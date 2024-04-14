package update_banner

type HandlerRequest struct {
	TagIds    *[]int64                `json:"tag_ids"`
	FeatureId *int64                  `json:"feature_id"`
	BannerId  *int64                  `json:"banner_id"`
	Content   *map[string]interface{} `json:"content"`
	IsActive  *bool                   `json:"is_active"`
	Version   *int64                  `json:"version"`
}

type HandlerResponse struct {
	Status int                   `json:"status"`
	Error  *HandlerResponseError `json:"content,omitempty"`
}

type HandlerResponseError struct {
	Error error `json:"error"`
}
