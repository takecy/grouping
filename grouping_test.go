package grouping

import (
	"testing"

	"github.com/matryer/is"
)

type elemTest struct {
	name  string
	ratio int
}

func (e *elemTest) GetName() string { return e.name }
func (e *elemTest) GetRatio() int   { return e.ratio }
func (e *elemTest) SetRatio(r int)  { e.ratio = r }

func TestNew(t *testing.T) {
	is := is.New(t)

	{
		// no elems
		group := GroupDefinition{}

		_, err := New(group)
		t.Logf("err: %+v", err)
		is.True(err != nil)
	}

	{
		// total ratio is over 100
		group := GroupDefinition{
			Elems: []Elementer{
				&elemTest{name: "01", ratio: 10},
				&elemTest{name: "02", ratio: 20},
				&elemTest{name: "03", ratio: 80},
			},
		}

		_, err := New(group)
		t.Logf("err: %+v", err)
		is.True(err != nil)
	}

	{
		// valid, total ratio is less than 100, but not specific default elem
		group := GroupDefinition{
			Elems: []Elementer{
				&elemTest{name: "01", ratio: 10},
				&elemTest{name: "02", ratio: 20},
				&elemTest{name: "03", ratio: 60},
			},
		}

		_, err := New(group)
		t.Logf("err: %+v", err)
		is.True(err != nil)
	}

	{
		// valid, total ratio is 100
		group := GroupDefinition{
			Elems: []Elementer{
				&elemTest{name: "01", ratio: 10},
				&elemTest{name: "03", ratio: 70},
				&elemTest{name: "02", ratio: 20},
			},
		}

		g, err := New(group)
		is.NoErr(err)
		is.True(g != nil)

		cli, ok := g.(*client)
		is.True(ok)

		for i, e := range cli.def.Elems {
			t.Logf("cli.def.Elems: i:e - %+v:%+v", i, e)
		}
		for k, v := range cli.srcElem {
			t.Logf("cli.srcElem: k:v - %+v:%+v", k, v)
		}

		is.Equal(len(cli.def.Elems), 3)
		is.Equal(len(cli.srcElem), 3)
		is.Equal(cli.def.Elems[0].GetName(), "01")
		is.Equal(cli.def.Elems[1].GetName(), "02")
		is.Equal(cli.def.Elems[2].GetName(), "03")
		is.Equal(cli.def.Elems[0].GetRatio(), 10)
		is.Equal(cli.def.Elems[1].GetRatio(), 30)
		is.Equal(cli.def.Elems[2].GetRatio(), 100)
	}

	{
		// valid, total ratio is less than 100, and specific default elem
		group := GroupDefinition{
			Elems: []Elementer{
				&elemTest{name: "01", ratio: 10},
				&elemTest{name: "03", ratio: 60},
				&elemTest{name: "02", ratio: 20},
			},
			DefaultElem: &elemTest{name: ""},
		}

		g, err := New(group)
		is.NoErr(err)
		is.True(g != nil)

		cli, ok := g.(*client)
		is.True(ok)

		for i, e := range cli.def.Elems {
			t.Logf("cli.def.Elems: i:e - %+v:%+v", i, e)
		}
		for k, v := range cli.srcElem {
			t.Logf("cli.srcElem: k:v - %+v:%+v", k, v)
		}

		is.Equal(len(cli.def.Elems), 3)
		is.Equal(len(cli.srcElem), 3)
		is.Equal(cli.def.Elems[0].GetName(), "01")
		is.Equal(cli.def.Elems[1].GetName(), "02")
		is.Equal(cli.def.Elems[2].GetName(), "03")
		is.Equal(cli.def.Elems[0].GetRatio(), 10)
		is.Equal(cli.def.Elems[1].GetRatio(), 30)
		is.Equal(cli.def.Elems[2].GetRatio(), 90)
	}

	{
		// valid, same ratio
		group := GroupDefinition{
			Elems: []Elementer{
				&elemTest{name: "02", ratio: 50},
				&elemTest{name: "01", ratio: 50},
			},
			DefaultElem: &elemTest{name: ""},
		}

		g, err := New(group)
		is.NoErr(err)
		is.True(g != nil)

		cli, ok := g.(*client)
		is.True(ok)

		for i, e := range cli.def.Elems {
			t.Logf("cli.def.Elems: i:e - %+v:%+v", i, e)
		}
		for k, v := range cli.srcElem {
			t.Logf("cli.srcElem: k:v - %+v:%+v", k, v)
		}

		is.Equal(len(cli.def.Elems), 2)
		is.Equal(len(cli.srcElem), 2)
		is.Equal(cli.def.Elems[0].GetName(), "01")
		is.Equal(cli.def.Elems[1].GetName(), "02")
		is.Equal(cli.def.Elems[0].GetRatio(), 50)
		is.Equal(cli.def.Elems[1].GetRatio(), 100)
	}
}

func TestNewWithHashFunc(t *testing.T) {
	is := is.New(t)

	{
		// hash function is nil
		group := GroupDefinition{
			Elems: []Elementer{
				&elemTest{name: "01", ratio: 10},
				&elemTest{name: "02", ratio: 20},
				&elemTest{name: "03", ratio: 80},
			},
		}

		_, err := NewWithHashFunc(group, nil)
		t.Logf("err: %+v", err)
		is.True(err != nil)
	}

	{
		// valid hash function
		group := GroupDefinition{
			Elems: []Elementer{
				&elemTest{name: "01", ratio: 10},
				&elemTest{name: "03", ratio: 70},
				&elemTest{name: "02", ratio: 20},
			},
		}

		g, err := NewWithHashFunc(group, func(seed string) uint32 { return uint32(10) })
		is.NoErr(err)
		is.True(g != nil)

		cli, ok := g.(*client)
		is.True(ok)

		for i, e := range cli.def.Elems {
			t.Logf("cli.def.Elems: i:e - %+v:%+v", i, e)
		}
		for k, v := range cli.srcElem {
			t.Logf("cli.srcElem: k:v - %+v:%+v", k, v)
		}

		is.Equal(len(cli.def.Elems), 3)
		is.Equal(len(cli.srcElem), 3)
		is.Equal(cli.def.Elems[0].GetName(), "01")
		is.Equal(cli.def.Elems[1].GetName(), "02")
		is.Equal(cli.def.Elems[2].GetName(), "03")
		is.Equal(cli.def.Elems[0].GetRatio(), 10)
		is.Equal(cli.def.Elems[1].GetRatio(), 30)
		is.Equal(cli.def.Elems[2].GetRatio(), 100)
	}
}

func TestGrouper_GetGroup(t *testing.T) {
	is := is.New(t)

	{
		// too short hash func
		group := GroupDefinition{
			Elems: []Elementer{
				&elemTest{name: "01", ratio: 10},
				&elemTest{name: "03", ratio: 70},
				&elemTest{name: "02", ratio: 20},
			},
		}
		// too short result
		hashFunc := func(seed string) uint32 { return uint32(1) }

		g, err := NewWithHashFunc(group, hashFunc)
		is.NoErr(err)
		is.True(g != nil)

		srcID := "id01"
		moreSeeds := "more01"
		_, err = g.GetGroup(srcID, moreSeeds)
		t.Logf("err: %+v", err)
		is.True(err != nil)
	}

	{
		// valid group
		group := GroupDefinition{
			Elems: []Elementer{
				&elemTest{name: "01", ratio: 10},
				&elemTest{name: "03", ratio: 70},
				&elemTest{name: "02", ratio: 20},
			},
		}
		// valid hash func
		hashFunc := func(seed string) uint32 { return uint32(10) }

		g, err := NewWithHashFunc(group, hashFunc)
		is.NoErr(err)
		is.True(g != nil)

		srcID := "id01"
		moreSeeds := "more01"
		elem, err := g.GetGroup(srcID, moreSeeds)
		is.NoErr(err)
		is.Equal(elem.GetName(), "02")
	}

	{
		// valid group, total ratio less than 100
		group := GroupDefinition{
			Elems: []Elementer{
				&elemTest{name: "01", ratio: 10},
				&elemTest{name: "03", ratio: 60},
				&elemTest{name: "02", ratio: 20},
			},
			DefaultElem: &elemTest{name: "default"},
		}
		// valid hash func
		// not match elems
		hashFunc := func(seed string) uint32 { return uint32(90) }

		g, err := NewWithHashFunc(group, hashFunc)
		is.NoErr(err)
		is.True(g != nil)

		srcID := "id01"
		moreSeeds := "more01"
		elem, err := g.GetGroup(srcID, moreSeeds)
		is.NoErr(err)
		is.Equal(elem.GetName(), "default")
	}
}

func Test_simpleHash(t *testing.T) {
	is := is.New(t)

	{
		// letter only
		seed := "abcde"
		h := simpleHash(seed)
		is.Equal(h, uint32(1956368136))
	}

	{
		// number only
		seed := "12345"
		h := simpleHash(seed)
		is.Equal(h, uint32(1136836824))
	}

	{
		// symbol
		seed := "++--//"
		h := simpleHash(seed)
		is.Equal(h, uint32(3244193459))
	}

	{
		// emoji
		seed := "🎧"
		h := simpleHash(seed)
		is.Equal(h, uint32(1524062823))
	}

	{
		// mix
		seed := "12345abcde++--//🎧"
		h := simpleHash(seed)
		is.Equal(h, uint32(876805499))
	}
}
