package services

import (
	"github.com/Ramshackle-Jamathon/DaaS/db"
	"log"
)

var PageSpeed PageSpeedService
var Result ResultService
var ServerLoad ServerLoadService
var Color ColorService
var Typography TypographyService

func init() {
	PageSpeed = PageSpeedService{}
	ServerLoad = ServerLoadService{}
	Typography = TypographyService{}

	Color = ColorService{}
	//go Color.Main();
	
	Result = ResultService{db.Session}
	log.Printf("Services Initialized");
}
