package delete_by_feature_tag

import (
	"context"
	"github.com/Gervva/avito_test_task/internal/model"
)

type BannerService interface {
	DeleteByFeatureTag(ctx context.Context, req *model.DeleteByFeatureTagReq) error
}