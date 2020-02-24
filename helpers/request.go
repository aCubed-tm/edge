package helpers

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/render"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetJsonFromRequestBody(r *http.Request, v interface{}) error {
	if r.Method == "GET" {
		return errors.New("cannot call GetJsonFromRequestBody on a GET request")
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(bodyBytes, &v)
	if err != nil {
		return err
	}
	return nil
}

func WriteSuccess(_ http.ResponseWriter, _ *http.Request) {
	log.Printf("Returning success without payload")
	// nothing to do, using this method to log and possibly extend in the future
}

func WriteSuccessJson(w http.ResponseWriter, r *http.Request, v interface{}) {
	log.Printf("Returning success: %v", v)
	var resp struct {
		Value interface{} `json:"data"`
	}
	resp.Value = v
	render.JSON(w, r, resp)
}

func WriteErrorJson(w http.ResponseWriter, r *http.Request, e error) {
	log.Printf("Returning error: %v", e.Error())
	var resp struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	resp.Error.Message = e.Error()
	render.JSON(w, r, resp)
}

func GetJwtToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("couldn't find authorization header")
	}

	if !strings.HasPrefix(strings.ToLower(header), "bearer") {
		return "", errors.New("authorization header didn't start with 'bearer'")
	}

	return header[7:], nil
}
