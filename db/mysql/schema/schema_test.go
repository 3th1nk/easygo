package schema

import (
	"testing"
)

func TestCustomTypePattern(t *testing.T) {
	for _, str := range []string{
		"{type:abc123a-b_c}",
		"{type: abc123a-b_c}",
		"{type :abc123a-b_c}",
		"{type : abc123a-b_c}",
		"{type:  abc123a-b_c}",
		"{type  :abc123a-b_c}",
		"{type  :  abc123a-b_c}",
		"{type:*abc123a-b_c}",
		"{type: *abc123a-b_c}",
		"{type :*abc123a-b_c}",
		"{type : *abc123a-b_c}",
		"{type:*ABC123a-b_c}",
		"{type: *ABC123a-b_c}",
		"{type :*ABC123a-b_c}",
		"{type : *ABC123a-b_c}",
	} {
		matches := customTypePattern.FindAllStringSubmatch(str, -1)
		if len(matches) == 0 {
			t.Errorf("no matches")
		} else {
			t.Logf("str=%v, type=%v", str, matches[0][1])
		}
	}
}
