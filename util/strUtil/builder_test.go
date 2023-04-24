package strUtil

import "testing"

func TestStringBuilder_Append(t *testing.T) {
	sb := NewBuilder()
	for i := 0; i < 1000; i++ {
		if i == 0 {
			sb.Append(Rand(6))
		} else {
			sb.Append(" " + Rand(6))
		}
	}
	t.Log(sb.String())
}
