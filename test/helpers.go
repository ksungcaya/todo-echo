package test

import (
	"encoding/json"
	"net/http/httptest"
)

// GetResponseData is a helper function to get data from JSON response
func GetResponseData(response *httptest.ResponseRecorder) map[string]interface{} {
	return GetResponse(response, "data")
}

// GetResponseErrors is a helper function to get errors from JSON response
func GetResponseErrors(response *httptest.ResponseRecorder) map[string]interface{} {
	return GetResponse(response, "errors")
}

// GetResponse is a helper function to get a key from JSON response
func GetResponse(response *httptest.ResponseRecorder, key string) map[string]interface{} {
	var data map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &data)

	return data[key].(map[string]interface{})
}
