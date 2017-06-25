package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
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
	ID            bson.ObjectId `bson:"_id,omitempty" json:"category_id"`
	Name          string        `bson:"name" json:"name" validate:"required"`
	StoreId       bson.ObjectId `bson:"store_id" json:"store_id"`
	Enabled       bool          `bson:"enabled" json:"enabled"`
	SortOrder     uint16        `bson:"sort_order" json:"sort_order"`
	Products      []Product     `bson:"products" json:"products" validate:"required"`
	PreviousTnxOp bson.ObjectId `bson:"previoustnxop" json:"-"`
}

type CategoryOrder struct {
	SID  bson.ObjectId   `json:"store_id" validate:"required"`
	CIDS []bson.ObjectId `json:"category_ids" validate:"required"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
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
				"products.sort_order": 1,
			},
		},
		bson.M{
			"$group": bson.M{
				"_id":        "$_id",
				"products":   bson.M{"$push": "$products"},
				"store_id":   bson.M{"$first": "$store_id"},
				"enabled":    bson.M{"$first": "$enabled"},
				"name":       bson.M{"$first": "$name"},
				"sort_order": bson.M{"$first": "$sort_order"},
			},
		},
		bson.M{
			"$sort": bson.M{
				"sort_order": 1,
			},
		},
	}
	c := cat.DB.C(CategoryCollectionName).With(cat.DBSession)
	pipe := c.Pipe(pipeline)
	resp := []ReadOnlyStoreCategory{}
	if err := pipe.All(&resp); err != nil {
		return err, resp
	}

	if enabled_only_products {
		for cat_index, _ := range resp {
			removed := 0
			for prod_index, _ := range resp[cat_index].Products {
				// filter all disabled products
				new_pindex := prod_index - removed
				if !resp[cat_index].Products[new_pindex].Enabled {
					if len(resp[cat_index].Products) <= new_pindex {
						resp[cat_index].Products = resp[cat_index].Products[:new_pindex]
					} else {
						resp[cat_index].Products = append(resp[cat_index].Products[:new_pindex], resp[cat_index].Products[new_pindex+1:]...)
						removed += 1
					}
				}
			}
		}
	}
	return nil, resp
}

func (cat *Category) AddStoreCategory() error {
	c := cat.DB.C(CategoryCollectionName).With(cat.DBSession)

	// might want to avoid lookup and just set sort index to staticaly
	// defined cap on number of categories
	cur_seq_len, _ := c.Find(bson.M{"store_id": cat.StoreId}).Count()
	cat.ID = bson.NewObjectId()
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
	log.Println(cat)
	log.Println(info)
	return err
}

func (cat *CategoryOrder) ReorderStoreCategories() error {
	c := cat.DB.C(CategoryCollectionName).With(cat.DBSession)
	cat_count, err := c.Find(bson.M{
		"$and": []bson.M{
			bson.M{
				"store_id": cat.SID,
				"enabled":  true,
			}, bson.M{
				"_id": bson.M{
					"$in": cat.CIDS,
				},
			},
		},
	}).Count()
	if err != nil || cat_count != len(cat.CIDS) {
		return errors.New("All active category ids must be provided. You may not include any categories not currently enabled.")
	}

	previousTnx := ReadOnlyStoreCategory{}
	c.Find(bson.M{"_id": cat.CIDS[0]}).One(&previousTnx)

	op_id := bson.NewObjectId()
	operations := make([]txn.Op, len(cat.CIDS))
	for index, cid := range cat.CIDS {
		operations[index] = txn.Op{
			C:  CategoryCollectionName,
			Id: cid,
			Update: bson.M{
				"$set": bson.M{
					"sort_order":    index,
					"previoustnxop": op_id,
				},
			},
		}
	}
	runner := txn.NewRunner(c)
	run_err := runner.Run(operations, op_id, nil)
	c.Remove(bson.M{"_id": previousTnx.PreviousTnxOp})
	return run_err
}
