package litmus

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/RoaringBitmap/roaring"
)

type Context struct {
	options          []Option
	externalCheckers []CheckerFunction
	externalMeta     MetaInterface
	internalMeta     internalMeta
	internalCheckers map[string]checkerFunc
}

func NewContext(opts []Option) (*Context, error) {
	ctx := &Context{
		options:          opts,
		externalCheckers: make([]CheckerFunction, 0),
		externalMeta:     nil,
		internalMeta: internalMeta{
			DateMap:     make(map[string]time.Time),
			RangeBitMap: make(map[string]*roaring.Bitmap),
			ArrayMap:    make([]map[string]map[interface{}]uint8, 0),
		},
		internalCheckers: nil,
	}

	if err := ctx.init(); err != nil {
		return nil, err
	}

	return ctx, nil
}

// AddChecker adds external checker function to check request with selectors
func (c *Context) AddChecker(checkers ...CheckerFunction) {
	c.externalCheckers = checkers
}

// AddMeta adds external meta information to use in external checker method
func (c *Context) AddMeta(meta MetaInterface) {
	c.externalMeta = meta
}

func (c *Context) init() error {
	if len(c.options) == 0 {
		return errors.New("must send at least one Option")
	}

	c.initInternalChecker()

	for _, option := range c.options {

		arrayMap := make(map[string]map[interface{}]uint8)

		selector := option.Selector
		valueOf := reflect.Indirect(reflect.ValueOf(selector))

		var typeOf reflect.Type

		switch k := reflect.TypeOf(selector); {
		case k.Kind() == reflect.Ptr:
			typeOf = k.Elem()
		default:
			typeOf = k
		}

		for i := 0; i < valueOf.NumField(); i++ {
			field := typeOf.Field(i)
			tagMeta, f := field.Tag.Lookup("meta")
			if !f {
				continue
			}

			switch tagMeta {
			case "DateMap":
				raw := valueOf.Field(i).String()
				v, err := ParseDate(raw)
				if err != nil {
					return fmt.Errorf("failed to parse Date from %v -> %v as %v", option.Key, field.Name, raw)
				}
				c.internalMeta.DateMap[raw] = v
			case "RangeBitMap":
				raw := valueOf.Field(i).String()
				v, err := ToRoarBitMap(raw)
				if err != nil {
					return fmt.Errorf("failed to convert to BitMap from %v -> %v as %v", option.Key, field.Name, raw)
				}
				c.internalMeta.RangeBitMap[raw] = v
			case "ArrayMap":
				raw := valueOf.Field(i).Interface()
				v, err := ToArrayMap(raw)
				if err != nil {
					return fmt.Errorf("failed to convert to ArrayMap from %v -> %v as %v", option.Key, field.Name, raw)
				}
				arrayMap[field.Name] = v
			default:
				return fmt.Errorf("invalid meta found: %v", tagMeta)
			}
		}

		c.internalMeta.ArrayMap = append(c.internalMeta.ArrayMap, arrayMap)
	}
	return nil
}

// GetResolver will first check externally provided Checker functions
// then will check internal Checker method to match with selector
func (c *Context) GetResolver(request RequestInterface, logger LoggerInterface) (ResolverInterface, bool) {

	for optionIndex, option := range c.options {
		optionSuccess := true
		for _, checker := range c.externalCheckers {
			if !checker(request, option.Selector, c.externalMeta, logger) {
				optionSuccess = false
				break
			}
		}

		valueOf := reflect.Indirect(reflect.ValueOf(request))
		typeOf := reflect.TypeOf(request)

		for i := 0; i < valueOf.NumField(); i++ {
			field := typeOf.Elem().Field(i)
			tagChecker, f := field.Tag.Lookup("checker")
			if !f {
				continue
			}

			value := valueOf.Field(i)
			if fn, f := c.internalCheckers[tagChecker]; f {
				if !fn(field.Tag.Get("selector"), value, optionIndex) {
					optionSuccess = false
					break
				}
			}
		}
		if optionSuccess {
			return option.Resolver, true
		}
	}
	return nil, false
}
