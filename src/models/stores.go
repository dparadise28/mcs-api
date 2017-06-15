package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"reflect"
	"strings"
)

var StoreCollectionName = "Stores"

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
	TaxRate         float64            `json:"tax_rate" validate:"required"`
	Delivery        StoreDelivery      `json:"delivery"`
	Distance        float64            `bson:"distance,omitempty" json:"distance,omitempty"`
	WorkingHours    WeeklyWorkingHours `json:"working_hours"`
	LongDescription string             `json:"long_description"`
	// this field has has a fulltext index for
	// full text search so must so we must
	// ensure its length for now to avoid index
	// bloating untill switching to a more robust
	// search solution or building one
	ShortDescription string `json:"short_description" validate:"max=50"`

	CategoryNames []string        `bson:"c_names" json:"category_names"`
	CategoryIDs   []bson.ObjectId `bson:"c_ids" json:"category_ids"`
	CTree         []StoreCategory `bson:"-" json:"category_tree" validate:"required"`

	products []Product
	//PlatformCategories []string        `json:"platform_categories"`

	DB        *mgo.Database
	DBSession *mgo.Session
}

func (s *Store) PrepStoreEntitiesForInsert() {
	s.ID = bson.NewObjectId()
	s.CategoryNames = []string{}
	for category_index, _ := range s.CTree {
		c_id := bson.NewObjectId()
		s.CTree[category_index].ID = c_id
		s.CTree[category_index].StoreId = s.ID

		s.CategoryIDs = append(s.CategoryIDs, c_id)
		s.CategoryNames = append(s.CategoryNames, s.CTree[category_index].Name)
		for product_index, _ := range s.CTree[category_index].Products {
			p_id := bson.NewObjectId()
			s.CTree[category_index].Products[product_index].ID = p_id
			s.CTree[category_index].Products[product_index].StoreID = s.ID
			s.CTree[category_index].Products[product_index].CategoryID = c_id

			s.CTree[category_index].ProductIDS = append(s.CTree[category_index].ProductIDS, p_id)
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
