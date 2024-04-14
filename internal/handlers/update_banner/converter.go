package update_banner

import (
	"encoding/json"

	"github.com/Gervva/avito_test_task/internal/model"
)

func HandlerRequestToServiceUpdateBanner(b *HandlerRequest) *model.UpdateBannerReq {
	content, err := json.Marshal(*b.Content)
	if err != nil {
		return nil
	}

	return &model.UpdateBannerReq{
		TagIds:    b.TagIds,
		FeatureId: b.FeatureId,
		GroupId:   b.BannerId,
		Content:   &content,
		IsActive:  b.IsActive,
		Version:   b.Version,
	}
}
