package get_banner_version

import (
	"context"
	"github.com/Gervva/avito_test_task/internal/model"
)

type BannerService interface {
	GetBannerVersion(ctx context.Context, req *model.GetBannerVersionReq) (*model.GetBannerVersionResp, error)
}