package litmus

import (
	"time"

	"github.com/RoaringBitmap/roaring"
)

type RequestInterface interface{}
type SelectorInterface interface{}
type ResolverInterface interface{}
type MetaInterface interface{}
type LoggerInterface interface{}

// Option should have a Selector and a Resolver.
// Selector will help find the best matched selector.
// Resolver will help implement the Selected option.
type Option struct {
	Key      string            `json:"key"`
	Selector SelectorInterface `json:"selector"`
	Resolver ResolverInterface `json:"resolver"`
}

type CheckerFunction func(RequestInterface, SelectorInterface, MetaInterface, LoggerInterface) bool

type internalMeta struct {
	DateMap     map[string]time.Time
	RangeBitMap map[string]*roaring.Bitmap
	ArrayMap    []map[string]map[interface{}]uint8
}
