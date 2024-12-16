package hash

import (
	"testing"
)

func TestMakeSHA256(t *testing.T) {
	hash, _ := MakeSHA256([]byte(`Hello world!`), `key`)
	if hash != `852d2fec4bda6add8f12c5c1dff8420510ac5b85ef432140c7097aaee3c270ca` {
		t.Error("hash error")
	}
	t.Log(hash)
}

//852d2fec4bda6add8f12c5c1dff8420510ac5b85ef432140c7097aaee3c270ca
