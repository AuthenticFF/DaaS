package services

import (
	"github.com/Ramshackle-Jamathon/DaaS/db"
	"log"
)

var PageSpeed PageSpeedService
var Result ResultService
var ServerLoad ServerLoadService
var Color ColorService

func init() {
	PageSpeed = PageSpeedService{}
	ServerLoad = ServerLoadService{}
	Color = ColorService{}
	Result = ResultService{db.Session}
	log.Printf("Services Initialized");
}
