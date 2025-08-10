package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/sunr3d/subscription-aggregator/internal/httpx"
	"github.com/sunr3d/subscription-aggregator/internal/interfaces/infra"
	"github.com/sunr3d/subscription-aggregator/models"
)

type Handler struct {
	db     infra.Database
	logger *zap.Logger
}

func New(db infra.Database, logger *zap.Logger) *Handler {
	return &Handler{
		db:     db,
		logger: logger,
	}
}

func (h *Handler) RegisterCRUDL(mux *http.ServeMux) {
	mux.HandleFunc("POST /subscriptions", h.createHandler)
	mux.HandleFunc("GET /subscriptions/{id}", h.getHandler)
	mux.HandleFunc("PUT /subscriptions/{id}", h.updateHandler)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.deleteHandler)
	mux.HandleFunc("GET /subscriptions", h.listHandler)
}

func (h *Handler) createHandler(w http.ResponseWriter, r *http.Request) {
	var req createSubscriptionReq

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		httpx.HttpError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}

	if err := validateCreateSubscription(req); err != nil {
		httpx.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	start, _ := time.Parse("01-2006", req.StartDate)
	start = start.UTC()

	var endPtr *time.Time
	if strings.TrimSpace(req.EndDate) != "" {
		end, _ := time.Parse("01-2006", req.EndDate)
		tt := end.UTC()
		endPtr = &tt
	}

	id, err := h.db.Create(r.Context(), models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   start,
		EndDate:     endPtr,
	})
	if err != nil {
		h.logger.Error("ошибка при создании подписки", zap.Error(err))
		httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	resp := subscriptionRes{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   start.Format("01-2006"),
	}
	if endPtr != nil {
		resp.EndDate = endPtr.Format("01-2006")
	}

	if err := httpx.WriteJSON(w, http.StatusCreated, resp); err != nil {
		h.logger.Error("не удалось записать ответ", zap.Error(err))
	}
}

func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Ручка получения подписки
}

func (h *Handler) updateHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Ручка обновления подписки
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Ручка удаления подписки
}

func (h *Handler) listHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Ручка получения списка подписок
}
