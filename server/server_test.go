package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
)

type TestData struct {
	Endpoint string
	Body     string
}

const baseUrl = "http://localhost:8080"

func getGoodRequest(test TestData) (res *MetricResSchema, err error) {
	url := baseUrl + test.Endpoint

	response, err := http.Post(url, "application/json", strings.NewReader(test.Body))

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusAccepted {
		err = json.NewDecoder(response.Body).Decode(&res)

		if err != nil {
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("want status code %q but instead got %d", 201, response.StatusCode)
	}

	return res, nil
}

func getBadRequest(test TestData) (res *ErrorSchema, err error) {
	url := baseUrl + test.Endpoint

	response, err := http.Post(url, "application/json", strings.NewReader(test.Body))

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusBadRequest {
		err = json.NewDecoder(response.Body).Decode(&res)

		if err != nil {
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("want status code %q but instead got %d", 400, response.StatusCode)
	}

	return res, nil
}

func TestScreen(t *testing.T) {
	req := TestData{
		Endpoint: "/screen-resize",
		Body: `{
      "eventType": "screenResize",
      "websiteUrl": "http://localhost:8080",
      "sessionId": "123123-123123-123123123",
      "resizeFrom": {
          "width": "1920",
          "height": "1080"
        },
        "resizeTo": {
          "width": "1000",
          "height": "100"
        }
      }`,
	}

	_, err := getGoodRequest(req)

	if err != nil {
		t.Errorf("problem sending request to %s, %v", req.Endpoint, err)
		return
	}
}

func TestBadScreen(t *testing.T) {
	req := TestData{
		Endpoint: "/screen-resize",
		Body: `{
      "eventType": "screenResize",
      "websiteUrl": "http://localhost:8080",
      "sessionId": "123123-123123-123123123",
      "copyPaste": {
        "formId": "abc",
        "pasted": true
      }
    }`,
	}

	_, err := getBadRequest(req)

	if err != nil {
		t.Errorf("problem sending request to %s, %v", req.Endpoint, err)
		return
	}
}

func TestCopy(t *testing.T) {
	req := TestData{
		Endpoint: "/copy-paste",
		Body: `{
      "eventType": "copyAndPaste",
      "websiteUrl": "http://localhost:8080",
      "sessionId": "123123-123123-123123123",
      "copyPaste": {
        "formId": "abc",
        "pasted": true
      }
      }`,
	}

	_, err := getGoodRequest(req)

	if err != nil {
		t.Errorf("problem sending request to %s, %v", req.Endpoint, err)
		return
	}
}

func TestBadCopy(t *testing.T) {
	req := TestData{
		Endpoint: "/copy-paste",
		Body: `{
      "eventType": "copyAndPaste",
      "websiteUrl": "http://localhost:8080",
      "sessionId": "123123-123123-123123123",
      "copyPste": {
        "formId": "abc",
        "pasted": true
      }
      }`,
	}

	_, err := getBadRequest(req)

	if err != nil {
		t.Errorf("problem sending request to %s, %v", req.Endpoint, err)
		return
	}
}

func TestSubmit(t *testing.T) {
	req := TestData{
		Endpoint: "/submit",
		Body: `{
      "eventType": "submit",
      "websiteUrl": "http://localhost:8080",
      "sessionId": "123123-123123-123123123",
      "timeTaken": 2345
      }`,
	}

	_, err := getGoodRequest(req)

	if err != nil {
		t.Errorf("problem sending request to %s, %v", req.Endpoint, err)
		return
	}
}

func TestBadSubmit(t *testing.T) {
	req := TestData{
		Endpoint: "/submit",
		Body: `{
      "eventType": "copyAndPaste",
      "websiteUrl": "http://localhost:8080",
      "sessionId": "123123-123123-123123123",
      "timeTaken": 2345
      }`,
	}

	_, err := getBadRequest(req)

	if err != nil {
		t.Errorf("problem sending request to %s, %v", req.Endpoint, err)
		return
	}
}

func TestConcurrent(t *testing.T) {
	reqs := []TestData{
		{
			Endpoint: "/screen-resize",
			Body: `{
        "eventType": "screenResize",
        "websiteUrl": "http://localhost:8080",
        "sessionId": "123123-123123-123123123",
        "resizeFrom": {
            "width": "1920",
            "height": "1080"
          },
        "resizeTo": {
          "width": "1000",
          "height": "100"
        }
      }`,
		},
		{
			Endpoint: "/screen-resize",
			Body: `{
        "eventType": "screenResize",
        "websiteUrl": "http://localhost:8080",
        "sessionId": "123123-123123-123123123",
        "resizeFrom": {
            "width": "1920",
            "height": "1080"
          },
        "resizeTo": {
          "width": "1000",
          "height": "100"
        }
      }`,
		},
		{
			Endpoint: "/submit",
			Body: `{
        "eventType": "submit",
        "websiteUrl": "http://localhost:8080",
        "sessionId": "123123-123123-123123123",
        "timeTaken": 2345
        }`,
		},
		{
			Endpoint: "/copy-paste",
			Body: `{
        "eventType": "copyAndPaste",
        "websiteUrl": "http://localhost:8080",
        "sessionId": "123123-123123-123123123",
        "copyPaste": {
          "formId": "abc",
          "pasted": true
        }
        }`,
		},
	}

	wg := sync.WaitGroup{}

	for _, req := range reqs {
		wg.Add(1)

		go func(req TestData) {
			_, err := getGoodRequest(req)

			if err != nil {
				t.Errorf("problem sending request to %s, %v", req.Endpoint, err)
			}

			fmt.Printf("POST %v sent", req.Endpoint)
			wg.Done()
		}(req)

		wg.Wait()
	}
}

//List of tests:
/**
- Many same requests to the same endpoint will not introduce race condition
- Many different requests to the B, C endpoint will not cause race condition with
  many same requests going to endpoint A. (different requests async, same requests mutex)
- a lot of requests
**/
