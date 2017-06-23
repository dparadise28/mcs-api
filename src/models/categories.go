package models

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var CategoryCollectionName = "Categories"

type Category struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"category_id"`
	Name      string        `bson:"name" json:"name" validate:"required"`
	StoreId   bson.ObjectId `bson:"store_id" json:"store_id" validate:"required"`
	SortOrder uint16        `bson:"sort_order" json:"sort_order"`
	Enabled   bool          `bson:"enabled" json:"enabled"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

// helper model for unpacking to avoid writing boilerplate
type StoreCategory struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"category_id"`
	Name      string        `bson:"name" json:"name" validate:"required"`
	StoreId   bson.ObjectId `bson:"store_id" json:"store_id"`
	Enabled   bool          `bson:"enabled" json:"enabled"`
	SortOrder uint16        `bson:"sort_order" json:"sort_order"`
	Products  []Product     `bson:"-" json:"products" validate:"required"`
}

type ReadOnlyStoreCategory struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"category_id"`
	Name      string        `bson:"name" json:"name" validate:"required"`
	StoreId   bson.ObjectId `bson:"store_id" json:"store_id"`
	Enabled   bool          `bson:"enabled" json:"enabled"`
	SortOrder uint16        `bson:"sort_order" json:"sort_order"`
	Products  []Product     `bson:"products" json:"products" validate:"required"`
}

func (cat *Category) RetrieveFullCategoriesByStoreID(id string, enabled_only_categories bool, enabled_only_products bool) (error, []ReadOnlyStoreCategory) {
	enabled := []bool{true}
	if !enabled_only_categories {
		enabled = []bool{true, false}
	}
	pipeline := []bson.M{
		bson.M{
			"$match": bson.M{
				"store_id": bson.ObjectIdHex(id),
				"enabled": bson.M{
					"$in": enabled,
				},
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path": "$products",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"from":         ProductCollectionName,
				"localField":   "_id",
				"foreignField": "category_id",
				"as":           "products",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path": "$products",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$sort": bson.M{
				"sort_order":          1,
				"products.sort_order": 1,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":        "$_id",
				"products":   bson.M{"$push": "$products"},
				"name":       bson.M{"$first": "$name"},
				"sort_order": bson.M{"$first": "$sort_order"},
			},
		},
	}
	c := cat.DB.C(CategoryCollectionName).With(cat.DBSession)
	pipe := c.Pipe(pipeline)
	resp := []ReadOnlyStoreCategory{}
	if err := pipe.All(&resp); err != nil {
		return err, resp
	}
	if !enabled_only_products {
		for cat_index, _ := range resp {
			for prod_index, _ := range resp[cat_index].Products {
				if !resp[cat_index].Products[prod_index].Enabled {
					// filter all disabled products
					resp[cat_index].Products = append(resp[cat_index].Products[:prod_index], resp[cat_index].Products[prod_index+1:]...)
				}
			}
		}
	}
	return nil, resp
}

func (cat *Category) AddStoreCategory() error {
	//var cur_seq_len uint16
	c := cat.DB.C(CategoryCollectionName).With(cat.DBSession)

	cur_seq_len, _ := c.Find(bson.M{"store_id": cat.StoreId}).Count()
	cat.ID = bson.NewObjectId()
	cat.Enabled = true
	cat.SortOrder = uint16(cur_seq_len) + 1
	if insert_err := c.Insert(&cat); insert_err != nil {
		return insert_err
	}
	return nil
}

func (cat *Category) UpdateStoreCategoryName() {
	c := cat.DB.C(CategoryCollectionName).With(cat.DBSession)
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{"name": cat.Name},
		},
	}
	info, _ := c.Find(bson.M{
		"_id":      cat.ID,
		"store_id": cat.StoreId,
	}).Apply(change, cat)
	log.Println(info)
}

func (cat *Category) ActivateStoreCategory() error {
	c := cat.DB.C(CategoryCollectionName).With(cat.DBSession)
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{"enabled": cat.Enabled},
		},
	}
	info, err := c.Find(bson.M{
		"_id":      cat.ID,
		"store_id": cat.StoreId,
	}).Apply(change, cat)
	log.Println(info)
	return err
}
