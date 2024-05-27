package influxdb

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/3th1nk/easygo/util/convertor"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	client http.Client
)

func init() {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.IdleConnTimeout = time.Minute
	t.MaxIdleConnsPerHost = 100
	t.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client.Transport = t
	client.Timeout = time.Minute
}

func doRequest(method string, url, token string, body interface{}, r interface{}) ([]byte, error) {
	if method == "" {
		method = http.MethodPost
	}

	var b io.Reader = nil
	if body != nil {
		switch t := body.(type) {
		case string:
			b = strings.NewReader(t)
		case []byte:
			b = bytes.NewReader(t)
		case io.Reader:
			b = t
		default:
			bodyStr, err := convertor.ToString(body)
			if err != nil {
				return nil, fmt.Errorf("invalid request body: %v", err.Error())
			}
			b = bytes.NewBufferString(bodyStr)
		}
	}

	req, err := http.NewRequest(method, url, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "text/plain; charset=UTF-8")
	if token != "" {
		req.Header.Set("Authorization", "Token "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("[%d]: %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if r != nil {
		return respBody, jsonUtil.Unmarshal(respBody, r)
	}

	return respBody, nil
}
