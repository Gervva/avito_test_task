package model

import "github.com/google/uuid"

type Banner struct {
	TagIds    *[]int64
	FeatureId *int64
	IsActive  *bool
	Content   *[]byte
}

type GetBannerReq struct {
	TagId     *int64
	FeatureId *int64
	Limit     *int64
	Offset    *int64
}

type GetBannerResp struct {
	GroupId   int64
	TagIds    []int64
	FeatureId int64
	Content   []byte
	Version   int64
	IsActive  bool
	CreatedAt string
	UpdatedAt string
}

type GetUserBannerReq struct {
	TagId           *int64
	FeatureId       *int64
	UseLastRevision bool
	IsAdmin         *bool
}

type GetUserBannerResp struct {
	Content  []byte
	IsActive bool
}

type UpdateBannerReq struct {
	TagIds    *[]int64
	FeatureId *int64
	GroupId   *int64
	Content   *[]byte
	IsActive  *bool
	Version   *int64
}

type UpdateBannerResp struct {
	Content []byte
}

type GetAllVersionsReq struct {
	GroupId *int64
}

type GetAllVersionsResp struct {
	GroupId   *int64
	TagIds    *[]int64
	FeatureId *int64
	Content   *[]byte
	IsActive  *bool
	Version   *int64
	CreatedAt *string
	UpdatedAt *string
}

type GetBannerVersionReq struct {
	GroupId *int64
	Version *int64
}

type GetBannerVersionResp struct {
	GroupId   *int64
	TagIds    *[]int64
	FeatureId *int64
	Content   *[]byte
	IsActive  *bool
	Version   *int64
	CreatedAt *string
	UpdatedAt *string
}

type DeleteByFeatureTagReq struct {
	FeatureId *int64
	TagId     *int64
}

type BannerWithPK struct {
	Id        uuid.UUID
	GroupId   *int64
	TagIds    *[]int64
	FeatureId *int64
	Content   *[]byte
	IsActive  *bool
	Version   *int64
	CreatedAt *string
	UpdatedAt *string
}
