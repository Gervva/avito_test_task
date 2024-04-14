package get_banner_version

import (
	"encoding/json"

	"github.com/Gervva/avito_test_task/internal/model"
)

func HandlerRequestToServiceGetBannerVersion(b *HandlerRequest) *model.GetBannerVersionReq {
	return &model.GetBannerVersionReq{
		GroupId: b.GroupId,
		Version: b.Version,
	}
}

func ModelToHandlerResp(b *model.GetBannerVersionResp) *HandlerResponseContent {
	var content map[string]interface{}
	json.Unmarshal(*b.Content, &content)

	banner := HandlerResponseContent{
		GroupId: b.GroupId,
		FeatureId: b.FeatureId,
		TagIds: b.TagIds,
		IsActive: b.IsActive,
		Content: &content,
		Version: b.Version,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}

	return &banner
}
