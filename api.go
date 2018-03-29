package train

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"

	"github.com/y4v8/errors"
)

type Api struct {
	baseUri string
}

type JsonResult struct {
	Data  json.RawMessage `json:"data"`
	Error *int            `json:"error,omitempty"`
}

func NewApi() *Api {
	api := &Api{
		baseUri: "https://booking.uz.gov.ua/ru/",
	}
	return api
}

func (a *Api) requestDataArray(method string, uri string, param interface{}, data interface{}) error {
	res, err := a.requestResult(method, uri, param)
	if err != nil {
		return errors.Wrap(err)
	}

	err = json.Unmarshal(res, data)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a *Api) requestDataObject(method string, uri string, param interface{}, data interface{}) error {
	res, err := a.requestResult(method, uri, param)
	if err != nil {
		return errors.Wrap(err)
	}

	var jsonRes JsonResult
	err = json.Unmarshal(res, &jsonRes)
	if err != nil {
		return errors.Wrap(err)
	}

	if jsonRes.Error != nil {
		return errors.New(string(jsonRes.Data))
	}

	err = json.Unmarshal(jsonRes.Data, data)
	if err != nil {
		return errors.Wrap(err)
	}

	return nil
}

func (a *Api) requestResult(method string, uri string, param interface{}) ([]byte, error) {
	v, err := query.Values(param)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	var resp *http.Response
	if method == "POST" {
		resp, err = http.PostForm(a.baseUri+uri, v)
	} else {
		resp, err = http.Get(a.baseUri + uri + "?" + v.Encode())
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.Wrap(err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("StatusCode: %v, %v", resp.StatusCode, resp.Status)
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err)
	}

	return res, nil
}
