package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type imagesEvents struct {
	eventsWriter *kafka.Writer
	logger       *logrus.Logger
}

func NewImagesEvents(cfg KafkaConfig, logger *logrus.Logger) *imagesEvents {
	w := &kafka.Writer{
		Addr:   kafka.TCP(cfg.Brokers...),
		Logger: logger,
	}
	w.AllowAutoTopicCreation = true
	return &imagesEvents{eventsWriter: w, logger: logger}
}

const (
	DeleteImageRequestTopic = "delete_image_request"
)

func (e *imagesEvents) Shutdown() error {
	return e.eventsWriter.Close()
}

func (e *imagesEvents) DeleteImageRequest(ctx context.Context, id string, category string) error {
	body, err := json.Marshal(deleteImageRequest{ID: id, Category: category})
	if err != nil {
		e.logger.Fatal(err)
	}
	return e.eventsWriter.WriteMessages(ctx, kafka.Message{
		Topic: DeleteImageRequestTopic,
		Key:   []byte(fmt.Sprint("image_", id)),
		Value: body,
	})
}
