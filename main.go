package main

import (
	"flag"
	"net/http"
	"os"

	loghttp "github.com/motemen/go-loghttp"
	"github.com/sirupsen/logrus"

	"github.com/radiorabe/annotation-agent/annotator"
	"github.com/radiorabe/annotation-agent/annotator/acr"
	"github.com/radiorabe/annotation-agent/annotator/archive"
	"github.com/radiorabe/annotation-agent/queue"
)

var pubsub *queue.PubSub

func main() {
	amqpDSN := flag.String("amqp-dsn", getenvDefault("AMQP_DSN", "amqp://pubsub:pubsub@rabbitmq:5672/example-pubsub"), "AMQP DSN")
	amqpTopic := flag.String("amqp-topic", getenvDefault("AMQP_TOPIC", "global_pubsub"), "AMQP topic")

	raarKey := flag.String("raar-key", getenvDefault("RAAR_KEY", "archive.created"), "Key that trigggers annotating of archive records.")
	raarUser := flag.String("raar-user", getenvDefault("RAAR_USER", ""), "RAAR username (must be able to download flac)")
	raarPass := flag.String("raar-pass", getenvDefault("RAAR_PASS", ""), "RAAR password")
	raarDownloadPrefix := flag.String("raar-download-prefix", getenvDefault("RAAR_DOWNLOAD_PREFIX", "https://archiv.rabe.ch"), "")

	acrKey := flag.String("acr-key", getenvDefault("ACR_KEY", "acr.created"), "Key that triggers annotating of ACRCloud records.")
	acrHostname := flag.String("acr-hostname", getenvDefault("ACR_HOSTNAME", "acrcloud.api.rabe.ch"), "")
	acrAnnotationContainer := flag.String("acr-annotation-container", getenvDefault("ACR_ANNOTATION_CONTAINER", "acr"), "")

	peaksStorageBucket := flag.String("peaks-storage-bucket", getenvDefault("PEAKS_STORAGE_BUCKET", "peaks"), "In which storage bucket to store peaks dat files.")
	peaksPublicURL := flag.String("peaks-public-url", getenvDefault("PEAKS_PUBLIC_URL", "https://peaks.api.rabe.ch/v1/audiowaveform/"), "Public URL prefixfor peaks data.")
	peaksAnnotationContainer := flag.String("peaks-annotation-container", getenvDefault("PEAKS_ANNOTATION_COTAINER", "peaks"), "")

	sonicAnnotatorTransform := flag.String("sonic-annotator-transform", getenvDefault("SONIC_ANNOTATOR_TRANSFORM", "/etc/annotation-agent/n3/segmentation.n3"), "")
	sonicAnnotatorAnnotationContainer := flag.String("sonc-annotator-annotation-container", getenvDefault("SONIC_ANNOTATOR_ANNOTATION_CONTAINER", "speech-music"), "")

	storageEndpoint := flag.String("storage-endpoint", getenvDefault("STORAGE_ENDPOINT", "minio:9000"), "Endpoint of minio compatible storage.")
	storageAccessKey := flag.String("storage-access-key", getenvDefault("STORAGE_ACCESS_KEY", "minio"), "Access key for accessing minio.")
	storageAccessSecret := flag.String("storage-access-secret", getenvDefault("STORAGE_ACCESS_SECRET", "minio123"), "Secret key for accessing minio.")
	storageUseSSL := flag.Bool("storage-use-ssl", true, "Use SSL when connecting to minio.")

	wapEndpoint := flag.String("wap-endpoint", getenvDefault("WAP_ENDPOINT", "http://elucidate:8080/annotation/w3c/"), "Elucidate server endpoint.")

	logDebug := flag.Bool("log-debug", false, "debug logs")
	logTrace := flag.Bool("log-trace", false, "trace logs")
	flag.Parse()

	if debug := os.Getenv("ANNOTATION_AGENT_LOG_DEBUG"); debug != "" {
		t := true
		logDebug = &t
	}
	if trace := os.Getenv("ANNOTATION_AGENT_LOG_TRACE"); trace != "" {
		t := true
		logDebug = &t
	}
	if useSSL := os.Getenv("STORAGE_USE_SSL"); useSSL != "false" {
		f := false
		storageUseSSL = &f
	}

	logging(*logDebug, *logTrace)

	pubsub = queue.NewPubSub(*amqpTopic)
	if err := pubsub.Init(*amqpDSN); err != nil {
		logrus.WithError(err).Fatal()
	}
	defer pubsub.Close()

	wapClient := annotator.NewClient(*wapEndpoint)

	pubsub.RegisterObserver(*acrKey, acr.NewAnnotator(
		wapClient,
		&acr.AnnotatorOptions{
			AcrHostname:         *acrHostname,
			AnnotationContainer: *acrAnnotationContainer,
		},
	))
	pubsub.RegisterObserver(*raarKey, archive.NewAnnotator(
		wapClient,
		&archive.AnnotatorOptions{
			RAARUsername:       *raarUser,
			RAARPassword:       *raarPass,
			RAARDownloadPrefix: *raarDownloadPrefix,

			StorageEndpoint:     *storageEndpoint,
			StorageAccessKey:    *storageAccessKey,
			StorageAccessSecret: *storageAccessSecret,
			StorageUseSSL:       *storageUseSSL,

			PeaksStorageBucket:       *peaksStorageBucket,
			PeaksStoragePublicURL:    *peaksPublicURL,
			PeaksAnnotationContainer: *peaksAnnotationContainer,

			SonicAnnotatorTransform:           *sonicAnnotatorTransform,
			SonicannotatorAnnotationContainer: *sonicAnnotatorAnnotationContainer,
		},
	))
	/*
		@TODO create enrich annnotation annotator to replace unknown shows from acr annotations
		annotationKey := flag.String("annotation-key", getenvDefault("ANNOTATION_KEY", "annotation.created"), "")
		pubsub.RegisterObserver(*annotationKey, annotation.NewAnnotator(
			wapClient,
			&annotation.AnnotatorOptions{},
		}))
	*/

	waitForWork(pubsub) // this is blocking!
}

func logging(debug bool, trace bool) {
	if trace {
		logrus.SetLevel(logrus.TraceLevel)
	} else if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	http.DefaultTransport = &loghttp.Transport{
		LogRequest: func(req *http.Request) {
			logrus.WithFields(logrus.Fields{
				"request": req,
				"method":  req.Method,
				"url":     req.URL,
			}).Tracef("%s %s", req.Method, req.URL)
		},
		LogResponse: func(resp *http.Response) {
			logrus.WithFields(logrus.Fields{
				"request":    resp.Request,
				"statuscode": resp.StatusCode,
				"status":     resp.Status,
				"url":        resp.Request.URL,
				"body":       resp.Body,
			}).Tracef("%d %s", resp.StatusCode, resp.Request.URL)
		},
		Transport: http.DefaultTransport,
	}
}

func waitForWork(pubsub *queue.PubSub) {
	if err := pubsub.QueueBindObservedQueues(); err != nil {
		logrus.WithError(err).Fatal("Failed to bind a queue")
	}

	if err := pubsub.Consume(); err != nil {
		logrus.WithError(err).Fatal("Failed to bind a queue")
	}

	forever := make(chan bool)

	go pubsub.Dispatch()

	logrus.Infof("Waiting for messages.")
	<-forever
}

// getenvDefault ...
func getenvDefault(env string, defaults string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaults
	}
	return e
}
