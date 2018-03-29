package train

import (
	"github.com/y4v8/errors"
)

type ParamsTrainSearch struct {
	From   string `url:"from"`
	To     string `url:"to"`
	Date   string `url:"date"`
	Time   string `url:"time"`
	GetTpl string `url:"get_tpl"`
}

type DataTrainSearch struct {
	List     []Train `json:"list"`
	TplPage  *string `json:"tplPage,omitempty"`
	TplTrain *string `json:"tplTrain,omitempty"`
}
type Train struct {
	Num            string      `json:"num"`
	Category       int         `json:"category"`
	TravelTime     string      `json:"travelTime"`
	From           StationFrom `json:"from"`
	To             StationTo   `json:"to"`
	Types          []TrainType `json:"types"`
	Child          Child       `json:"child"`
	AllowStud      int         `json:"allowStud"`
	AllowBooking   int         `json:"allowBooking"`
	AllowRoundtrip int         `json:"allowRoundtrip"`
	IsEurope       int         `json:"isEurope"`
}
type StationTo struct {
	Code         string `json:"code"`
	Station      string `json:"station"`
	StationTrain string `json:"stationTrain"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	SortTime     int    `json:"sortTime"`
}
type StationFrom struct {
	StationTo
	SrcDate string `json:"srcDate"`
}
type TrainType struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Letter string `json:"letter"`
	Places int    `json:"places"`
}
type Child struct {
	MinDate string `json:"minDate"`
	MaxDate string `json:"maxDate"`
}

func (a *Api) TrainSearch(param ParamsTrainSearch) (*DataTrainSearch, error) {
	var data DataTrainSearch

	err := a.requestDataObject("POST", "train_search/", param, &data)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &data, nil
}
