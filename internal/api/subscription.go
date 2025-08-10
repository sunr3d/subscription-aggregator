package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/sunr3d/subscription-aggregator/internal/httpx"
	infraErr "github.com/sunr3d/subscription-aggregator/internal/infra"
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
	mux.HandleFunc("PATCH /subscriptions/{id}", h.updateHandler)
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

	if err := httpx.WriteJSON(w, http.StatusCreated, map[string]int{"id": id}); err != nil {
		switch {
		case errors.Is(err, httpx.ErrJSONMarshal):
			h.logger.Error("не удалось сериализовать JSON", zap.Error(err))
			httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		case errors.Is(err, httpx.ErrWriteBody):
			h.logger.Warn("клиент закрыл соединение, ответ не был отправлен", zap.Error(err))
		}
		return
	}
}

func (h *Handler) getHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		httpx.HttpError(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	dataItem, err := h.db.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, infraErr.ErrNotFound) {
			httpx.HttpError(w, http.StatusNotFound, "Подписка не найдена")
			return
		}
		h.logger.Error("Ошибка GetByID()",
			zap.Int("id", id),
			zap.Error(err),
		)
		httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	resp := subscriptionRes{
		ID:          dataItem.ID,
		ServiceName: dataItem.ServiceName,
		Price:       dataItem.Price,
		UserID:      dataItem.UserID,
		StartDate:   dataItem.StartDate.UTC().Format("01-2006"),
	}
	if dataItem.EndDate != nil {
		resp.EndDate = dataItem.EndDate.UTC().Format("01-2006")
	}

	if err := httpx.WriteJSON(w, http.StatusOK, resp); err != nil {
		if errors.Is(err, httpx.ErrJSONMarshal) {
			h.logger.Error("не удалось сериализовать JSON", zap.Error(err))
			httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		} else if errors.Is(err, httpx.ErrWriteBody) {
			h.logger.Warn("клиент закрыл соединение, ответ не был отправлен", zap.Error(err))
		}
	}
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
