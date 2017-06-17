package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"reflect"
	"strings"
)

const StoreCollectionName = "Stores"

var MAX_DISTANCE = 1609.34 * 2 // max distance is static for now (2 mi)

type OpenHours struct {
	Hours  Hours `json:"hours"`
	IsOpen bool  `json:"open"`
}

type WeeklyWorkingHours struct {
	Sun OpenHours `bson:"sun" json:"sunday"`
	Mon OpenHours `bson:"mon" json:"monday"`
	Tue OpenHours `bson:"tue" json:"tuesday"`
	Wed OpenHours `bson:"wed" json:"wednesday"`
	Thu OpenHours `bson:"thu" json:"thursday"`
	Fri OpenHours `bson:"fri" json:"friday"`
	Sat OpenHours `bson:"sat" json:"saturday"`
}

type StoreDelivery struct {
	// we might want to offer a variable delivery
	// fee model at some point but this will do
	// for now. That can be split out form the base
	// store model
	Fee       uint   `bson:"fee,omitempty" json:"delivery_fee"`
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
	} `bson:"pickup_items" json:"pickup_items"`
}

type Store struct {
	ID              bson.ObjectId      `bson:"_id,omitempty" json:"store_id"`
	Name            string             `bson:"name" json:"name"`
	Image           string             `json:"image"`
	Phone           string             `json:"phone"`
	Pickup          StorePickup        `json:"pickup"`
	Address         Address            `json:"address"`
	TaxRate         float64            `bson:"tax_rate" json:"tax_rate" validate:"required"`
	Delivery        StoreDelivery      `json:"delivery"`
	Distance        float64            `bson:"distance,omitempty" json:"distance,omitempty"`
	WorkingHours    WeeklyWorkingHours `bson:"working_hours" json:"working_hours"`
	LongDescription string             `bson:"long_desc" json:"long_description"`
	// this field has has a fulltext index for
	// full text search so must so we must
	// ensure its length for now to avoid index
	// bloating untill switching to a more robust
	// search solution or building one
	ShortDescription string `bson:"short_desc" json:"short_description" validate:"max=50"`

	CategoryNames []string        `bson:"c_names" json:"category_names"`
	CTree         []StoreCategory `bson:"-" json:"category_tree" validate:"required"`

	products []Product `bson:"-" json:"-"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (s *Store) PrepStoreEntitiesForInsert() {
	s.ID = bson.NewObjectId()
	s.CategoryNames = []string{}
	for category_index, _ := range s.CTree {
		c_id := bson.NewObjectId()

		s.CTree[category_index].ID = c_id
		s.CTree[category_index].StoreId = s.ID
		s.CTree[category_index].SortOrder = uint16(category_index)

		s.CategoryNames = append(s.CategoryNames, s.CTree[category_index].Name)
		for product_index, _ := range s.CTree[category_index].Products {
			p_id := bson.NewObjectId()

			s.CTree[category_index].Products[product_index].ID = p_id
			s.CTree[category_index].Products[product_index].StoreID = s.ID
			s.CTree[category_index].Products[product_index].CategoryID = c_id
			s.CTree[category_index].Products[product_index].SortOrder = uint16(product_index)

			s.products = append(s.products, s.CTree[category_index].Products[product_index])
		}
	}
}

func I(array interface{}) []interface{} {
	v := reflect.ValueOf(array)
	t := v.Type()

	if t.Kind() != reflect.Slice {
		log.Panicf("`array` should be %s but got %s", reflect.Slice, t.Kind())
	}

	result := make([]interface{}, v.Len(), v.Len())

	for i := 0; i < v.Len(); i++ {
		result[i] = v.Index(i).Interface()
	}

	return result
}

func (s *Store) Insert() error {
	c := s.DB.C(StoreCollectionName).With(s.DBSession)
	s.PrepStoreEntitiesForInsert()
	s.Address.Location.Type = "Point"
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
		if strings.Count(insert_err.Error(), "ObjectId") == 1 {
			log.Println("during products")
			log.Println(insert_err)
			return nil
		}
		return insert_err
	}
	return s.InsertStoreProducts()
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
				"categories":        bson.M{"$push": "$categories"},
				"working_hours":     bson.M{"$first": "$working_hours"},
				"image":             bson.M{"$first": "$image"},
				"delivery":          bson.M{"$first": "$delivery"},
				"phone":             bson.M{"$first": "$phone"},
				"tax_rate":          bson.M{"$first": "$tax_rate"},
				"address":           bson.M{"$first": "$address"},
				"long_description":  bson.M{"$first": "$long_desc"},
				"distance":          bson.M{"$first": "$distance"},
				"name":              bson.M{"$first": "$name"},
				"pickup":            bson.M{"$first": "$pickup"},
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

func (s *Store) FindStoresByLocation(long float64, lat float64, maxDist float64, time int) (error, []bson.M) {
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
	stores := []bson.M{}
	err := c.Find(query).All(&stores)
	return err, stores
}
