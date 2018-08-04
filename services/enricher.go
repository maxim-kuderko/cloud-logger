package services

import (
	"github.com/maxim-kuderko/cloud-logger/enrichers"
	"fmt"
)

type Enricher struct {
	eReg *enrichers.Registry
}

func NewEnricher() *Enricher {
	rg := registerEnrichers()
	return &Enricher{eReg: rg}
}

func registerEnrichers() *enrichers.Registry {
	r := &enrichers.Registry{}
	r.Add("set-body", enrichers.SetBody)
	r.Add("add-headers", enrichers.AddHeaders)
	return r
}

func (e *Enricher) Do(data map[string][]byte, params []map[string]string) ([]byte, error) {
	var output []byte
	for _, enrParams := range params {
		enr, ok := e.eReg.Get(enrParams[`name`])
		if !ok {
			return nil, fmt.Errorf("enricher %s not found", enrParams[`name`])
		}
		out, err := enr(output, data, enrParams)
		if err != nil {
			return nil, fmt.Errorf("enricher error %s", err)
		}
		output = out
	}
	return output, nil
}
