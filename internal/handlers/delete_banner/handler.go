package delete_banner

import (
	"context"
	"encoding/json"
	"strings"
	"strconv"
	"errors"

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
	path := strings.Trim(r.URL.Path, "/")
  	pathParts := strings.Split(path, "/")
	
	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseContent{
			Error: err,
		})
		return
	}

	response := h.handle(r.Context(), int64(id))

	if response.Status == http.StatusNotFound {
		w.WriteHeader(response.Status)
	} else if response.Status != http.StatusOK {
		w.Header().Set("Content-Type", handlers.ContentTypeJSON)
		http.Error(w, response.Content.Error.Error(), response.Status)
	}

	w.WriteHeader(response.Status)
}

func (h Handler) handle(ctx context.Context, id int64) HandlerResponse {
	err := h.bunnerService.DeleteBanner(ctx, &id)
	if err != nil {
		h.logger.ErrorContext(ctx, "error while delete banner", "error", err, "banner_id", id)
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

	return HandlerResponse{Status: http.StatusOK}
}
