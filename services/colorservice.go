package services

import (
	"github.com/Ramshackle-Jamathon/DaaS/models"
    "log"
    "github.com/conformal/gotk3/gtk"
    "github.com/conformal/gotk3/glib"
    "github.com/sourcegraph/go-webkit2/webkit2"
    "github.com/sqs/gojs"
    "runtime"
    "image"
    "image/jpeg"
    "bytes"
    "encoding/base64"
    "io/ioutil"
    "path/filepath"
    "fmt"
    "strings"
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
        log.Println("Load failed.")
    })



    //loading html snipit
    absPath, _ := filepath.Abs("/services/assets/typography.html")
    testHTML, err := ioutil.ReadFile(absPath)
    if err != nil {
        return result, err
    }
    js := fmt.Sprint("document.body.insertAdjacentHTML( 'afterbegin', '",string(testHTML), "' );")
    js = strings.Replace(js,"\n","<br>",-1)



    webView.Connect("load-changed", func(_ *glib.Object, i int) {
        loadEvent := webkit2.LoadEvent(i)
        switch loadEvent {
            case webkit2.LoadFinished:
                log.Println("Load finished.")
                log.Printf("Title: %q\n", webView.Title())
                log.Printf("URI: %s\n", webView.URI())
                webView.RunJavaScript(js, func(val *gojs.Value, err error) {
                    if err != nil {
                        log.Println(err)
                    } else {
                        log.Printf("Hostname (from JavaScript): %q\n", val)
                    }
                })  
                webView.GetSnapshot(func(img *image.RGBA, err error) {
                    if err != nil {
                        log.Printf("GetSnapshot error: %q", err)
                    }
                    if img.Pix == nil {
                        log.Printf("!img.Pix")
                    }
                    if img.Stride == 0 || img.Rect.Max.X == 0 || img.Rect.Max.Y == 0 {
                        log.Printf("!img.Stride or !img.Rect.Max.X or !img.Rect.Max.Y")
                    }
                    buf := new(bytes.Buffer)
                    jpeg.Encode(buf, img, &jpeg.Options{100})

                    imgBase64Str := base64.StdEncoding.EncodeToString(buf.Bytes())

                    result.Image = imgBase64Str
                    gtk.MainQuit()
                })

        }
    })

    glib.IdleAdd(func() bool {
        webView.LoadURI("https://www.google.com/")
        return false
    })


    // Create a new toplevel window, set its title, and connect it to the
    // "destroy" signal to exit the GTK main loop when it is destroyed.
    win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
    if err != nil {
        log.Fatal("Unable to create window:", err)
    }
    win.SetTitle("Simple Example")
    win.Connect("destroy", func() {
        gtk.MainQuit()
    })
    // Create a new label widget to show in the window.
    l, err := gtk.LabelNew("Hello, gotk3!")
    if err != nil {
        log.Fatal("Unable to create label:", err)
    }

    // Add the label to the window.
    win.Add(l)

    // Set the default window size.
    win.SetDefaultSize(1920, 1080)

    // Recursively show all widgets contained in this window.
    win.ShowAll()

    gtk.Main()

    // output:
    // Load finished.
    // Title: "Google"
    // URI: https://www.google.com/
    // Hostname (from JavaScript): "www.google.com"

    return result, nil 

}