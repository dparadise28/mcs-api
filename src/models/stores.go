package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

const StoreCollectionName = "Stores"

var MAX_DISTANCE = 1609.34 * 2 // max distance is static for now (2 mi)

type OpenHours struct {
	Hours  Hours `json:"hours"  validate:"dive"`
	IsOpen bool  `json:"open" validate:"required"`
}

type WeeklyWorkingHours struct {
	Sun OpenHours `bson:"sun" json:"sunday" validate:"required,dive"`
	Mon OpenHours `bson:"mon" json:"monday" validate:"required,dive"`
	Tue OpenHours `bson:"tue" json:"tuesday" validate:"required,dive"`
	Wed OpenHours `bson:"wed" json:"wednesday" validate:"required,dive"`
	Thu OpenHours `bson:"thu" json:"thursday" validate:"required,dive"`
	Fri OpenHours `bson:"fri" json:"friday" validate:"required,dive"`
	Sat OpenHours `bson:"sat" json:"saturday" validate:"required,dive"`
}

type StoreDelivery struct {
	// we might want to offer a variable delivery
	// fee model at some point but this will do
	// for now. That can be split out form the base
	// store model
	Fee       uint32 `bson:"fee,omitempty" json:"delivery_fee"`
	Offered   bool   `bson:"offered" json:"service_offered" validate:"required"`
	MaxDist   int    `bson:"max_distance,omitempty" json:"delivery_distance"`
	MinTime   uint8  `bson:"min_time,omitempty" json:"maximum_time_to_delivery"`
	MaxTime   uint8  `bson:"max_time,omitempty" json:"minimum_time_to_delivery"`
	MinAmount uint16 `bson:"min_amount,omitempty" json:"delivery_minimum"`
}

type StorePickup struct {
	Offered         bool  `json:"offered" validate:"required"`
	MinTime         uint8 `bson:"min_time" json:"minimum_time_to_pickup" validate:"max=90,min=1"`
	MaxTime         uint8 `bson:"max_time" json:"maximum_time_to_pickup" validate:"max=90,min=1"`
	PickupItemCount struct {
		Min uint8  `json:"min" validate:"max=255,min=1"`
		Max uint32 `json:"max" validate:"max=4294967295,min=1"`
	} `bson:"pickup_items" json:"pickup_items,dive"`
}

type LegalEntity struct {
	BillingAddress Address `json:"billing_address" bson:"billing_address" validate:"required,dive"`
	BusinessTaxID  string  `json:"business_tax_id" bson:"-"`
	BusinessName   string  `json:"legal_business_name" bson:"legal_business_name" validate:"required"`
	PersonalID     string  `json:"personal_id" bson:"-"`
	SSNLast4       string  `json:"last_4_ssn" bson:"-" validate:"max=4,min=4"`
	Owner          struct {
		First string `validate:"required"`
		Last  string `validate:"required"`
		DOB   struct {
			Day   uint8  `json:"day" validate:"required"`
			Month uint8  `json:"month" validate:"required"`
			Year  uint16 `json:"year" validate:"required"`
		} `validate:"required,dive"`
	} `validate:"required,dive"`
}

type StorePaymentDetails struct {
	StoreID            bson.ObjectId `bson:"store_id,omitempty" json:"store_id"`
	LegalEntity        LegalEntity   `json:"legal_entity" bson:"legal_entity" validate:"dive"`
	BusinessType       string        `json:"business_type" bson:"business_type" validate:"required"`
	StripeAccountID    string        `bson:"stripe_custom_account_id" json:"-"`
	AcceptsCCPayment   bool          `bson:"cc_payment_available" json:"cc_payment_available"`
	AcceptsCashPayment bool          `bson:"cash_payment_available" json:"cash_payment_available"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type Store struct {
	ID              bson.ObjectId      `bson:"_id,omitempty" json:"store_id"`
	Name            string             `bson:"name" json:"name"`
	Image           string             `json:"image"`
	Phone           string             `json:"phone" validate:"required"`
	Email           string             `json:"email" validate:"required,email"`
	Pickup          StorePickup        `json:"pickup" validate:"required,dive"`
	Address         Address            `json:"address" validate:"required,dive"`
	TaxRate         float64            `bson:"tax_rate" json:"tax_rate" validate:"required"`
	Enabled         bool               `bson:"enabled" json:"enabled"`
	Delivery        StoreDelivery      `json:"delivery" validate:"dive"`
	Distance        float64            `bson:"distance,omitempty" json:"distance"`
	WorkingHours    WeeklyWorkingHours `bson:"working_hours" json:"working_hours" validate:"required,dive"`
	LongDescription string             `bson:"long_desc" json:"long_description" validate:"max=200"`
	// this field has has a fulltext index for
	// full text search so must so we must
	// ensure its length for now to avoid index
	// bloating untill switching to a more robust
	// search solution or building one
	ShortDescription string              `bson:"short_desc" json:"short_description" validate:"max=50"`
	PaymentDetails   StorePaymentDetails `bson:"payment_details" json:"payment_details" validate:"-"`
	ReviewScore      int64               `bson:"review_score" json:"review_score"`
	ReviewCount      int64               `bson:"review_count" json:"review_count"`

	CategoryNames []string        `bson:"c_names" json:"category_names"`
	CTree         []StoreCategory `bson:"-" json:"categories" validate:"required,dive"`

	products []Product `bson:"-" json:"-"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (s *Store) PrepStoreEntitiesForInsert() error {
	s.PaymentDetails.AcceptsCCPayment = false
	s.PaymentDetails.AcceptsCashPayment = true
	s.Enabled = false

	s.ID = bson.NewObjectId()
	s.CategoryNames = []string{}
	for category_index, _ := range s.CTree {
		c_id := bson.NewObjectId()

		s.CTree[category_index].ID = c_id
		s.CTree[category_index].StoreId = s.ID
		s.CTree[category_index].Enabled = true
		s.CTree[category_index].SortOrder = uint16(category_index)

		s.CategoryNames = append(s.CategoryNames, s.CTree[category_index].Name)
		for product_index, _ := range s.CTree[category_index].Products {
			if s.CTree[category_index].Products[product_index].AssetID.Hex() == "" {
				return errors.New("Must provide a valid asset_id for every product.")
			}
			p_id := bson.NewObjectId()

			s.CTree[category_index].Products[product_index].ID = p_id
			s.CTree[category_index].Products[product_index].Enabled = true
			s.CTree[category_index].Products[product_index].StoreID = s.ID
			s.CTree[category_index].Products[product_index].CategoryID = c_id
			s.CTree[category_index].Products[product_index].SortOrder = uint16(product_index)

			s.products = append(s.products, s.CTree[category_index].Products[product_index])
		}
	}
	return nil
}

func (s *Store) Insert() error {
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	if err := s.PrepStoreEntitiesForInsert(); err != nil {
		return err
	}
	s.Address.Location.Type = "Point"
	s.Address.Location.Coordinates = []float64{s.Address.Longitude, s.Address.Latitude}
	if insert_err := c.Insert(&s); insert_err != nil {
		// dont move on if store record is faulty
		if strings.Count(insert_err.Error(), "ObjectId") == 1 {
			log.Println("during store")
			log.Println(insert_err)
			return nil
		}
		return insert_err
	}
	return s.InsertStoreCategories()
}

func (s *Store) InsertStoreCategories() error {
	c := s.DB.C(CategoryCollectionName).With(s.DBSession)
	if insert_err := c.Insert(I(s.CTree)...); insert_err != nil {
		if strings.Count(insert_err.Error(), "ObjectId") == 1 {
			log.Println("during categories")
			log.Println(insert_err)
			return nil
		}
		return insert_err
	}
	return s.InsertStoreProducts()
}

func (s *Store) InsertStoreProducts() error {
	c := s.DB.C(ProductCollectionName).With(s.DBSession)
	if insert_err := c.Insert(I(s.products)...); insert_err != nil {
		log.Println("during products")
		log.Println(insert_err)
		if strings.Count(insert_err.Error(), "ObjectId") == 1 {
			return nil
		}
		return insert_err
	}
	return s.InsertStoreProducts()
}

func (s *Store) RetrieveStoreByID(id string) error {
	query := bson.M{
		"_id": bson.ObjectIdHex(id),
	}
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	err := c.Find(query).One(s)
	return err
}

func (s *Store) RetrieveStoreByOID() error {
	query := bson.M{
		"_id": store.ID,
	}
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	err := c.Find(query).One(s)
	return err
}

func (s *Store) RetrieveFullStoreByID(id string) (error, bson.M) {
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"_id": bson.ObjectIdHex(id),
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         CategoryCollectionName,
				"localField":   "_id",
				"foreignField": "store_id",
				"as":           "categories",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path": "$categories",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$sort": bson.M{
				"categories.sort_order": 1,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         ProductCollectionName,
				"localField":   "categories._id",
				"foreignField": "category_id",
				"as":           "categories.products",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":               "$_id",
				"name":              bson.M{"$first": "$name"},
				"image":             bson.M{"$first": "$image"},
				"phone":             bson.M{"$first": "$phone"},
				"pickup":            bson.M{"$first": "$pickup"},
				"address":           bson.M{"$first": "$address"},
				"delivery":          bson.M{"$first": "$delivery"},
				"tax_rate":          bson.M{"$first": "$tax_rate"},
				"distance":          bson.M{"$first": "$distance"},
				"categories":        bson.M{"$push": "$categories"},
				"working_hours":     bson.M{"$first": "$working_hours"},
				"payment_details":   bson.M{"$first": "$payment_details"},
				"long_description":  bson.M{"$first": "$long_desc"},
				"short_description": bson.M{"$first": "$short_desc"},
			},
		},
		bson.M{
			"$sort": bson.M{
				"categories.sort_order":          1,
				"categories.products.sort_order": 1,
			},
		},
	}
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	pipe := c.Pipe(pipeline)
	resp := []bson.M{}
	err := pipe.All(&resp)
	return err, resp[0]
}

func (s *Store) FindStoresByLocation(long float64, lat float64, maxDist float64, time int) (error, []Store) {
	query := bson.M{
		"address.location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
				"$maxDistance": maxDist,
			},
		},
	}
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	stores := []Store{}
	err := c.Find(query).All(&stores)
	for index, store := range stores {
		stores[index].Distance = Distance(
			store.Address.Latitude,
			store.Address.Longitude,
			lat, long,
		)
	}
	return err, stores
}

func (s *Store) UpdateStoreInfo() {
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{
				"name":          s.Name,
				"image":         s.Image,
				"email":         s.Email,
				"phone":         s.Phone,
				"pickup":        s.Pickup,
				"address":       s.Address,
				"delivery":      s.Delivery,
				"tax_rate":      s.TaxRate,
				"distance":      s.Distance,
				"long_desc":     s.LongDescription,
				"short_desc":    s.ShortDescription,
				"working_hours": s.WorkingHours,
			},
		},
	}
	info, _ := c.Find(bson.M{
		"_id": s.ID,
	}).Apply(change, s)
	log.Println(info)
}

func (s *Store) AddStoreCCPaymentMethod() {
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	s.PaymentDetails.AcceptsCCPayment = true
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{
				"payment_details": s.PaymentDetails,
			},
		},
	}
	info, _ := c.Find(bson.M{
		"_id": s.ID,
	}).Apply(change, s)
	log.Println(info)
}
