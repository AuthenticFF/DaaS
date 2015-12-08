package services

import (
	"github.com/AuthenticFF/DaaS/models"
    "time"

    vegeta "github.com/tsenart/vegeta/lib"
)

var serverLoadService IServerLoadService

type IServerLoadService interface {
	GetData() (models.Result, error)
}

type ServerLoadService struct {

}

func (s *ServerLoadService) GetData(result models.Result) (models.Result, error){

    rate := uint64(10) // per second
    duration := 4 * time.Second
    targeter := vegeta.NewStaticTargeter(vegeta.Target{
        Method: "GET",
        URL:    result.Url,
    })
    attacker := vegeta.NewAttacker()


    var metrics vegeta.Metrics
    for res := range attacker.Attack(targeter, rate, duration) {
        metrics.Add(res)
    }
    metrics.Close();


    return result, nil 

}