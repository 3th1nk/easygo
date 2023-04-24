package internal

import "testing"

func TestShortStr(t *testing.T) {
	for _, s := range []string{
		"aaa",
		"aaa aaa",
		"aaa aaa aaa",
		"测试一下",
	} {
		t.Log(ShortStr(s, 5))
	}

	str := "aaa aaa aaa"
	for i := -20; i < 20; i += 5 {
		t.Logf("%d: %s", i, ShortStr(str, i))
	}
}
