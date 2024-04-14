package banner

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Gervva/avito_test_task/internal/model"
	bannerDatabase "github.com/Gervva/avito_test_task/internal/storage/database"
)

type Service struct {
	bannerRepo Repository
}

type Response struct {
	Status int
}

func New(bannerRepo Repository) Service {
	return Service{bannerRepo: bannerRepo}
}

func (s Service) AddBanner(ctx context.Context, banner *model.Banner) (*int64, error) {
	if banner.FeatureId == nil || *banner.FeatureId <= 0 {
		return nil, fmt.Errorf("feature_id cannot be empty")
	}

	if banner.TagIds == nil || len(*banner.TagIds) == 0 {
		return nil, fmt.Errorf("number of tag_ids must be more than 0")
	}
	for _, id := range *banner.TagIds {
		if id <= 0 {
			return nil, fmt.Errorf("tag_id must be more than 0")
		}
	}

	// if banner.Content == nil || !json.Valid(*banner.Content) {
	if banner.Content == nil {
		return nil, fmt.Errorf("invalid json object in content")
	}

	if banner.FeatureId == nil || *banner.FeatureId <= 0 {
		return nil, fmt.Errorf("is_active cannot be empty")
	}

	id, err := s.bannerRepo.AddBanner(ctx, *banner)

	return id, err
}

func (s Service) DeleteBanner(ctx context.Context, id *int64) error {
	if id == nil || *id <= 0 {
		return fmt.Errorf("id must be more than 0")
	}

	err := s.bannerRepo.DeleteBanner(ctx, *id)

	return err
}

func (s Service) GetBanner(ctx context.Context, req *model.GetBannerReq) (*[]model.GetBannerResp, error) {
	if req.TagId != nil && *req.TagId <= 0 {
		return nil, fmt.Errorf("tag_id must be more than 0")
	}

	if req.FeatureId != nil && *req.FeatureId <= 0 {
		return nil, fmt.Errorf("feature_id must be more than 0")
	}

	if req.Limit != nil && *req.Limit <= 0 {
		return nil, fmt.Errorf("limit must be more than 0")
	}

	if req.Offset != nil && *req.Offset < 0 {
		return nil, fmt.Errorf("offset must be more than or equal 0")
	}

	resp, err := s.bannerRepo.GetBanner(ctx, *req)

	return resp, err
}

func (s Service) GetUserBanner(ctx context.Context, req *model.GetUserBannerReq) (*model.GetUserBannerResp, error) {
	if req.TagId != nil && *req.TagId <= 0 {
		return nil, fmt.Errorf("banner_id must be more than 0")
	}

	if req.FeatureId != nil && *req.FeatureId <= 0 {
		return nil, fmt.Errorf("tag_id must be more than 0")
	}

	resp, err := s.bannerRepo.GetUserBanner(ctx, *req)

	// если баннер неактивен и юзер - не админ, то возвращаем ошибку
	if resp != nil && !resp.IsActive && !*req.IsAdmin {
		return nil, bannerDatabase.ErrBannerNotExist
	}

	return resp, err
}

func (s Service) UpdateBanner(ctx context.Context, req *model.UpdateBannerReq) error {
	if req.GroupId == nil && *req.GroupId <= 0 {
		return fmt.Errorf("incorrect banner_id")
	}

	if req.TagIds != nil && len(*req.TagIds) == 0 {
		if len(*req.TagIds) == 0 {
			return fmt.Errorf("tag_ids can not be empty")
		}
		
		for _, id := range *req.TagIds {
			if id <= 0 {
				return fmt.Errorf("tag_id must be more than 0")
			}
		}
	}

	if req.FeatureId != nil && *req.FeatureId <= 0 {
		return fmt.Errorf("feature_id must be more than 0")
	}

	if req.Content != nil && !json.Valid(*req.Content) {
		return fmt.Errorf("invalid content")
	}

	if req.Version == nil || *req.Version <= 0 {
		return fmt.Errorf("version must be more than 0")
	}

	err := s.bannerRepo.UpdateBanner(ctx, *req)

	return err
}

func (s Service) GetBannerVersion(ctx context.Context, req *model.GetBannerVersionReq) (*model.GetBannerVersionResp, error) {
	if *req.GroupId <= 0 {
		return nil, fmt.Errorf("banner_id must be more than 0")
	}

	if *req.Version <= 0 {
		return nil, fmt.Errorf("version must be more than 0")
	}

	banner, err := s.bannerRepo.GetBannerVersion(ctx, req)

	return banner, err
}

func (s Service) GetAllVersions(ctx context.Context, req *model.GetAllVersionsReq) (*[]model.GetAllVersionsResp, error) {
	if *req.GroupId <= 0 {
		return nil, fmt.Errorf("banner_id must be more than 0")
	}

	banners, err := s.bannerRepo.GetAllVersions(ctx, req)

	return banners, err
}

func (s Service) DeleteByFeatureTag(ctx context.Context, req *model.DeleteByFeatureTagReq) error {
	if req.FeatureId == nil && req.TagId == nil {
		return fmt.Errorf("feature_id and tag_id can not be empty at the same time")
	}

	if req.FeatureId != nil && *req.FeatureId <= 0 {
		return fmt.Errorf("feature_id must be more than 0")
	}

	if req.TagId != nil && *req.TagId <= 0 {
		return fmt.Errorf("tag_id must be more than 0")
	}

	err := s.bannerRepo.DeleteByFeatureTag(ctx, req)

	return err
}
