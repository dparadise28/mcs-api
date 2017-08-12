package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var CartCollectionName = "Carts"

var CartStates = map[string]int{
	"ABANDONED": 0,
	"ACTIVE":    1,
	"COMPLETED": 2,
}

type Totals struct {
	Subtotal uint32  `bson:"-" json:"subtotal"`
	CashTip  bool    `bson:"cash_tip" json:"cash_tip"`
	Total    float64 `bson:"-" json:"total"`
	Tax      float64 `bson:"-" json:"tax"`
	Tip      uint32  `bson:"tip" json:"tip"`
}

type Cart struct {
	ID           bson.ObjectId `bson:"_id" json:"id"`
	Totals       Totals        `bson:"-" json:"totals"`
	UserID       bson.ObjectId `bson:"user_id" json:"user_id"`
	StoreID      bson.ObjectId `bson:"store_id" json:"store_id"`
	Products     []CartProduct `bson:"products" json:"products"`
	StoreName    string        `bson:"store_name" json:"store_name"`
	CartState    int           `bson:"cart_state" json:"cart_state"`
	DateCreated  time.Time     `bson:"created_at" json:"created_at"`
	LastUpdated  time.Time     `bson:"last_updated" json:"last_updated"`
	DeliveryFee  uint32        `bson:"delivery_fee" json:"delivery_fee"`
	StoreTaxRate float64       `bson:"tax_rate" json:"tax_rate"`

	//StoreInfo Store `bson:"-" json:"-"`
	IsNew bool `bson:"-" json:"-"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (cart *Cart) UpdateProductQuantityQueries(p CartProduct) []bson.M {
	pushQuery := bson.M{
		"$set": bson.M{
			"last_updated": cart.LastUpdated,
		},
		"$push": bson.M{
			"products": p,
		},
	}
	pullQuery := bson.M{
		"$pull": bson.M{
			"products": bson.M{
				"_id": p.ID,
			},
		},
		"$set": bson.M{
			"last_updated": cart.LastUpdated,
			"cart_state":   CartStates["ACTIVE"],
		},
	}
	// supper hacky (generally want to come up with a better appraoch but will do for now)
	if cart.IsNew {
		pullQuery = bson.M{
			"$pull": bson.M{
				"products": bson.M{
					"_id": p.ID,
				},
			},
			"$set": bson.M{
				"last_updated": cart.LastUpdated,
				"delivery_fee": cart.DeliveryFee,
				"tax_rate":     cart.StoreTaxRate,
				"cart_state":   CartStates["ACTIVE"],
			},
		}
	}
	return []bson.M{pullQuery, pushQuery}
}

func (cart *Cart) RunUpsertQueries(queries []bson.M) error {
	c := cart.DB.C(CartCollectionName).With(cart.DBSession)
	for _, query := range queries {
		change := mgo.Change{
			ReturnNew: true,
			Upsert:    cart.IsNew,
			Remove:    false,
			Update:    query,
		}
		_, err := c.Find(bson.M{
			"_id":        cart.ID,
			"user_id":    cart.UserID,
			"store_id":   cart.StoreID,
			"cart_state": CartStates["ACTIVE"],
		}).Apply(change, cart)
		if err != nil {
			return err
		}
	}
	return nil
}

func (cart *Cart) AbandonCart() error {
	c := cart.DB.C(CartCollectionName).With(cart.DBSession)
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{
				"last_updated": time.Now(),
				"cart_state":   CartStates["ABANDONED"],
			},
		},
	}
	_, err := c.Find(bson.M{
		"_id":        cart.ID,
		"user_id":    cart.UserID,
		"cart_state": CartStates["ACTIVE"],
	}).Apply(change, cart)
	if err != nil {
		return err
	}
	return nil
}

func (cart *Cart) UpdateCartTotals() {
	for _, product := range cart.Products {
		cart.Totals.Subtotal += product.PriceCents * uint32(product.Quantity)
	}
	cart.Totals.Tax = float64(cart.Totals.Subtotal) * cart.StoreTaxRate / 100.00
	cart.Totals.Total = cart.Totals.Tax + float64(cart.Totals.Subtotal) + float64(cart.DeliveryFee)
}

func (cart *Cart) UpdateProductQuantity(id bson.ObjectId, instructions string, quantity uint16) error {
	p_collection := cart.DB.C(ProductCollectionName).With(cart.DBSession)
	var p CartProduct
	pQuery := bson.M{
		"_id":      id,
		"store_id": cart.StoreID,
	}
	if err := p_collection.Find(pQuery).One(&p); err != nil || p.PriceCents == 0 {
		return err
	}
	p.Instructions = instructions
	p.Quantity = quantity

	queries := cart.UpdateProductQuantityQueries(p)
	if err := cart.RunUpsertQueries(queries); err != nil {
		return err
	}
	cart.UpdateCartTotals()
	return nil
}

func (cart *Cart) RetrieveUserCartsByStatus() ([]Cart, error) {
	c := cart.DB.C(CartCollectionName).With(cart.DBSession)

	carts := []Cart{}
	err := c.Find(bson.M{
		"user_id":    cart.UserID,
		"cart_state": cart.CartState,
	}).All(&carts)
	for cartIndex, _ := range carts {
		carts[cartIndex].UpdateCartTotals()
	}
	return carts, err
}

func (cart *Cart) ActiveUserCartCountForStore() (int, error) {
	c := cart.DB.C(CartCollectionName).With(cart.DBSession)
	return c.Find(bson.M{
		"user_id":    cart.UserID,
		"store_id":   cart.StoreID,
		"cart_state": CartStates["ACTIVE"],
	}).Count()
}

func (cart *Cart) ReActivateCart() error {
	c := cart.DB.C(CartCollectionName).With(cart.DBSession)
	c.Find(bson.M{
		"_id":        cart.ID,
		"user_id":    cart.UserID,
		"cart_state": CartStates["COMPLETED"],
	}).One(cart)
	if len(cart.Products) == 0 {
		return errors.New(
			"The cart you have selected is either empty or could not be found",
		)
	}
	if cartCount, countErr := cart.ActiveUserCartCountForStore(); countErr != nil {
		return countErr
	} else if cartCount != 0 {
		return errors.New(
			"It appears as though you have an active cart for" +
				" the store associated with the cart attepting" +
				" to be re-activated. Please drop that cart and" +
				" try again.")
	}
	cart.ID = bson.NewObjectId()
	cart.CartState = CartStates["ACTIVE"]
	return c.Insert(cart)
}
