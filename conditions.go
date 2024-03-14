package conditions

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

func (c *Conditions) Check(instance any, condition any) bool {
	// Recursively check conditions if it's a slice, treating it as an AND condition
	if conditions, ok := condition.([]any); ok {
		for _, cond := range conditions {
			if !c.Check(instance, cond) {
				return false
			}
		}
		return true
	}

	// Handle condition maps
	if condMap, ok := condition.(map[string]any); ok {
		for key, value := range condMap {
			valueKind := reflect.ValueOf(value).Kind()
			if operator, exists := stringToSimpleOperator[key]; exists {
				return c.checkSimpleOperator(operator, value, instance)
			} else if operator, exists := stringToLogicOperator[key]; exists {
				return c.checkLogicOperator(operator, value, instance)
			} else if valueKind == reflect.Map || valueKind == reflect.Struct {
				result, err := c.checkCommonOperator(key, value, instance)
				return err == nil && result
			} else {
				leftSide := c.getValueByTemplate(key, instance)
				rightSide := c.getValueByTemplate(value, instance)

				// Check for equality
				return reflect.DeepEqual(leftSide, rightSide)
			}
		}
	}

	return false
}

func (c *Conditions) checkSimpleOperator(operator SimpleOperatorsEnum, value any, instance any) bool {
	fact := c.getValueByTemplate(value, instance)
	switch operator {
	case NULL:
		return fact == nil
	case DEFINED:
		return fact != nil
	case UNDEFINED:
		return fact == nil
	case EXIST:
		return fact != nil
	case EMPTY:
		v := reflect.ValueOf(fact)
		switch v.Kind() {
		case reflect.Array, reflect.Slice, reflect.String, reflect.Map:
			return v.Len() == 0
		default:
			// TODO check for empty struct
			return false
		}
	case BLANK:
		if fact == nil || reflect.ValueOf(fact).IsNil() {
			return true
		}
		v := reflect.ValueOf(fact)
		switch v.Kind() {
		case reflect.Array, reflect.Slice, reflect.String:
			return v.Len() == 0
		case reflect.Bool:
			return !v.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return v.Int() == 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return v.Uint() == 0
		case reflect.Float32, reflect.Float64:
			return v.Float() == 0
		case reflect.Map:
			return v.Len() == 0
		default:
			return false // throw an error
		}
	case TRULY:
		return fact == true
	case FALSY:
		return fact == false
	default:
		return false
	}
}

func (c *Conditions) checkCommonOperator(key string, value any, instance any) (bool, error) {
	fact := c.getValueByTemplate(key, instance)

	// Ensure that value is a map containing our conditions.
	conditionMap, ok := value.(map[string]any)
	if !ok {
		return false, fmt.Errorf("expected condition to be a map, got %T", value)
	}

	for operator, conditionValue := range conditionMap {
		switch CommonOperatorsEnum(operator) {
		case EQ:
			if !reflect.DeepEqual(fact, conditionValue) {
				return false, nil
			}
		case NE:
			if reflect.DeepEqual(fact, conditionValue) {
				return false, nil
			}
		case LT, GT, LTE, GTE:
			result, err := compareNumbersOrDates(fact, conditionValue)
			if err != nil {
				return false, err
			}

			switch CommonOperatorsEnum(operator) {
			case LT:
				return result == -1, nil
			case GT:
				return result == 1, nil
			case LTE:
				return result <= 0, nil
			case GTE:
				return result >= 0, nil
			}
		case IN:
			list, ok := reflect.ValueOf(conditionValue).Interface().([]string)
			if !ok {
				return false, fmt.Errorf("expected value to be a string array")
			}
			return contains(list, fact.(string)), nil
		case NI:
			list, ok := reflect.ValueOf(conditionValue).Interface().([]string)
			if !ok {
				return false, fmt.Errorf("expected value to be a string array")
			}
			return !contains(list, fact.(string)), nil
		case RE:
			pattern, ok := conditionValue.(string)
			if !ok {
				return false, fmt.Errorf("expected string for $re operator, got %T", conditionValue)
			}
			re, err := regexp.Compile(pattern)
			if err != nil {
				return false, err
			}
			str, ok := fact.(string)
			if !ok {
				return false, fmt.Errorf("expected string for regex match, got %T", fact)
			}
			if !re.MatchString(str) {
				return false, nil
			}
		case SW:
			prefix, ok := conditionValue.(string)
			if !ok {
				return false, fmt.Errorf("expected string for $sw operator, got %T", conditionValue)
			}
			str, ok := fact.(string)
			if !ok {
				return false, fmt.Errorf("expected string for instance value, got %T", fact)
			}
			return strings.HasPrefix(str, prefix), nil
		case EW:
			prefix, ok := conditionValue.(string)
			if !ok {
				return false, fmt.Errorf("expected string for $ew operator, got %T", conditionValue)
			}
			str, ok := fact.(string)
			if !ok {
				return false, fmt.Errorf("expected string for instance value, got %T", fact)
			}
			return strings.HasSuffix(str, prefix), nil
		case INCL, HAS:
			return isInCollection(fact, conditionValue), nil
		case EXCL:
			return !isInCollection(fact, conditionValue), nil
		case POWER:
			// TODO: update for uint64 case
			numValue, ok := toFloat64(fact)
			if !ok {
				return false, fmt.Errorf("expected numeric type for instance value, got %T", fact)
			}
			powerValue, ok := toFloat64(conditionValue)
			if !ok {
				return false, fmt.Errorf("expected numeric type for condition value, got %T", conditionValue)
			}
			num := int(numValue)
			power := int(powerValue)

			return (num & power) != 0, nil
		case BETWEEN:
			// First, ensure conditionValue can be treated as a slice of any.
			val := reflect.ValueOf(conditionValue)
			if val.Kind() != reflect.Slice || val.Len() != 2 {
				return false, fmt.Errorf("expected condition to be a slice with exactly two elements")
			}

			// Extract the lower and upper bounds as interface{}.
			lowerBound := val.Index(0).Interface()
			upperBound := val.Index(1).Interface()

			// Perform the comparisons.
			compLower, err := compareNumbersOrDates(fact, lowerBound)
			if err != nil || compLower == -1 {
				return false, err // If fact is less than the lower bound, or an error occurred.
			}

			compUpper, err := compareNumbersOrDates(fact, upperBound)
			if err != nil || compUpper == 1 {
				return false, err // If fact is greater than the upper bound, or an error occurred.
			}

			return true, nil
		case SOME:
			factVal := reflect.ValueOf(fact)
			if factVal.Kind() != reflect.Slice {
				return false, fmt.Errorf("expected a slice for instance value under key %s, got %T", key, fact)
			}

			conditionVal := reflect.ValueOf(conditionValue)
			if conditionVal.Kind() != reflect.Slice {
				return false, fmt.Errorf("expected a slice for $some operator, got %T", conditionValue)
			}

			// Iterate over each item in the conditionValue slice.
			for i := 0; i < conditionVal.Len(); i++ {
				conditionItem := conditionVal.Index(i).Interface()

				// Check if conditionItem is in factVal slice.
				for j := 0; j < factVal.Len(); j++ {
					factItem := factVal.Index(j).Interface()

					if reflect.DeepEqual(factItem, conditionItem) {
						return true, nil
					}
				}
			}

			return false, nil // No matching elements found.
		case EVERY, NOONE:
			factVal := reflect.ValueOf(fact)
			if factVal.Kind() != reflect.Slice {
				return false, fmt.Errorf("expected a slice for instance value under key %s, got %T", key, fact)
			}

			conditionVal := reflect.ValueOf(conditionValue)
			if conditionVal.Kind() != reflect.Slice {
				return false, fmt.Errorf("expected a slice for $every operator, got %T", conditionValue)
			}

			// Iterate over each item in the conditionValue slice.
			for i := 0; i < conditionVal.Len(); i++ {
				conditionItem := conditionVal.Index(i).Interface()

				// Check if conditionItem is or is not (for $noone) in factVal slice.
				for j := 0; j < factVal.Len(); j++ {
					factItem := factVal.Index(j).Interface()
					eq := reflect.DeepEqual(factItem, conditionItem)

					if CommonOperatorsEnum(operator) == EVERY && !eq {
						return false, nil
					} else if CommonOperatorsEnum(operator) == NOONE && eq {
						return false, nil
					}
				}
			}

			// All elements in conditionValue exist or don't exist (for $noone) in the factVal slice.
			return true, nil
		default:
			return false, fmt.Errorf("unhandled operator %s", operator)
		}
	}

	return true, nil
}

func (c *Conditions) checkLogicOperator(operator LogicOperatorsEnum, value any, instance any) bool {
	// Convert value to a slice of conditions
	var conditions []map[string]any

	// Use reflection to handle both []any and []map[string]any cases gracefully
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Slice {
		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()

			// Attempt to assert each item's type to a map[string]any
			if cond, ok := item.(map[string]any); ok {
				conditions = append(conditions, cond)
			} else {
				// Handle error or unexpected type
				fmt.Printf("Unexpected type in conditions slice: got %T\n", item)
				return false
			}
		}
	} else if v, ok := value.(map[string]any); ok {
		// If a single map is provided, convert it into an array of maps,
		// each containing one key-value pair from the original map.
		for key, val := range v {
			singleCondition := map[string]any{key: val}
			conditions = append(conditions, singleCondition)
		}
	} else {
		fmt.Printf("Unexpected type for value: got %T\n", value)
		return false
	}

	switch operator {
	case OR:
		for _, cond := range conditions {
			if c.Check(instance, cond) {
				return true
			}
		}
		return false
	case XOR:
		trueCount := 0
		for _, cond := range conditions {
			if c.Check(instance, cond) {
				trueCount++
			}
		}
		return trueCount == 1
	case AND:
		for _, cond := range conditions {
			if !c.Check(instance, cond) {
				return false
			}
		}
		return true
	case NOT:
		for _, cond := range conditions {
			if c.Check(instance, cond) {
				return false
			}
		}
		return true
	default:
		return false // throw Unrecognized operator
	}
}

// getValueByTemplate fetches the value specified by a template string or returns the direct value.
func (c *Conditions) getValueByTemplate(value any, instance any) any {
	valueStr, ok := value.(string)

	if !ok {
		return value
	}

	if strings.HasPrefix(valueStr, "~~") {
		return c.getTemplateString(valueStr[2:], instance)
	} else if strings.HasPrefix(valueStr, "{{") && strings.HasSuffix(valueStr, "}}") {
		valueStr = strings.TrimSpace(valueStr[2 : len(valueStr)-2])
		return c.getValueByChain(valueStr, instance)
	}

	return c.getValueByChain(valueStr, instance)
}

// getTemplateString processes a template string with placeholders, replacing them with actual values from the instance.
func (c *Conditions) getTemplateString(value string, instance any) string {
	re := regexp.MustCompile(`\{\{[-a-zA-Z0-9_]+\}\}`)
	matches := re.FindAllString(value, -1)
	for _, match := range matches {
		placeholder := match[2 : len(match)-2] // Trim off the {{ and }}
		replacement := c.getValueByChain(placeholder, instance)
		replacementStr, ok := replacement.(string)
		if !ok && replacement != nil {
			// Attempt to convert basic types to strings.
			switch reflect.TypeOf(replacement).Kind() {
			case reflect.Int, reflect.Int64, reflect.Float64:
				replacementStr = fmt.Sprintf("%v", replacement)
			default:
				panic("Bad type of hard string")
			}
		}
		value = strings.Replace(value, match, replacementStr, 1)
	}
	return value
}

// getValueByChain retrieves a value from an instance based on a "dot" path (e.g., "a.b.c").
func (c *Conditions) getValueByChain(param string, instance any) any {
	chain := strings.Split(param, ".")
	var result any

	for _, step := range chain {
		instanceValue := reflect.ValueOf(instance)

		if instanceValue.Kind() == reflect.Pointer {
			instanceValue = instanceValue.Elem()
		}

		switch instanceValue.Kind() {
		case reflect.Map:
			instanceValue := instanceValue.MapIndex(reflect.ValueOf(step))
			if !instanceValue.IsValid() {
				return nil // Key not found in map
			}
			instance = instanceValue.Interface()
		case reflect.Struct:
			instanceValue := instanceValue.FieldByName(step)
			if !instanceValue.IsValid() {
				return nil // Field not found in struct
			}
			instance = instanceValue.Interface()
		default:
			return nil // Not a map or struct
		}
		result = instance
	}
	return result
}
