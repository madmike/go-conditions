package conditions

import "testing"

func TestNestedLogicalOperators(t *testing.T) {
	cond := NewConditions()
	instance := map[string]any{
		"person": map[string]any{
			"age": 25,
			"name": map[string]any{
				"first": "John",
				"last":  "Doe",
			},
			"tags": []string{"student", "active"},
		},
		"status": "active",
	}

	condition := map[string]any{
		"$and": []any{
			map[string]any{"{{person.age}}": map[string]any{"$gte": 18}},
			map[string]any{
				"$or": []map[string]any{
					{"{{person.name.first}}": map[string]any{"$eq": "Jane"}},
					{"{{person.name.last}}": map[string]any{"$eq": "Doe"}},
				},
			},
			map[string]any{"{{person.tags}}": map[string]any{"$has": "student"}},
		},
	}

	want := true
	got := cond.Check(instance, condition)
	if got != want {
		t.Errorf("Check() = %v, want %v", got, want)
	}
}

func TestComplexConditionsWithBetweenAndSome(t *testing.T) {
	cond := NewConditions()
	instance := map[string]any{
		"metrics": map[string]any{
			"temperature": 72,
			"humidity":    40,
			"readings":    []int{100, 200, 300},
		},
		"identifiers": []string{"X123", "Y456"},
	}

	condition := map[string]any{
		"$and": []any{
			map[string]any{"{{metrics.temperature}}": map[string]any{"$between": []int{70, 75}}},
			map[string]any{"{{metrics.humidity}}": map[string]any{"$lt": 50}},
			map[string]any{
				"$or": []map[string]any{
					{"{{identifiers}}": map[string]any{"$some": "X123"}},
					{"{{metrics.readings}}": map[string]any{"$every": []int{100, 200}}},
				},
			},
		},
	}

	want := true
	got := cond.Check(instance, condition)
	if got != want {
		t.Errorf("Check() = %v, want %v", got, want)
	}
}
