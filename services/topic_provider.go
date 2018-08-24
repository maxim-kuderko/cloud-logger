package services

import (
	"github.com/maxim-kuderko/cloud-logger/initializers"
	"github.com/maxim-kuderko/storage-buffer"
	"time"
	"sync"
)

type TopicProvider struct {
	db          *initializers.Db
	maxTopicAge time.Duration
	cache       map[string]*TopicOptionWrapper
	m           sync.RWMutex
}

type TopicOptionWrapper struct {
	opt       *storage_buffer.TopicOptions
	updatedAt time.Time
	m         sync.RWMutex
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
	tow.m.RUnlock()
	defer tow.m.RUnlock()
	if time.Now().Sub(tow.updatedAt) >= maxTopicAge {
		return true
	}
	return false
}

func NewTopicProvider(db *initializers.Db, maxTopicAge time.Duration) *TopicProvider {
	return &TopicProvider{
		db:          db,
		maxTopicAge: maxTopicAge,
	}
}

func (tp *TopicProvider) Provide(topicName string) (*storage_buffer.TopicOptions, error) {
	tow, err := tp.loadOrStore(topicName)
	if err != nil && tow == nil{
		return nil, err
	}
	return tow.options(), err
}

func (tp *TopicProvider) loadOrStore(topicName string) (*TopicOptionWrapper, error) {
	v, ok := tp.safeRead(topicName)
	if !ok {
		v, err := tp.safeInitTow(topicName)
		if err != nil {
			return v, err
		}
	}
	if v.shouldUpdate(tp.maxTopicAge) {
		to, err := tp.getFromDb(topicName)
		if err != nil {
			return v, err
		}
		v.update(to)
	}
	return v, nil
}

func (tp *TopicProvider) safeRead(key string) (*TopicOptionWrapper, bool) {
	tp.m.RLock()
	defer tp.m.RUnlock()
	v, ok := tp.cache[key]
	return v, ok

}

func (tp *TopicProvider) safeInitTow(key string) (*TopicOptionWrapper, error){
	tp.m.Lock()
	defer tp.m.Unlock()
	v, ok := tp.cache[key]
	if ok {
		return v, nil
	}
	to, err := tp.getFromDb(key)
	if err != nil {
		return nil, err
	}
	v = &TopicOptionWrapper{
		opt:       to,
		updatedAt: time.Now(),
	}
	tp.cache[key] = v
	return v, nil
}

func (tp *TopicProvider) getFromDb(key string) (*storage_buffer.TopicOptions, error) {

}
