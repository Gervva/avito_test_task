package get_user_banner

import (
	"context"

	"github.com/Gervva/avito_test_task/internal/model"
)

type BannerService interface {
	GetUserBanner(ctx context.Context, banner *model.GetUserBannerReq) (*model.GetUserBannerResp, error)
}
