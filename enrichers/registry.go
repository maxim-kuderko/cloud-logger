package enrichers

type Registry map[string]Enricher

type Enricher func(output []byte, data map[string][]byte, params map[string]string) ([]byte, error)

func (r Registry) Add(name string, enricher Enricher) {
	r[name] = enricher
}

func (r Registry) Get(name string) (Enricher, bool) {
	e, ok := r[name]
	return e, ok
}
