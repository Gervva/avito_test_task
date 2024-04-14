package delete_banner

type HandlerRequest struct {
	Id int64 `json:"id"`
}

// type HandlerResponse struct {
// 	Status int                   `json:"status"`
// 	Msg    string                `json:"message,omitempty"`
// 	Error  *HandlerResponseError `json:"error,omitempty"`
// }

type HandlerResponse struct {
	Status  int             `json:"status"`
	Content ResponseContent `json:"content"`
}

type ResponseContent struct {
	Error    error `json:"error,omitempty"`
}

type HandlerResponseError struct {
	Message string `json:"message"`
}