package conditions

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
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
				if !c.checkCommonOperator(operator, value, instance) {
					return false
				}
			} else if operator, exists := stringToLogicOperator[key]; exists {
				// Implement your checkCommonOperator logic here
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

	default:
		return false, fmt.Errorf("unhandled operator %s", operator)
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
