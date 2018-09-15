# grouping
`grouping` is simple id grouping package in golang. Useful for AB testing.

![](https://img.shields.io/badge/golang-1.11.0-blue.svg?style=flat-square)

## Usage

### Basic usage

see [example](./example/simple/main.go)

```
package main

import (
	"fmt"

	"github.com/takecy/grouping"
)

// SimpleElem implements `grouping.Elementer`
type SimpleElem struct {
	name  string
	ratio int
}

func (e *SimpleElem) GetName() string { return e.name }
func (e *SimpleElem) GetRatio() int   { return e.ratio }
func (e *SimpleElem) SetRatio(r int)  { e.ratio = r }

func main() {
	group := grouping.GroupDefinition{
		// case: A+B=100
		Elems: []grouping.Elementer{
			&SimpleElem{name: "group-A", ratio: 20},
			&SimpleElem{name: "group-B", ratio: 80},
		},
	}

	g, err := grouping.New(group)
	if err != nil {
		panic(err)
	}

	//
	// The same result will be obtained no matter how many times it is executed.
	//
	testName := "welcome_content_test"

	userID1 := "user-001"
	elem1, err := g.GetGroup(userID1, testName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("elem1: %v\n", elem1.GetName()) // group-A

	userID2 := "user-002"
	elem2, err := g.GetGroup(userID2, testName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("elem2: %v\n", elem2.GetName()) // group-B
}
```

### Advanced usage

#### Default specification

see [example](./example/default_elem/main.go)

```
	group := grouping.GroupDefinition{
		// this means 20% of all
		Elems: []grouping.Elementer{
			&SimpleElem{name: "group-A", ratio: 20},
		},
		// not match
		DefaultElem: &SimpleElem{name: "group-default"},
	}

	g, err := grouping.New(group)
	if err != nil {
		panic(err)
	}
```


#### Customize hash function

You can customize function for generate hash.

see [example](./example/hash_func/main.go)

```
	group := grouping.GroupDefinition{
		// case: A+B=100
		Elems: []grouping.Elementer{
			&SimpleElem{name: "group-A", ratio: 30},
			&SimpleElem{name: "group-B", ratio: 70},
		},
	}
	// replace hash function
	hashFunc := func(seed string) uint32 {
		return uint32(len(seed))
	}

	g, err := grouping.NewWithHashFunc(group, hashFunc)
	if err != nil {
		panic(err)
	}
```


## LICENSE
[MIT](./LICENSE)