package services

import (
	"DaaS/models"
    "gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

var resultService IResultService

type IResultService interface {
	NewResult(newResult models.Result) (models.Result, error)
}

type ResultService struct {
	session *mgo.Session
}

func (s *ResultService) NewResult(newResult models.Result) (models.Result, error) {
    // Add an Id
    newResult.Id = bson.NewObjectId()

    s.session.DB("Daas").C("result").Insert(newResult)

	return newResult, nil
}
