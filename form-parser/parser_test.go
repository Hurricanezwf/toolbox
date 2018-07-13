package formparser

import (
	"fmt"
	"reflect"
	"testing"
)

type Hello struct {
	A int                `a:"a"`
	B string             `a:"b"`
	C int64              `a:"c"`
	D float64            `a:"d"`
	E []int              `a:"e"`
	F Persion            `a:"f"`
	G bool               `a:"g"`
	H []*Persion         `a:"h"`
	I map[string]*string `a:"i"`
}

type Persion struct {
	CPU *string `a:"cpu"`
}

func BenchmarkParse(b *testing.B) {
	h := Hello{
		A: 1,
		B: "BB",
		C: 2,
		D: 3.14,
		E: []int{2, 0, 32},
		F: Persion{CPU: StringPtr("1核")},
		G: true,
		H: []*Persion{
			&Persion{CPU: StringPtr("2核")},
			&Persion{CPU: StringPtr("3核")},
			&Persion{CPU: StringPtr("4核")},
		},
		I: map[string]*string{
			"m1": StringPtr("m1"),
			"m2": StringPtr("m2"),
		},
	}

	p := NewFormParser("a", "-")
	for i := 0; i < b.N; i++ {
		_, err := p.parse(reflect.ValueOf(h))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}
	p.Debug(reflect.ValueOf(h))
}
