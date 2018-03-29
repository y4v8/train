package train

import (
	"github.com/y4v8/errors"
)

type ParamsTrainWagon struct {
	From          string  `url:"from"`
	To            string  `url:"to"`
	Train         string  `url:"train"`
	Date          string  `url:"date"`
	WagonNum      int     `url:"wagon_num"`
	WagonType     string  `url:"wagon_type"`
	WagonClass    string  `url:"wagon_class"`
	CachedScheme1 *string `url:"cached_scheme[],omitempty"`
	CachedScheme2 *string `url:"cached_scheme[],omitempty"`
}

type DataTrainWagon struct {
	Places   map[string][]string `json:"places"`
	SchemeID string              `json:"schemeId"`
}

func (a *Api) TrainWagon(param ParamsTrainWagon) (*DataTrainWagon, error) {
	var data DataTrainWagon

	err := a.requestDataObject("POST", "train_wagon/", param, &data)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return &data, nil
}
