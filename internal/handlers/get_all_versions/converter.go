package get_all_versions

import (
	"encoding/json"

	"github.com/Gervva/avito_test_task/internal/model"
)

func HandlerRequestToServiceGetAllVersions(b *HandlerRequest) *model.GetAllVersionsReq {
	return &model.GetAllVersionsReq{
		GroupId: b.GroupId,
	}
}

func ModelToHandlerResp(banners *[]model.GetAllVersionsResp) *[]HandlerResponseContent {
	bs := make([]HandlerResponseContent, 0, len(*banners))
	for _, b := range *banners {
		var content map[string]interface{}
		json.Unmarshal(*b.Content, &content)

		bs = append(bs, HandlerResponseContent{
			GroupId: b.GroupId,
			FeatureId: b.FeatureId,
			TagIds: b.TagIds,
			IsActive: b.IsActive,
			Content: &content,
			Version: b.Version,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}

	return &bs
}
