package train

import (
	"github.com/y4v8/errors"
)

type ParamsRoute struct {
	From  string `url:"routes[0][from]"`
	To    string `url:"routes[0][to]"`
	Date  string `url:"routes[0][date]"`
	Train string `url:"routes[0][train]"`
}

type DataRoute struct {
	Tpl    string       `json:"tpl"`
	Routes []TrainRoute `json:"routes"`
}

type TrainRoute struct {
	Train string    `json:"train"`
	List  []Station `json:"list"`
}

type Station struct {
	Name          string  `json:"name"`
	Code          int     `json:"code"`
	ArrivalTime   string  `json:"arrivalTime"`
	DepartureTime string  `json:"departureTime"`
	Distance      string  `json:"distance"`
	Lat           float64 `json:"lat"`
	Long          float64 `json:"long"`
}

func (a *Api) Route(param ParamsRoute) (*DataRoute, error) {
	var data DataRoute

	err := a.requestDataObject("POST", "route/", param, &data)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &data, nil
}
