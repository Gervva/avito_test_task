package get_user_banner

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
	TagId := queryParams.Get("tag_id")

	if TagId != "" {
		id, err := strconv.ParseInt(TagId, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
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
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(HandlerResponseError{
				Error: errors.New("incorrect feature_id"),
			})
			return
		}

		request.FeatureId = &id
	}
	UseLastRevision := queryParams.Get("use_last_revision")

	if UseLastRevision != "" {
		ulr, err := strconv.ParseBool(UseLastRevision)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(HandlerResponseError{
				Error: errors.New("incorrect use_last_revision"),
			})
			return
		}

		request.UseLastRevision = ulr
	} else {
		request.UseLastRevision = false
	}
	isAdmin := queryParams.Get("is_admin")

	if isAdmin != "" {
		is_admin, err := strconv.ParseBool(isAdmin)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(HandlerResponseError{
				Error: errors.New("incorrect use_last_revision"),
			})
			return
		}

		request.IsAdmin = &is_admin
	}

	response := h.handle(r.Context(), request)

	if response.Error == nil {
		w.Header().Set("Content-Type", handlers.ContentTypeJSON)
		json.NewEncoder(w).Encode(response.Content)
	} else if response.Status == http.StatusNotFound {
		w.WriteHeader(response.Status)
	} else {
		http.Error(w, response.Error.Error.Error(), response.Status)
	}
}

func (h Handler) handle(ctx context.Context, req HandlerRequest) HandlerResponse {
	content, err := h.bunnerService.GetUserBanner(ctx, HandlerRequestToServiceGetUserBanner(&req))
	if err != nil {
		h.logger.ErrorContext(ctx, "error while delete banner", "error", err, "request", req)
		if errors.Is(err, databaseRepo.ErrBannerNotExist) {
			return HandlerResponse{
				Status: http.StatusNotFound,
				Error: &HandlerResponseError{
					Error: err,
				},
			}
		}

		return HandlerResponse{
			Status: http.StatusInternalServerError,
			Error: &HandlerResponseError{
				Error: err,
			},
		}
	}

	var c map[string]interface{}
	err = json.Unmarshal(content.Content, &c)
	if err != nil {
		return HandlerResponse{
			Status: http.StatusInternalServerError,
			Error: &HandlerResponseError{
				Error: err,
			},
		}
	}

	return HandlerResponse{Status: http.StatusOK, Content: c}
}
