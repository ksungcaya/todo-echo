package controllers

// ResponseData is a struct for response with data.
type ResponseData struct {
	Data interface{} `json:"data"`
}

// NewResponseData creates a response with "data" as parent node.
func NewResponseData(data interface{}) ResponseData {
	r := ResponseData{Data: data}
	return r
}
