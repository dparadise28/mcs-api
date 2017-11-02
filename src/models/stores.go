package models

import (
	"encoding/json"
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
	Offered   bool   `bson:"offered" json:"service_offered"`
	MaxDist   int    `bson:"max_distance,omitempty" json:"delivery_distance"`
	MinTime   uint8  `bson:"min_time,omitempty" json:"minimum_time_to_delivery"`
	MaxTime   uint8  `bson:"max_time,omitempty" json:"maximum_time_to_delivery"`
	MinAmount uint16 `bson:"min_amount,omitempty" json:"delivery_minimum"`
}

type StorePickup struct {
	Offered         bool  `json:"offered"`
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
	CategoryIds      []bson.ObjectId     `bson:"category_ids" json:"category_ids"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type StoreReturn struct {
	ID               bson.ObjectId       `bson:"_id,omitempty" json:"store_id"`
	Name             string              `bson:"name" json:"name"`
	Image            string              `json:"image"`
	Phone            string              `json:"phone" validate:"required"`
	CTree            []StoreCategory     `bson:"categories" json:"categories" validate:"required,dive"`
	Email            string              `json:"email" validate:"required,email"`
	Pickup           StorePickup         `json:"pickup" validate:"required,dive"`
	Address          Address             `json:"address" validate:"required,dive"`
	TaxRate          float64             `bson:"tax_rate" json:"tax_rate" validate:"required"`
	Enabled          bool                `bson:"enabled" json:"enabled"`
	Delivery         StoreDelivery       `json:"delivery" validate:"dive"`
	Distance         float64             `bson:"distance,omitempty" json:"distance"`
	ReviewScore      int64               `bson:"review_score" json:"review_score"`
	ReviewCount      int64               `bson:"review_count" json:"review_count"`
	WorkingHours     WeeklyWorkingHours  `bson:"working_hours" json:"working_hours" validate:"required,dive"`
	PaymentDetails   StorePaymentDetails `bson:"payment_details" json:"payment_details" validate:"-"`
	LongDescription  string              `bson:"long_desc" json:"long_description" validate:"max=200"`
	ShortDescription string              `bson:"short_desc" json:"short_description" validate:"max=50"`
	CategoryIds      []bson.ObjectId     `bson:"category_ids" json:"category_ids"`
}

type StoreInfo struct {
	ID               bson.ObjectId      `bson:"_id,omitempty" json:"store_id"`
	Name             string             `bson:"name" json:"name"`
	Image            string             `json:"image"`
	Phone            string             `json:"phone" validate:"required"`
	Email            string             `json:"email" validate:"required,email"`
	Pickup           StorePickup        `json:"pickup" validate:"required,dive"`
	Address          Address            `json:"address" validate:"required,dive"`
	TaxRate          float64            `bson:"tax_rate" json:"tax_rate" validate:"required"`
	Delivery         StoreDelivery      `json:"delivery" validate:"dive"`
	WorkingHours     WeeklyWorkingHours `bson:"working_hours" json:"working_hours" validate:"required,dive"`
	LongDescription  string             `bson:"long_desc" json:"long_description" validate:"max=200"`
	ShortDescription string             `bson:"short_desc" json:"short_description" validate:"max=50"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (s *Store) PrepStoreEntitiesForInsert() error {
	s.PaymentDetails.AcceptsCCPayment = false
	s.PaymentDetails.AcceptsCashPayment = true
	s.Enabled = false

	s.ID = bson.NewObjectId()
	s.Address.ID = bson.NewObjectId()
	s.Address.UserID = s.ID
	s.PaymentDetails.LegalEntity.BillingAddress.ID = s.ID
	s.PaymentDetails.LegalEntity.BillingAddress.UserID = s.ID
	s.PaymentDetails.StoreID = s.ID
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
		log.Println(insert_err)
		st, _ := json.Marshal(s)
		log.Println(string(st))
		// dont move on if store record is faulty
		if strings.Count(insert_err.Error(), "ObjectId") == 1 {
			log.Println("during store")
			log.Println(insert_err)
			return nil
		}
		return insert_err
	}
	return nil // s.InsertStoreCategories()
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
		"_id": s.ID,
	}
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	err := c.Find(query).One(s)
	return err
}

func (s *Store) AddCategories() error {
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{
				"category_ids": s.CategoryIds,
			},
		},
	}
	_, err := c.Find(bson.M{
		"_id": s.ID,
	}).Apply(change, s)
	return err
}

func (s *Store) FindStoresByLocation(long float64, lat float64, maxDist float64, time int) (error, []Store) {
	query := bson.M{
		"enabled": true,
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

func (s *StoreInfo) UpdateStoreInfo() error {
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	s.Address.Location.Type = "Point"
	s.Address.Location.Coordinates = []float64{s.Address.Longitude, s.Address.Latitude}
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
				"long_desc":     s.LongDescription,
				"short_desc":    s.ShortDescription,
				"working_hours": s.WorkingHours,
			},
		},
	}
	_, err := c.Find(bson.M{
		"_id": s.ID,
	}).Apply(change, s)
	//log.Println(info)
	return err
}

func (s *Store) AddStoreCCPaymentMethod() error {
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
	_, err := c.Find(bson.M{
		"_id": s.ID,
	}).Apply(change, s)
	return err
}
