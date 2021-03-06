package litmus

import (
	"reflect"
	"strings"
	"time"
)

// func(selectorFieldName, requestFieldValue, OptionIndex) bool
type checkerFunc func(string, reflect.Value, int) bool

func (c *Context) initInternalChecker() {
	checker := make(map[string]checkerFunc)
	checker["Equal"] = c.equal
	checker["ExistsInArray"] = c.existsInArray
	checker["TimeBetween"] = c.timeBetween
	checker["TimeAndHourBetween"] = c.timeAndHourBetween
	checker["ExistsInRange"] = c.existsInRange
	checker["EndsWith"] = c.endsWith
	c.internalCheckers = checker
}

// equal method checks equality of request and selector value. (r == s)
// still now, this method only checks single value.
func (c *Context) equal(field string, requestValue reflect.Value, optionIndex int) bool {
	selectorI := c.options[optionIndex].Selector
	selectorValue := reflect.Indirect(reflect.ValueOf(selectorI)).FieldByName(field)

	kind := requestValue.Kind()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return requestValue.Int() == selectorValue.Int()
	case reflect.String:
		return requestValue.String() == selectorValue.String()
	case reflect.Bool:
		return requestValue.Bool() == selectorValue.Bool()
	}

	return false
}

// existsInArray checks existence of request value in array generated by ArrayMap meta converter
// if selector has "*" value, any request value will be matched.
func (c *Context) existsInArray(field string, requestValue reflect.Value, optionIndex int) bool {

	data, found := c.internalMeta.ArrayMap[optionIndex][field]
	if !found {
		return false
	}

	if _, found := data["*"]; found {
		return true
	}

	kind := requestValue.Kind()

	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if _, f := data[requestValue.Int()]; f {
			return true
		}
	case reflect.String:
		if _, f := data[requestValue.String()]; f {
			return true
		}
	}

	return false
}

// timeBetween method checks if provided request time is in between two selector time
// request 'field' must have a tag selector pointing to selector fields like 'StartTime,EndTime'
func (c *Context) timeBetween(field string, requestValue reflect.Value, optionIndex int) bool {
	selectorI := c.options[optionIndex].Selector

	timestampFields := strings.Split(field, ",")
	if len(timestampFields) != 2 {
		return false
	}

	startTime := reflect.Indirect(reflect.ValueOf(selectorI)).FieldByName(timestampFields[0])
	endTime := reflect.Indirect(reflect.ValueOf(selectorI)).FieldByName(timestampFields[1])

	var reqTime time.Time
	switch v := requestValue.Interface().(type) {
	case time.Time:
		reqTime = v
	case int64:
		reqTime = time.Unix(v, 0)
	}

	return c.checkTimeBetween(reqTime, startTime, endTime)
}

// timeAndHourBetween method checks if provided request time is in between two selector time
// and also matches with provided hours.
// request 'field' must have a tag selector pointing to selector fields like 'StartTime,EndTime'
func (c *Context) timeAndHourBetween(field string, requestValue reflect.Value, optionIndex int) bool {
	selectorI := c.options[optionIndex].Selector

	fields := strings.Split(field, ",")
	if len(fields) != 3 {
		return false
	}

	startTime := reflect.Indirect(reflect.ValueOf(selectorI)).FieldByName(fields[0])
	endTime := reflect.Indirect(reflect.ValueOf(selectorI)).FieldByName(fields[1])
	activeHours := reflect.Indirect(reflect.ValueOf(selectorI)).FieldByName(fields[2])

	var reqTime time.Time
	switch v := requestValue.Interface().(type) {
	case time.Time:
		reqTime = v
	case int64:
		reqTime = time.Unix(v, 0)
	}

	if !c.checkTimeBetween(reqTime, startTime, endTime) {
		return false
	}

	if activeHours.String() == "*" {
		return true
	}

	hour := reqTime.Hour()
	if v, found := c.internalMeta.RangeBitMap[activeHours.String()]; !found {
		return false
	} else {
		return v.ContainsInt(hour)
	}

}

// existsInRange checks existence of request value in BitMap generated by RangeBitMap meta converter
// if selector has "*" value, any request value will be matched.
func (c *Context) existsInRange(field string, requestValue reflect.Value, optionIndex int) bool {
	selectorI := c.options[optionIndex].Selector

	value := int(requestValue.Int())
	selectorValue := reflect.Indirect(reflect.ValueOf(selectorI)).FieldByName(field)

	if selectorValue.String() == "*" {
		return true
	}
	if v, found := c.internalMeta.RangeBitMap[selectorValue.String()]; !found {
		return false
	} else {
		return v.ContainsInt(value)
	}
}

// endsWith method checks if request value ends with a specific digit
// still now, only integer is consider for this Checker
func (c *Context) endsWith(field string, requestValue reflect.Value, optionIndex int) bool {
	value := int(requestValue.Int())
	lastDigit := value % 10
	return c.existsInRange(field, reflect.ValueOf(lastDigit), optionIndex)
}

// internal method //

// checkTimeBetween is internal method, not exposed as Checker
func (c *Context) checkTimeBetween(reqTime time.Time, startTime, endTime reflect.Value) bool {

	selectorStartTime := c.internalMeta.DateMap[startTime.String()]
	selectorEndTime := c.internalMeta.DateMap[endTime.String()]

	if reqTime.After(selectorEndTime) || reqTime.Before(selectorStartTime) {
		return false
	}
	return true
}
