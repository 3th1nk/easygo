package charset

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding/simplifiedchinese"
	"testing"
)

func TestToUTF8(t *testing.T) {
	s := "abc123中文繁體"
	s2, err := simplifiedchinese.GB18030.NewEncoder().String(s)
	assert.NoError(t, err)
	s3, err := simplifiedchinese.GB18030.NewDecoder().String(s2)
	assert.NoError(t, err)

	assert.Equal(t, true, IsGB18030(s2))
	assert.Equal(t, s, ToUTF8(s2))
	assert.Equal(t, s, s3)

	s = "abc123\xff中文繁體"
	assert.Equal(t, false, IsUTF8(s))

	s = "电信"
	assert.Equal(t, true, IsUTF8(s))
	assert.Equal(t, true, IsGB18030(s))
	assert.Equal(t, s, ToUTF8(s))
}
