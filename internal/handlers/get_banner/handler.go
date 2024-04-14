package get_banner

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/Gervva/avito_test_task/internal/handlers"

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
	w.Header().Set("Content-Type", handlers.ContentTypeJSON)
	var request HandlerRequest

	queryParams := r.URL.Query()

	TagId := queryParams.Get("tag_id")

	if TagId != "" {
		id, err := strconv.ParseInt(TagId, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(HandlerResponseError{
				Error: errors.New("incorrect tag_id"),
			})
			return
		}

		request.TagId = &id
	}

	FeatureId := queryParams.Get("feature_id")

	if FeatureId != "" {
		id, err := strconv.ParseInt(FeatureId, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(HandlerResponseError{
				Error: errors.New("incorrect feature_id"),
			})
			return
		}

		request.FeatureId = &id
	}

	Limit := queryParams.Get("limit")

	if Limit != "" {
		limit, err := strconv.ParseInt(Limit, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(HandlerResponseError{
				Error: errors.New("incorrect limit"),
			})
			return
		}

		request.Limit = &limit
	}
	Offset := queryParams.Get("offset")

	if Offset != "" {
		offset, err := strconv.ParseInt(Offset, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(HandlerResponseError{
				Error: errors.New("incorrect offset"),
			})
			return
		}

		request.Offset = &offset
	}

	response := h.handle(r.Context(), request)

	if response.Error != nil {
		http.Error(w, response.Error.Error.Error(), response.Status)
	} else {
		w.WriteHeader(response.Status)
		json.NewEncoder(w).Encode(response.Content)
	}
}

func (h Handler) handle(ctx context.Context, req HandlerRequest) HandlerResponse {
	banners, err := h.bunnerService.GetBanner(ctx, HandlerRequestToServiceGetBanner(&req))
	if err != nil {
		h.logger.ErrorContext(ctx, "error while delete banner", "error", err, "request", req)

		return HandlerResponse{
			Status: http.StatusInternalServerError,
			Error: &HandlerResponseError{
				Error: err,
			},
		}
	}

	return HandlerResponse{Status: http.StatusOK, Content: ModelToHandlerResp(banners)}
}
