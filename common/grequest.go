package common

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	errorCode    = 500
	errorHeaders = http.Header{}
	errorBody    = func(a interface{}) string { return fmt.Sprintf(`{"msg":"%s"}`, a) }
)

// getHttpClient get Http client
func getHttpClient(verify bool) *http.Client {
	client := &http.Client{}
	if verify {
		tlsVerify := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		client = &http.Client{Transport: tlsVerify}
	}
	return client
}

// setHeaders for http request
func setHeaders(request *http.Request, header map[string]string) *http.Request {
	for key, value := range header {
		request.Header.Set(key, value)
	}
	return request
}

// response for http request
func response(request *http.Request, err error, header map[string]string, verify bool) (int, http.Header, string) {
	// 如果异常，返回
	if err != nil {
		return errorCode, errorHeaders, errorBody(err)
	}
	// set headers
	setHeaders(request, header)
	if request.Method == "POST" || request.Method == "PUT" || request.Method == "PATCH" {
		if header["Content-Type"] == "" {
			header["Content-Type"] = "application/x-www-form-urlencoded"
		}
	}
	client := getHttpClient(verify)
	resp, e := client.Do(request)
	if e != nil {
		return errorCode, errorHeaders, errorBody(e)
	}
	// close stream
	defer resp.Body.Close()
	data, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return errorCode, errorHeaders, errorBody(e)
	}
	return resp.StatusCode, resp.Header, string(data)
}

// Get method for http request
func Get(url string, header map[string]string, verify bool) (int, http.Header, string) {
	request, e := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	return response(request, e, header, verify)
}
