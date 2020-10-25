package annotator

// @TODO bearer token

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/sirupsen/logrus"
)

// ClientInterface ...
type ClientInterface interface {
	Post(anno *Annotation, collection string) (string, error)
	CreateCollection(collection string) error
	SearchByTargetID(id string, generator string) (*SingleTypeAnnotationContainer, error)
}

// Client ...
type Client struct {
	endpoint string

	client *http.Client

	log *logrus.Entry

	collectionCreated map[string]bool
}

// NewClient ...
func NewClient(endpoint string) ClientInterface {
	return &Client{
		endpoint: endpoint,

		client: &http.Client{},

		log: logrus.
			WithField("system", "wap-client").
			WithField("endpoint", endpoint),

		collectionCreated: make(map[string]bool),
	}
}

// Post ...
func (c *Client) Post(anno *Annotation, collection string) (string, error) {
	log := c.log.WithField("collection", collection)

	json, err := json.Marshal(anno)
	if err != nil {
		log.WithError(err).Error("Failed to marshal JSON.")
		return "", err
	}
	log.WithField("json", string(json)).Debug("POSTing annotation.")
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s/", c.endpoint, collection), bytes.NewBuffer(json))
	if err != nil {
		log.WithError(err).Error("Failed to create POST request.")
		return "", err
	}
	req.Header.Set("Content-Type", ContentTypeJSONLD)

	resp, err := c.client.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to post annotation.")
		return "", err
	}
	defer resp.Body.Close()

	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		log.WithError(err).Error("Failed to read response data.")
		return "", err
	}

	location, err := resp.Location()
	if err != nil {
		log.WithError(err).Error("Failed to find location in response.")
		return "", err
	}

	log.WithField("status", resp.Status).WithField("id", location.String()).Info("Created annotation")
	return location.String(), nil
}

// CreateCollection ...
func (c *Client) CreateCollection(collection string) error {
	if c.collectionCreated[collection] {
		return nil
	}
	collectionURL := fmt.Sprintf("%s%s/", c.endpoint, collection)

	log := c.log.WithField("collection", collection)
	log.Debug("Checking annotation collection.")

	req, err := http.NewRequest("HEAD", collectionURL, nil)
	if err != nil {
		log.WithError(err).Error("Failed to head annotation collection.")
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		log.WithError(err).Error("Failed to do head annotation collection.")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		ac := &AnnotationContainer{
			Context: []string{"http://www.w3.org/ns/anno.jsonld", "http://www.w3.org/ns/ldp.jsonld"},
			Type:    []string{"BasicContainer", "AnnotationCollection"},
			Label:   collection,
		}
		json, err := json.Marshal(ac)
		if err != nil {
			log.WithError(err).Error("Failed to marshal JSON.")
			return err
		}
		log.WithField("json", string(json)).Info("Creating annotation collection.")
		req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(json))
		req.Header.Set("Content-Type", ContentTypeJSONLD)
		req.Header.Set("Slug", collection)
		if err != nil {
			log.WithError(err).Error("Failed to post annotation collection")
			return err
		}
		resp, err = c.client.Do(req)
		if err != nil {
			log.WithError(err).Error("Failed to do post annotation collection.")
			return err
		}
		log.WithField("status", resp.Status).Info("Created Collection.")
		defer resp.Body.Close()
	}
	c.collectionCreated[collection] = true
	return nil
}

// SearchByTargetID ...
func (c *Client) SearchByTargetID(id string, generator string) (*SingleTypeAnnotationContainer, error) {
	searchURL := fmt.Sprintf(
		"%sservices/search/target?fields=id&value=%s&generator=%s",
		c.endpoint,
		url.QueryEscape(id),
		url.QueryEscape(generator),
	)

	log := c.log.WithField("url", searchURL)
	log.Debug("Searching by target ID.")

	resp, err := http.Get(searchURL)
	if err != nil {
		log.WithError(err).Error()
		return nil, err
	}
	defer resp.Body.Close()

	annos := &SingleTypeAnnotationContainer{}
	if err := json.NewDecoder(resp.Body).Decode(annos); err != nil {
		log.WithField("json", resp.Body).WithError(err).Error()
		return nil, err
	}
	log.WithField("total", annos.Total).Debug("Got search by target ID results.")
	return annos, nil
}
