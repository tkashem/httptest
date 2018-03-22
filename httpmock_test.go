package httpmock_test

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/RueLaLa/chad-test/httpmock"
)

func do(t *testing.T, method string, url string, body string) *http.Response {
	request, error := http.NewRequest(method, url, bytes.NewBuffer([]byte(body)))
	if error != nil {
		t.Logf("Error: %s", error)
		t.Fail()
	}

	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		t.Logf("Error: %s", error)
		t.FailNow()
	}

	return response
}

func TestHTTPS(t *testing.T) {
	// Get an instance of the mocked http server
	server := httpmock.New(true)
	defer server.Close()

	call := server.Setup()
	defer call.Remove()

	call.Expect("GET", "/v1/token", http.Header{
		"foo": {"bar"}}).ReplyWith(http.Header{
		"foo": {"bar"},
	}, `{"id": 1}`)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	_, error := client.Get("https://127.0.0.1:8000")

	if error != nil {
		t.Logf("Error: %s", error)
		t.Fail()
	}
}

func TestHTTP(t *testing.T) {
	// Get an instance of the mocked http server
	server := httpmock.New(false)
	defer server.Close()

	call1 := server.Setup()
	defer call1.Remove()

	call1.Expect("POST", "/v1/fetch/order", http.Header{}).ReplyWith(http.Header{
		"foo": {"bar"},
	}, `{"id": 1}`)

	call2 := server.Setup()
	defer call2.Remove()

	call2.Expect("POST", "/v1/fetch/item", http.Header{}).ReplyWith(http.Header{
		"foo": {"bar"},
	}, `{"id": 2}`)

	response1 := do(t, "POST", "http://127.0.0.1:8000/v1/fetch/order", `{"id": 1}`)
	response2 := do(t, "POST", "http://127.0.0.1:8000/v1/fetch/item", `{"id": 2}`)

	defer response1.Body.Close()
	body1, _ := ioutil.ReadAll(response1.Body)

	defer response2.Body.Close()
	body2, _ := ioutil.ReadAll(response2.Body)

	if string(body1) != `{"id": 1}` {
		t.Logf("Error: body returned [%s] does not match expected", string(body1))
		t.Fail()
	}

	if string(body2) != `{"id": 2}` {
		t.Logf("Error: body returned [%s] does not match expected", string(body2))
		t.Fail()
	}
}

func TestHTTPPath(t *testing.T) {
	server := httpmock.New(false)
	defer server.Close()

	call := server.Setup()
	defer call.Remove()
	call.Expect("POST", "/v1/fetch/order", http.Header{}).ReplyWith(http.Header{}, `{"id": 1}`)

	response := do(t, "POST", "http://127.0.0.1:8000/v1/fetch/order?foo=1", `{"id": 1}`)
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	if string(body) != `{"id": 1}` {
		t.Logf("Error: body returned [%s] does not match expected", string(body))
		t.Fail()
	}
}
