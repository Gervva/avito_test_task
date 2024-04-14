package delete_by_feature_tag

import (
	"context"
	"encoding/json"

	"errors"
	"strconv"

	"github.com/Gervva/avito_test_task/internal/handlers"
	databaseRepo "github.com/Gervva/avito_test_task/internal/storage/database"

	"log/slog"
	"net/http"
)

type Handler struct {
	bunnerService BannerService
	logger        *slog.Logger
}

func New(s BannerService, l *slog.Logger) Handler {
	return Handler{bunnerService: s, logger: l}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request HandlerRequest

	queryParams := r.URL.Query()
	FeatureId := queryParams.Get("feature_id")

	if FeatureId != "" {
		featureId, err := strconv.ParseInt(FeatureId, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseContent{
				Error: err,
			})
			return
		}

		request.FeatureId = &featureId
	}

	TagId := queryParams.Get("tag_id")

	if TagId != "" {
		tagId, err := strconv.ParseInt(TagId, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ResponseContent{
				Error: err,
			})
			return
		}

		request.TagId = &tagId
	}

	response := h.handle(r.Context(), request)

	if response.Status != http.StatusNoContent {
		w.Header().Set("Content-Type", handlers.ContentTypeJSON)
		http.Error(w, response.Content.Error.Error(), response.Status)
	}
	
	w.WriteHeader(response.Status)
}

func (h Handler) handle(ctx context.Context, req HandlerRequest) HandlerResponse {
	err := h.bunnerService.DeleteByFeatureTag(ctx, HandlerRequestToServiceDeleteByFeatureTagReq(&req))
	if err != nil {
		h.logger.ErrorContext(ctx, "error while delete banner", "error", err, "request", req)
		if errors.Is(err, databaseRepo.ErrBannerNotExist) {
			return HandlerResponse{
				Status: http.StatusNotFound,
				Content: ResponseContent{
					Error: err,
				},
			}
		}

		return HandlerResponse{
			Status: http.StatusInternalServerError,
			Content: ResponseContent{
				Error: err,
			},
		}
	}

	return HandlerResponse{Status: http.StatusNoContent}
}
