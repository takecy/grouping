package grouping

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
)

type (
	// Elementer is interface
	Elementer interface {
		GetName() string
		GetRatio() int
		SetRatio(int)
	}

	// Grouper is interface
	Grouper interface {
		GetGroup(srcID string, moreSeeds ...string) (Elementer, error)
	}

	// GroupDefinition is struct for group definition
	GroupDefinition struct {
		// Elems is the group to be distributed
		Elems []Elementer
		// DefaultElem is default element when a seed nothing matches elements.
		DefaultElem Elementer
	}

	// client is struct
	client struct {
		srcElem map[string]Elementer
		def     GroupDefinition
	}

	// HashFunc returns hash from seed
	HashFunc func(seed string) uint32
)

// hashFunc is function for generate hash
var hashFunc HashFunc

// New returns Grouper interface
func New(group GroupDefinition) (Grouper, error) {
	return NewWithHashFunc(group, simpleHash)
}

// NewWithHashFunc returns Grouper interface
func NewWithHashFunc(group GroupDefinition, hf HashFunc) (Grouper, error) {
	if len(group.Elems) == 0 {
		return nil, fmt.Errorf("elements required")
	}
	if hf == nil {
		return nil, fmt.Errorf("hash function required")
	}

	hashFunc = hf

	srcElems := make(map[string]Elementer)
	// convert to map
	for _, elem := range group.Elems {
		srcElems[elem.GetName()] = elem
	}

	// sort asc by ratio and name
	sort.Slice(group.Elems, func(i, j int) bool {
		return group.Elems[i].GetRatio() <= group.Elems[j].GetRatio() &&
			group.Elems[i].GetName() <= group.Elems[j].GetName()
	})

	// total ratio
	sumNum := 0
	for i := range group.Elems {
		sumNum += group.Elems[i].GetRatio()
	}
	if sumNum > 100 {
		return nil, fmt.Errorf("total ratio must be less than or equal 100")
	}
	if sumNum < 100 {
		if group.DefaultElem == nil {
			return nil, fmt.Errorf("when total ratio is less than 100, default elem is required")
		}
	}

	tmpNum := 0
	for i := range group.Elems {
		if i == 0 {
			tmpNum = group.Elems[i].GetRatio()
			continue
		}
		rt := group.Elems[i].GetRatio()
		rt += tmpNum
		tmpNum = rt
		group.Elems[i].SetRatio(rt)
	}

	cli := &client{
		srcElem: srcElems,
		def:     group,
	}

	return cli, nil
}

// GetGroup is implementation of `Grouper` interface
func (cli *client) GetGroup(srcID string, moreSeeds ...string) (elem Elementer, err error) {
	for i := range moreSeeds {
		srcID += "_" + moreSeeds[i]
	}

	numHash := hashFunc(srcID)
	strHash := fmt.Sprintf("%d", numHash)
	if len(strHash) < 2 {
		err = fmt.Errorf("too short hash: %s", strHash)
		return
	}

	// extract the last two digits
	cutNumStr := strHash[len(strHash)-2 : len(strHash)]
	// errors never occur
	seedRatio, _ := strconv.Atoi(cutNumStr)

	// retrieve from elem
	for _, el := range cli.def.Elems {
		if seedRatio < el.GetRatio() {
			elem = cli.srcElem[el.GetName()]
			break
		}
	}

	// not match from elems
	if elem == nil {
		elem = cli.def.DefaultElem
	}

	return
}

// simpleHash returns hash from string
func simpleHash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
