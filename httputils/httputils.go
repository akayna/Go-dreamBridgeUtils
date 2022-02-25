package httputils

import (
	"io/ioutil"
	"net/http"
)

func RequestBodyToString(req *http.Request) (string, error) {

	bodyString, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	}
	return string(bodyString), nil

}
