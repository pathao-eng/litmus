package litmus

import (
	"strconv"
	"strings"
	"time"

	"github.com/RoaringBitmap/roaring"
)

const (
	layoutISO = "2006-01-02 15:04 -0700"
)

func ParseDate(dateTimeStr string) (time.Time, error) {
	return time.Parse(layoutISO, dateTimeStr)
}

// ToArrayMap converts string/array into ArrayMap
func ToArrayMap(value interface{}) (map[interface{}]uint8, error) {
	// 	data := reflect.ValueOf(value)
	dataMap := make(map[interface{}]uint8)
	switch value.(type) {
	case string:
		data := value.(string)
		dataA := strings.Split(data, ",")
		for _, d := range dataA {
			d = strings.Trim(d, " ")
			dataMap[d] = 1
		}
	case []string:
		data := value.([]string)
		for _, d := range data {
			dataMap[d] = 1
		}
	case []int:
		data := value.([]int)
		for _, d := range data {
			dataMap[int64(d)] = 1
		}
	case []int8:
		data := value.([]int8)
		for _, d := range data {
			dataMap[int64(d)] = 1
		}
	case []int16:
		data := value.([]int16)
		for _, d := range data {
			dataMap[int64(d)] = 1
		}
	case []int32:
		data := value.([]int32)
		for _, d := range data {
			dataMap[int64(d)] = 1
		}
	case []int64:
		data := value.([]int64)
		for _, d := range data {
			dataMap[int64(d)] = 1
		}
	}

	return dataMap, nil
}

// ToRoarBitMap converts string into Bitmap
func ToRoarBitMap(matchStr string) (*roaring.Bitmap, error) {
	rangeBitmap := roaring.NewBitmap()
	if matchStr == "*" {
		return rangeBitmap, nil
	}
	splits := strings.Split(matchStr, ",")
	var err error
	for _, aSplit := range splits {
		ranges := strings.Split(aSplit, "-")
		if len(ranges) == 1 {
			hr, err := strconv.ParseUint(ranges[0], 10, 32)
			if err != nil {
				break
			}
			rangeBitmap.Add(uint32(hr))
		} else {
			end, err := strconv.Atoi(ranges[1])
			if err != nil {
				break
			}
			start, err := strconv.Atoi(ranges[0])
			if err != nil {
				break
			}
			if start > end {
				start, end = end, start
			}
			rangeBitmap.AddRange(uint64(start), uint64(end+1))
		}
	}
	return rangeBitmap, err
}
