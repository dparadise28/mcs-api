package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"
	"log"
)

var ProductCollectionName = "Products"

type Product struct {
	ID             bson.ObjectId `bson:"_id" json:"product_id"`
	Image          bool          `bson:"image" json:"image"`
	Enabled        bool          `bson:"enabled" json:"enabled"`
	StoreID        bson.ObjectId `bson:"store_id" json:"store_id"`
	SortOrder      uint16        `bson:"sort_order" json:"sort_order"`
	PriceCents     uint32        `bson:"price_cents" json:"price_cents" validate:"required"`
	CategoryID     bson.ObjectId `bson:"category_id" json:"category_id"`
	Description    string        `bson:"desc" json:"description"`
	ProductTitle   string        `bson:"title" json:"title" validate:"required"`
	DisplayPrice   string        `bson:"-" json:"display_price"`
	NewCategoryID  bson.ObjectId `bson:"-" json:"new_category_id"`
	ProductRatings struct {
		ReviewCount           uint64  `bson:"review_count" json:"total_reviews"`
		ReviewPercentageScore float64 `bson:"pct_score" json:"review_percent"`
	}

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type ReadOnlyProduct struct {
	ID             bson.ObjectId `bson:"_id" json:"product_id"`
	Image          bool          `bson:"image" json:"image"`
	Enabled        bool          `bson:"enabled" json:"enabled"`
	StoreID        bson.ObjectId `bson:"store_id" json:"store_id"`
	SortOrder      uint16        `bson:"sort_order" json:"sort_order"`
	PriceCents     uint32        `bson:"price_cents" json:"price_cents" validate:"required"`
	CategoryID     bson.ObjectId `bson:"category_id" json:"category_id"`
	Description    string        `bson:"desc" json:"description"`
	ProductTitle   string        `bson:"title" json:"title" validate:"required"`
	DisplayPrice   string        `bson:"-" json:"display_price"`
	NewCategoryID  bson.ObjectId `bson:"-" json:"new_category_id"`
	PreviousTnxOp  bson.ObjectId `bson:"previoustnxop" json:"-"`
	ProductRatings struct {
		ReviewCount           uint64  `bson:"review_count" json:"total_reviews"`
		ReviewPercentageScore float64 `bson:"pct_score" json:"review_percent"`
	}

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type CartProduct struct {
	//StoreID      bson.ObjectId `bson:"-" json:"store_id"`
	ID           bson.ObjectId `bson:"_id" json:"product_id"`
	Image        bool          `bson:"image" json:"image"`
	Quantity     uint16        `bson:"qty" json:"quantity"`
	PriceCents   uint32        `bson:"price_cents" json:"price_cents"`
	ProductTitle string        `bson:"title" json:"title"`
	Instructions string        `bson:"instructions" json:"instructions"`
}

type ProductOrder struct {
	SID  bson.ObjectId   `json:"store_id" validate:"required"`
	CID  bson.ObjectId   `json:"category_id" validate:"required"`
	PIDS []bson.ObjectId `json:"product_ids" validate:"required"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type CartRequest struct {
	SID          bson.ObjectId `json:"store_id" validate:"required"`
	PID          bson.ObjectId `json:"product_id" validate:"required"`
	CID          bson.ObjectId `json:"cart_id"`
	QTY          uint16        `json:"quantity" validate:"required"`
	IsNewCart    bool          `json:"is_new_cart"`
	Instructions string        `json:"instructions"`
}

func (p *Product) AddProductToStoreCategory() error {
	c := p.DB.C(ProductCollectionName).With(p.DBSession)

	// might want to avoid lookup and just set sort index to staticaly
	// defined cap on number of categories
	cur_seq_len, _ := c.Find(bson.M{
		"store_id":    p.StoreID,
		"category_id": p.CategoryID,
	}).Count()
	p.ID = bson.NewObjectId()
	p.SortOrder = uint16(cur_seq_len) + 1
	if insert_err := c.Insert(&p); insert_err != nil {
		return insert_err
	}
	return nil
}

func (p *Product) ActivateStoreProduct() error {
	c := p.DB.C(ProductCollectionName).With(p.DBSession)
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{"enabled": p.Enabled},
		},
	}
	info, err := c.Find(bson.M{
		"_id":         p.ID,
		"store_id":    p.StoreID,
		"category_id": p.CategoryID,
	}).Apply(change, p)
	log.Println(info)
	return err
}

func (p *Product) UpdateStoreProduct() {
	c := p.DB.C(ProductCollectionName).With(p.DBSession)
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{
				"desc":        p.Description,
				"title":       p.ProductTitle,
				"price_cents": p.PriceCents,
				"category_id": p.NewCategoryID,
			},
		},
	}
	info, _ := c.Find(bson.M{
		"_id":         p.ID,
		"store_id":    p.StoreID,
		"category_id": p.CategoryID,
	}).Apply(change, p)
	log.Println(info)
}

func (p *ProductOrder) ReorderStoreProducts() error {
	c := p.DB.C(ProductCollectionName).With(p.DBSession)
	p_count, err := c.Find(bson.M{
		"$and": []bson.M{
			bson.M{
				"category_id": p.CID,
				"store_id":    p.SID,
				"enabled":     true,
			}, bson.M{
				"_id": bson.M{
					"$in": p.PIDS,
				},
			},
		},
	}).Count()
	if err != nil || p_count != len(p.PIDS) {
		return errors.New("All active product ids must be provided. You may not include any products not currently enabled.")
	}

	previousTnx := ReadOnlyProduct{}
	c.Find(bson.M{
		"store_id":    p.SID,
		"category_id": p.CID,
		"previoustnxop": bson.M{
			"$exists": true,
		},
	}).One(&previousTnx)

	op_id := bson.NewObjectId()
	operations := make([]txn.Op, len(p.PIDS))
	for index, id := range p.PIDS {
		operations[index] = txn.Op{
			C:  ProductCollectionName,
			Id: id,
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
