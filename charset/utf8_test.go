package charset

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestToUTF8(t *testing.T) {
	var builder strings.Builder
	for _, r := range []rune{0x5e, 0x2d, 0x2d, 0x2d, 0x2d, 0x2d, 0x65e0, 0x6cd5, 0x8bc6, 0x522b, 0x7684, 0x5173, 0x952e, 0x5b57, 0x3a} {
		builder.WriteRune(r)
	}
	s := builder.String()
	t.Log(s)
	assert.Equal(t, "^-----无法识别的关键字:", ToUTF8(s))
}
