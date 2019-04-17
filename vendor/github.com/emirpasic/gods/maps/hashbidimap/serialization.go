// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hashbidimap

import (
	"encoding/json"
	"github.com/emirpasic/gods/containers"
)

func assertSerializationImplementation() {
	var _ containers.JSONSerializer = (*Map)(nil)
	var _ containers.JSONDeserializer = (*Map)(nil)
}

// ToJSON outputs the JSON representation of the map.
func (m *Map) ToJSON() ([]byte, error) {
	return m.forwardMap.ToJSON()
}

// FromJSON populates the map from the input JSON representation.
func (m *Map) FromJSON(data []byte) error {
	elements := make(map[string]interface{})
	err := json.Unmarshal(data, &elements)
	if err == nil {
		m.Clear()
		for key, value := range elements {
			m.Put(key, value)
		}
	}
	return err
}
