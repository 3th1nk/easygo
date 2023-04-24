package mapUtil

import (
	"github.com/3th1nk/easygo/util/strUtil"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestKeysValues(t *testing.T) {
	for i := 0; i < 5; i++ {
		m := map[string]int{
			"b": 2,
			"a": 1,
			"c": 3,
		}

		t.Log("keys: " + strings.Join(StringKeys(m), ","))

		sortedKeys := SortedStringKeys(m)
		if len(sortedKeys) != 3 || sortedKeys[0] != "a" || sortedKeys[1] != "b" || sortedKeys[2] != "c" {
			t.Errorf("assert faild: %v", sortedKeys)
		}

		t.Log("values: " + strUtil.JoinInt(IntValues(m), ","))

		sortedValues := SortedIntValues(m)
		if len(sortedValues) != 3 || sortedValues[0] != 1 || sortedValues[1] != 2 || sortedValues[2] != 3 {
			t.Errorf("assert faild: %v", sortedValues)
		}

		time.Sleep(time.Duration(rand.Int63n(10)))
	}
}
