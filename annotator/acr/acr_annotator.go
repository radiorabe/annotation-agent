package acr

import (
	"fmt"
	"time"

	"github.com/radiorabe/acr-webhook-receiver/models"
	"github.com/sirupsen/logrus"

	"github.com/radiorabe/annotation-agent/annotator"
)

// Annotator implementation for ACRCloud
type Annotator struct {
	bodyFormat string
	bodyType   []string

	unknownShowURN                string
	clocktimeAnnotationURNPattern string

	client              ClientInterface
	annotationContainer string

	wapClient annotator.ClientInterface
}

// AnnotatorOptions ...
type AnnotatorOptions struct {
	AcrHostname string

	AnnotationContainer string
}

// NewAnnotator gets an Annotator
func NewAnnotator(wapClient annotator.ClientInterface, options *AnnotatorOptions) *Annotator {
	return &Annotator{
		bodyFormat: "application/json",
		bodyType:   []string{"Dataset"},

		unknownShowURN:                "urn:annotation.api.rabe.ch:annotation-agent:show:unknown",
		clocktimeAnnotationURNPattern: "urn:annotation.api.rabe.ch:annotation-agent:date?t=clocktime:%s",

		client: GetClient(
			options.AcrHostname,
		),
		annotationContainer: options.AnnotationContainer,

		wapClient: wapClient,
	}
}

// CreateAnnotations based on an ACRCloud record.
func (a *Annotator) CreateAnnotations(url string) ([]string, error) {
	return a.FromRecord(url, a.client.GetRecord(url))
}

// FromRecord creates an annotation and returns it's URL
func (a *Annotator) FromRecord(uri string, record *models.Result) ([]string, error) {
	anno := annotator.NewAnnotation(a.wapClient)

	ts, err := time.Parse("2006-01-02 15:04:05", *record.Result.Data.Metadata.TimestampUtc)
	if err != nil {
		logrus.WithError(err).Fatal()
	}
	anno.Generated = ts
	anno.Body = a.getAnnotationBody(uri)

	anno.Target = [2]annotator.AnnotationContent{
		{
			ID: a.unknownShowURN,
		},
		{
			ID:    fmt.Sprintf(a.clocktimeAnnotationURNPattern, ts),
			Value: ts.Format(time.RFC1123),
		},
	}
	location, err := anno.Post(a.annotationContainer)

	return []string{location}, err
}

func (a *Annotator) getAnnotationBody(uri string) *annotator.AnnotationContent {
	return &annotator.AnnotationContent{
		ID:     uri,
		Type:   a.bodyType,
		Format: a.bodyFormat,
	}
}
