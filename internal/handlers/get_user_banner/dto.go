package get_user_banner

type HandlerRequest struct {
	TagId           *int64 `json:"tag_id"`
	FeatureId       *int64 `json:"feature_id"`
	UseLastRevision bool   `json:"use_last_revision"`
	IsAdmin         *bool  `json:"is_admin"`
}

type HandlerResponse struct {
	Status  int                    `json:"status"`
	Content map[string]interface{} `json:"content,omitempty"`
	Error   *HandlerResponseError  `json:"error,omitempty"`
}

type HandlerResponseError struct {
	Error error `json:"error"`
}
