package annotator

import (
	"time"

	"github.com/sirupsen/logrus"
)

// ContentTypeJSONLD ...
const ContentTypeJSONLD = "application/ld+json; profile=\"http://www.w3.org/ns/anno.jsonld\""

// Annotation W3C Annotation Model representation
type Annotation struct {
	Context    string      `json:"@context"`
	Type       string      `json:"type"`
	Creator    string      `json:"creator"`
	Created    time.Time   `json:"created"`
	Generator  string      `json:"generator"`
	Generated  time.Time   `json:"generated"`
	Body       interface{} `json:"body"`
	Target     interface{} `json:"target"`
	Motivation string      `json:"motivation"`

	client ClientInterface
}

// AnnotationContentArray ...
type AnnotationContentArray []AnnotationContent

// AnnotationContent ...
type AnnotationContent struct {
	ID         string             `json:"id,omitempty"`
	Type       []string           `json:"type,omitempty"`
	ConformsTo string             `json:"conformsTo,omitempty"`
	Format     string             `json:"format,omitempty"`
	Value      string             `json:"value,omitempty"`
	Selector   *AnnotationContent `json:"selector,omitempty"`
	Generator  string             `json:"generator,omitempty"`
}

// AnnotationContainer ...
type AnnotationContainer struct {
	Context []string `json:"@context"`
	Type    []string `json:"type"`
	Label   string   `json:"label"`
	Total   int      `json:"total,omitempty"`
}

// SingleTypeAnnotationContainer ...
type SingleTypeAnnotationContainer struct {
	Context []string `json:"@context"`
	Type    string   `json:"type"`
	Label   string   `json:"label"`
	Total   int      `json:"total,omitempty"`
}

// NewAnnotation ...
func NewAnnotation(client ClientInterface) *Annotation {

	return &Annotation{
		Context:   "http://www.w3.org/ns/anno.jsonld",
		Type:      "Annotation",
		Creator:   "https://github.com/radiorabe/annotation-agent",
		Created:   time.Now().UTC(),
		Generator: "https://github.com/radiorabe/annotation-agent",
		Generated: time.Now().UTC(),

		client: client,
	}
}

// Post annotation to a WAP server.
func (a *Annotation) Post(collection string) (string, error) {
	log := logrus.
		WithField("collection", collection)

	if err := a.client.CreateCollection(collection); err != nil {
		log.WithError(err).Error("Failed to create collection.")
	}
	return a.client.Post(a, collection)
}
