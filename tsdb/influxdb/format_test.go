package influxdb

import (
	"fmt"
	"testing"
)

func TestEscapeTagValue(t *testing.T) {
	s := `"\n ,=n1'\`
	s2 := EscapeTagValue(s)
	s3 := EscapeCondValue(s)
	t.Log(s)
	t.Log(fmt.Sprintf(`insert measurement1,tag1=%v value=1`, s2))
	t.Log(fmt.Sprintf(`select * from measurement1 where tag1='%v'`, s3))
}

func TestEscapeFieldValue(t *testing.T) {
	s := `"\n ,=n1'\`
	s2 := EscapeFieldValue(s)
	s3 := EscapeCondValue(s, true)
	t.Log(s)
	t.Log(fmt.Sprintf(`insert measurement1,tag1=abc value="%v"`, s2))
	t.Log(fmt.Sprintf(`select * from measurement1 where tag1='abc' and value='%v'`, s3))
}
