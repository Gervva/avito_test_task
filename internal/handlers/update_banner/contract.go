package update_banner

import (
	"context"

	"github.com/Gervva/avito_test_task/internal/model"
)

type BannerService interface {
	UpdateBanner(ctx context.Context, banner *model.UpdateBannerReq) error
}
