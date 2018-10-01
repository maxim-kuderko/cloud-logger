package enrichers

import (
	"encoding/json"
	"fmt"
)

type EnricherManager struct {
	eReg *Registry
}

func NewEnricherManager() *EnricherManager {
	rg := registerEnrichers()
	return &EnricherManager{eReg: rg}
}

func registerEnrichers() *Registry {
	r := &Registry{}
	r.Add("add-body-json", SetBody)
	r.Add("add-headers", AddHeaders)
	return r
}

func (e *EnricherManager) Do(data map[string][]byte, params []map[string]string) (*enrichedData, error) {
	var output enrichedData
	for _, enrParams := range params {
		enr, ok := e.eReg.Get(enrParams[`name`])
		if !ok {
			return nil, fmt.Errorf("enricher %s not found", enrParams[`name`])
		}
		_, err := enr(&output, data, enrParams)
		if err != nil {
			return nil, fmt.Errorf("enricher error %s", err)
		}
	}
	return &output, nil
}

type enrichedData struct {
	Body    string `json:"body"`
	Headers string `json:"headers"`
	Partitions []string
}

func (ed *enrichedData) Serialize(serializer string) ([]byte, error){
	switch serializer {
	case `json`:
		return json.Marshal(ed)
	}
	return nil, fmt.Errorf("unknown serializer %s", serializer)
}
