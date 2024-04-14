package get_all_versions

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
	logger         *slog.Logger
}

func New(s BannerService, l *slog.Logger) Handler {
	return Handler{bunnerService: s, logger: l}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", handlers.ContentTypeJSON)
	var request HandlerRequest

	queryParams := r.URL.Query()
	GroupId := queryParams.Get("banner_id")

	if GroupId != "" {
		id, err := strconv.ParseInt(GroupId, 10, 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(HandlerResponseError{
				Error: errors.New("banner_id cannot be empty"),
			})
			return
		}

		request.GroupId = &id
	}

	response := h.handle(r.Context(), request)

	if response.Content == nil {
		http.Error(w, response.Error.Error.Error(), response.Status)
	} else {
		w.WriteHeader(response.Status)
		json.NewEncoder(w).Encode(response.Content)
	}
}

func (h Handler) handle(ctx context.Context, req HandlerRequest) HandlerResponse {
	banners, err := h.bunnerService.GetAllVersions(ctx, HandlerRequestToServiceGetAllVersions(&req))
	if err != nil {
		h.logger.ErrorContext(ctx, "error while get all bunner versions", "error", err, "request", req)
		if errors.Is(err, databaseRepo.ErrBannerNotExist) {
			return HandlerResponse{
				Status: http.StatusInternalServerError,
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

	return HandlerResponse{Status: http.StatusOK, Content: ModelToHandlerResp(banners)}
}
