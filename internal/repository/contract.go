package repository

import (
	"context"

	"github.com/Gervva/avito_test_task/internal/model"
	"github.com/Gervva/avito_test_task/internal/model/cache"
)

type CacheManager interface {
	GetUserBanner(ctx context.Context, req model.GetUserBannerReq) (*model.GetUserBannerResp, error)
	AddUserBanner(ctx context.Context, banner cache.AddBannerReq) error
}

type DBManager interface {
	AddBanner(ctx context.Context, banner model.Banner) (*int64, error)
	DeleteBanner(ctx context.Context, id int64) error
	GetBanner(ctx context.Context, req model.GetBannerReq) (*[]model.GetBannerResp, error)
	GetUserBanner(ctx context.Context, req model.GetUserBannerReq) (*model.GetUserBannerResp, error)
	UpdateBanner(ctx context.Context, req model.UpdateBannerReq) error
	GetBannerVersion(ctx context.Context, req *model.GetBannerVersionReq) (*model.GetBannerVersionResp, error)
	GetAllVersions(ctx context.Context, req *model.GetAllVersionsReq) (*[]model.GetAllVersionsResp, error)
	DeleteByFeatureTag(ctx context.Context, req *model.DeleteByFeatureTagReq) error
}
