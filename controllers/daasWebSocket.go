package controllers

import (
	"github.com/AuthenticFF/DaaS/services"
    "golang.org/x/net/websocket"
    "net/http"
    "fmt"

)

type daasWebSocketController struct {
	serverLoadService services.ServerLoadService
	resultService services.ResultService
}
func (c *daasWebSocketController) Init() {
	http.Handle("/socket/server", websocket.Handler(socketHandler))
}

func socketHandler(ws *websocket.Conn) {
    var err error

    for {
        var reply string

        if err = websocket.Message.Receive(ws, &reply); err != nil {
            fmt.Println("Can't receive")
            break
        }

        fmt.Println("Received back from client: " + reply)

        msg := "Received:  " + reply
        fmt.Println("Sending to client: " + msg)

        if err = websocket.Message.Send(ws, msg); err != nil {
            fmt.Println("Can't send")
            break
        }
    }
}