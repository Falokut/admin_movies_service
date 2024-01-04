package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type moviesEvents struct {
	eventsWriter *kafka.Writer
	logger       *logrus.Logger
}

func NewMoviesEvents(cfg KafkaConfig, logger *logrus.Logger) *moviesEvents {
	w := &kafka.Writer{
		Addr:   kafka.TCP(cfg.Brokers...),
		Logger: logger,
	}
	w.AllowAutoTopicCreation = true
	return &moviesEvents{eventsWriter: w, logger: logger}
}

const (
	movieDeletedTopic = "movie_deleted"
)

func (e *moviesEvents) Shutdown() error {
	return e.eventsWriter.Close()
}

func (e *moviesEvents) MovieDeleted(ctx context.Context, id int32) error {
	body, err := json.Marshal(movieDeletedEvent{ID: id})
	if err != nil {
		e.logger.Fatal(err)
	}
	return e.eventsWriter.WriteMessages(ctx, kafka.Message{
		Topic: movieDeletedTopic,
		Key:   []byte(fmt.Sprint("movie_", id)),
		Value: body,
	})
}
