package common

// Response for SUCCESS
type Response struct {
	Code   int         `json:"code"`
	Detail string      `json:"detail"`
	Data   interface{} `json:"data"`
}

// SUCCESS Response
func SUCCESS(d interface{}) *Response {
	return &Response{
		Code:   0,
		Detail: "Success",
		Data:   d,
	}
}
