package conditions

import (
	"testing"
)

// TestLogicOperators tests the logic operators like $and, $or, $not, $xor.
func TestLogicOperators(t *testing.T) {
	cond := NewConditions()

	tests := []struct {
		name      string
		instance  any
		condition map[string]any
		want      bool
	}{
		{
			name: "Test $and operator true",
			instance: map[string]any{
				"age":  25,
				"name": "John",
			},
			condition: map[string]any{
				"$and": []map[string]any{
					{"{{age}}": map[string]any{"$gt": 20}},
					{"{{name}}": map[string]any{"$eq": "John"}},
				},
			},
			want: true,
		},
		{
			name: "Test $or operator true",
			instance: map[string]any{
				"age":  18,
				"name": "Jane",
			},
			condition: map[string]any{
				"$or": []map[string]any{
					{"{{age}}": map[string]any{"$lt": 20}},
					{"{{name}}": map[string]any{"$eq": "John"}},
				},
			},
			want: true,
		},
		{
			name: "Test $not operator true",
			instance: map[string]any{
				"age": 18,
			},
			condition: map[string]any{
				"$not": []map[string]any{
					{"{{age}}": map[string]any{"$gt": 20}},
				},
			},
			want: true,
		},
		{
			name: "Test $xor operator true",
			instance: map[string]any{
				"age":  22,
				"name": "Jane",
			},
			condition: map[string]any{
				"$xor": []map[string]any{
					{"{{age}}": map[string]any{"$gt": 20}},
					{"{{name}}": map[string]any{"$eq": "John"}},
				},
			},
			want: true,
		},
		// Add more tests to cover false scenarios and edge cases...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cond.Check(tt.instance, tt.condition); got != tt.want {
				t.Errorf("%s: Check() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
