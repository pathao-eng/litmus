# Litmus

<p align="center">
  <img alt="litmus"  src="/docs/images/banner.png">
</p>

Litmus is a coloring compound that under goes a color changes when added to an acid or base. If we have lots of solution either acid or base, we have no problem to find out which one is acid with litmus. Litmus paper is treated with natural dyes that reacts with acid/base to change color. 

Imagine a ***litmus*** paper that can identify numbers of solution. Instead of natural dyes, this ***litmus*** is treated with variety of indicators. Each combination of indicators reacts with a particular solution to make color.

Now, lets relate this imaginary ***litmus*** with our concerned **Litmus**.

First, we need to treat our **Litmus** with different indicators. This indicators are called here **Request** object.

Second, each solution (**Option**) is a mixture of multiple compounds. These compounds are called **Selector**. 

Third, instead of color, after **Litmus** test, solution will return a set of identifier. This set is called here **Resolver**.

That means, when we are testing each **Option** with **Request** object, if this object returns true for all **Selector** compounds, **Litmus** will return **Resolver** like color in litmus paper

Now main points are like:

An **Option** is formed with two parts.

1. First, **Selector** that holds some rules which must be true for this Option.
	
2. Second, **Resolver** that holds pre-defined knowledge about this Option. 

Lets say, we have multiple options to select from. According to different conditions, one option will be selected based on it's rules.

```
Option {
    Selector
    Resolver
}
```

To check with Litmus, we need another object called **Request**. This Request object holds all conditions which are used to filter Options.


Litmus process is done by some `Checker` methods those are used to check conditions.

## Checker

All we need to mention the `Checker` in the Tag of request fields. For example,

```
type Request struct {
    HexID   string  `checker:"ExistsInArray"`
}
```

Here, `ExistsInArray` Checker method is responsible to check `HexID` with Options' Selector.

> Note: All fields with 'checker' Tag requires a 'selector' Tag that
> will point to field(s) of selector object.

```
type Request struct {
    HexID   string  `selector:"HexIDs" checker:"ExistsInArray"`
}
```

We have now following `Checker` methods

### Equal

This method checks equality of request and selector value.
Still now, it only checks single value of

* Bool
* Int, Int8, Int16, Int32, Int64
* String

### ExistsInArray

This method checks existence of request value in array generated by ArrayMap meta converter.
We will learn about `meta converter` later.

This method works on following types of request:

* Int, Int8, Int16, Int32, Int64
* String

### TimeBetween

This method checks if provided request time is in between two selector time.
This also requires two fields in `selector` tag.

```
type Request struct {
    CreatedAt   time.Time   `selector:"StartTime,EndTime" checker:"TimeBetween"`
}
```

Above example implies that, each selector must have two fields named `StartTime` & `EndTime`.

### TimeAndHourBetween

This method checks if provided request time is in between two selector time and also matched with hour range.
This also requires three fields in `selector` tag.

```
type Request struct {
    CreatedAt   time.Time   `selector:"StartTime,EndTime,ActiveHours" checker:"TimeAndHourBetween"`
}
```

Above example implies that, each selector must have three fields named `StartTime`, `EndTime` & `ActiveHours`.

### ExistsInRange

This method checks existence of request value in BitMap generated by RangeBitMap meta converter
If selector has "*" value, any request value will be matched.

### EndsWith

This method checks if request value ends with a specific digit mentioned in selector.

## Meta Converter

Most of the Checker methods require some conversion in their data.

We have following Meta converter

### RangeBitMap

RangeBitMap must be applied to `String` data.

Example: `1-5`, `1-5,8,9-10`, `*`

Above String will be converted into BitMap and checked with Integer value.

```
type Selector struct {
    Days string  `meta:"RangeBitMap"`
}
```

```
type Request struct {
    Day int `selector:"Days" checker:"ExistsInRange"`
}
```

> Note: RangeBitMap is required for ExistsInRange checker.


### ArrayMap

ArrayMap converts list of items into array.

Followings are possible for now:

* []Int, []Int8, []Int16, []Int32, []Int64
* []String
* string ( comma-separated)

If selector has "*" value, any request value will be considered matched.

> Note: ArrayMap is required for ExistsInArray checker.

### DateMap

DateMap can be applied on String value of `time.Time`

> Note: DateMap is required for TimeBetween & TimeAndHourBetween


## Examples

```
type Selector struct {
    StartTime   string  `meta:"DateMap"`
    EndTime     string  `meta:"DateMap"`
    ActiveHours string  `meta:"RangeBitMap"`
}

type Request struct {
    CreatedAt   time.Time   `selector:"StartTime,EndTime,ActiveHours" checker:"TimeAndHourBetween"`
}
```

If we want to check if `CreatedAt` is in between `StartTime` & `EndTime` and also matches `ActiveHours`,
we need to use `TimeAndHourBetween` checker.

This requires `StartTime` & `EndTime` to have `DateMap` meta converter. And also it requires `ActiveHours` to have
`RangeBitMap` meta converter.

```json
{
  "StartTime":"2019-05-01 00:15 +0600",
  "EndTime":"2099-05-08 23:59 +0600",
  "ActiveHours":"0-19"
}
```


## How to use

```go
package main
import (
	"github.com/pathao-eng/litmus"
)

type Option struct {
	Key      string   `json:"key"`
	Selector Selector `json:"selector"`
	Resolver Resolver `json:"resolver"`
}

type Selector struct{}

type Resolver struct{}

type Request struct{}

func main() {
	var options []Option
	// read your options into this array

	var litmusOptions []litmus.Option
	for _, option := range options {
		litmusOptions = append(litmusOptions, litmus.Option{
			Key:      option.Key,
			Selector: option.Selector,
			Resolver: option.Resolver,
		})
	}

	litmusCtx, _ := litmus.NewContext(litmusOptions)

	req := &Request{}
	resolver, ok := litmusCtx.GetResolver(req, nil)
}
```
