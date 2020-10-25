package queue

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/radiorabe/annotation-agent/annotator"
)

// PubSub implements a RabbitMQ publisher and receiver
type PubSub struct {
	exchange   string
	maxThreads int32

	log *logrus.Entry

	conn *amqp.Connection
	ch   *amqp.Channel
	q    amqp.Queue

	msgs      <-chan amqp.Delivery
	observers map[string]annotator.Annotator
}

// NewPubSub creates a PubSub instance
func NewPubSub(exchange string) *PubSub {
	return &PubSub{
		exchange:   exchange,
		maxThreads: 5,
		observers:  make(map[string]annotator.Annotator),

		log: logrus.WithField("system", "pubsub"),
	}
}

// Init pubsub connection
func (p *PubSub) Init(dsn string) error {

	if err := p.Bind(dsn); err != nil {
		p.log.WithError(err).Error("Failed to connect and bind to RabbitMQ")
		return err
	}

	if err := p.ExchangeDeclare(); err != nil {
		p.log.WithError(err).Warning("Failed to declare an exchange")
		return err
	}

	if err := p.QueueDeclare(); err != nil {
		p.log.WithError(err).Warning("Failed to declare a queue")
		return err
	}
	return nil
}

// RegisterObserver ...
func (p *PubSub) RegisterObserver(key string, observer annotator.Annotator) {
	p.log.WithField("key", key).Debug("Registered Observer.")
	p.observers[key] = observer
}

// Dispatch incoming messages
func (p *PubSub) Dispatch() {
	type messageID struct{}
	var threads int32
	var messages int
	atomic.StoreInt32(&threads, 0)
	for d := range p.msgs {
		ctx := context.WithValue(context.Background(), &messageID{}, d.MessageId)
		log := p.log.
			WithField("routing_key", d.RoutingKey).
			WithField("exchange", d.Exchange).
			WithField("stage", "dispatch").
			WithField("type", "receive").
			WithField("body", string(d.Body)).
			WithContext(ctx)

		// throttle to maxThreads by rejecting messages and requeuing them
		if atomic.LoadInt32(&threads) >= p.maxThreads {
			if err := d.Nack(false, true); err != nil {
				log.WithError(err).Error("Failed to nack/requeue message.")
			}
			continue
		}
		messages++
		log = log.WithField("count", messages)

		if err := d.Ack(false); err != nil {
			log.WithError(err).Error("Failed to ack message.")
		}

		go func(d amqp.Delivery) {
			log.WithField("thread_count", atomic.AddInt32(&threads, 1)).Info("Dispatching message.")
			var a annotator.Annotator
			if p.observers[d.RoutingKey] != nil {
				a = p.observers[d.RoutingKey]
			} else {
				a = &annotator.DefaultAnnotator{}
				log.Trace(d)
			}
			uris, _ := a.CreateAnnotations(string(d.Body))
			for _, uri := range uris {
				if err := p.Publish(ctx, uri, "created"); err != nil {
					log.WithError(err).Error()
				}
			}
			log.WithField("thread_count", atomic.AddInt32(&threads, -1)).Info("Handled message.")
		}(d)
	}
}

// Publish a message to the exchange
func (p *PubSub) Publish(ctx context.Context, uri string, action string) error {
	return p.PublishWithKey(ctx, uri, action, "annotation")
}

// PublishWithKey a message to the exchange with a non-default key (ie. not annotation.<action>)
func (p *PubSub) PublishWithKey(ctx context.Context, uri string, action string, key string) error {
	routingKey := fmt.Sprintf("%s.%s", key, action)

	log := p.log.
		WithField("routing_key", routingKey).
		WithField("exchange", p.exchange).
		WithField("type", "publish").
		WithField("body", uri).
		WithContext(ctx)

	if err := p.ch.Publish(
		p.exchange, // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(uri),
		}); err != nil {

		log.Warning("Failed to publish a message")
		return err
	}

	log.Info("Published message.")
	return nil
}

// Bind to a amqp dsn.
func (p *PubSub) Bind(dsn string) error {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		p.log.WithError(err).Warning("Failed to connect to RabbitMQ")
		return err
	}
	p.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		p.log.WithError(err).Warning("Failed to open a channel")
		return err
	}
	if err = ch.Confirm(
		false, // nowait
	); err != nil {
		p.log.WithError(err).Warning("Failed to confirm channel")
		return err
	}
	p.ch = ch

	return nil
}

// Close connection and channels.
func (p *PubSub) Close() {
	defer p.conn.Close()
	defer p.ch.Close()
}

// ExchangeDeclare used for global message subbing
func (p *PubSub) ExchangeDeclare() error {
	err := p.ch.ExchangeDeclare(
		p.exchange, // name
		"topic",    // type
		false,      // durable
		true,       // auto-deleted
		false,      // internal
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		p.log.WithError(err).Warning("Failed to declare an exchange")
	}
	return err
}

// QueueDeclare used for global messages subbing
func (p *PubSub) QueueDeclare() error {
	q, err := p.ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		p.log.WithError(err).Warning("Failed to declare a queue")
		return err
	}
	p.q = q
	return nil
}

// QueueBindObservedQueues ...
func (p *PubSub) QueueBindObservedQueues() error {
	var queues []string
	for queue := range p.observers {
		p.log.WithField("queue", queue).Debug("Binding observable queue.")
		queues = append(queues, queue)
	}
	return p.QueueBind(queues)
}

// QueueBind to the queue we want to ingest
func (p *PubSub) QueueBind(queues []string) error {
	for _, queue := range queues {
		if err := p.ch.QueueBind(
			p.q.Name,   // queue name
			queue,      // routing key
			p.exchange, // exchange
			false,
			nil); err != nil {

			p.log.WithError(err).Warning("Failed to bind a queue")
			return err
		}
	}
	return nil
}

// Consume messages and put the on msgs channel for later dispatching
func (p *PubSub) Consume() error {
	msgs, err := p.ch.Consume(
		p.q.Name, // queue
		"",       // consumer
		false,    // auto ack
		false,    // exclusive
		false,    // no local
		false,    // no wait
		nil,      // args
	)
	if err != nil {
		p.log.WithError(err).Warning("Failed to consume a queue")
		return err
	}
	p.msgs = msgs
	return nil
}
