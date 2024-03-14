package conditions

import (
	"testing"
)

func TestCheckSimpleOperators(t *testing.T) {
	cond := NewConditions()

	tests := []struct {
		name      string
		condition map[string]any
		instance  any
		want      bool
	}{
		{
			name: "Test NULL operator true",
			condition: map[string]any{
				"$null": "value",
			},
			instance: map[string]any{
				"value": nil,
			},
			want: true,
		},
		{
			name: "Test DEFINED operator false",
			condition: map[string]any{
				"$defined": "missingValue",
			},
			instance: map[string]any{},
			want:     false,
		},
		{
			name: "Test UNDEFINED operator true",
			condition: map[string]any{
				"$undefined": "undefinedValue",
			},
			instance: map[string]any{
				"value": "something",
			},
			want: true,
		},
		{
			name: "Test EXIST operator true",
			condition: map[string]any{
				"$exist": "value",
			},
			instance: map[string]any{
				"value": "I exist",
			},
			want: true,
		},
		{
			name: "Test EMPTY operator true",
			condition: map[string]any{
				"$empty": "emptyArray",
			},
			instance: map[string]any{
				"emptyArray": []int{},
			},
			want: true,
		},
		{
			name: "Test BLANK operator true for nil value",
			condition: map[string]any{
				"$blank": "nilValue",
			},
			instance: map[string]any{
				"nilValue": nil,
			},
			want: true,
		},
		{
			name: "Test TRULY operator true",
			condition: map[string]any{
				"$truly": "trueValue",
			},
			instance: map[string]any{
				"trueValue": true,
			},
			want: true,
		},
		{
			name: "Test FALSY operator true",
			condition: map[string]any{
				"$falsy": "falseValue",
			},
			instance: map[string]any{
				"falseValue": false,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cond.Check(tt.instance, tt.condition); got != tt.want {
				t.Errorf("%s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
