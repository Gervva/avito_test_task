package get_banner

type HandlerRequest struct {
	TagId     *int64 `json:"tag_id"`
	FeatureId *int64 `json:"feature_id"`
	Limit     *int64 `json:"limit"`
	Offset    *int64 `json:"offset"`
}

type HandlerResponse struct {
	Status  int                       `json:"status"`
	Content *[]HandlerResponseContent `json:"content,omitempty"`
	Error   *HandlerResponseError     `json:"error,omitempty"`
}

type HandlerResponseContent struct {
	GroupId   int64                  `json:"banner_id"`
	TagIds    []int64                `json:"tag_ids"`
	FeatureId int64                  `json:"feature_id"`
	Content   map[string]interface{} `json:"content"`
	IsActive  bool                   `json:"is_active"`
	Version   int64                  `json:"version"`
	CreatedAt string                 `json:"created_at"`
	UpdatedAt string                 `json:"updated_at"`
}

type HandlerResponseError struct {
	Error error `json:"error"`
}
