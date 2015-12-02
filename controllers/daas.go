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
    "sync"
	"github.com/asaskevich/govalidator"
)
///docker rmi -f $(docker images | grep "^<none>" | awk "{print $3}")
type daasController struct {
	pageSpeedService services.PageSpeedService
	serverLoadService services.ServerLoadService
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

    wg.Add(2)
    var result models.Result
    var pageResult models.Result
    //var snapshotResult models.Result
   // var typographyResult models.Result
    var serverResult models.Result
    var pageErr error
    //var snapshotErr error
    //var typographyErr error
    var serverErr error



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
		go func() {
			pageResult, pageErr = c.pageSpeedService.GetData(result)
	        log.Println("Pagespeed complete.")
        	wg.Done()
    	}()

		//page snapshot
		/*go func(snapshotResult models.Result, snapshotErr error) {
			snapshotResult, snapshotErr = c.colorService.GetData(result)
	        log.Println("page snapshot complete.")
        	wg.Done()
    	}(snapshotResult, snapshotErr)*/

		//typography
		/*typographyResult, err := c.colorService.GetTypography(snapshotResult)
		if err != nil {
			return nil, ServerError(err)
		}
        log.Println("color analysis complete.")*/


		//vegeta server test
		go func() {
			serverResult, serverErr = c.serverLoadService.GetData(result)
	        log.Println("server load complete.")
        	wg.Done()
    	}()

    			//page color test
		result, err = c.colorService.GetData(result)
		if err != nil {
			return nil, ServerError(err)
		}

	    wg.Wait()
	    log.Println("synched")
	    result.PageData = pageResult.PageData
	    result.ServerData = serverResult.ServerData

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

    //will break if docker container is not updating dns properly
	_, err = net.LookupHost(urlObject.Host)
	if err != nil {
		return urlString, err
	}
    log.Println("validation Host exists.")


	return urlObject.String(), nil
}


