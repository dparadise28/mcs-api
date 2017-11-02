package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
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

	CIDS      []bson.ObjectId `bson:"-" json:"-"`
	CID       bson.ObjectId   `bson:"-" json:"-"`
	SID       bson.ObjectId   `bson:"-" json:"-"`
	DB        *mgo.Database   `bson:"-" json:"-"`
	DBSession *mgo.Session    `bson:"-" json:"-"`
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

func (cat *Category) RetrieveT1StoreCategoryIdsFromProducts() {
	c := cat.DB.C(ProductCollectionName).With(cat.DBSession)
	c.Find(bson.M{"store_id": cat.SID}).Distinct("category_id", &cat.CIDS)
}

func (cat *Category) FindStoreTier2Categories() []Category {
	categories := []Category{}
	c := cat.DB.C(TemplateCategoryCollectionName).With(cat.DBSession)

	// grab t1 ids through all the products in the store
	cat.RetrieveT1StoreCategoryIdsFromProducts()
	c.Find(bson.M{
		"parent": cat.CID,
		"_id": bson.M{
			"$in": cat.CIDS,
		},
	}).All(&categories)
	return categories
}

func (cat *Category) FindStoreTier1Categories() []Category {
	categories := []Category{}
	c := cat.DB.C(TemplateCategoryCollectionName).With(cat.DBSession)

	// grab t1 ids through all the products in the store
	cat.RetrieveT1StoreCategoryIdsFromProducts()
	c.Find(bson.M{
		"_id": bson.M{
			"$in": cat.CIDS,
		},
	}).Distinct("parent", &cat.CIDS)
	c.Find(bson.M{
		"_id": bson.M{
			"$in": cat.CIDS,
		},
	}).All(&categories)
	return categories
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
