package archive

/**
Archive Record Annotator
========================
Annotate records from RAAR based on scanning flac files with various tools.

To initialize the RAAR client:
	curl https://archiv.rabe.ch/api | python -mjson.tool > annotator/archive/raar/swagger.json
	patch annotator/archive/raar/swagger.json hack/generate_archive_client.patch
	go generate annotator/archive/archive_annotator.go
*/

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/radiorabe/annotation-agent/annotator"
	"github.com/radiorabe/annotation-agent/annotator/archive/raar/models"
)

//go:generate ../../hack/generate_archive_client.sh

// Annotator implementation for RAAR
type Annotator struct {
	raarAllowedCodec   string
	raarDownLoadPrefix string

	audiowaveformBin         string
	peaksContentType         string
	peaksPublicURL           string
	peaksAnnotationContainer string

	sonicannotatorBin                 string
	sonicannotatorTransform           string
	sonicannotatorAnnotationContainer string
	sonicannotatorMotivation          string
	sonicannotatorType                []string
	sonicannotatorFormat              string

	client  LoggedInClientInterface
	storage StorageInterface

	wapClient  annotator.ClientInterface
	bodyType   []string
	bodyFormat string

	log *logrus.Entry
}

// AnnotatorOptions ...
type AnnotatorOptions struct {
	RAARUsername       string
	RAARPassword       string
	RAARDownloadPrefix string

	StorageEndpoint     string
	StorageAccessKey    string
	StorageAccessSecret string
	StorageUseSSL       bool

	PeaksStorageBucket       string
	PeaksStoragePublicURL    string
	PeaksAnnotationContainer string

	SonicAnnotatorTransform           string
	SonicannotatorAnnotationContainer string
}

// NewAnnotator gets an Annotator
func NewAnnotator(wapClient annotator.ClientInterface, options *AnnotatorOptions) *Annotator {
	return &Annotator{
		raarAllowedCodec:   "flac",
		raarDownLoadPrefix: options.RAARDownloadPrefix,

		audiowaveformBin:         "audiowaveform",
		peaksContentType:         "application/binary",
		peaksPublicURL:           options.PeaksStoragePublicURL,
		peaksAnnotationContainer: options.PeaksAnnotationContainer,

		sonicannotatorBin:                 "sonic-annotator",
		sonicannotatorTransform:           options.SonicAnnotatorTransform,
		sonicannotatorAnnotationContainer: options.SonicannotatorAnnotationContainer,
		sonicannotatorMotivation:          "classifying",
		sonicannotatorType:                []string{"TextualBody", "Dataset"},
		sonicannotatorFormat:              "text/plain",

		wapClient:  wapClient,
		bodyFormat: "text/plain",
		bodyType:   []string{"TextualBody", "Dataset"},

		client: GetClient().Login(
			options.RAARUsername,
			options.RAARPassword,
		),
		storage: NewStorage(
			options.StorageEndpoint,
			options.StorageAccessKey,
			options.StorageAccessSecret,
			options.StorageUseSSL,
			options.PeaksStorageBucket,
		).Init(),
		log: logrus.WithField("system", "archive-annotator"),
	}
}

// CreateAnnotations based on RAAR data and files.
func (a *Annotator) CreateAnnotations(url string) ([]string, error) {
	return a.FromRecord(url, a.client.GetRecord(url))
}

// FromRecord ...
func (a *Annotator) FromRecord(uri string, record *models.Broadcast) ([]string, error) {
	log := a.log.WithField("uri", uri)

	if annos, err := a.wapClient.SearchByTargetID(uri, "urn:annotation-agent"); err == nil && annos.Total > 0 {
		log.Warning("Already handled.")
		return nil, nil
	} else if err != nil {
		log.WithError(err).Panic("Failed to check if annotations exist.")
	}

	file := NewFile(record.ID)
	if file.HasDownload() {
		log.Warning("Unprocessed file exists.")
	}
	files := a.client.GetFiles(record)
	for _, apiFile := range files {
		if apiFile.Attributes.Codec == a.raarAllowedCodec {
			err := file.Download(a.fullyQualifyURL(apiFile.Links.Download)) // blocks until file is downloaded and ready
			if err != nil {
				log.WithError(err).Error("Failed to download file.")
				return nil, err
			}
		}
	}
	annos := []string{}
	if !file.HasPeaks() {
		annos = append(annos, a.LoadPeaks(uri, files, file.Path(), file.PeakPath())...)
	}
	if !file.HasSegments() {
		annos = append(annos, a.LoadSegments(uri, files, record, file.Path(), file.SegmentsPath())...)
	}
	anno := annotator.NewAnnotation(a.wapClient)
	anno.Body = &annotator.AnnotationContent{
		Type:      a.bodyType,
		Format:    a.bodyFormat,
		Value:     "Processed",
		Generator: "urn:annotation-agent",
	}
	anno.Target = &annotator.AnnotationContent{
		ID:        uri,
		Generator: "urn:annotation-agent",
	}
	annoLocation, err := anno.Post("annotation-agent")
	if err != nil {
		log.WithError(err).Error("Failed to create annotation")
		return nil, err
	}
	annos = append(annos, annoLocation)
	// @TODO delete downloaded file (with keep option for debuging)!
	return annos, nil
}

// LoadPeaks ...
func (a *Annotator) LoadPeaks(uri string, files []*models.AudioFile, filePath string, peakPath string) []string {
	log := a.log.
		WithField("id", uri).
		WithField("bin", a.audiowaveformBin)

	bin, err := exec.LookPath(a.audiowaveformBin)
	if err != nil {
		log.WithError(err).Fatal()
	}
	log.Infof("Scanning file.")

	cmd := exec.Command(
		bin,
		"--input-filename", filePath,
		"--output-filename", peakPath,
		"--output-format", "dat",
	)
	if err := cmd.Run(); err != nil {
		log.WithError(err).Fatal()
	}
	log.Infof("Scanned file.")

	storeID, err := a.storage.Store(
		uri,
		peakPath,
		a.peaksContentType,
	)
	if err != nil {
		log.WithError(err).Fatal("Failed to store scan results.")
	}
	storeURL := fmt.Sprintf("%s%s", a.peaksPublicURL, storeID)

	anno := annotator.NewAnnotation(a.wapClient)
	anno.Target = a.getTargets(uri, files, "urn:audiowavefile")
	anno.Body = &annotator.AnnotationContent{
		ID:        storeURL,
		Type:      []string{a.peaksContentType},
		Generator: "urn:audiowavefile",
	}

	location, err := anno.Post(a.peaksAnnotationContainer)
	if err != nil {
		log.WithError(err).Error("Failed to store peaks annotation.")
	}

	return []string{location}
}

// LoadSegments ...
func (a *Annotator) LoadSegments(uri string, files []*models.AudioFile, record *models.Broadcast, filePath string, segmentPath string) []string {
	log := a.log.
		WithField("filePath", filePath).
		WithField("segmentPath", segmentPath).
		WithField("bin", a.sonicannotatorBin).
		WithField("transform", a.sonicannotatorTransform)

	bin, err := exec.LookPath(a.sonicannotatorBin)
	if err != nil {
		logrus.WithError(err).Fatal()
	}

	cmd := exec.Command(
		bin,
		"--transform", a.sonicannotatorTransform,
		filePath,
		"--writer", "csv",
		"--csv-one-file", segmentPath,
		"--csv-omit-filename",
		"--csv-fill-ends",
	)
	log = log.
		WithField("bin", bin).
		WithField("args", cmd.Args)

	log.Infof("Scanning file.")
	if err := cmd.Run(); err != nil { // will block until command is ran
		logrus.WithError(err).Fatal()
	}
	log.Infof("Scanned file.")

	csvFile, err := os.Open(segmentPath)
	if err != nil {
		log.WithError(err).Fatal()
	}
	r := csv.NewReader(csvFile)

	results := []string{}
	for {
		csvRecord, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.WithError(err).Fatal()
		}
		start, err := strconv.ParseFloat(csvRecord[0], 64)
		if err != nil {
			log.WithError(err).Fatal()
		}
		end, err := strconv.ParseFloat(csvRecord[1], 64)
		if err != nil {
			log.WithError(err).Fatal()
		}
		if end == 0 {
			// we assume that it must be that last line that has no end since it lasts until the end of the file
			s, err := time.Parse(time.RFC3339, record.Attributes.FinishedAt.String())
			if err != nil {
				log.WithError(err).Fatal()
			}
			e, err := time.Parse(time.RFC3339, record.Attributes.FinishedAt.String())
			if err != nil {
				log.WithError(err).Fatal()
			}
			end = s.Sub(e).Seconds()
		}
		recordType := csvRecord[3]

		anno := annotator.NewAnnotation(a.wapClient)
		anno.Motivation = a.sonicannotatorMotivation
		anno.Target = a.fragmentTargets(a.getTargets(uri, files, "urn:sonic-annotator"), start, end)
		anno.Body = &annotator.AnnotationContent{
			Type:      a.sonicannotatorType,
			Format:    a.sonicannotatorFormat,
			Value:     recordType,
			Generator: "urn:sonic-annotator",
		}

		location, err := anno.Post(a.sonicannotatorAnnotationContainer)
		if err != nil {
			log.WithError(err).Fatal()
		}
		results = append(results, location)
	}

	return results
}

func (a *Annotator) getTargets(uri string, files []*models.AudioFile, generator string) []annotator.AnnotationContent {
	var targets []annotator.AnnotationContent
	targets = append(targets, annotator.AnnotationContent{
		ID: uri,
	})
	for _, apiFile := range files {
		content := annotator.AnnotationContent{
			ID:        a.fullyQualifyURL(apiFile.Links.Play),
			Type:      []string{"Sound"},
			Format:    fmt.Sprintf("audio/%s", apiFile.Attributes.Codec),
			Generator: generator,
		}
		targets = append(targets, content)
	}
	return targets
}

func (a *Annotator) fragmentTargets(targets []annotator.AnnotationContent, start float64, end float64) []annotator.AnnotationContent {
	for i, target := range targets {
		target.Selector = &annotator.AnnotationContent{
			Type:       []string{"FragmentSelector"},
			ConformsTo: "http://www.w3.org/TR/media-frags/",
			Value:      fmt.Sprintf("t=%f,%f", start, end),
		}
		targets[i] = target
	}
	return targets
}

func (a *Annotator) fullyQualifyURL(path string) string {
	// @TODO get rid of this method by figuring out to make raar return fully qualified links
	return fmt.Sprintf("%s%s", a.raarDownLoadPrefix, path)
}
