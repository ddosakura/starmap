package check

import (
	"encoding/json"
)

// Model for check
type Model map[string]interface{}

// BaseLoader for entity
type BaseLoader struct {
	entity Model
	Err    error
}

// Load for entity
func (b *BaseLoader) Load(es ...interface{}) error {
	if b.Err != nil {
		return b.Err
	}
	for _, e := range es {
		if bs, err := json.Marshal(b.entity); err != nil {
			return ErrAuthLoadCrash
		} else if err = json.Unmarshal(bs, e); err != nil {
			return ErrAuthLoadCrash
		}
	}
	return nil
}
