package hasher_test

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"http-response-hasher/hasher"
	"net/http"
	"strings"
	"testing"
)

func mockResponseAndMD5Hash(body string) (*http.Response, []byte) {
	data := []byte(body)
	content := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
		"Content-Length: %d\r\n\r\n%s", len(data), data)

	br := bufio.NewReader(strings.NewReader(content))
	res, _ := http.ReadResponse(br, &http.Request{Method: "GET"})

	h := md5.New()
	h.Write(data)
	return res, h.Sum(nil)
}

func TestHashHttpResponse(t *testing.T) {
	data, err := hasher.HashHttpResponse(nil)
	if err == nil {
		t.Error("HashHttpResponse(nil) should return an error")
	}
	if data != nil {
		t.Error("HashHttpResponse(nil) should not return hasing result")
	}

	res, expHash := mockResponseAndMD5Hash("test response body")
	data, err = hasher.HashHttpResponse(res)
	if err != nil {
		t.Errorf("HashHttpResponse(\"test response body\") should not return an error: %s", err)
	}
	if data == nil {
		t.Error("HashHttpResponse(\"test response body\") should return hashing result")
	}
	if bytes.Compare(data, expHash) != 0 {
		t.Error("HashHttpResponse(\"test response body\") returned an unexpected hashing result")
	}
}

func TestHashHttpResponseToString(t *testing.T) {
	data, err := hasher.HashHttpResponseToString(nil)
	if err == nil {
		t.Error("HashHttpResponseToString(nil) should return an error")
	}
	if data != "" {
		t.Error("HashHttpResponseToString(nil) should not return hasing result")
	}

	res, expHash := mockResponseAndMD5Hash("test response body")
	expHashStr := hasher.HashToStr(expHash)
	data, err = hasher.HashHttpResponseToString(res)
	if err != nil {
		t.Errorf("HashHttpResponseToString(\"test response body\") should not return an error: %s", err)
	}
	if data == "" {
		t.Error("HashHttpResponseToString(\"test response body\") should return hashing result")
	}
	if data != expHashStr {
		t.Error("HashHttpResponseToString(\"test response body\") returned an unexpected hashing result")
	}
}

func TestHashIsMD5(t *testing.T) {
	data := []byte("some test data")
	h := md5.New()
	h.Write(data)

	expHash := h.Sum(nil)
	hash := hasher.ComputeHash(data)

	if bytes.Compare(expHash, hash) != 0 {
		t.Error("ComputeHash(...) returned unexpected hashing result")
	}
}

func TestHashToStr(t *testing.T) {
	text := "some test data"
	data := []byte(text)
	h := md5.New()
	h.Write(data)
	hash := h.Sum(nil)

	expHashStr := hex.EncodeToString(hash)
	hashStr := hasher.HashToStr(hash)

	if expHashStr != hashStr {
		t.Error("HashToStr(...) returned unexpected result")
	}
}
