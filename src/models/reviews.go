package models

import (
	"encoding/json"
	"errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

const (
	ReviewCollectionName = "Reviews"

	PlatformReview = "platform"
	ProductReview  = "product"
	StoreReview    = "store"
	OrderReview    = "order"
	PageSize       = 100
)

type Review struct {
	ID         bson.ObjectId `bson:"_id" json:"id"`
	ReviewedOn time.Time     `bson:"reviewed_on" json:"reviewed_on"`
	ReviewFor  string        `bson:"review_for" json:"review_for"`
	ProductId  bson.ObjectId `bson:"product_id,omitempty" json:"product_id"`
	Username   string        `bson:"username" json:"username"`
	Comment    string        `bson:"comment" json:"comment" validate:"max=200"`
	StoreId    bson.ObjectId `bson:"store_id,omitempty" json:"store_id"`
	OrderId    bson.ObjectId `bson:"order_id,omitempty,omitempty" json:"order_id"`
	UserId     bson.ObjectId `bson:"user_id" json:"user_id"`
	Score      uint8         `bson:"score" json:"score" validate:"required,max=5"`

	DB          *mgo.Database `bson:"-" json:"-"`
	DBSession   *mgo.Session  `bson:"-" json:"-"`
	CurrentPage int           `bson:"-" json:"-"`
}

func (r *Review) AddReview() error {
	var user User
	user.DB = r.DB
	user.ID = r.UserId
	user.DBSession = r.DBSession
	r.ReviewedOn = time.Now()
	if err := user.GetById(); err != nil {
		return errors.New(
			"We could not find the user associated with this request.",
		)
	}
	r.UserId = user.ID
	r.Username = strings.Split(user.Email, "@")[0]
	c := r.DB.C(ReviewCollectionName).With(r.DBSession)
	r.ID = bson.NewObjectId()
	if r.ReviewFor == StoreReview || r.ReviewFor == OrderReview {
		if err := r.UpdateStoreReview(); err != nil {
			return err
		}
	}
	buff, _ := json.Marshal(r)
	log.Println(string(buff))

	return c.Insert(r)
}

func (r *Review) UpdateStoreReview() error {
	c := r.DB.C(StoreCollectionName).With(r.DBSession)
	var store Store
	change := mgo.Change{
		ReturnNew: false,
		Upsert:    false,
		Remove:    false,
		Update: bson.M{
			"$inc": bson.M{
				"review_score": r.Score,
				"review_count": 1,
			},
		},
	}
	log.Println(r.StoreId)
	_, err := c.Find(bson.M{
		"_id": r.StoreId,
	}).Apply(change, &store)
	return err
}

func (r *Review) StoreReviews() (error, []Review) {
	c := r.DB.C(ReviewCollectionName).With(r.DBSession)
	reviews := []Review{}
	err := c.Find(bson.M{
		"store_id":   r.StoreId,
		"review_for": StoreReview,
	}).Skip(
		PageSize * (r.CurrentPage - 1),
	).Sort("-$natural").All(&reviews)
	return err, reviews
}

func (r *Review) PlatformReviews() (error, []Review) {
	c := r.DB.C(ReviewCollectionName).With(r.DBSession)
	reviews := []Review{}
	err := c.Find(bson.M{
		"review_for": PlatformReview,
	}).Skip(
		PageSize * (r.CurrentPage - 1),
	).Sort("-$natural").All(&reviews)
	return err, reviews
}
