package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type consumeRequest struct {
	T string `json:"t"`
}

type consumeResponse struct {
	M   string `json:"m"` //will send all messages in a single string
	Err string `json:"err, omitempty"`
}

func makeConsumeEndpoint(svc KafkaService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(consumeRequest)
		v, err := svc.Consume(ctx, req.T)
		if err != nil {
			return consumeResponse{v, err.Error()}, nil
		}
		return consumeResponse{v, ""}, nil
	}
}

func decodeConsumeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request consumeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func decodeConsumeResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response consumeResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}

	r.Body = ioutil.NopCloser(&buf)
	return nil
}
