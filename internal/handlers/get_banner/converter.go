package get_banner

import (
	"encoding/json"

	"github.com/Gervva/avito_test_task/internal/model"
)

func HandlerRequestToServiceGetBanner(b *HandlerRequest) *model.GetBannerReq {
	return &model.GetBannerReq{
		TagId:       b.TagId,
		FeatureId:   b.FeatureId,
		Limit:       b.Limit,
		Offset:      b.Offset,
	}
}

func ModelToHandlerResp(banners *[]model.GetBannerResp) *[]HandlerResponseContent {
	bs := make([]HandlerResponseContent, 0, len(*banners))
	for _, b := range *banners {
		var content map[string]interface{}
		json.Unmarshal(b.Content, &content)

		bs = append(bs, HandlerResponseContent{
			GroupId: b.GroupId,
			FeatureId: b.FeatureId,
			TagIds: b.TagIds,
			IsActive: b.IsActive,
			Content: content,
			Version: b.Version,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}

	return &bs
}
