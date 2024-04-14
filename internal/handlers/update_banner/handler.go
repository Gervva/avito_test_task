package update_banner

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Gervva/avito_test_task/internal/handlers"
	databaseRepo "github.com/Gervva/avito_test_task/internal/storage/database"

	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	bunnerService BannerService
	logger        *slog.Logger
}

func New(s BannerService, l *slog.Logger) Handler {
	return Handler{bunnerService: s, logger: l}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	pathParts := strings.Split(path, "/")

	defer func() { _ = r.Body.Close() }()

	Bannerid, err := strconv.ParseInt(pathParts[1], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HandlerResponseError{
			Error: errors.New("incorrect banner_id"),
		})
		return
	}

	var request HandlerRequest

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(HandlerResponseError{
			Error: err,
		})
		return
	}
	request.BannerId = &Bannerid

	response := h.handle(r.Context(), request)

	if response.Error != nil {
		w.Header().Set("Content-Type", handlers.ContentTypeJSON)
		http.Error(w, response.Error.Error.Error(), response.Status)
	}

	w.WriteHeader(response.Status)
}

func (h Handler) handle(ctx context.Context, req HandlerRequest) HandlerResponse {
	reqBody := HandlerRequestToServiceUpdateBanner(&req)
	if reqBody == nil {
		return HandlerResponse{
			Status: http.StatusBadRequest,
			Error: &HandlerResponseError{
				Error: errors.New("invalid request body"),
			},
		}
	}

	err := h.bunnerService.UpdateBanner(ctx, reqBody)
	if err != nil {
		h.logger.ErrorContext(ctx, "error while delete banner", "error", err, "request", req)
		if errors.Is(err, databaseRepo.ErrBannerNotExist) {
			return HandlerResponse{
				Status: http.StatusNotFound,
			}
		}

		return HandlerResponse{
			Status: http.StatusInternalServerError,
			Error: &HandlerResponseError{
				Error: err,
			},
		}
	}

	return HandlerResponse{Status: http.StatusOK}
}
