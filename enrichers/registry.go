package enrichers

type Registry map[string]Enricher

type Enricher func(output *enrichedData , data map[string][]byte, params map[string]string) (*enrichedData, error)

func (r Registry) Add(name string, enricher Enricher) {
	r[name] = enricher
}

func (r Registry) Get(name string) (Enricher, bool) {
	e, ok := r[name]
	return e, ok
}
