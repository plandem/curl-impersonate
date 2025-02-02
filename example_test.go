package curl_test

import (
	"encoding/json"
	"fmt"
	"github.com/plandem/curl-impersonate"
	"github.com/plandem/curl-impersonate/presets"
	"io"
	"log"
)

func Example_http() {
	c := curl.New()
	_, _, err := c.Request("http://httpbin.org/ip1")
	fmt.Println(err)

	r, _, _ := c.Request("http://httpbin.org/ip")
	fmt.Println(r.StatusCode)
	//Output:
	//HTTP Error. Not Found (404)
	//200
}

func Example() {
	// settings as part of creation
	c := curl.New(
		curl.Binary("curl-impersonate"),
		curl.Preset(presets.Random),
		curl.Flag("location", true),
		curl.Header("Accept", "application/json"),
		curl.Header("x-header1", "1"),
		curl.Header("X-HEADER2", "2"),
	)

	// settings via Set
	c.Set(
		curl.Header("Accept", "application/json"),
		curl.Header("Cache-Control", "no-cache, no-store, must-revalidate, max-age=0"),
		curl.Header("Pragma", "no-cache"),
	)

	// single header setter
	c.SetHeader("Accept", "application/json")

	// single flag setter
	c.SetFlag("location", true)

	resp, headers, err := c.Request("http://httpbin.org/ip")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// getting response's header
	etag := resp.Header.Get("ETag")
	fmt.Println(etag)

	// getting body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	type JSONResponse struct {
		Origin string `json:"origin"`
	}

	jsonData := JSONResponse{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(headers)
}
