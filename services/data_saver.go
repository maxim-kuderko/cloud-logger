package services

import (
	"github.com/maxim-kuderko/storage-buffer"
	"log"
)

type DataSaver struct {
	provider *TopicProvider
	stg      *storage_buffer.Collection
	en       *Enricher
}

func NewDataSaver(provider *TopicProvider, enricher *Enricher) *DataSaver {
	return &DataSaver{
		provider: provider,
		stg:      storage_buffer.NewCollection(1024 * 1024 * 1024 * 10),
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
	log.Println(string(data))
	_, err = ds.stg.Write(t, data)
	if err != nil {
		return err
	}
	return nil
}


func (ds *DataSaver) Close(){
	ds.stg.Shutdown()
}