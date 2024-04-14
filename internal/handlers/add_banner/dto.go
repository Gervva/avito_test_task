package add_banner

type HandlerRequest struct {
	TagIds    *[]int64                `json:"tag_ids"`
	FeatureId *int64                  `json:"feature_id"`
	Content   *map[string]interface{} `json:"content"`
	IsActive  *bool                   `json:"is_active"`
}

type HandlerResponse struct {
	Status  int             `json:"status"`
	Content ResponseContent `json:"content"`
}

type ResponseContent struct {
	BannerId int64 `json:"banner_id,omitempty"`
	Error    error `json:"error,omitempty"`
}

// type HandlerResponseError struct {
// 	Error string `json:"error"`
// }
