package main

/*
 Triggers loading of yesterdays broadcasts into the system.
*/

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/radiorabe/annotation-agent/annotator/archive"
	"github.com/radiorabe/annotation-agent/queue"
)

func main() {
	raarUser := flag.String("raar-user", getenvDefault("RAAR_USER", ""), "RAAR username (must be able to download flac)")
	raarPass := flag.String("raar-pass", getenvDefault("RAAR_PASS", ""), "RAAR password")
	raarDownloadPrefix := flag.String("raar-download-prefix", getenvDefault("RAAR_DOWNLOAD_PREFIX", "https://archiv.rabe.ch"), "")

	amqpDSN := flag.String("amqp-dsn", getenvDefault("AMQP_DSN", "amqp://pubsub:pubsub@rabbitmq:5672/example-pubsub"), "AMQP DSN")
	amqpTopic := flag.String("amqp-topic", getenvDefault("AMQP_TOPIC", "global_pubsub"), "AMQP topic")
	flag.Parse()

	pubsub := queue.NewPubSub(*amqpTopic)
	if err := pubsub.Init(*amqpDSN); err != nil {
		logrus.WithError(err).Fatal()
	}
	defer pubsub.Close()

	for _, bcast := range archive.GetClient().Login(*raarUser, *raarPass).GetBroadcasts(time.Now().AddDate(0, 0, -1)) {
		if err := pubsub.PublishWithKey(context.TODO(), fmt.Sprintf("%s%s", *raarDownloadPrefix, bcast.Links.Self), "imported", "archive"); err != nil {
			logrus.WithError(err).Fatal()
		}
	}
}

// getenvDefault ...
func getenvDefault(env string, defaults string) string {
	e := os.Getenv(env)
	if e == "" {
		return defaults
	}
	return e
}
