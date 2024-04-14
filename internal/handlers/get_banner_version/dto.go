package get_banner_version

type HandlerRequest struct {
	GroupId *int64 `json:"id"`
	Version *int64 `json:"version"`
}

type HandlerResponse struct {
	Status  int
	Content *HandlerResponseContent
	Error   *HandlerResponseError
}

type HandlerResponseContent struct {
	GroupId   *int64                  `json:"banner_id"`
	TagIds    *[]int64                `json:"tag_ids"`
	FeatureId *int64                  `json:"feature_id"`
	Content   *map[string]interface{} `json:"content"`
	IsActive  *bool                   `json:"is_active"`
	Version   *int64                  `json:"version"`
	CreatedAt *string                 `json:"created_at"`
	UpdatedAt *string                 `json:"updated_at"`
}

type HandlerResponseError struct {
	Error error `json:"error"`
}
