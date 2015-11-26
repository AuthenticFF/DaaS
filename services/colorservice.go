package services

import (
	"DaaS/models"
    "fmt"
    "github.com/conformal/gotk3/gtk"
    "github.com/conformal/gotk3/glib"
    "github.com/sourcegraph/go-webkit2/webkit2"
    "github.com/sqs/gojs"
    "runtime"
)

var colorService IColorService

type IColorService interface {
	GetData() (models.Result, error)
}

type ColorService struct {

}

func (s *ColorService) GetData(result models.Result) (models.Result, error){

    runtime.LockOSThread()
    gtk.Init(nil)

    webView := webkit2.NewWebView()
    defer webView.Destroy()

    webView.Connect("load-failed", func() {
        fmt.Println("Load failed.")
    })
    webView.Connect("load-changed", func(_ *glib.Object, i int) {
        loadEvent := webkit2.LoadEvent(i)
        switch loadEvent {
        case webkit2.LoadFinished:
            fmt.Println("Load finished.")
            fmt.Printf("Title: %q\n", webView.Title())
            fmt.Printf("URI: %s\n", webView.URI())
            webView.RunJavaScript("window.location.hostname", func(val *gojs.Value, err error) {
                if err != nil {
                    fmt.Println("JavaScript error.")
                } else {
                    fmt.Printf("Hostname (from JavaScript): %q\n", val)
                }
                gtk.MainQuit()
            })
        }
    })

    glib.IdleAdd(func() bool {
        webView.LoadURI("https://www.google.com/")
        return false
    })

    gtk.Main()

    // output:
    // Load finished.
    // Title: "Google"
    // URI: https://www.google.com/
    // Hostname (from JavaScript): "www.google.com"

    return result, nil 

}