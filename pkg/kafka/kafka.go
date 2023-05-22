package kafka

import (
	"context"
	"errors"
	"sync"

	"code.cloudfoundry.org/bytefmt"
	"github.com/google/wire"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kzerolog"
	"github.com/valyala/bytebufferpool"
)

type Options struct {
	Urls     []string
	ClientID string `koanf:"client_id"`
}

func NewOptions(k *koanf.Koanf) (*Options, error) {
	o := new(Options)
	if err := k.Unmarshal("kafka", o); err != nil {
		return nil, err
	}
	return o, nil
}

type Client struct {
	client *kgo.Client
	logger zerolog.Logger
}

func NewClient(o Options, logger zerolog.Logger) (*Client, func(), error) {
	l := logger.Level(zerolog.WarnLevel)
	client, err := kgo.NewClient(kgo.SeedBrokers(o.Urls...),
		kgo.ClientID(o.ClientID),
		kgo.ProducerBatchCompression(kgo.GzipCompression()),
		kgo.RecordRetries(10),
		kgo.MaxBufferedRecords(100000),
		kgo.ProducerBatchMaxBytes(10*bytefmt.MEGABYTE),
		kgo.MaxProduceRequestsInflightPerBroker(1000000),
		kgo.DisableIdempotentWrite(),
		kgo.WithLogger(kzerolog.New(&l)))
	if err != nil {
		return nil, nil, err
	}
	return &Client{
		client: client,
		logger: logger,
	}, func() { client.Close() }, nil
}

var p = sync.Pool{New: func() any { return &kgo.Record{} }}

func (c *Client) Produce(ctx context.Context, topic string, fn func(*bytebufferpool.ByteBuffer) error) {
	b := bytebufferpool.Get()
	defer bytebufferpool.Put(b)
	if err := fn(b); err != nil {
		c.logger.Err(err).Str("topic", topic).Msg("produce data error")
		return
	}
	record := p.Get().(*kgo.Record)
	if b.Len() > len(record.Value) {
		record.Value = append(record.Value, make([]byte, b.Len()-len(record.Value))...)
	} else {
		record.Value = record.Value[:b.Len()]
	}
	copy(record.Value, b.Bytes())
	c.client.Produce(ctx, record, func(r *kgo.Record, err error) {
		p.Put(r)
		if err != nil && !errors.Is(err, context.Canceled) {
			c.logger.Err(err).Str("topic", r.Topic).Time("timestamp", r.Timestamp).Msg("record produce error")
		}
	})
}

func (c *Client) Consume(ctx context.Context, topics []string, fn func(*kgo.Record)) {
	c.client.AddConsumeTopics(topics...)
	for {
		fetches := c.client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				c.logger.Error().
					Err(err.Err).
					Str("topic", err.Topic).
					Int32("partition", err.Partition).
					Msg("consume error")
			}
		}
		fetches.EachPartition(func(p kgo.FetchTopicPartition) {
			p.EachRecord(fn)
		})
	}
}

var ProviderSet = wire.NewSet(NewOptions, NewClient)
