package registry

import (
	"github.com/maxim-kuderko/cloud-logger/initializers"
	"github.com/maxim-kuderko/cloud-logger/services"
)

// This is a static registry, initialized at boot
type Registry struct {
	Ds *services.DataSaver
}

func NewRegistry() *Registry {
	db := initializers.NewDb()
	ds, en := initServices(db)
	return &Registry{
		Ds: ds,
	}
}

func initServices(db *initializers.Db) (*services.DataSaver, *services.Enricher) {
	tp := services.NewTopicProvider(db)
	en := services.NewEnricher()
	ds := services.NewDataSaver(tp)
	return ds, en
}
