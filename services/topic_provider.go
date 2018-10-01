package services

import (
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/maxim-kuderko/cloud-logger/initializers"
	"github.com/maxim-kuderko/storage-buffer"
	"github.com/maxim-kuderko/storage-buffer/storage"
	"io"
	"sync"
	"time"
)

type TopicProvider struct {
	db          *initializers.Db
	maxTopicAge time.Duration
	cache       map[string]*TopicOptionWrapper
	query       *sql.Stmt
	m           sync.RWMutex
}

type TopicOptionWrapper struct {
	opt           *storage_buffer.TopicOptions
	partitionsDef []string
	callbackUrl   string
	updatedAt     time.Time
	m             sync.RWMutex
}

func (tow *TopicOptionWrapper) options() *storage_buffer.TopicOptions {
	tow.m.RLock()
	defer tow.m.RUnlock()
	return tow.opt
}

func (tow *TopicOptionWrapper) update(opt *storage_buffer.TopicOptions) {
	tow.m.Lock()
	defer tow.m.Unlock()
	tow.opt = opt
	tow.updatedAt = time.Now()
}

func (tow *TopicOptionWrapper) shouldUpdate(maxTopicAge time.Duration) bool {
	tow.m.RLock()
	defer tow.m.RUnlock()
	if time.Now().Sub(tow.updatedAt) >= maxTopicAge {
		return true
	}
	return false
}

func NewTopicProvider(db *initializers.Db, maxTopicAge time.Duration) *TopicProvider {
	stmt, _ := db.Prepare("SELECT `id`, `name`, `max_length`, `max_size`, `interval`, `storage_driver`, `storage_creds`, `partitions`, `updated_at`")
	return &TopicProvider{
		db:          db,
		maxTopicAge: maxTopicAge,
		cache:       make(map[string]*TopicOptionWrapper),
		query:       stmt,
	}
}

func (tp *TopicProvider) Provide(topicId string) (*storage_buffer.TopicOptions, error) {
	tow, err := tp.loadOrStore(topicId)
	if err != nil && tow == nil {
		return nil, err
	}
	return tow.options(), err
}

func (tp *TopicProvider) loadOrStore(topicId string) (*TopicOptionWrapper, error) {
	var err error
	v, ok := tp.safeRead(topicId)
	if !ok || v.shouldUpdate(tp.maxTopicAge) {
		v, err = tp.safeInitTow(topicId)
		if err != nil {
			return v, err
		}
	}
	return v, nil
}

func (tp *TopicProvider) safeRead(key string) (*TopicOptionWrapper, bool) {
	tp.m.RLock()
	defer tp.m.RUnlock()
	v, ok := tp.cache[key]
	return v, ok

}

func (tp *TopicProvider) safeInitTow(key string) (*TopicOptionWrapper, error) {
	tp.m.Lock()
	defer tp.m.Unlock()
	v, ok := tp.cache[key]
	if ok && !v.shouldUpdate(tp.maxTopicAge) {
		return v, nil
	}
	to, err := tp.getFromDb(key)
	if err != nil {
		return v, err
	}
	tp.cache[key] = to
	return to, nil
}

func (tp *TopicProvider) getFromDb(key string) (*TopicOptionWrapper, error) {
	res, err := tp.query.Query(key)
	if err != nil {
		return nil, err
	}
	optw := TopicOptionWrapper{
		opt: &storage_buffer.TopicOptions{},
	}
	var t, creds []byte
	res.Scan()
	return &TopicOptionWrapper{

	}, nil
}

func (tp *TopicProvider) defineStorageDriver(t []byte, creds []byte, name string, callback func(resp *s3manager.UploadOutput, err error)) (func(partition []string) io.WriteCloser, error) {
	switch string(t) {
	case `s3`:
		s := storage.NewS3Loader(name, ``, ``,``,`json`,``,``,callback)
		return s.S3Store, nil
	}
	return nil, fmt.Errorf("unkwon storage driver %s", t)
}
