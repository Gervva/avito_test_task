package get_user_banner

import (
	"github.com/Gervva/avito_test_task/internal/model"
)

func HandlerRequestToServiceGetUserBanner(b *HandlerRequest) *model.GetUserBannerReq {
	return &model.GetUserBannerReq{
		TagId:           b.TagId,
		FeatureId:       b.FeatureId,
		UseLastRevision: b.UseLastRevision,
		IsAdmin:         b.IsAdmin,
	}
}
