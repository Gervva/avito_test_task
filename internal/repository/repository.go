package repository

import (
	"context"

	"github.com/Gervva/avito_test_task/internal/model"
	cacheModel "github.com/Gervva/avito_test_task/internal/model/cache"
)

type DataManager struct {
	db DBManager
	cache CacheManager
}

func New(db DBManager, cache CacheManager) *DataManager {
	return &DataManager{
		db: db,
		cache: cache,
	}
}

func (r DataManager) AddBanner(ctx context.Context, banner model.Banner) (*int64, error) {
	id, err := r.db.AddBanner(ctx, banner)

	return id, err
}

func (r DataManager) DeleteBanner(ctx context.Context, id int64) error {
	err := r.db.DeleteBanner(ctx, id)

	return err
}

func (r DataManager) GetBanner(ctx context.Context, req model.GetBannerReq) (*[]model.GetBannerResp, error) {
	banners, err := r.db.GetBanner(ctx, req)

	return banners, err
}

func (r DataManager) GetUserBanner(ctx context.Context, req model.GetUserBannerReq) (*model.GetUserBannerResp, error) {
	var banner *model.GetUserBannerResp
	var err error

	if req.UseLastRevision {
		banner, err = r.db.GetUserBanner(ctx, req)
	} else {
		banner, err = r.cache.GetUserBanner(ctx, req)

		if err != nil {
			banner, err = r.db.GetUserBanner(ctx, req)
			if err != nil {
				return banner, err
			}

			cacheBanner := cacheModel.AddBannerReq{
				FeatureId: *req.FeatureId,
				TagId: *req.TagId,
				Content: banner.Content,
				IsActive: banner.IsActive,
			}
			r.cache.AddUserBanner(ctx, cacheBanner)
		}
	}

	return banner, err
}

func (r DataManager) UpdateBanner(ctx context.Context, req model.UpdateBannerReq) error {
	err := r.db.UpdateBanner(ctx, req)

	return err
}

func (r DataManager) GetBannerVersion(ctx context.Context, req *model.GetBannerVersionReq) (*model.GetBannerVersionResp, error) {
	banner, err := r.db.GetBannerVersion(ctx, req)

	return banner, err
}

func (r DataManager) GetAllVersions(ctx context.Context, req *model.GetAllVersionsReq) (*[]model.GetAllVersionsResp, error) {
	banner, err := r.db.GetAllVersions(ctx, req)

	return banner, err
}

func (r DataManager) DeleteByFeatureTag(ctx context.Context, req *model.DeleteByFeatureTagReq) error {
	err := r.db.DeleteByFeatureTag(ctx, req)

	return err
}
