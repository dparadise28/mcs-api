package models

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
	"time"
)

const (
	AssetCollectionName    = "Assets"
	DefaultS3UploadTimeout = 60 * time.Second
)

var (
	AWS_SECRET_ACCESS_KEY = ""
	AWS_ACCESS_KEY_ID     = ""
	AWS_SESSION_TOKEN     = ""
	AWS_REGION            = ""

	// the two variables below dictate the path to assets.
	// example:
	//		AWS_S3_BUCKET_NAME = mcs-images
	//		AWS_S3_BUCKET_KEY = production/images/products
	//
	//		final asset destination = mcs-images/production/images/products/asset_id
	AWS_S3_BUCKET_NAME      = ""
	AWS_S3_BUCKET_KEY       = ""
	AWS_S3_STORE_BUCKET_KEY = "stores"
)

type Asset struct {
	ID                 bson.ObjectId      `bson:"_id" json:"asset_id"`
	Link               string             `bson:"link" json:"link" validate:"required"`
	Size               string             `bson:"size" json:"size"`
	Title              string             `bson:"title" json:"title" validate:"required"`
	Details            []NutritionDetails `json:"details" bson:"details"`
	Nutrition          Nutrition          `bson:"nutrition" json:"nutrition"`
	PriceCents         uint32             `bson:"price_cents" json:"price_cents" validate:"required"`
	DisplayTitle       string             `bson:"display_title" json:"display_title"`
	TemplateCategoryID bson.ObjectId      `bson:"template_category_id,omitempty" json:"template_category_id" validate:"-"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type NewAsset struct {
	ID                 bson.ObjectId      `bson:"_id" json:"asset_id"`
	Size               string             `bson:"size" json:"size"`
	Image              string             `bson:"image" json:"image" validate:"required"`
	Title              string             `bson:"title" json:"title"`
	Details            []NutritionDetails `json:"details" bson:"details"`
	Nutrition          Nutrition          `bson:"nutrition" json:"nutrition"`
	PriceCents         uint32             `bson:"price_cents" json:"price_cents"`
	DisplayTitle       string             `bson:"display_title" json:"display_title" validate:"required"`
	TemplateCategoryID bson.ObjectId      `bson:"template_category_id,omitempty" json:"template_category_id,omitempty" validate:"-"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type AutocompleteAsset struct {
	ID                 bson.ObjectId `bson:"_id" json:"asset_id"`
	Link               string        `bson:"link" json:"image" validate:"required"`
	Title              string        `bson:"title" json:"title" validate:"required"`
	DisplayTitle       string        `bson:"display_title" json:"label"`
	TemplateCategoryID bson.ObjectId `bson:"template_category_id,omitempty" json:"template_category_id"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type TemplateAsset struct {
	ID                 bson.ObjectId `bson:"_id" json:"asset_id"`
	Size               string        `bson:"size" json:"size"`
	Link               string        `bson:"link" json:"image" validate:"required"`
	Title              string        `bson:"title" json:"title" validate:"required"`
	PriceCents         uint32        `bson:"price_cents" json:"price_cents"`
	DisplayTitle       string        `bson:"display_title" json:"label"`
	TemplateCategoryID bson.ObjectId `bson:"template_category_id" json:"template_category_id"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (ta *TemplateAsset) RetrieveTemplateCategoryAssets(cid bson.ObjectId, pg int) []TemplateAsset {
	var assets []TemplateAsset
	c := ta.DB.C(AssetCollectionName).With(ta.DBSession)
	c.Find(bson.M{
		"template_category_id": cid,
	}).Sort("$natural").Limit(100).Skip(100 * pg).All(&assets)
	if assets == nil {
		assets = []TemplateAsset{}
	}
	return assets
}

func (a *Asset) RetrieveTemplateAssetById(id bson.ObjectId) {
	c := a.DB.C(AssetCollectionName).With(a.DBSession)
	c.Find(bson.M{
		"_id": id,
	}).One(&a)
}

func (a *Asset) SearchForAsset(queryTerm string) []AutocompleteAsset {
	var assets []AutocompleteAsset
	c := a.DB.C(AssetCollectionName).With(a.DBSession)
	c.Find(bson.M{
		"$text": bson.M{
			"$search": queryTerm,
		},
	}).Select(bson.M{
		"score": bson.M{
			"$meta": "textScore",
		},
	}).Sort("$textScore:score").Limit(100).All(&assets)
	if assets == nil {
		assets = []AutocompleteAsset{}
	}
	return assets
}

func (a *NewAsset) UploadAsset() (error, Asset) {
	var asset Asset

	a.ID = bson.NewObjectId()
	image, decode_err := base64.StdEncoding.DecodeString(
		a.Image[strings.IndexByte(a.Image, ',')+1:],
	)
	if decode_err != nil {
		return decode_err, asset
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AWS_REGION),
		Credentials: credentials.NewStaticCredentials(
			AWS_ACCESS_KEY_ID,
			AWS_SECRET_ACCESS_KEY,
			AWS_SESSION_TOKEN,
		),
	})
	if err != nil {
		return err, asset
	}

	svc := s3.New(sess)
	ctx, cancelFn := context.WithTimeout(context.Background(), DefaultS3UploadTimeout)
	defer cancelFn()
	bucketKey := AWS_S3_BUCKET_KEY
	if a.TemplateCategoryID.Hex() == "" {
		asset.Size = a.Size
		asset.Details = a.Details
		asset.Nutrition = a.Nutrition
		asset.PriceCents = a.PriceCents
		asset.TemplateCategoryID = a.TemplateCategoryID
		bucketKey = AWS_S3_STORE_BUCKET_KEY
	}
	if _, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(AWS_S3_BUCKET_NAME),
		Key:         aws.String(bucketKey + "/" + a.ID.Hex() + ".jpeg"),
		ACL:         aws.String("public-read"),
		Body:        strings.NewReader(string(image)),
		ContentType: aws.String("image/jpeg"),
	}); err != nil {
		return err, asset
	}

	c := a.DB.C(AssetCollectionName).With(a.DBSession)

	asset.ID = a.ID
	asset.Title = strings.ToLower(a.DisplayTitle)
	asset.DisplayTitle = a.DisplayTitle
	asset.Link = "https://s3.amazonaws.com/" +
		AWS_S3_BUCKET_NAME + "/" +
		bucketKey + "/" + asset.ID.Hex() + ".jpeg"
	if insert_err := c.Insert(&asset); insert_err != nil {
		log.Println("insert err:", insert_err)
		return insert_err, asset
	}
	return nil, asset
}
