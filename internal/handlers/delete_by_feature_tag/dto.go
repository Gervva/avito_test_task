package delete_by_feature_tag

type HandlerRequest struct {
	FeatureId *int64 `json:"feature_id"`
	TagId     *int64 `json:"tag_id"`
}

type HandlerResponse struct {
	Status  int             `json:"status"`
	Content ResponseContent `json:"content"`
}

type ResponseContent struct {
	Error error `json:"error,omitempty"`
}
