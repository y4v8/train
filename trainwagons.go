package train

import (
	"github.com/y4v8/errors"
)

type ParamsTrainWagons struct {
	From        string `url:"from"`
	To          string `url:"to"`
	Date        string `url:"date"`
	Train       string `url:"train"`
	WagonTypeID string `url:"wagon_type_id"`
	GetTpl      string `url:"get_tpl"`
}

type ResultTrainWagons struct {
	Data DataTrainWagons `json:"data"`
}
type DataTrainWagons struct {
	Types   []WagonType `json:"types"`
	Wagons  []Wagon     `json:"wagons"`
	TplPage *string     `json:"tplPage,omitempty"`
}
type WagonType struct {
	TypeID    string `json:"type_id"`
	Title     string `json:"title"`
	Letter    string `json:"letter"`
	Free      int    `json:"free"`
	Cost      int    `json:"cost"`
	IsOneCost bool   `json:"isOneCost"`
}
type Wagon struct {
	Num           int            `json:"num"`
	TypeID        string         `json:"type_id"`
	Type          string         `json:"type"`
	Class         string         `json:"class"`
	Railway       int            `json:"railway"`
	Free          int            `json:"free"`
	ByWishes      bool           `json:"byWishes"`
	HasBedding    bool           `json:"hasBedding"`
	Services      []string       `json:"services"`
	Prices        map[string]int `json:"prices"`
	ReservedPrice int            `json:"reservedPrice"`
	AllowBonus    bool           `json:"allowBonus"`
}

func (a *Api) TrainWagons(param ParamsTrainWagons) (*DataTrainWagons, error) {
	var data DataTrainWagons

	err := a.requestDataObject("POST", "train_wagons/", param, &data)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &data, nil
}
