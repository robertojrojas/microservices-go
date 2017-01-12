package services

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	catModels "github.com/robertojrojas/microservices-go/pets/cats-service/models"
	"golang.org/x/net/context"
)

const (
	CatServiceKey = "CatService"
)

type CatService struct {
	URL string
}

func (service *CatService) RPC(rpcRequest *RPCRequest) (rpcResponse *RPCResponse, err error) {

	log.Printf("[%T] Using ServiceAddress  %s\n", service, service.URL)
	req, err := http.NewRequest("GET", service.URL, nil)
	if err != nil {
		return
	}
	rpcResponse = &RPCResponse{
		Key: CatServiceKey,
	}
	err = httpDo(rpcRequest.Ctx, req, func(resp *http.Response, err error) error {
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		cats := []*catModels.Cat{}
		err = json.Unmarshal(data, &cats)
		if err != nil {
			return err
		}
		rpcResponse.Data = cats

		return nil
	})

	return

}

// httpDo issues the HTTP request and calls f with the response. If ctx.Done is
// closed while the request or f is running, httpDo cancels the request, waits
// for f to exit, and returns ctx.Err. Otherwise, httpDo returns f's error.
func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}
