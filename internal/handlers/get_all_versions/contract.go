package get_all_versions

import (
	"context"
	"github.com/Gervva/avito_test_task/internal/model"
)

type BannerService interface {
	GetAllVersions(ctx context.Context, banner *model.GetAllVersionsReq) (*[]model.GetAllVersionsResp, error)
}