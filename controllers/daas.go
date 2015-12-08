package controllers

import (
	"github.com/AuthenticFF/DaaS/models"
	"github.com/AuthenticFF/DaaS/services"

    "net"
	"net/url"
	"net/http"
	"errors"
    "log"
    "sync"

	"github.com/julienschmidt/httprouter"
	"github.com/asaskevich/govalidator"
)
type daasController struct {
	pageSpeedService services.PageSpeedService
	resultService services.ResultService
	colorService services.ColorService
	typographyService services.TypographyService
}

func (c *daasController) Init(router *httprouter.Router) *httprouter.Router {
	router.GET("/api/daas", PublicRoute(Daas.Analyze))

	return router

}
/*
Main API method
*/
func (c *daasController) Analyze(writer http.ResponseWriter, req *http.Request, params httprouter.Params) (interface{}, httpStatus) {
    var wg sync.WaitGroup
    var result models.Result
    var pageResult models.Result
    var pageErr error



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
    	wg.Add(1)
		go func() {
			pageResult, pageErr = c.pageSpeedService.GetData(result)
	        log.Println("Pagespeed complete.")
        	wg.Done()
    	}()


		//typography
		/*typographyResult, err := c.colorService.GetTypography(snapshotResult)
		if err != nil {
			return nil, ServerError(err)
		}
        log.Println("color analysis complete.")*/


		//vegeta server test
		/*go func() {
			serverResult, serverErr = c.serverLoadService.GetData(result)
	        log.Println("server load complete.")
        	wg.Done()
    	}()*/

	    wg.Wait()
	    log.Println("synched")

	   	if pageErr != nil {
			return nil, ServerError(err)
		}
	    result.PageData = pageResult.PageData

		//store in MongoDB
		storedResult, err := c.resultService.NewResult(result)
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

	//DNS lookup 
	urlObject, err := url.Parse(urlString)
	if err != nil {
		return urlString, err
	}

	urlObject.Scheme = "http"
	urlObject, err = url.Parse(urlObject.String())
	if err != nil {
		return urlString, err
	}

    //will break if docker container is not updating dns properly
	_, err = net.LookupHost(urlObject.Host)
	if err != nil {
		return urlString, err
	}


	return urlObject.String(), nil
}


