package train

import (
	"github.com/y4v8/errors"
)

type ParamsStation struct {
	Term string `url:"term"`
}

type DataStation struct {
	Title  string `json:"title"`
	Region string `json:"region"`
	Value  int    `json:"value"`
}

func (a *Api) Station(param ParamsStation) ([]DataStation, error) {
	var data []DataStation

	err := a.requestDataArray("GET", "train_search/station/", param, &data)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return data, nil
}
