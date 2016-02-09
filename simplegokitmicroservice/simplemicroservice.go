package main

// Adapted from the original GoKit examples available at https://github.com/go-kit/kit/blob/master/examples

import (
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/go-kit/kit/endpoint"
	"net/http"
	"log"
	httptransport "github.com/go-kit/kit/transport/http"
	"strings"
	"errors"
)


type SimpleStringManipulatingMicroservice interface {
	Uppercase(string) (string, error)
}

type simpleStringManipulatingMicroservice struct{}

func (simpleStringManipulatingMicroservice) Uppercase(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return strings.ToUpper(s), nil
}

func main() {
	ctx := context.Background()
	svc := simpleStringManipulatingMicroservice{}

	uppercaseHandler := httptransport.NewServer(
		ctx,
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	http.Handle("/uppercase", uppercaseHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}


func makeUppercaseEndpoint(svc SimpleStringManipulatingMicroservice) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(uppercaseRequest)
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return uppercaseResponse{v, err.Error()}, nil
		}
		return uppercaseResponse{v, ""}, nil
	}
}

func decodeUppercaseRequest(r *http.Request) (interface{}, error) {
	var request uppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

type uppercaseRequest struct {
	S string `json:"s"`
}

type uppercaseResponse struct {
	V   string `json:"v"`
	Err string `json:"err,omitempty"` // errors don't define JSON marshaling
}

var ErrEmpty = errors.New("empty string")