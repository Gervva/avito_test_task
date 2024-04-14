package add_banner

import (
	"encoding/json"

	"github.com/Gervva/avito_test_task/internal/model"
)

func ToBannerFromHandler(h HandlerRequest) *model.Banner {
	content, err := json.Marshal(*h.Content)
	if err != nil {
		return nil
	}

	return &model.Banner{
		TagIds:    h.TagIds,
		FeatureId: h.FeatureId,
		IsActive:  h.IsActive,
		Content:   &content,
	}
}
