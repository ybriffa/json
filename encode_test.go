package json

import (
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	var tests = []struct {
		input       interface{}
		expected    []byte
		shouldCrash bool
	}{
		//0 no crash, not a map
		{
			input:       42,
			expected:    []byte(`42`),
			shouldCrash: false,
		},

		//1 empty map
		{
			input:       map[string]int{},
			expected:    []byte(`{}`),
			shouldCrash: false,
		},

		//2 int map
		{
			input: map[int]int{
				42: 4,
				13: 2,
				5:  1,
				32: 3,
			},
			expected:    []byte(`{"5":1,"13":2,"32":3,"42":4}`),
			shouldCrash: false,
		},

		//3 string map
		{
			input: map[string]int{
				"CCCD":   4,
				"BBB":    2,
				"AAAAAA": 1,
				"CCCCC":  3,
			},
			expected:    []byte(`{"AAAAAA":1,"BBB":2,"CCCCC":3,"CCCD":4}`),
			shouldCrash: false,
		},

		//4 float map
		{
			input: map[float64]int{
				2.42: 3,
				1.56: 2,
				0.5:  1,
				3.32: 4,
			},
			expected:    []byte(`{"0.5":1,"1.56":2,"2.42":3,"3.32":4}`),
			shouldCrash: false,
		},
	}

	for i, test := range tests {
		result, err := Marshal(test.input)
		if (err != nil) != test.shouldCrash {
			t.Fatalf("test #%d failed : expected to crash `%t` and got %s", i, test.shouldCrash, err)
		}

		if test.shouldCrash {
			continue
		}

		if !reflect.DeepEqual(result, test.expected) {
			t.Fatalf("test #%d failed :\nexpected : `%s`\n     got : `%s`", i, string(test.expected), string(result))
		}
	}
}
