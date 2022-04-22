package textual

import (
	"testing"
)

func TestTrimSpace(t *testing.T) {
	list := []string{"  aaa ", " bbb  "}
	TrimSpace(list)
	t.Logf("%#v\n", list)
}

func TestSortStingNumber(t *testing.T) {
	list := []string{"234", "345", "123"}
	res := SortStingNumber(list, false)
	t.Logf("%#v\n", res)
}
