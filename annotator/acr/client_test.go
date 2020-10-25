package acr

import (
	"fmt"
	"testing"

	"github.com/radiorabe/acr-webhook-receiver/client/api"
	"github.com/radiorabe/acr-webhook-receiver/models"
)

type SpyClient struct {
	Called bool
	ID     int64
}

func (s *SpyClient) GetResult(params *api.GetResultParams) (*api.GetResultOK, error) {
	s.Called = true
	s.ID = params.ResultID
	r := api.NewGetResultOK()
	r.Payload = &models.Result{
		ID: params.ResultID,
	}
	return r, nil
}

func TestGetClient(T *testing.T) {
	c := GetClient("acrcloud.api.example.org")

	_, _ = c.(ClientInterface)
}

func TestGetRecord(T *testing.T) {
	s := SpyClient{}
	c := &Client{
		client: &s,
	}
	fmt.Println(c.GetRecord("http://acrcloud.api.example.com/v1/api/record/123456789"))
	fmt.Println(s.Called)
	fmt.Println(s.ID)
	// Output:
	// 123456789
	// true
	// 123456789
}
