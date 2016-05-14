package stockfighter

import (
	"github.com/franela/goreq"
	"fmt"
	"io/ioutil"
)

type Api struct {
	config *config
}

type HeartBeat struct {
	Ok    bool
	Error string
}

func InitApi(config *config) *Api {
	return &Api{
		config: config,
	}

}

func (s *Api) HeartBeat() (*HeartBeat, error) {
	response, err := s.GetRequest("ob/api/heartbeat")

	if err != nil {
		return nil, err
	}

	var hb *HeartBeat
	jsonErr := response.Body.FromJsonTo(&hb)

	if jsonErr == nil {
		return hb, nil
	} else {
		return nil, jsonErr
	}
}

func (s *Api) GetRequest(url string) (*goreq.Response, error) {
	return s.Request("GET", url)
}

func (s *Api) Request(method string, path string) (*goreq.Response, error) {
	req := goreq.Request{
		Uri: fmt.Sprintf("https://api.stockfighter.io/%s", path),
		Accept: "application/json",
		ContentType: "application/json",
	}

	if response, err := req.Do(); err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	} else {
		return response, nil
	}
}

func (s *Api) RequestContent(method string, path string) (string, error) {
	req := goreq.Request{
		Uri: fmt.Sprintf("https://api.stockfighter.io/%s", path),
		Accept: "application/json",
		ContentType: "application/json",
	}

	if response, err := req.Do(); err != nil {
		fmt.Printf("%+v\n", err)
		return "", err
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return "", err
		}
		return string(contents), nil
	}
}