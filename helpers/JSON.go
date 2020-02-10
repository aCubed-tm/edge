package helpers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/render"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

func GetJsonFromPostRequest(r *http.Request, v interface{}) error {
	if r.Method != "POST" {
		return errors.New(fmt.Sprintf("expected a POST request, received %s", r.Method))
	}

	bodyBytes, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(bodyBytes, &v)
	if err != nil {
		return err
	}
	return nil
}

func WriteSuccess(w http.ResponseWriter, r *http.Request) {
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

func RunGrpc(ip string, f func(context.Context, *grpc.ClientConn) (interface{}, error)) (interface{}, error) {
	log.Printf("Starting gRPC connection to %s", ip)
	conn, err := grpc.Dial(ip, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, errors.New(fmt.Sprintf("did not connect: %v", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	ret, err := f(ctx, conn)

	cancel()
	_ = conn.Close()

	return ret, err
}
