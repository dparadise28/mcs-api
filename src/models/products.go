package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math"
	"strings"
)

var ProductCollectionName = "Products"

type Product struct {
	ID             bson.ObjectId `bson:"_id" json:"product_id"`
	Size           string        `bson:"size" json:"size" validate:"required"`
	Image          string        `bson:"image" json:"image" validate:"required"`
	AssetID        bson.ObjectId `bson:"asset_id" json:"asset_id" validate:"required"`
	Enabled        bool          `bson:"enabled" json:"enabled" validate:"required"`
	StoreID        bson.ObjectId `bson:"store_id" json:"store_id"`
	TaxRate        float64       `bson:"tax_rate" json:"tax_rate" validate:"required"`
	PriceCents     uint32        `bson:"price_cents" json:"price_cents" validate:"required"`
	CategoryID     bson.ObjectId `bson:"category_id" json:"category_id" validate:"required"`
	ProductTitle   string        `bson:"title" json:"title" validate:"required"`
	ProductRatings struct {
		ReviewCount           uint64  `bson:"review_count" json:"total_reviews"`
		ReviewPercentageScore float64 `bson:"pct_score" json:"review_percent"`
	}
	// DisplayPrice   string        `bson:"-" json:"display_price"`
	// NewCategoryID  bson.ObjectId `bson:"-" json:"new_category_id"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type NewProduct struct {
	Size         string        `bson:"size" json:"size" validate:"required"`
	Image        string        `bson:"image" json:"image" validate:"required"`
	AssetID      bson.ObjectId `bson:"asset_id" json:"asset_id" validate:"required"`
	Enabled      bool          `bson:"enabled" json:"enabled" validate:"required"`
	TaxRate      float64       `bson:"tax_rate" json:"tax_rate" validate:"required"`
	PriceCents   uint32        `bson:"price_cents" json:"price_cents" validate:"required"`
	CategoryID   bson.ObjectId `bson:"category_id" json:"category_id" validate:"required"`
	ProductTitle string        `bson:"title" json:"title" validate:"required"`
}

type CartProduct struct {
	ID           bson.ObjectId `bson:"_id" json:"product_id"`
	Size         string        `bson:"size" json:"size"`
	Image        string        `bson:"image" json:"image" validate:"required"`
	AssetID      bson.ObjectId `bson:"asset_id" json:"asset_id" validate:"required"`
	TaxRate      float64       `bson:"tax_rate" json:"tax_rate" validate:"required"`
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
	QTY          uint16        `json:"quantity" validate:"gte=0"`
	Instructions string        `json:"instructions"`
}

type PaginatedProducts struct {
	Metadata PaginationMetadata `bson:"-" json:"metadata"`
	Results  []Product          `bson:"-" json:"results"`

	Size      int           `bson:"-" json:"-"`
	SID       bson.ObjectId `bson:"-" json:"-"`
	CID       bson.ObjectId `bson:"-" json:"-"`
	PG        int           `bson:"-" json:"-"`
	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (p *Product) AddProducts(products []NewProduct, sid bson.ObjectId) (error, []Product) {
	newProducts := []Product{}
	cmap := map[bson.ObjectId]bool{}
	for _, np := range products {
		cmap[np.CategoryID] = true
		newProducts = append(newProducts, Product{
			ID:           bson.NewObjectId(),
			Size:         np.Size,
			Image:        np.Image,
			AssetID:      np.AssetID,
			Enabled:      np.Enabled,
			StoreID:      sid,
			TaxRate:      np.TaxRate,
			PriceCents:   np.PriceCents,
			CategoryID:   np.CategoryID,
			ProductTitle: np.ProductTitle,
		})
	}

	cids := []bson.ObjectId{}
	for k, _ := range cmap {
		cids = append(cids, k)
	}
	cat := p.DB.C(TemplateCategoryCollectionName).With(p.DBSession)
	ccount, err := cat.Find(bson.M{
		"_id": bson.M{
			"$in": cids,
		},
	}).Count()
	if err != nil || ccount != len(cids) {
		return errors.New("Please only include available categories."), []Product{}
	}
	c := p.DB.C(ProductCollectionName).With(p.DBSession)
	if insert_err := c.Insert(I(newProducts)...); insert_err != nil {
		if strings.Count(insert_err.Error(), "ObjectId") == 1 {
			return nil, newProducts
		}
		return insert_err, []Product{}
	}
	return nil, newProducts
}

func (p *PaginatedProducts) RetrieveStoreProductsByCategory() {
	c := p.DB.C(ProductCollectionName).With(p.DBSession)
	query := bson.M{
		"category_id": p.CID,
	}
	count, err := c.Find(query).Count()
	if err != nil {
		p.Results = []Product{}
	}
	if p.PG*p.Size < count && err == nil {
		c.Find(query).Sort("$natural").Limit(p.Size).Skip(p.Size * p.PG).All(&p.Results)
		if p.Results == nil {
			p.Results = []Product{}
		}
	} else {
		p.Results = []Product{}
	}

	p.Metadata.Page = p.PG
	p.Metadata.PerPage = p.Size
	p.Metadata.PageSize = len(p.Results)
	p.Metadata.PageCount = int(math.Ceil(float64(count) / float64(p.Size)))
	p.Metadata.TotalCount = count
}

func (p *Product) AddProductToStoreCategory() error {
	c := p.DB.C(ProductCollectionName).With(p.DBSession)
	p.ID = bson.NewObjectId()
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
				"title":       p.ProductTitle,
				"size":        p.Size,
				"price_cents": p.PriceCents,
				// "category_id": p.NewCategoryID,
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
