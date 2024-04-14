package database

import (
	"github.com/Gervva/avito_test_task/internal/model"
)

func UpdateRow(newB model.UpdateBannerReq, oldB model.BannerWithPK) model.UpdateBannerReq {
	var res model.UpdateBannerReq

	if newB.FeatureId != nil {
		res.FeatureId = newB.FeatureId
	} else {
		res.FeatureId = oldB.FeatureId
	}
	
	if newB.TagIds != nil {
		res.TagIds = newB.TagIds
	} else {
		res.TagIds = oldB.TagIds
	}
	
	if newB.Content != nil {
		res.Content = newB.Content
	} else {
		res.Content = oldB.Content
	}
	
	if newB.IsActive != nil {
		res.IsActive = newB.IsActive
	} else {
		res.IsActive = oldB.IsActive
	}

	res.GroupId = oldB.GroupId
	
	return res
}