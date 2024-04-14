package cache

type AddBannerReq struct {
	FeatureId int64
	TagId     int64
	Content   []byte
	IsActive  bool
}

type Banner struct {
	Content  []byte `json:"content"`
	IsActive bool   `json:"is_active"`
}
