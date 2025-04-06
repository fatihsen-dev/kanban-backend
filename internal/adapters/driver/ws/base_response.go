package ws

type BaseResponse struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}
