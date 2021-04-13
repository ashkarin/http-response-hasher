package hasher

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
)

func ComputeHash(data []byte) []byte {
	h := md5.New()
	h.Write(data)
	return h.Sum(nil)
}

func HashToStr(hdata []byte) string {
	return hex.EncodeToString(hdata)
}

// HashHttpResponse returns a result of ComputeHash
// computed for the response body
func HashHttpResponse(res *http.Response) ([]byte, error) {
	if res == nil {
		return nil, fmt.Errorf("No http.Response given")
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return ComputeHash(data), nil
}

// HashHttpResponseToString returns a result of ComputeHash
// conputed for the response body and encoded as to a string
func HashHttpResponseToString(res *http.Response) (string, error) {
	hdata, err := HashHttpResponse(res)
	if err != nil {
		return "", err
	}
	return HashToStr(hdata), nil
}
