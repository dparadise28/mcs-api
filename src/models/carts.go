package models

import (
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	// "log"
	"strconv"
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
	IsNew    bool `bson:"-" json:"-"`
	ApplyFee bool `bson:"-" json:"-"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (cart *Cart) UpdateProductQuantityQueries(p CartProduct) []bson.M {
	pushQuery := bson.M{
		"$set": bson.M{
			"last_updated": time.Now(),
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
			"last_updated": time.Now(),
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
				"last_updated": time.Now(),
				"delivery_fee": cart.DeliveryFee,
				"created_at":   time.Now(),
				"store_name":   cart.StoreName,
				"cart_state":   CartStates["ACTIVE"],
				"tax_rate":     cart.StoreTaxRate,
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

func (cart *Cart) CompleteCart() error {
	c := cart.DB.C(CartCollectionName).With(cart.DBSession)
	change := mgo.Change{
		ReturnNew: true,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$set": bson.M{
				"last_updated": time.Now(),
				"cart_state":   CartStates["COMPLETED"],
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
	cart.Totals.Total = cart.Totals.Tax + float64(cart.Totals.Subtotal)
	if cart.ApplyFee {
		cart.Totals.Total = float64(cart.DeliveryFee)
	}
}

func FormatPriceCents(n int64) string {
	in := strconv.FormatInt(n, 10)
	neg := false
	out := make([]byte, len(in)+(len(in)-2+int(in[0]/'0'))/3)
	if in[0] == '-' {
		in, neg = in[1:], true
	}
	if len(in) == 1 {
		return "$" + string(out) + "0.0" + in
	}
	if len(in) == 2 {
		return "$" + string(out) + "0." + in
	}
	cents := in[len(in)-2:]
	in = in[:len(in)-2]
	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			break
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
	if neg {
		return "-$" + string(out) + "." + cents
	} else {
		return "$" + string(out) + "." + cents
	}
}

func (cart *Cart) GetCartProductsOrderMarkup() string {
	markup := ""
	for _, p := range cart.Products {
		if p.Quantity > 0 {
			markup += `<br><br>
			<div class="column-left"><b>` + strconv.Itoa(int(p.Quantity)) + ` </b> x </b></div>
			<div class="column-center" align="left">` + p.ProductTitle + `</div>
			<div class="column-right">` + FormatPriceCents(int64(p.PriceCents)) + `</div>`
		}
	}
	return markup
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

func (cart *Cart) GetCartsById() error {
	c := cart.DB.C(CartCollectionName).With(cart.DBSession)

	if err := c.Find(bson.M{
		"_id": cart.ID,
	}).One(cart); err != nil {
		return err
	}
	cart.UpdateCartTotals()
	return nil
}

func (cart *Cart) GetActiveCartsById() error {
	c := cart.DB.C(CartCollectionName).With(cart.DBSession)

	if err := c.Find(bson.M{
		"_id":        cart.ID,
		"cart_state": CartStates["ACTIVE"],
	}).One(cart); err != nil {
		return err
	}
	cart.UpdateCartTotals()
	return nil
}

func (cart *Cart) ActiveUserCartCountForStore() (int, error) {
	return cart.DB.C(CartCollectionName).With(cart.DBSession).Find(bson.M{
		"user_id":    cart.UserID,
		"store_id":   cart.StoreID,
		"cart_state": CartStates["ACTIVE"],
	}).Count()
}

func (cart *Cart) ReActivateCart() error {
	// copy vars here so theyre not overriden on fetch (need better abstraction; ok for now)
	db := cart.DB
	sess := cart.DBSession

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
	cart.DB, cart.DBSession = db, sess
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
