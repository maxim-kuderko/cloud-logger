package services

import (
	"github.com/maxim-kuderko/cloud-logger/enrichers"
	"github.com/maxim-kuderko/storage-buffer"
	"time"
)

type DataSaver struct {
	provider *TopicProvider
	stg      *storage_buffer.Collection
	en       *enrichers.EnricherManager
}

func NewDataSaver(provider *TopicProvider, enricher *enrichers.EnricherManager) *DataSaver {
	return &DataSaver{
		provider: provider,
		stg:      storage_buffer.NewCollection(1024 * 1024 * 1024 * 10, time.Minute * 2),
		en:       enricher,
	}
}

func (ds *DataSaver) Write(topic string, headers []byte, body []byte) error {
	// find topic
	t, err := ds.provider.Provide(topic)
	if err != nil {
		return err
	}
	// enrich
	data, err := ds.en.Do(map[string][]byte{"body": body, "headers": headers}, []map[string]string{{"name": "set-body", "opt": "val"}})
	if err != nil {
		return err
	}
	// push to storage
	b, err := data.Serialize(`json`)
	if err != nil{
		return err
	}
	_, err = ds.stg.Write(t,data.Partitions, b)
	if err != nil {
		return err
	}
	return nil
}


func (ds *DataSaver) Close(){
	ds.stg.Shutdown()
}