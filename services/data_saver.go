package services

import "github.com/maxim-kuderko/storage-buffer"

type DataSaver struct {
	provider *TopicProvider
	stg      *storage_buffer.Collection
	en       *Enricher
}

func NewDataSaver(provider *TopicProvider, enricher *Enricher) *DataSaver {
	return &DataSaver{
		provider: provider,
		stg:      storage_buffer.NewCollection(1024*1024*1024*10, provider.provide),
		en:       enricher,
	}
}

func (ds *DataSaver) Write(topic string, headers map[string]string, data []byte) error {
	// find topic
	// enrich
	// push to storage
	return nil
}
