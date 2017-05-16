package json

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sort"
)

func Marshal(v interface{}) ([]byte, error) {
	reflected := reflect.ValueOf(v)

	for reflected.Type().Kind() == reflect.Ptr {
		reflected = reflect.Indirect(reflected)
	}

	if reflected.Type().Kind() != reflect.Map {
		return json.Marshal(v)
	}

	reflectedKeys := reflected.MapKeys()
	if len(reflectedKeys) == 0 {
		return []byte("{}"), nil
	}
	keys := reflect.MakeSlice(reflect.SliceOf(reflectedKeys[0].Type()), len(reflectedKeys), len(reflectedKeys))

	for i, reflectedKey := range reflectedKeys {
		keys.Index(i).Set(reflectedKey)
	}

	switch k := keys.Interface().(type) {
	case []int:
		sort.Ints(k)
	case []string:
		sort.Strings(k)
	case []float64:
		sort.Float64s(k)
	case sort.Interface:
		sort.Sort(k)
	}

	var buffer bytes.Buffer
	buffer.Write([]byte("{"))
	for i := 0; i < keys.Len(); i++ {
		if i > 0 {
			buffer.Write([]byte(","))
		}
		keyMarshalled, err := json.Marshal(keys.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		buffer.Write([]byte(`"`))
		buffer.Write(bytes.Replace(keyMarshalled, []byte(`"`), []byte(`\\"`), -1))
		buffer.Write([]byte(`"`))
		buffer.Write([]byte(":"))
		valueMarshalled, err := Marshal(reflected.MapIndex(keys.Index(i)).Interface())
		if err != nil {
			return nil, err
		}
		buffer.Write(valueMarshalled)

	}
	buffer.Write([]byte("}"))

	return buffer.Bytes(), nil
}

func MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	// carbon copy of json.MarshalIndent,
	// just using our function instead of json.Marshal
	b, err := Marshal(v)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, b, prefix, indent)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
