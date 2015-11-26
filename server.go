package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/Ramshackle-Jamathon/DaaS/controllers"
	"github.com/Ramshackle-Jamathon/DaaS/db"
	"html/template"
	"io/ioutil"	
	"log"
	"net/http"
	"os"
)

var PathSeperator = string(os.PathSeparator)
var Templates = template.New("")
var htmlTemplates = []string{"frontend/httpdocs/index.html"}

func main() {
	CompileRootTemplates()
	router := httprouter.New()
	defer db.Session.Close();
	// USER APP PAGE ROUTING
	// Probably Could Use Some Work
	router.GET("/", RenderTemplate("frontend/httpdocs/index"))

	// USER APP STATIC FILES
	// Probably Could Use Some Work
	router.ServeFiles("/assets/*filepath", http.Dir("frontend/httpdocs/assets/"))
	router = controllers.Init(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9091"
	}
	log.Printf("Authentic Form & Function (& Framework) listening on %s", port)

	http.ListenAndServe(":"+port, router)
}

func CompileRootTemplates() {
	for _, html := range htmlTemplates {
		filetext, err := ioutil.ReadFile(html)
		if err != nil {
			log.Fatal("Root app html view failed to compile: " + err.Error())
		}
		text := string(filetext)
		Templates.New(html).Parse(text)
	}
}

type Page struct {
	Title string
}

func RenderTemplate(template string) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		defer req.Body.Close()
		w.Header().Set("X-Powered-By", "Authentic F&F")
		page := Page{Title: "Demo"}
		err := Templates.ExecuteTemplate(w, template+".html", page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
