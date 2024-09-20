package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type SelectorSpec struct {
	// Selector of an attribute from the contextual properties provided by kuadrant
	// during request and connection processing
	Selector string `json:"selector"`
}

type StaticSpec struct {
	Key string `json:"key"`
}

type Static struct {
	Static StaticSpec `json:"static"`
}

type Selector struct {
	Selector SelectorSpec `json:"selector"`
}

type DataType struct {
	Value interface{}
}

func (d *DataType) UnmarshalJSON(data []byte) error {
	//fmt.Println("UnmarshalJSON")
	//fmt.Println(string(data))

	types := []interface{}{
		&Static{},
		&Selector{},
	}

	var err error

	for idx := range types {
		//fmt.Printf("unmarshalling %+v\n", types[idx])
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.DisallowUnknownFields() // Force errors
		err = dec.Decode(&types[idx])
		if err == nil {
			d.Value = types[idx]
			return nil
		}
	}

	return err
}

type Conf struct {
	Name string     `json:"name"`
	Data []DataType `json:"data"`
}

func main() {

	tests := []struct {
		name    string
		valid   bool
		jsonStr string
	}{
		{
			"validemptyJson",
			true,
			`
{
	"name": "validJson",
	"data": []
}
`,
		},
		{
			"validJson",
			true,
			`
{
	"name": "validJson",
	"data": [
		{
			"static": { "key": "keyA" }
		},
		{
			"selector": { "selector": "selectorA" }
		},
		{
			"static": { "key": "keyB" }
		},
		{
			"selector": { "selector": "selectorB" }
		}
	]
}
`,
		},
		{
			"invalidJsonA",
			false,
			`
{
	"name": "invalidJsonA",
	"data": [
		{
			"other": { "key": "keyA" }
		},
		{
			"selector": { "selector": "selectorA" }
		}
	]
}
`,
		},
		{
			"invalidJsonB",
			false,
			`
{
	"name": "invalidJsonB",
	"data": [
		{
			"static": { "key": "keyA" },
			"selector": { "selector": "selectorA" }
		}
	]
}
`,
		},
	}

	for _, test := range tests {
		var conf Conf
		err := json.Unmarshal([]byte(test.jsonStr), &conf)
		fmt.Printf("testing %s: %t, err: %v\n", test.name, test.valid == (err == nil), err)
	}
}
