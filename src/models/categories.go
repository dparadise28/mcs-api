package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	/*
		"encoding/json"
		"gopkg.in/mgo.v2/txn"
		"log"
	*///
)

var CategoryCollectionName = "Categories"
var TemplateCategoryCollectionName = "TemplateCategories"

type Category struct {
	ID       bson.ObjectId   `bson:"_id,omitempty" json:"category_id"`
	Name     string          `bson:"name" json:"name" validate:"required"`
	Root     bson.ObjectId   `bson:"root" json:"root"`
	Parent   bson.ObjectId   `bson:"parent" json:"parent"`
	Enabled  bool            `bson:"enabled" json:"enabled"`
	Children []bson.ObjectId `bson:"children" json:"children"`
	// TemplateCategoryID bson.ObjectId   `bson:"template_category_id" json:"template_category_id"`
	//	SortOrder uint16          `bson:"sort_order" json:"sort_order"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

// helper model for unpacking to avoid writing boilerplate
type StoreCategory struct {
	Name     string          `bson:"name" json:"name" validate:"required"`
	Children []StoreCategory `bson:"children" json:"children" validate:"dive"`
}

type StoreCategoryIds struct {
	// represents the active store category ids from the template
	// only categories in the template can be added to the store
	TemplateCategoryIds []bson.ObjectId `bson:"template_category_ids" json:"template_category_ids" validate:"required"`

	Store     Store         `bson:"-" json:"-" validate:"-"`
	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (cat *StoreCategoryIds) AddStoreCategories(storeID bson.ObjectId) error {
	c := cat.DB.C(TemplateCategoryCollectionName).With(cat.DBSession)
	count, _ := c.Find(bson.M{
		"_id": bson.M{
			"$in": cat.TemplateCategoryIds,
		},
	}).Count()
	if count != len(cat.TemplateCategoryIds) {
		return errors.New("Please only include available template categories")
		// , []Category{}
	}

	var store Store
	store.ID = storeID
	store.DB = cat.DB
	store.DBSession = cat.DBSession
	store.CategoryIds = cat.TemplateCategoryIds

	return store.AddCategories()
}

func (cat *Category) FindStoreTemplateTier1Categories() []Category {
	categories := []Category{}
	c := cat.DB.C(TemplateCategoryCollectionName).With(cat.DBSession)
	c.Find(bson.M{
		"enabled": true,
		"children": bson.M{
			"$gt": []bson.ObjectId{},
		},
	}).All(&categories)
	return categories
}

func (cat *Category) FindStoreTemplateTier2Categories(cid bson.ObjectId) []Category {
	categories := []Category{}
	c := cat.DB.C(TemplateCategoryCollectionName).With(cat.DBSession)
	c.Find(bson.M{
		"enabled": true,
		"parent":  cid,
	}).All(&categories)
	return categories
}

func (cat *Category) AddStoreCategories(categories []StoreCategory, storeID bson.ObjectId, template bool) (error, []Category) {
	cats := []Category{}
	collectionMap := map[bool]string{
		true:  TemplateCategoryCollectionName,
		false: CategoryCollectionName,
	}
	for _, catT1 := range categories {
		newParent := Category{
			ID:       bson.NewObjectId(),
			Name:     catT1.Name,
			Root:     storeID,
			Parent:   storeID,
			Enabled:  true,
			Children: []bson.ObjectId{},
		}
		for _, catT2 := range catT1.Children {
			newID := bson.NewObjectId()
			newParent.Children = append(newParent.Children, newID)
			cats = append(cats, Category{
				ID:       newID,
				Name:     catT2.Name,
				Root:     storeID,
				Parent:   newParent.ID,
				Enabled:  true,
				Children: []bson.ObjectId{},
			})
		}
		cats = append(cats, newParent)
	}

	c := cat.DB.C(collectionMap[template]).With(cat.DBSession)
	if insert_err := c.Insert(I(cats)...); insert_err != nil {
		if strings.Count(insert_err.Error(), "ObjectId") == 1 {
			return nil, cats
		}
		return insert_err, []Category{}
	}
	return nil, cats
}

/*
type ReadOnlyStoreCategory struct {
	ID                 bson.ObjectId   `bson:"_id,omitempty" json:"category_id"`
	Name               string          `bson:"name" json:"name" validate:"required"`
	Root               bson.ObjectId   `bson:"root" json:"root"`
	Parent             bson.ObjectId   `bson:"parent" json:"parent"`
	Enabled            bool            `bson:"enabled" json:"enabled"`
	Children           []bson.ObjectId `bson:children" json:"children"`
	// Products           []Product       `bson:"products" json:"products" validate:"required"`
	// TemplateCategoryID bson.ObjectId   `bson:"parent" json:"parent"`
	//	StoreId            bson.ObjectId   `bson:"store_id" json:"store_id"`
	//	SortOrder          uint16          `bson:"sort_order" json:"sort_order"`
}

type CategoryOrder struct {
	SID  bson.ObjectId   `json:"store_id" validate:"required"`
	CIDS []bson.ObjectId `json:"category_ids" validate:"required"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (cat *Category) AddStoreCategoriesFromTemplate(categories []Category, storeID bson.ObjectId) (error, []Category) {
	c := cat.DB.C(CategoryCollectionName).With(cat.DBSession)
	ct := cat.DB.C(TemplateCategoryCollectionName).With(cat.DBSession)

	cats := []Category{}
	idMap := map[bson.ObjectId]bson.ObjectId{}
	idList := []bson.ObjectId{}
	for _, catT1 := range categories {
		idMap[catT1.ID] = bson.NewObjectId()
		idList = append(idList, catT1.ID)
		parent := catT1.Parent
		if catT1.Root == catT1.Parent {
			parent = storeID
		}
		cats = append(cats, Category{
			ID:                 idMap[catT1.ID],
			Name:               catT1.Name,
			Root:               storeID,
			Parent:             parent, // catT1.Parent,
			Enabled:            false,
			Children:           catT1.Children,
			TemplateCategoryID: catT1.ID,
		})
	}
	count, _ := ct.Find(bson.M{
		"_id": bson.M{
			"$in": idList,
		},
	}).Count()
	if count != len(idList) {
		return errors.New("Please only include available template categories"), []Category{}
	}
	for catT1Index, catT1 := range categories {
		if catT1.Root != catT1.Parent {
			cats[catT1Index].Parent = idMap[cats[catT1Index].Parent]
		} else {
			for subCatIndex, subCatId := range cats[catT1Index].Children {
				cats[catT1Index].Children[subCatIndex] = idMap[subCatId]
			}
		}
	}
	if insert_err := c.Insert(I(cats)...); insert_err != nil {
		if strings.Count(insert_err.Error(), "ObjectId") == 1 {
			return nil, cats
		}
		return insert_err, []Category{}
	}
	return nil, cats
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
	// cur_seq_len, _ := c.Find(bson.M{"store_id": cat.StoreId}).Count()
	cat.ID = bson.NewObjectId()
	// cat.SortOrder = uint16(cur_seq_len) + 1
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
		"_id":  cat.ID,
		"root": cat.Root,
	}).Apply(change, cat)
	log.Println(info)
}

func (cat *Category) FindStoreTier1Categories(storeID bson.ObjectId, enabled bool) []Category {
	categories := []Category{}
	c := cat.DB.C(CategoryCollectionName).With(cat.DBSession)
	c.Find(bson.M{
		"root":    storeID,
		"enabled": enabled,
		"children": bson.M{
			"$gt": []bson.ObjectId{},
		},
	}).All(&categories)
	return categories
}

func (cat *Category) FindStoreTier2Categories(enabled bool, cid, storeID bson.ObjectId) []Category {
	categories := []Category{}
	c := cat.DB.C(CategoryCollectionName).With(cat.DBSession)
	c.Find(bson.M{
		"root":    storeID,
		"parent":  cid,
		"enabled": enabled,
	}).All(&categories)
	return categories
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
		"_id":  cat.ID,
		"root": cat.Root,
	}).Apply(change, cat)
	//log.Println(cat)
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
	// c.Remove(bson.M{"_id": previousTnx.PreviousTnxOp})
	return run_err
}
*/
