package services

import (
	"DaaS/models"
	"net/http"
    "bytes"
    "time"
    "io/ioutil"
	"encoding/json"
)

var pageSpeedService IPageSpeedService

type IPageSpeedService interface {
	GetData() (models.Result, error)
}

type PageSpeedService struct {

}

func (s *PageSpeedService) GetData(result models.Result) (models.Result, error) {

	timeout := time.Duration(8 * time.Second)
	client := http.Client{
	    Timeout: timeout,
	}

    var buffer bytes.Buffer

    buffer.WriteString("https://www.googleapis.com/pagespeedonline/v2/runPagespeed?screenshot=true&url=")
	buffer.WriteString(result.Url)
    var queryUrl = buffer.String()
	response, err := client.Get(queryUrl)

    if err != nil {
        return result , err
    } else {
        defer response.Body.Close()
        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
        	return result , err
        }
        if err := json.Unmarshal(contents, &result.PageData); err != nil {
        	return result , err
	    }
    }

	return result, err 
}