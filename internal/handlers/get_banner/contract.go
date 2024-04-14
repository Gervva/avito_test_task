package get_banner

import (
	"context"

	"github.com/Gervva/avito_test_task/internal/model"
)

type BannerService interface {
	GetBanner(ctx context.Context, banner *model.GetBannerReq) (*[]model.GetBannerResp, error)
}
