package events

import "context"

type KafkaConfig struct {
	Brokers []string
}

type movieDeletedEvent struct {
	ID int32 `json:"movie_id"`
}

type deleteImageRequest struct {
	ID, Category string
}

//go:generate mockgen -source=events.go -destination=mocks/events.go
type MoviesEventsMQ interface {
	MovieDeleted(ctx context.Context, id int32) error
}

type ImagesEventsMQ interface {
	DeleteImageRequest(ctx context.Context, id string, category string) error
}
