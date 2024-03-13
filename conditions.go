package conditions

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func (c *Conditions) Check(instance interface{}, condition interface{}) bool {
	// Recursively check conditions if it's a slice, treating it as an AND condition
	if conditions, ok := condition.([]interface{}); ok {
		for _, cond := range conditions {
			if !c.Check(instance, cond) {
				return false
			}
		}
		return true
	}

	// Handle condition maps
	if condMap, ok := condition.(map[string]interface{}); ok {
		for key, value := range condMap {
			if operator, exists := stringToSimpleOperator[key]; exists {
				if !c.checkSimpleOperator(operator, value, instance) {
					return false
				}
			} else if operator, exists := stringToCommonOperator[key]; exists {
				result, err := c.checkCommonOperator(operator, value, instance)
				if err != nil {
					// Handle the error here
					return false
				}
				if !result {
					return false
				}
			} else if operator, exists := stringToLogicOperator[key]; exists {
				if !c.checkLogicOperator(operator, value, instance) {
					return false
				}
			} else {
				leftSide := c.getValueByTemplate(key, instance)
				rightSide := c.getValueByTemplate(value, instance)

				// Check for equality
				if !reflect.DeepEqual(leftSide, rightSide) {
					return false
				}
			}
		}
	}

	return false // Placeholder return
}

func (c *Conditions) checkSimpleOperator(operator SimpleOperatorsEnum, value interface{}, instance interface{}) bool {
	fact := c.getValueByTemplate(value, instance)

	switch operator {
	case NULL:
		return fact == nil
	case DEFINED:
		return !reflect.ValueOf(fact).IsNil()
	// case UNDEFINED:
	// 	return reflect.ValueOf(fact).IsNil()
	case EXIST:
		return fact != nil && !reflect.ValueOf(fact).IsNil()
	case EMPTY:
		v := reflect.ValueOf(fact)
		switch v.Kind() {
		case reflect.Array, reflect.Slice, reflect.String, reflect.Map:
			return v.Len() == 0
		default:
			// Go doesn't have a direct equivalent to JavaScript's broad object type;
			// structs could be checked for being "empty" in a different way, if needed.
			return false // Or throw an error, as per your application's needs
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
			return false // Or throw an error
		}
	case TRULY:
		return fact == true
	case FALSY:
		return fact == false
	default:
		return false
	}
}

// checkCommonOperator evaluates the instance against a common operator condition.
func (c *Conditions) checkCommonOperator(operator CommonOperatorsEnum, value interface{}, instance interface{}) (bool, error) {
	//var side interface{}

	fact := c.getValueByTemplate(value, instance)

	switch operator {
	case EQ:
		return reflect.DeepEqual(fact, c.getValueByTemplate(value, instance)), nil
	case NE:
		return !reflect.DeepEqual(fact, c.getValueByTemplate(value, instance)), nil
	case LT, GT, LTE, GTE:
		_, err := c.compareNumbersOrDates(fact, c.getValueByTemplate(value, instance), operator)
		if err != nil {
			return false, err
		}

		return true, nil
	case RE:
		pattern, ok := c.getValueByTemplate(value, instance).(string)
		if !ok {
			return false, nil // or log an error
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return false, nil // or log an error
		}
		str, ok := fact.(string)
		if !ok {
			return false, nil // or log an error
		}
		return re.MatchString(str), nil
	case IN:
		list, ok := c.getValueByTemplate(value, instance).([]interface{})
		if !ok {
			return false, nil // or log an error
		}
		return contains(list, fact.(string)), nil

	case NI:
		list, ok := c.getValueByTemplate(value, instance).([]interface{})
		if !ok {
			return false, nil // or log an error
		}
		return !contains(list, fact.(string)), nil
	case SW:
		prefix, ok := c.getValueByTemplate(value, instance).(string)
		if !ok {
			return false, nil // or log an error
		}
		str, ok := fact.(string)
		if !ok {
			return false, nil // or log an error
		}
		return strings.HasPrefix(str, prefix), nil

	case EW:
		suffix, ok := c.getValueByTemplate(value, instance).(string)
		if !ok {
			return false, nil // or log an error
		}
		str, ok := fact.(string)
		if !ok {
			return false, nil // or log an error
		}
		return strings.HasSuffix(str, suffix), nil
	case INCL:
		element := c.getValueByTemplate(value, instance)
		return isInCollection(fact, element), nil

	case EXCL:
		element := c.getValueByTemplate(value, instance)
		return !isInCollection(fact, element), nil
	case POWER:
		num, ok := fact.(int) // Assuming fact is an int for simplicity
		if !ok {
			return false, nil // or log an error
		}
		power, ok := c.getValueByTemplate(value, instance).(int)
		if !ok {
			return false, nil // or log an error
		}
		return (num & power) != 0, nil
	case BETWEEN:
		rangeSlice, ok := c.getValueByTemplate(value, instance).([]interface{})
		if !ok || len(rangeSlice) != 2 {
			return false, nil // Incorrect format
		}

		lowerBound, upperBound := rangeSlice[0], rangeSlice[1]

		compLower, err := attemptCompare(fact, lowerBound)
		if err != nil || compLower == -1 {
			return false, nil
		}

		compUpper, err := attemptCompare(fact, upperBound)
		if err != nil || compUpper == 1 {
			return false, nil
		}

		return true, nil
	case SOME:
		arrCond, ok := value.([]interface{})
		if ok {
			return false, fmt.Errorf("Bad fact type for $some operator")
		}

		for _, item := range arrCond {
			if c.Check(item, value) {
				return true, nil
			}
		}
		return false, nil
	case EVERY:
		arrCond, ok := value.([]interface{})
		if ok {
			return false, fmt.Errorf("Bad fact type for $some operator")
		}

		for _, item := range arrCond {
			if !c.Check(item, value) {
				return false, nil
			}
		}
		return true, nil
	case NOONE:
		arrCond, ok := value.([]interface{})
		if ok {
			return false, fmt.Errorf("Bad fact type for $some operator")
		}

		for _, item := range arrCond {
			if c.Check(item, value) {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, fmt.Errorf("unhandled operator %s", operator)
	}
}

// checkLogicOperator evaluates the logical operation on a set of conditions.
func (c *Conditions) checkLogicOperator(operator LogicOperatorsEnum, value interface{}, instance interface{}) bool {
	// Convert value to a slice of conditions
	var conditions []map[string]interface{}
	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			if cond, ok := item.(map[string]interface{}); ok {
				conditions = append(conditions, cond)
			}
		}
	case map[string]interface{}:
		conditions = append(conditions, v)
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
		return false // Unrecognized operator
	}
}

// Helper function to check if a slice contains a specific key
func contains[T comparable](s []T, e string) bool {
	for _, a := range s {
		if reflect.DeepEqual(a, e) {
			return true
		}
	}
	return false
}

func isInCollection(collection interface{}, element interface{}) bool {
	val := reflect.ValueOf(collection)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(val.Index(i).Interface(), element) {
				return true
			}
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			if reflect.DeepEqual(val.MapIndex(key).Interface(), element) {
				return true
			}
		}
	}
	return false
}

// attemptCompare tries to compare two values, accommodating for different types.
// Returns 0 if equal, -1 if v1 < v2, 1 if v1 > v2, and an error if incomparable.
func attemptCompare(v1, v2 interface{}) (int, error) {
	switch v1Typed := v1.(type) {
	case float64, float32, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		f1, _ := toFloat64(v1)
		f2, _ := toFloat64(v2)

		if f1 < f2 {
			return -1, nil
		} else if f1 > f2 {
			return 1, nil
		}
		return 0, nil
	case time.Time:
		t2, ok := v2.(time.Time)
		if !ok {
			return 0, fmt.Errorf("cannot compare time.Time with non-time.Time")
		}
		if v1Typed.Before(t2) {
			return -1, nil
		} else if v1Typed.After(t2) {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type for comparison")
	}
}

func toFloat64(v interface{}) (float64, bool) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(rv.Uint()), true
	case reflect.Float32, reflect.Float64:
		return rv.Float(), true
	default:
		return 0, false
	}
}

// Assuming `check` is defined on a `Conditions` struct as before.

// checkSome evaluates if some elements in a slice meet the provided condition.
func (c *Conditions) checkSome(instance []interface{}, condition interface{}) bool {
	for _, item := range instance {
		if c.Check(item, condition) {
			return true
		}
	}
	return false
}

// checkEvery evaluates if every element in a slice meets the provided condition.
func (c *Conditions) checkEvery(instance []interface{}, condition interface{}) bool {
	for _, item := range instance {
		if !c.Check(item, condition) {
			return false
		}
	}
	return true
}

// checkNoOne evaluates if no elements in a slice meet the provided condition.
func (c *Conditions) checkNoOne(instance []interface{}, condition interface{}) bool {
	return !c.checkSome(instance, condition)
}

// getValueByTemplate fetches the value specified by a template string or returns the direct value.
func (c *Conditions) getValueByTemplate(value interface{}, instance interface{}) interface{} {
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
	return value
}

// getTemplateString processes a template string with placeholders, replacing them with actual values from the instance.
func (c *Conditions) getTemplateString(value string, instance interface{}) string {
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
func (c *Conditions) getValueByChain(param string, instance interface{}) interface{} {
	chain := strings.Split(param, ".")
	var result interface{}

	for _, step := range chain {
		instanceValue := reflect.ValueOf(instance)
		if instanceValue.Kind() == reflect.Pointer {
			instanceValue = instanceValue.Elem()
		}

		switch instanceValue.Kind() {
		case reflect.Map:
			instance = instanceValue.MapIndex(reflect.ValueOf(step)).Interface()
		case reflect.Struct:
			instance = instanceValue.FieldByName(step).Interface()
		default:
			return nil // Not a map or struct
		}
		result = instance
	}
	return result
}

func (c *Conditions) compareNumbersOrDates(fact interface{}, side interface{}, operator CommonOperatorsEnum) (bool, error) {
	// Type assert and compare
	// This is a simplified example; you'll need to handle different types appropriately
	factValue, ok := fact.(float64) // Assuming you've normalized numbers to float64
	if !ok {
		return false, fmt.Errorf("fact is not a number")
	}
	sideValue, ok := side.(float64)
	if !ok {
		return false, fmt.Errorf("side is not a number")
	}

	switch operator {
	case LT:
		return factValue < sideValue, nil
	case GT:
		return factValue > sideValue, nil
	case LTE:
		return factValue <= sideValue, nil
	case GTE:
		return factValue >= sideValue, nil
	}

	return false, fmt.Errorf("invalid comparison operator")
}
