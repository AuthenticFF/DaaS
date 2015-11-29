package controllers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/Ramshackle-Jamathon/DaaS/models"
	"github.com/Ramshackle-Jamathon/DaaS/services"
    "net"
	"net/url"
	"net/http"
	"errors"
    "log"
	"github.com/asaskevich/govalidator"
)

type daasController struct {
	pageSpeedService services.PageSpeedService
	serverLoadService services.ServerLoadService
	resultService services.ResultService
	colorService services.ColorService
}

func (c *daasController) Init(router *httprouter.Router) *httprouter.Router {
	router.GET("/api/daas", PublicRoute(Daas.Analyze))

	return router

}
/*
Main API method
*/
func (c *daasController) Analyze(writer http.ResponseWriter, req *http.Request, params httprouter.Params) (interface{}, httpStatus) {
    result := models.Result{}

	result.Url = req.URL.Query().Get("url")
	//check for param
	if len(result.Url) != 0 {

		//url validation
		var err error
		result.Url, err = c.validateURL(result.Url)
		if err != nil {
			return nil, ServerError(err)
		}
        log.Println("URL valid.")
		//google pagespeed test
		pageResult, err := c.pageSpeedService.GetData(result)
		if err != nil {
			return nil, ServerError(err)
		}
        log.Println("Pagespeed complete.")

		//vegeta server test
		serverResult, err := c.serverLoadService.GetData(pageResult)
		if err != nil {
			return nil, ServerError(err)
		}
        log.Println("server load complete.")

		//page color test
		colorResult, err := c.colorService.GetData(serverResult)
		if err != nil {
			return nil, ServerError(err)
		}
        log.Println("color analysis complete.")

		//store in MongoDB
		storedResult, err := c.resultService.NewResult(colorResult)
		if err != nil {
			return nil, ServerError(err)
		}
		return storedResult, StatusOk(http.StatusOK)
	}
	return nil, ServerError(errors.New("No URL Provided"))

}

/*
URL validation/correction and DNS resolution
*/
func (c *daasController) validateURL(urlString string) (string, error) {

	//check for valid URL format
	validURL := govalidator.IsURL(urlString)
	if validURL == false {
		return urlString, errors.New("Invalid URL Format")
	}
    log.Println("validation is url.")

	//DNS lookup 
	urlObject, err := url.Parse(urlString)
	if err != nil {
		return urlString, err
	}
    log.Println("validation DNS success.")

	urlObject.Scheme = "http"
	urlObject, err = url.Parse(urlObject.String())
	if err != nil {
		return urlString, err
	}
    log.Println("validation Scheme.")

    //temporarily unreliable
	_, err = net.LookupHost(urlObject.Host)
	if err != nil {
		return urlString, err
	}
    log.Println("validation Host exists.")

	return urlObject.String(), nil
}


