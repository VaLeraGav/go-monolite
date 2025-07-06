package testinit

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func SendRequest(t *testing.T, url, method, body string) *http.Response {
	t.Helper()

	var buf *bytes.Buffer
	if body != "" {
		buf = bytes.NewBufferString(body)
	} else {
		buf = &bytes.Buffer{}
	}

	req, err := http.NewRequest(method, url, buf)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return resp
}

func DecodeJSON(t *testing.T, reader io.Reader, out any) {
	t.Helper()
	body, err := io.ReadAll(reader)
	require.NoError(t, err)
	t.Logf("Response body: %s", string(body))
	err = json.Unmarshal(body, out)
	require.NoError(t, err)
}

func MarshalUnmarshal(t *testing.T, in interface{}, out any) {
	t.Helper()
	bytes, err := json.Marshal(in)
	require.NoError(t, err)
	err = json.Unmarshal(bytes, out)
	require.NoError(t, err)
}
