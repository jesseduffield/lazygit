// Copyright (c) 2015, Emir Pasic. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package treebidimap

import (
	"encoding/json"
	"github.com/emirpasic/gods/containers"
	"github.com/emirpasic/gods/utils"
)

func assertSerializationImplementation() {
	var _ containers.JSONSerializer = (*Map)(nil)
	var _ containers.JSONDeserializer = (*Map)(nil)
}

// ToJSON outputs the JSON representation of the map.
func (m *Map) ToJSON() ([]byte, error) {
	elements := make(map[string]interface{})
	it := m.Iterator()
	for it.Next() {
		elements[utils.ToString(it.Key())] = it.Value()
	}
	return json.Marshal(&elements)
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
