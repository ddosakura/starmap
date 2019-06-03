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
func (b *BaseLoader) Load(e interface{}) error {
	if b.Err != nil {
		return b.Err
	}
	if bs, err := json.Marshal(b.entity); err != nil {
		return ErrAuthLoadCrash
	} else if err = json.Unmarshal(bs, e); err != nil {
		return ErrAuthLoadCrash
	}
	return nil
}
