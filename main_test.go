package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
)

func TestEmptyUrl(t *testing.T) {
	url := os.Getenv("GOJUMPCLOUD_URL")
	if len(url) == 0 {
		t.Skip("Provide GOJUMPCLOUD_URL env variable")
	}

	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(body))
}

func postForm(t *testing.T, u string) string {
	resp, err := http.PostForm(u, url.Values{"password": {"angryMonkey"}})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	return string(body)
}

func postMultiForm(t *testing.T, u string, ch chan<- struct{}) {
	resp, err := http.PostForm(u, url.Values{"password": {"angryMonkey"}})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	n, err := io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Read %d bytes", n)

	ch <- struct{}{}
}

func TestSinglePostForm(t *testing.T) {
	jcURL := os.Getenv("GOJUMPCLOUD_URL")
	if len(jcURL) == 0 {
		t.Skip("Provide GOJUMPCLOUD_URL env variable")
	}

	body := postForm(t, jcURL)
	t.Log(body)
}

func TestMultiPostForm(t *testing.T) {
	// note, this is lame load test
	// it simply fires n requests and waits for each to finish

	jcURL := os.Getenv("GOJUMPCLOUD_URL")
	if len(jcURL) == 0 {
		t.Skip("Provide GOJUMPCLOUD_URL env variable")
	}

	sJcCount := os.Getenv("GOJUMPCLOUD_COUNT")
	if len(sJcCount) == 0 {
		t.Skip("Provide GOJUMPCLOUD_COUNT env variable")
	}
	jcCount, err := strconv.Atoi(sJcCount)
	if err != nil {
		t.Fatal("Invalid GOJUMPCLOUD_COUNT env value")
	}

	ch := make(chan struct{})
	for i := 0; i < jcCount; i++ {
		go postMultiForm(t, jcURL, ch)
		t.Logf("Posted %s [%d]\n", jcURL, i)
	}

	for i := 0; i < jcCount; i++ {
		<-ch
		t.Logf("Received %d\n", i)
	}
}

func TestHashEncode(t *testing.T) {
	// example encoding via openssl.exe

	// $ echo -n angryMonkey | openssl sha1 -sha512 -binary | openssl enc -base64
	// ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0v
	// KbZPZklJz0Fd7su2A+gf7Q==

	// $ echo -n bruce | openssl sha1 -sha512 -binary | openssl enc -base64
	// 9IZoQRJ9m4fDMTCU8RihwMpxJgyr+DgYcWg0on2mgSFICOmonH7YGKQNDkun92aK
	// sJRKSO2OA8ELDen1CVCaqg==

	// $ echo -n topher | openssl sha1 -sha512 -binary | openssl enc -base64
	// +MzpZ9kU8vZ06dDMjgBO44s99oVkdbkBstEqDq9zhFZ9E/T+Q3DxPf98lk6ni4K0
	//  xdSrTmQ19bHup/o4LIXd5w==

	var tt = []struct {
		in       string
		expected string
	}{
		{
			in:       "bruce",
			expected: "9IZoQRJ9m4fDMTCU8RihwMpxJgyr+DgYcWg0on2mgSFICOmonH7YGKQNDkun92aKsJRKSO2OA8ELDen1CVCaqg==",
		},
		{
			in:       "topher",
			expected: "+MzpZ9kU8vZ06dDMjgBO44s99oVkdbkBstEqDq9zhFZ9E/T+Q3DxPf98lk6ni4K0xdSrTmQ19bHup/o4LIXd5w==",
		},
		{
			in:       "angryMonkey",
			expected: "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==",
		},
	}

	for _, test := range tt {
		p := hashEncodePassword(test.in)

		if test.expected != p {
			t.Fatalf("Expected %s, received %s", test.expected, p)
		}
	}
}
