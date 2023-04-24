package strUtil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSnakeToUpperKebab(t *testing.T) {
	for key, expect := range map[string]string{
		"sso_session_uid": "Sso-Session-Uid",
	} {
		if val := SnakeToUpperKebab(key); val != expect {
			t.Errorf("assert faild: expect %v, but %v", expect, val)
		}
	}
}

func TestCamelToLowerSnake(t *testing.T) {
	cases := []struct {
		str    string
		expect string
	}{
		{str: "ID", expect: "id"},
		{str: "intToStr", expect: "int_to_str"},
		{str: "IDAndName", expect: "id_and_name"},
		{str: "objIDAndName", expect: "obj_id_and_name"},
		{str: "HTTPPostBody", expect: "http_post_body"},
	}
	for _, c := range cases {
		assert.Equal(t, c.expect, CamelToSnake(c.str), c.str)
	}
}

func TestKebabToSnake(t *testing.T) {
	cases := []struct {
		str    string
		expect string
	}{
		{str: "ID", expect: "id"},
		{str: "int-To-Str", expect: "int_to_str"},
		{str: "ID-And-Name", expect: "id_and_name"},
		{str: "obj-ID-And-Name", expect: "obj_id_and_name"},
		{str: "HTTP-Post-Body", expect: "http_post_body"},
	}
	for _, c := range cases {
		assert.Equal(t, c.expect, KebabToSnake(c.str), c.str)
	}
}
