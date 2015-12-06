package services

import (
	"github.com/Ramshackle-Jamathon/DaaS/models"
    
    "log"
    "github.com/gotk3/gotk3/gtk"
    "github.com/gotk3/gotk3/glib"
    "github.com/Ramshackle-Jamathon/go-webkit2/webkit2"
    //"github.com/sqs/gojs"
    //"io/ioutil"
    //"path/filepath"
    "fmt"
    //"strings"
    "runtime"
    "image"
    "image/jpeg"
    "bytes"
    "encoding/base64"
)

var colorService IColorService

type IColorService interface {
	GetData() (models.Result, error)
}

type ColorService struct {

}

func init() {
    runtime.LockOSThread()

    webContext := webkit2.DefaultWebContext()
    webContext.SetCacheModel(webkit2.DocumentViewerCacheModel)

    gtk.Init(nil)
        
}

func (s *ColorService) Main() {
    defer func() {
        if r := recover(); r != nil {
                fmt.Println("goroutine paniced:", r)
        }
    }()
    for f := range mainfunc {
        f()
    }
}

// queue of work to run in main thread.
var mainfunc = make(chan func())

// do runs f on the main thread.
func do(f func()) {
    done := make(chan bool, 1)
    mainfunc <- func() {
        f()
        done <- true
    }
    <-done
}


func (s *ColorService) GetData(result models.Result) (models.Result, error){

    do(func() {

        webView := webkit2.NewWebView()
        defer webView.Destroy()

        webView.Connect("load-failed", func() {
            log.Println("Load failed.")
        })
        webView.Connect("load-changed", func(_ *glib.Object, i int) {
            loadEvent := webkit2.LoadEvent(i)

            switch loadEvent {
                case webkit2.LoadStarted:
                    fmt.Println("Load started.")
                    fmt.Printf("Title: %q\n", webView.Title())
                    fmt.Printf("URI: %s\n", webView.URI())

                case webkit2.LoadRedirected:
                    fmt.Println("Load redirected.")
                    fmt.Printf("Title: %q\n", webView.Title())
                    fmt.Printf("URI: %s\n", webView.URI())

                case webkit2.LoadCommitted:
                    fmt.Println("Load committed.")
                    fmt.Printf("Title: %q\n", webView.Title())
                    fmt.Printf("URI: %s\n", webView.URI())

                case webkit2.LoadFinished:
                    fmt.Println("Load finished.")
                    fmt.Printf("Title: %q\n", webView.Title())
                    fmt.Printf("URI: %s\n", webView.URI())
                    webView.GetSnapshot(func(img *image.RGBA, err error) {
                        if err != nil {
                            log.Printf("GetSnapshot error: %q", err)
                            gtk.MainQuit()
                        } else {
                            if img == nil{
                                log.Printf("!img")
                                gtk.MainQuit()
                            } else if img.Pix == nil {
                                log.Printf("!img.Pix")
                                gtk.MainQuit()
                            } else if img.Stride == 0 || img.Rect.Max.X == 0 || img.Rect.Max.Y == 0 {
                                log.Printf("!img.Stride or !img.Rect.Max.X or !img.Rect.Max.Y")
                                gtk.MainQuit()
                            } else {
                                buf := new(bytes.Buffer)
                                jpeg.Encode(buf, img, &jpeg.Options{100})
                                imgBase64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
                                result.Image = imgBase64Str
                                gtk.MainQuit()
                            }
                        }
                    })

            }

        })

        glib.IdleAdd(func() bool {
            webView.LoadURI(result.Url)
            return false
        })

        // Create a new toplevel window, set its title, and connect it to the
        // "destroy" signal to exit the GTK main loop when it is destroyed.
        win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
        if err != nil {
            log.Fatal("Unable to create window:", err)
        }

        // Set the default window size.
        win.SetDefaultSize(1280, 1200)

        // Recursively show all widgets contained in this window.
        win.ShowAll()

        gtk.Main()

    })
    return result, nil 


}



