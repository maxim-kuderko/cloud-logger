package services

import (
	"fmt"
	"github.com/maxim-kuderko/cloud-logger/initializers"
	"github.com/maxim-kuderko/storage-buffer"
)

type TopicProvider struct {
}

func NewTopicProvider(db *initializers.Db) *TopicProvider {
	return &TopicProvider{}
}

func (tp *TopicProvider) provide(topicName string) (*storage_buffer.TopicOptions, error) {
	return &storage_buffer.TopicOptions{}, fmt.Errorf("topic not found: %s", topicName)
}
