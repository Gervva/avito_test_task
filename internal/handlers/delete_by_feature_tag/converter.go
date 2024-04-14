package delete_by_feature_tag

import (
	"github.com/Gervva/avito_test_task/internal/model"
)

func HandlerRequestToServiceDeleteByFeatureTagReq(b *HandlerRequest) *model.DeleteByFeatureTagReq {
	return &model.DeleteByFeatureTagReq{
		FeatureId: b.FeatureId,
		TagId:     b.TagId,
	}
}
