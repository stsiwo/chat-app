package test

import (
	"github.com/stretchr/testify/assert"
	"testing"
  "fmt"
)

func TestPointerV1(t *testing.T) {

  a := 1
  var pa *int
  pa = &a

  fmt.Println(&a)
  fmt.Println(&pa)
  fmt.Println(*pa)
  fmt.Println(*&a)

	assert.Equal(t, 0, 1)
}

type Sample struct {

  child *Child
}

type Child struct {

  value string
}

func (c *Child) setValue(value string) {
   c.value = value
}

func changeChildValue(child Child) {
  child.value = "kaoru"
}


func TestStructWithPointer(t *testing.T) {

  sample := Sample{child: &Child{value: "satoshi"}}

  changeChildValue(*sample.child)

  fmt.Println(sample.child.value)

  assert.True(t, false)
}

//func TestReferenceVsPointer(t *testing.T) {
//
//}
