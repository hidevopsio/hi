package httpclient

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPClientDoSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en")
	response, err := client.Do(req)
	require.NoError(t, err, "should not have failed to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}

func TestHTTPClientGetSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Get(server.URL, headers)
	require.NoError(t, err, "should not have failed to make a GET request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))

	response, err = Get(server.URL, headers)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHTTPClientPostSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Post(server.URL, requestBody, headers)
	require.NoError(t, err, "should not have failed to make a POST request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))

	requestBody = bytes.NewReader([]byte(requestBodyString))
	response, err = Post(server.URL, requestBody, headers)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHTTPClientDeleteSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Delete(server.URL, headers, func(req *http.Request) {})
	require.NoError(t, err, "should not have failed to make a DELETE request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))

	response, err = Delete(server.URL, headers, func(req *http.Request) {})

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHTTPClientPutSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(100 * time.Millisecond))

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Put(server.URL, requestBody, headers, func(req *http.Request) {})
	require.NoError(t, err, "should not have failed to make a PUT request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))

	requestBody = bytes.NewReader([]byte(requestBodyString))
	response, err = Put(server.URL, requestBody, headers, func(req *http.Request) {})

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHTTPClientPatchSuccess(t *testing.T) {
	client := NewClient(WithHTTPTimeout(100 * time.Millisecond))

	requestBodyString := `{ "name": "heimdall" }`

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "en", r.Header.Get("Accept-Language"))

		rBody, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err, "should not have failed to extract request body")

		assert.Equal(t, requestBodyString, string(rBody))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	requestBody := bytes.NewReader([]byte(requestBodyString))

	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept-Language", "en")

	response, err := client.Patch(server.URL, requestBody, headers, func(req *http.Request) {})
	require.NoError(t, err, "should not have failed to make a PATCH request")

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))

	requestBody = bytes.NewReader([]byte(requestBodyString))
	response, err = Patch(server.URL, requestBody, headers, func(req *http.Request) {})

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "{ \"response\": \"ok\" }", respBody(t, response))
}

func TestHTTPClientGetRetriesOnFailure(t *testing.T) {
	count := 0
	noOfRetries := 3
	noOfCalls := noOfRetries + 1
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(100*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(NewRetrier(NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err, "should have failed to make GET request")

	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.Equal(t, "{ \"response\": \"something went wrong\" }", respBody(t, response))

	assert.Equal(t, noOfCalls, count)
}

func BenchmarkHTTPClientGetRetriesOnFailure(b *testing.B) {
	noOfRetries := 3
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(NewRetrier(NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	for i := 0; i < b.N; i++ {
		_, _ = client.Get(server.URL, http.Header{})
	}
}

func TestHTTPClientPostRetriesOnFailure(t *testing.T) {
	count := 0
	noOfRetries := 3
	noOfCalls := noOfRetries + 1
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(NewRetrier(NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Post(server.URL, strings.NewReader("a=1"), http.Header{}, func(req *http.Request) {})
	require.NoError(t, err, "should have failed to make GET request")

	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.Equal(t, "{ \"response\": \"something went wrong\" }", respBody(t, response))

	assert.Equal(t, noOfCalls, count)
}

func BenchmarkHTTPClientPostRetriesOnFailure(b *testing.B) {
	noOfRetries := 3
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(NewRetrier(NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	for i := 0; i < b.N; i++ {
		_, _ = client.Post(server.URL, strings.NewReader("a=1"), http.Header{})
	}
}

func TestHTTPClientGetReturnsNoErrorsIfRetriesFailWith5xx(t *testing.T) {
	count := 0
	noOfRetries := 2
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(noOfRetries),
		WithRetrier(NewRetrier(NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err)

	require.Equal(t, noOfRetries+1, count)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)
	require.Equal(t, "{ \"response\": \"something went wrong\" }", respBody(t, response))
}

func TestHTTPClientGetReturnsNoErrorsIfRetrySucceeds(t *testing.T) {
	count := 0
	countWhenCallSucceeds := 2
	backoffInterval := 1 * time.Millisecond
	maximumJitterInterval := 1 * time.Millisecond

	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithRetryCount(3),
		WithRetrier(NewRetrier(NewConstantBackoff(backoffInterval, maximumJitterInterval))),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		if count == countWhenCallSucceeds {
			w.Write([]byte(`{ "response": "success" }`))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{ "response": "something went wrong" }`))
		}
		count++
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{}, func(req *http.Request) {})
	require.NoError(t, err, "should not have failed to make GET request")

	require.Equal(t, countWhenCallSucceeds+1, count)
	require.Equal(t, http.StatusOK, response.StatusCode)
	require.Equal(t, "{ \"response\": \"success\" }", respBody(t, response))
}

func TestHTTPClientGetReturnsErrorOnClientCallFailure(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	server.URL = "" // Invalid URL to simulate client.Do error
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.Error(t, err, "should have failed to make GET request")

	require.Nil(t, response)

	assert.Equal(t, "Get \"\": unsupported protocol scheme \"\"", err.Error())
}

func TestHTTPClientGetReturnsNoErrorOn5xxFailure(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{ "response": "something went wrong" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	response, err := client.Get(server.URL, http.Header{})
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, response.StatusCode)



}

func TestHTTPClientGetReturnsErrorOnFailure(t *testing.T) {
	client := NewClient(WithHTTPTimeout(10 * time.Millisecond))

	response, err := client.Get("url_doenst_exist", http.Header{})
	require.EqualError(t, err, "Get \"url_doenst_exist\": unsupported protocol scheme \"\"")
	require.Nil(t, response)

}

type myHTTPClient struct {
	client http.Client
}

func (c *myHTTPClient) Do(request *http.Request) (*http.Response, error) {
	request.Header.Set("foo", "bar")
	return c.client.Do(request)
}

func TestCustomHTTPClientHeaderSuccess(t *testing.T) {
	client := NewClient(
		WithHTTPTimeout(10*time.Millisecond),
		WithHTTPClient(&myHTTPClient{
			client: http.Client{Timeout: 25 * time.Millisecond}}),
	)

	dummyHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "bar", r.Header.Get("foo"))
		assert.NotEqual(t, "baz", r.Header.Get("foo"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{ "response": "ok" }`))
	}

	server := httptest.NewServer(http.HandlerFunc(dummyHandler))
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	response, err := client.Do(req)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err)
	assert.Equal(t, "{ \"response\": \"ok\" }", string(body))
}

func respBody(t *testing.T, response *http.Response) string {
	if response.Body != nil {
		defer response.Body.Close()
	}

	respBody, err := ioutil.ReadAll(response.Body)
	require.NoError(t, err, "should not have failed to read response body")

	return string(respBody)
}

func TestHTTPClientOnFailure(t *testing.T) {

	c := NewClient()
	t.Run("should get error when the url is unsupported on Get", func(t *testing.T) {
		_, err := c.Get("::", nil)
		require.EqualError(t, err, "GET - request creation failed: parse \"::\": missing protocol scheme")
	})
	t.Run("should get error when the url is unsupported on Put", func(t *testing.T) {
		_, err := c.Put("::", nil, nil)
		require.EqualError(t, err, "PUT - request creation failed: parse \"::\": missing protocol scheme")

	})
	t.Run("should get error when the url is unsupported on Post", func(t *testing.T) {
		_, err := c.Post("::", nil, nil)
		require.EqualError(t, err, "POST - request creation failed: parse \"::\": missing protocol scheme")

	})
	t.Run("should get error when the url is unsupported on Delete", func(t *testing.T) {
		_, err := c.Delete("::", nil)
		require.EqualError(t, err, "DELETE - request creation failed: parse \"::\": missing protocol scheme")

	})
	t.Run("should get error when the url is unsupported on Patch", func(t *testing.T) {
		_, err := c.Patch("::", nil, nil)
		require.EqualError(t, err, "PATCH - request creation failed: parse \"::\": missing protocol scheme")

	})

}

func TestHTTPClientWhenHeaderIsNotNil(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hi"))
	}))

	server := httptest.NewServer(r)
	defer server.Close()
	c := NewClient()
	t.Run("should get error when the url is unsupported on Get", func(t *testing.T) {
		_, err := c.Get(server.URL, http.Header{})
		require.NoError(t, err, "should not have failed to make a GET request")
	})
	t.Run("should get error when the url is unsupported on Put", func(t *testing.T) {
		_, err := c.Put(server.URL, nil, http.Header{})
		require.NoError(t, err, "should not have failed to make a PUT request")
	})
	t.Run("should get error when the url is unsupported on Post", func(t *testing.T) {
		_, err := c.Post(server.URL, nil, http.Header{})
		require.NoError(t, err, "should not have failed to make a POST request")

	})
	t.Run("should get error when the url is unsupported on Delete", func(t *testing.T) {
		_, err := c.Delete(server.URL, http.Header{})
		require.NoError(t, err, "should not have failed to make a DELETE request")

	})
	t.Run("should get error when the url is unsupported on Patch", func(t *testing.T) {
		_, err := c.Patch(server.URL, nil, http.Header{})
		require.NoError(t, err, "should not have failed to make a PATCH request")

	})

}

func TestHTTPClientDoWhenReadAllFail(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("hi"))
	}))

	server := httptest.NewServer(r)
	defer server.Close()
	c := NewClient()

	mockReadCloser := mockReadCloser{}
	// if Read is called, it will return error
	mockReadCloser.On("Read", mock.AnythingOfType("[]uint8")).Return(0, fmt.Errorf("error reading"))
	// if Close is called, it will return error
	mockReadCloser.On("Close").Return(fmt.Errorf("error closing"))

	body := ioutil.NopCloser(mockReadCloser)

	t.Run("should get error when client do reading body with a ReadAll error", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, server.URL, body)

		_, err := c.Do(request)
		println(err.Error())
		require.EqualError(t, err, "error reading")

	})
}

type mockReadCloser struct {
	mock.Mock
}

func (m mockReadCloser) Read(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *mockReadCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}
