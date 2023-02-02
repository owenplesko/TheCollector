package riot

import (
	"encoding/json"
	"errors"
	"net/http"
)

func getJson(url string, target interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	return json.NewDecoder(res.Body).Decode(target)
}
