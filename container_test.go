package depcontainer_test

import (
	"fmt"
	"testing"

	"github.com/dripolles/depcontainer"
)

type A struct {
	X int
}

type B struct {
	Y int
}

type C struct {
	A *A
	B *B
}

func TestContainer_Basic(t *testing.T) {
	container := depcontainer.NewContainer()

	container.MustAdd(
		&depcontainer.ComponentDef{
			Name:        "a",
			IsSingleton: true,
			Builder: func(deps depcontainer.DepsMap) (interface{}, error) {
				return &A{X: 1}, nil
			},
		},
	).MustAdd(
		&depcontainer.ComponentDef{
			Name:        "b",
			IsSingleton: false,
			Builder: func(deps depcontainer.DepsMap) (interface{}, error) {
				return &B{Y: 2}, nil
			},
		},
	).MustAdd(
		&depcontainer.ComponentDef{
			Name:        "c",
			DependsOn:   []string{"a", "b"},
			IsSingleton: true,
			Builder: func(deps depcontainer.DepsMap) (interface{}, error) {
				return &C{
					A: deps["a"].(*A),
					B: deps["b"].(*B),
				}, nil
			},
		},
	)

	c := container.MustGet("c").(*C)
	fmt.Printf("C: %v\n", c)
}
