package util

import "testing"

func TestShort(t *testing.T) {
	for _, s := range []string{
		"aaa",
		"aaa aaa",
		"aaa aaa aaa",
		"测试一下",
	} {
		Println(ShortStr(s, 5))
	}
}

func TestShortN(t *testing.T) {
	str := "aaa aaa aaa"
	for i := -20; i < 20; i += 5 {
		Println("%d: %s", i, ShortStr(str, i))
	}
}

func ExampleShortStr() {
	for _, s := range []string{
		"aaa",
		"aaa aaa",
		"aaa aaa aaa",
	} {
		Println(ShortStr(s, 5))
	}

	// output:
	// aaa
	// aaa a...(2 more)
	// aaa a...(6 more)
}
