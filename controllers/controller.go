package controllers

import (	
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/AuthenticFF/DaaS/services"
	"net/http"
	"os"
	"log"
)

var Daas daasController
var DaasWebSocket daasWebSocketController

func Init(router *httprouter.Router) *httprouter.Router {

	Daas = daasController{services.PageSpeed, services.Result, services.Color, services.Typography}
	DaasWebSocket = daasWebSocketController{services.ServerLoad, services.Result}
	router = Daas.Init(router)
	DaasWebSocket.Init();
	log.Printf("Controllers Initialized");

	return router
}

type httpStatus struct {
	err    error
	status int
}

var (
	verifyKey = os.Getenv("JWT_PUB_KEY")
	signKey   = os.Getenv("JWT_PRIV_KEY")
)

func ServerError(err error) httpStatus {
	return httpStatus{err, http.StatusInternalServerError}
}

func StatusOk(status int) httpStatus {
	return httpStatus{nil, status}
}

type controllerRoute func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, httpStatus)

/* JSON REST utils */

func WriteResponse(w http.ResponseWriter, result interface{}, httpStatus httpStatus) {
	var responseBody string
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Powered-By", "Authentic F&F")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
	if httpStatus.err != nil {
		w.WriteHeader(httpStatus.status)
		jsonBody, _ := json.Marshal(httpStatus.err.Error())
		responseBody = string(jsonBody)
	} else {
		w.WriteHeader(httpStatus.status)
		jsonBody, _ := json.Marshal(result)
		responseBody = string(jsonBody)
	}
	fmt.Fprintf(w, responseBody)
}

func ResponseHandler(r controllerRoute) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		result, httpStatus := r(w, req, p)
		WriteResponse(w, result, httpStatus)
	}
}

func authenticate(token string) error {
	if token == "" {
		return errors.New("No authorization token present")
	}

	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(verifyKey), nil
	})

	var authError error
	if err == nil && tk.Valid {
		authError = nil
	} else {
		switch err.(type) {
		case nil:
			if !tk.Valid {
				authError = errors.New("Invalid token")
			}
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				authError = errors.New("Token has expired. Please login again.")
			default:
				authError = errors.New("Invalid token")
			}
		default:
			authError = errors.New("There was an error while trying to validate your request")
		}
	}
	return authError
}

func RestrictedHandler(r controllerRoute) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		/*token := req.Header.Get("Authorization")

		err := authenticate(token)
		if err != nil {
			WriteResponse(w, "", httpStatus{err, http.StatusUnauthorized})
			return
		}
		*/
		result, status := r(w, req, p)
		WriteResponse(w, result, status)
	}
}

func PublicRoute(r controllerRoute) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		result, status := r(w, req, p)
		WriteResponse(w, result, status)
	}
}
func RestrictedRoute(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		route := p.ByName("route")
		w.Header().Set("X-Powered-By", "Authentic F&F")
		if route != "login" {
			/*token := req.Header.Get("Authorization")
			err := authenticate(token)
			if err != nil {
				w.Header().Set("X-Authenticated", "Not Allowed")
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.Header().Set("X-Authenticated", "Allowed")
			}*/
			w.Header().Set("X-Authenticated", "Allowed")
		}
		fn(w, req, p)
	}
}

