package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/Gervva/avito_test_task/internal/model"
	"github.com/Gervva/avito_test_task/internal/model/cache"

	json "github.com/goccy/go-json"
)

var ErrMissCache = errors.New("miss cache")

type Repository struct {
	db *redis.Client
}

func NewRepository(db *redis.Client) *Repository {
	return &Repository{
		db: db,
	}
}

func (r Repository) GetUserBanner(ctx context.Context, req model.GetUserBannerReq) (*model.GetUserBannerResp, error) {
	key := fmt.Sprintf("%d:%d", *req.FeatureId, *req.TagId)

	var resp []byte
	var banner model.GetUserBannerResp

	err := r.db.Get(ctx, key).Scan(&resp)
	if err != nil {
		return nil, ErrMissCache
	}

	err = json.Unmarshal(resp, &banner)
	if err != nil {
		return nil, err
	}

	return &banner, nil
}

func (r Repository) AddUserBanner(ctx context.Context, req cache.AddBannerReq) error {
	key := fmt.Sprintf("%d,%d", req.FeatureId, req.TagId)
	value := cache.Banner{
		Content: req.Content,
		IsActive: req.IsActive,
	}

	jsonVal, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = r.db.Set(ctx, key, jsonVal, 5 * time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}
