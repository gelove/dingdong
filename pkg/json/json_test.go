package json

import (
	"testing"
)

func TestMustTransform(t *testing.T) {
	data := map[string]any{
		"a": map[string]any{
			"a1": 1,
			"a2": "2",
			"a3": true,
		},
		"b": []any{
			"b1",
			2,
			false,
		},
	}
	out := make(map[string]any)
	MustTransform(data, &out)
	t.Logf("%#v", out)
}
