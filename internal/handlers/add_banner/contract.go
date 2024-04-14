package add_banner

import (
	"context"
	"github.com/Gervva/avito_test_task/internal/model"
)

type BannerService interface {
	AddBanner(ctx context.Context, banner *model.Banner) (*int64, error)
}