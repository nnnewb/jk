// Code generated by jk -t Service generate all; DO NOT EDIT.

package service

import (
	"bytes"
	"context"
	"encoding/json"
	endpoint "github.com/go-kit/kit/endpoint"
	errors "github.com/juju/errors"
	"net/http"
	"net/url"
)

func MakeBuyRemoteEndpoint(host string, client *http.Client) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		buffer := bytes.NewBufferString("")
		u := url.URL{
			Host:   host,
			Path:   "/api/v1/service/buy",
			Scheme: "https",
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), buffer)
		if err != nil {
			return nil, errors.Trace(err)
		}

		response, err := client.Do(request)
		if err != nil {
			return nil, errors.Trace(err)
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return nil, errors.Errorf("call remote endpoint failed, http status %d %s", response.StatusCode, response.Status)
		}

		var resp BuyResponse
		err = json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return resp, nil
	}
}

func MakeBuyHandlerFunc(svc Service) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var payload BuyRequest
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
		resp, err := svc.Buy(req.Context(), payload)
		err = json.NewEncoder(wr).Encode(resp)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
	}
}

func MakeJoinRemoteEndpoint(host string, client *http.Client) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		buffer := bytes.NewBufferString("")
		u := url.URL{
			Host:   host,
			Path:   "/api/v1/service/join",
			Scheme: "https",
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), buffer)
		if err != nil {
			return nil, errors.Trace(err)
		}

		response, err := client.Do(request)
		if err != nil {
			return nil, errors.Trace(err)
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return nil, errors.Errorf("call remote endpoint failed, http status %d %s", response.StatusCode, response.Status)
		}

		var resp JoinResponse
		err = json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return resp, nil
	}
}

func MakeJoinHandlerFunc(svc Service) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var payload JoinRequest
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
		resp, err := svc.Join(req.Context(), payload)
		err = json.NewEncoder(wr).Encode(resp)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
	}
}

func MakeJoin2RemoteEndpoint(host string, client *http.Client) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		buffer := bytes.NewBufferString("")
		u := url.URL{
			Host:   host,
			Path:   "/api/v1/service/join-2",
			Scheme: "https",
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), buffer)
		if err != nil {
			return nil, errors.Trace(err)
		}

		response, err := client.Do(request)
		if err != nil {
			return nil, errors.Trace(err)
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return nil, errors.Errorf("call remote endpoint failed, http status %d %s", response.StatusCode, response.Status)
		}

		var resp Join2Response
		err = json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return resp, nil
	}
}

func MakeJoin2HandlerFunc(svc Service) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var payload Join2Request
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
		resp, err := svc.Join2(req.Context(), payload)
		err = json.NewEncoder(wr).Encode(resp)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
	}
}

func MakeJoin3RemoteEndpoint(host string, client *http.Client) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		buffer := bytes.NewBufferString("")
		u := url.URL{
			Host:   host,
			Path:   "/api/v1/service/join-3",
			Scheme: "https",
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), buffer)
		if err != nil {
			return nil, errors.Trace(err)
		}

		response, err := client.Do(request)
		if err != nil {
			return nil, errors.Trace(err)
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return nil, errors.Errorf("call remote endpoint failed, http status %d %s", response.StatusCode, response.Status)
		}

		var resp Join3Response
		err = json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return resp, nil
	}
}

func MakeJoin3HandlerFunc(svc Service) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var payload Join3Request
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
		resp, err := svc.Join3(req.Context(), payload)
		err = json.NewEncoder(wr).Encode(resp)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
	}
}

func MakeLowercaseRemoteEndpoint(host string, client *http.Client) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		buffer := bytes.NewBufferString("")
		u := url.URL{
			Host:   host,
			Path:   "/api/v1/service/lowercase",
			Scheme: "https",
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), buffer)
		if err != nil {
			return nil, errors.Trace(err)
		}

		response, err := client.Do(request)
		if err != nil {
			return nil, errors.Trace(err)
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return nil, errors.Errorf("call remote endpoint failed, http status %d %s", response.StatusCode, response.Status)
		}

		var resp LowercaseResponse
		err = json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return resp, nil
	}
}

func MakeLowercaseHandlerFunc(svc Service) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var payload LowercaseRequest
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
		resp, err := svc.Lowercase(req.Context(), payload)
		err = json.NewEncoder(wr).Encode(resp)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
	}
}

func MakeUppercaseRemoteEndpoint(host string, client *http.Client) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		buffer := bytes.NewBufferString("")
		u := url.URL{
			Host:   host,
			Path:   "/api/v1/service/uppercase",
			Scheme: "https",
		}
		request, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), buffer)
		if err != nil {
			return nil, errors.Trace(err)
		}

		response, err := client.Do(request)
		if err != nil {
			return nil, errors.Trace(err)
		}

		defer response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return nil, errors.Errorf("call remote endpoint failed, http status %d %s", response.StatusCode, response.Status)
		}

		var resp UppercaseResponse
		err = json.NewDecoder(response.Body).Decode(&resp)
		if err != nil {
			return nil, errors.Trace(err)
		}

		return resp, nil
	}
}

func MakeUppercaseHandlerFunc(svc Service) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		var payload UppercaseRequest
		err := json.NewDecoder(req.Body).Decode(&payload)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
		resp, err := svc.Uppercase(req.Context(), payload)
		err = json.NewEncoder(wr).Encode(resp)
		if err != nil {
			panic(errors.Errorf("unexpected unmarshal error %+v", err))
		}
	}
}

func NewClient(host string, client *http.Client) Service {
	return EndpointSet{
		BuyEndpoint:       MakeBuyRemoteEndpoint(host, client),
		Join2Endpoint:     MakeJoin2RemoteEndpoint(host, client),
		Join3Endpoint:     MakeJoin3RemoteEndpoint(host, client),
		JoinEndpoint:      MakeJoinRemoteEndpoint(host, client),
		LowercaseEndpoint: MakeLowercaseRemoteEndpoint(host, client),
		UppercaseEndpoint: MakeUppercaseRemoteEndpoint(host, client),
	}
}
func Register(svc Service, m *http.ServeMux) *http.ServeMux {
	m.HandleFunc("/api/v1/service/buy", MakeBuyHandlerFunc(svc))
	m.HandleFunc("/api/v1/service/join", MakeJoinHandlerFunc(svc))
	m.HandleFunc("/api/v1/service/join-2", MakeJoin2HandlerFunc(svc))
	m.HandleFunc("/api/v1/service/join-3", MakeJoin3HandlerFunc(svc))
	m.HandleFunc("/api/v1/service/lowercase", MakeLowercaseHandlerFunc(svc))
	m.HandleFunc("/api/v1/service/uppercase", MakeUppercaseHandlerFunc(svc))
	return m
}
