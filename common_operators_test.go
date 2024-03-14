package conditions

import (
	"testing"
)

func TestCheckCommonOperators(t *testing.T) {
	cond := NewConditions()

	tests := []struct {
		name      string
		condition map[string]any
		instance  any
		want      bool
	}{
		{
			name: "Test $eq operator true",
			condition: map[string]any{
				"{{age}}": map[string]any{"$eq": 30},
			},
			instance: map[string]any{
				"age": 30,
			},
			want: true,
		},
		{
			name: "Test $ne operator false",
			condition: map[string]any{
				"{{age}}": map[string]any{"$ne": 30},
			},
			instance: map[string]any{
				"age": 30,
			},
			want: false,
		},
		{
			name: "Test $lt operator true",
			condition: map[string]any{
				"{{age}}": map[string]any{"$lt": 31},
			},
			instance: map[string]any{
				"age": 30,
			},
			want: true,
		},
		{
			name: "Test $gt operator false",
			condition: map[string]any{
				"{{age}}": map[string]any{"$gt": 30},
			},
			instance: map[string]any{
				"age": 30,
			},
			want: false,
		},
		{
			name: "Test $lte operator true",
			condition: map[string]any{
				"{{age}}": map[string]any{"$lte": 30},
			},
			instance: map[string]any{
				"age": 30,
			},
			want: true,
		},
		{
			name: "Test $gte operator true",
			condition: map[string]any{
				"{{age}}": map[string]any{"$gte": 30},
			},
			instance: map[string]any{
				"age": 29,
			},
			want: false,
		},
		{
			name: "Test $in operator true",
			condition: map[string]any{
				"{{status}}": map[string]any{"$in": []string{"active", "inactive"}},
			},
			instance: map[string]any{"status": "active"},
			want:     true,
		},
		{
			name: "Test $ni operator true",
			condition: map[string]any{
				"{{tag}}": map[string]any{"$ni": []string{"go", "golang", "programming"}},
			},
			instance: map[string]any{"tag": "missing"},
			want:     true,
		},
		{
			name: "Test $re operator true",
			condition: map[string]any{
				"{{name}}": map[string]any{"$re": "^J.*"},
			},
			instance: map[string]any{
				"name": "John",
			},
			want: true,
		},
		{
			name: "Test $sw operator true",
			condition: map[string]any{
				"{{name}}": map[string]any{"$sw": "Jo"},
			},
			instance: map[string]any{
				"name": "John",
			},
			want: true,
		},
		{
			name: "Test $ew operator true",
			condition: map[string]any{
				"{{file}}": map[string]any{"$ew": ".txt"},
			},
			instance: map[string]any{
				"file": "document.txt",
			},
			want: true,
		},
		{
			name: "Test $incl operator true",
			condition: map[string]any{
				"{{values}}": map[string]any{"$incl": 3},
			},
			instance: map[string]any{
				"values": []int{1, 2, 3, 4},
			},
			want: true,
		},
		{
			name: "Test $excl operator true",
			condition: map[string]any{
				"{{values}}": map[string]any{"$excl": 5},
			},
			instance: map[string]any{
				"values": []int{1, 2, 3, 4},
			},
			want: true,
		},
		{
			name: "Test $has operator true",
			condition: map[string]any{
				"{{values}}": map[string]any{"$has": "key1"},
			},
			instance: map[string]any{
				"values": map[string]int{"key1": 1, "key2": 2},
			},
			want: true,
		},
		{
			name: "Test $power operator true",
			condition: map[string]any{
				"{{value}}": map[string]any{"$power": 2}, // Example checks if value has a bit set (binary AND operation), here checking for 2^1
			},
			instance: map[string]any{
				"value": 3, // Binary 11, has 2^1 set
			},
			want: true,
		},
		{
			name: "Test $between operator true",
			condition: map[string]any{
				"{{age}}": map[string]any{"$between": []int{25, 35}},
			},
			instance: map[string]any{
				"age": 30,
			},
			want: true,
		},
		{
			name: "Test $some operator true",
			condition: map[string]any{
				"{{values}}": map[string]any{"$some": []int{5, 6}},
			},
			instance: map[string]any{
				"values": []int{1, 5, 9},
			},
			want: true,
		},
		{
			name: "Test $every operator false",
			condition: map[string]any{
				"{{values}}": map[string]any{"$every": []int{1, 5, 9}},
			},
			instance: map[string]any{
				"values": []int{1, 5},
			},
			want: false,
		},
		{
			name: "Test $noone operator true",
			condition: map[string]any{
				"{{values}}": map[string]any{"$noone": []int{10, 11}},
			},
			instance: map[string]any{
				"values": []int{1, 5, 9},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cond.Check(tt.instance, tt.condition); got != tt.want {
				t.Errorf("Check() for %s = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
