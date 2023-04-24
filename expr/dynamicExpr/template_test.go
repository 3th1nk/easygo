package dynamicExpr

import (
	"github.com/Masterminds/sprig"
	"github.com/stretchr/testify/assert"
	"html/template"
	"strings"
	"testing"
)

func TestTemplate(t *testing.T) {
	tmpl, err := template.New("").Funcs(sprig.FuncMap()).Parse(`{"name":"{{.obj.name|upper|trimSuffix "C"}}"}`)
	assert.NoError(t, err)
	strBuf := &strings.Builder{}
	err = tmpl.Execute(strBuf, map[string]interface{}{
		"int": 1,
		"str": "abc",
		"obj": map[string]interface{}{
			"name":  "abc",
			"value": 123,
		},
	})
	assert.NoError(t, err)
	str := strBuf.String()
	assert.Equal(t, `{"name":"AB"}`, str)
}
