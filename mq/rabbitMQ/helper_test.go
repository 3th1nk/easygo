package rabbitMQ

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	rawUrl     = "amqp://gsuser:Gs@mq?>>@192.168.1.213:20002"
	rawUrl2    = "amqps://gsuser@192.168.1.213:20002"
	rawUrl3    = "amqps://gsuser:@192.168.1.213:20002"
	encodedUrl = "amqp://gsuser:Gs%40mq%3F%3E%3E@192.168.1.213:20002"
)

func TestParseUrl(t *testing.T) {
	var scheme, username, password, host string
	var err error

	scheme, username, password, host, err = ParseUrl(rawUrl)
	assert.NoError(t, err, rawUrl)
	assert.Equal(t, "amqp", scheme)
	assert.Equal(t, "gsuser", username)
	assert.Equal(t, "Gs@mq?>>", password)
	assert.Equal(t, "192.168.1.213:20002", host)

	scheme, username, password, host, err = ParseUrl(rawUrl2)
	assert.NoError(t, err, rawUrl2)
	assert.Equal(t, "amqps", scheme)
	assert.Equal(t, "gsuser", username)
	assert.Equal(t, "", password)
	assert.Equal(t, "192.168.1.213:20002", host)

	scheme, username, password, host, err = ParseUrl(rawUrl3)
	assert.NoError(t, err, rawUrl3)
	assert.Equal(t, "amqps", scheme)
	assert.Equal(t, "gsuser", username)
	assert.Equal(t, "", password)
	assert.Equal(t, "192.168.1.213:20002", host)

	scheme, username, password, host, err = ParseUrl(encodedUrl)
	assert.NoError(t, err, encodedUrl)
	assert.Equal(t, "amqp", scheme)
	assert.Equal(t, "gsuser", username)
	assert.Equal(t, "Gs@mq?>>", password)
	assert.Equal(t, "192.168.1.213:20002", host)
}

func TestEncodeUrl(t *testing.T) {
	assert.Equal(t, encodedUrl, AutoEncodeUrl(rawUrl))
	assert.Equal(t, rawUrl2, AutoEncodeUrl(rawUrl2))
	assert.Equal(t, rawUrl2, AutoEncodeUrl(rawUrl3))
	assert.Equal(t, encodedUrl, AutoEncodeUrl(encodedUrl))
}
