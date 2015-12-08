package services

import (
    "github.com/AuthenticFF/DaaS/models"
    "fmt"
    "log"

    "github.com/PuerkitoBio/goquery"

)

var typographyService ITypographyService

type ITypographyService interface {
    GetData() (models.Result, error)
}

type TypographyService struct {

}
func (s *TypographyService) GetData(result models.Result) (models.Result, error){

    doc, err := goquery.NewDocument(result.Url) 
    if err != nil {
        log.Fatal(err)
    }

    doc.Find(".reviews-wrap article .review-rhs").Each(func(i int, s *goquery.Selection) {
        band := s.Find("h3").Text()
        title := s.Find("i").Text()
        fmt.Printf("Review %d: %s - %s\n", i, band, title)
    })

    return result, nil 
}
