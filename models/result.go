package models

import (
    "gopkg.in/mgo.v2/bson"
	//"time"
)

type (  
    // Result represents the structure of our resource
	Result struct {
		Id          bson.ObjectId       `json:"id" bson:"_id"`
		Url        string    `json:"url" bson:"url"`
		PageData map[string]interface{} `json:"pagedata" bson:"pagedata"`
	}
)