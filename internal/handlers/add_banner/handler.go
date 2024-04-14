package add_banner

import (
	"context"
	"encoding/json"
	"errors"

	"log/slog"
	"net/http"

	"github.com/Gervva/avito_test_task/internal/handlers"
	databaseRepo "github.com/Gervva/avito_test_task/internal/storage/database"
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

	defer func() { _ = r.Body.Close() }()

	var request HandlerRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseContent{
			Error: err,
		})
		return
	}

	response := h.handle(r.Context(), request)

	if response.Content.Error != nil {
		http.Error(w, response.Content.Error.Error(), response.Status)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response.Content)
	}
}

func (h Handler) handle(ctx context.Context, request HandlerRequest) HandlerResponse {
	reqBody := ToBannerFromHandler(request)
	if reqBody == nil {
		return HandlerResponse{
			Status: http.StatusBadRequest,
			Content: ResponseContent{
				Error: errors.New("invalid request body"),
			},
		}
	}

	id, err := h.bunnerService.AddBanner(ctx, reqBody)
	if err != nil {
		h.logger.ErrorContext(ctx, "error while adding banner", "error", err, "request", request)
		if errors.Is(err, databaseRepo.ErrRepoDB) {
			return HandlerResponse{
				Status: http.StatusInternalServerError,
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

	return HandlerResponse{
		Status: http.StatusOK,
		Content: ResponseContent{
			BannerId: *id,
		},
	}
}
