package main

import (
	"testing"

	"github.com/sohaha/zlsgo/zstring"
)

func TestXxx(t *testing.T) {
	t.Log(zstring.String2Bytes("5734574.pdf"))
	t.Log(len(zstring.String2Bytes("5734574.pdf")))
	t.Log(len(zstring.Bytes2String([]byte{5, 5, 5, 5, 5})))
	t.Log((zstring.Bytes2String([]byte{5, 5, 5, 5, 5})))

	b, _ := zstring.Base64Decode(zstring.String2Bytes("S6B4iGEVQ1mY7TYBNxtHEw=="))
	t.Log(b)
	t.Log(len(b))

	r, err := zstring.AesDecrypt(b, "86fefr4ozq")
	t.Log(r, err)
	t.Log(len(r))
	t.Log(string(r))
}
