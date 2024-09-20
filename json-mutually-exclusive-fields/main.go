package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
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
	Value interface{} `json:,inline`
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
		err = dec.Decode(types[idx])
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

		if err == nil && test.valid {
			for idx := range conf.Data {
				switch val := conf.Data[idx].Value.(type) {
				case *Static:
					fmt.Printf("testing %s: is static key: %s\n", test.name, val.Static.Key)
				case *Selector:
					fmt.Printf("testing %s: is selector key: %s\n", test.name, val.Selector.Selector)
				default:
					panic("should not happen")
				}
			}
		}

		serialized, err := json.Marshal(conf)
		fmt.Printf("marshal error: %v\n", err)
		fmt.Println(string(serialized))

	}

	fmt.Println("== compare method")
	config := Conf{
		Name: "validJson",
		Data: []DataType{
			{Value: &Static{Static: StaticSpec{Key: "keyA"}}},
			{Value: &Selector{Selector: SelectorSpec{Selector: "selectorA"}}},
			{Value: &Static{Static: StaticSpec{Key: "keyB"}}},
			{Value: &Selector{Selector: SelectorSpec{Selector: "selectorB"}}},
		},
	}

	var unmarshaledConfig Conf
	configMarshaled := []byte(`
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
`)
	if err := json.Unmarshal(configMarshaled, &unmarshaledConfig); err != nil {
		panic(err)
	}

	if !cmp.Equal(unmarshaledConfig, config) {
		diff := cmp.Diff(unmarshaledConfig, config)
		fmt.Println("Diff")
		fmt.Println(diff)
	}

	fmt.Println("no diff 🎉")
}
