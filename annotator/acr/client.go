package acr

import (
	"net/url"
	"path"
	"strconv"

	"github.com/go-openapi/strfmt"
	"github.com/sirupsen/logrus"

	"github.com/radiorabe/acr-webhook-receiver/client"
	"github.com/radiorabe/acr-webhook-receiver/client/api"
	"github.com/radiorabe/acr-webhook-receiver/models"
)

// ClientInterface for clients
type ClientInterface interface {
	GetRecord(uri string) *models.Result
}

// ServiceInterface for ACRCloud API
type ServiceInterface interface {
	GetResult(params *api.GetResultParams) (*api.GetResultOK, error)
}

// Client concrete type
type Client struct {
	client ServiceInterface
}

var acrClient ClientInterface

// GetClient returns the acrClient singleton
func GetClient(host string) ClientInterface {
	if acrClient == nil {
		c := client.NewHTTPClientWithConfig(
			strfmt.Default,
			client.
				DefaultTransportConfig().
				WithHost(host),
		)
		acrClient = &Client{
			client: c.API,
		}
	}
	return acrClient
}

// GetRecord ...
func (c *Client) GetRecord(uri string) *models.Result {
	u, err := url.Parse(uri)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	id, err := strconv.Atoi(path.Base(u.Path))
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	result, err := c.client.GetResult(api.NewGetResultParams().WithResultID(int64(id)))
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	return result.GetPayload()
}
