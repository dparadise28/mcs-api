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
	AWS_S3_BUCKET_NAME = ""
	AWS_S3_BUCKET_KEY  = ""
)

type Asset struct {
	ID           bson.ObjectId `bson:"_id" json:"asset_id"`
	Link         string        `bson:"link" json:"link" validate:"required"`
	Title        string        `bson:"title" json:"title" validate:"required"`
	DisplayTitle string        `bson:"display_title" json:"display_title"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type NewAsset struct {
	ID           bson.ObjectId `bson:"_id" json:"asset_id"`
	Image        string        `bson:"image" json:"image" validate:"required"`
	Title        string        `bson:"title" json:"title"`
	DisplayTitle string        `bson:"display_title" json:"display_title" validate:"required"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

type AutocompleteAsset struct {
	ID           bson.ObjectId `bson:"_id" json:"asset_id"`
	Link         string        `bson:"link" json:"image" validate:"required"`
	Title        string        `bson:"title" json:"title" validate:"required"`
	DisplayTitle string        `bson:"display_title" json:"label"`

	DB        *mgo.Database `bson:"-" json:"-"`
	DBSession *mgo.Session  `bson:"-" json:"-"`
}

func (a *Asset) SearchForAsset(queryTerm string) []AutocompleteAsset {
	var assets []AutocompleteAsset
	c := a.DB.C(AssetCollectionName).With(a.DBSession)
	c.Find(bson.M{
		"$text": bson.M{
			"$search": queryTerm,
		},
	}).Limit(100).All(&assets)
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
	if _, err := svc.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(AWS_S3_BUCKET_NAME),
		Key:         aws.String(AWS_S3_BUCKET_KEY + "/" + a.ID.Hex()),
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
	asset.Link = "https://s3-" +
		AWS_REGION + ".amazonaws.com/" +
		AWS_S3_BUCKET_NAME + "/" +
		AWS_S3_BUCKET_KEY + "/" + asset.ID.Hex()
	if insert_err := c.Insert(&asset); insert_err != nil {
		return insert_err, asset
	}
	return nil, asset
}
