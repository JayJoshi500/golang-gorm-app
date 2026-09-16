package handlers

import (
	"log/slog"
	"net/http"

	"github.com/JayJoshi500/golang-gorm-app/internal/services"
	"github.com/JayJoshi500/golang-gorm-app/pkg/response"
)

type RestoHandler struct {
	service services.RestoService
	logger  *slog.Logger
}

func NewRestoHandler(service services.RestoService, logger *slog.Logger) *RestoHandler {
	return &RestoHandler{service: service, logger: logger}
}

func (h *RestoHandler) AvailableSlots(w http.ResponseWriter, r *http.Request) {
	availableSlots, err := h.service.GetAvailableSlots(r.Context())
	if err != nil {
		h.logger.Error("failed", "error", err)
		response.Error(w, http.StatusInternalServerError, "could not able to retrieve")
		return
	}

	response.Success(w, http.StatusCreated, availableSlots)
}

func (h *RestoHandler) TimingSlots(w http.ResponseWriter, r *http.Request) {
	availableSlots, err := h.service.GetTimingSlots(r.Context())
	if err != nil {
		h.logger.Error("failed", "error", err)
		response.Error(w, http.StatusInternalServerError, "could not able to retrieve")
		return
	}

	response.Success(w, http.StatusCreated, availableSlots)
}
