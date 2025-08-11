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
	"github.com/sunr3d/subscription-aggregator/internal/interfaces/services"
	"github.com/sunr3d/subscription-aggregator/models"
)

type Handler struct {
	svc    services.SubscriptionService
	logger *zap.Logger
}

func New(svc services.SubscriptionService, logger *zap.Logger) *Handler {
	return &Handler{
		svc:    svc,
		logger: logger,
	}
}

func (h *Handler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("POST /subscriptions", h.createHandler)
	mux.HandleFunc("GET /subscriptions/{id}", h.getHandler)
	mux.HandleFunc("PATCH /subscriptions/{id}", h.updateHandler)
	mux.HandleFunc("DELETE /subscriptions/{id}", h.deleteHandler)
	mux.HandleFunc("GET /subscriptions", h.listHandler)
	mux.HandleFunc("GET /subscriptions/total", h.totalCostHandler)
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
	start = start.Local()

	var endPtr *time.Time
	if strings.TrimSpace(req.EndDate) != "" {
		end, _ := time.Parse("01-2006", req.EndDate)
		tt := end.Local()
		endPtr = &tt
	}

	id, err := h.svc.Create(r.Context(), models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   start,
		EndDate:     endPtr,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidation):
			httpx.HttpError(w, http.StatusBadRequest, err.Error())
		default:
			h.logger.Error("ошибка при создании подписки", zap.Error(err))
			httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
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

	dataItem, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
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
		StartDate:   dataItem.StartDate.Local().Format("01-2006"),
	}
	if dataItem.EndDate != nil {
		resp.EndDate = dataItem.EndDate.Local().Format("01-2006")
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
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		httpx.HttpError(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var req updateSubscriptionReq

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		httpx.HttpError(w, http.StatusBadRequest, "Некорректный JSON")
		return
	}

	if err := validateUpdateSubscription(req); err != nil {
		httpx.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	dataItem, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			httpx.HttpError(w, http.StatusNotFound, "Подписка не найдена")
			return
		}
		h.logger.Error("Ошибка GetByID() при обновлении подписки",
			zap.Int("id", id),
			zap.Error(err),
		)
		httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	if req.ServiceName != nil {
		dataItem.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		dataItem.Price = *req.Price
	}
	if req.UserID != nil {
		dataItem.UserID = *req.UserID
	}
	if req.StartDate != nil {
		t, _ := time.Parse("01-2006", *req.StartDate)
		dataItem.StartDate = t.Local()
	}

	if req.EndDate != nil {
		if strings.TrimSpace(*req.EndDate) == "" {
			dataItem.EndDate = nil
		} else {
			t, _ := time.Parse("01-2006", *req.EndDate)
			tt := t.Local()
			dataItem.EndDate = &tt
		}
	}

	if err := h.svc.Update(r.Context(), dataItem); err != nil {
		switch {
		case errors.Is(err, services.ErrValidation):
			httpx.HttpError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, services.ErrNotFound):
			httpx.HttpError(w, http.StatusNotFound, "Подписка не найдена")
		default:
			h.logger.Error("Ошибка Update()", zap.Int("id", id), zap.Error(err))
			httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		httpx.HttpError(w, http.StatusBadRequest, "Некорректный ID")
		return
	}
	
	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, services.ErrNotFound) {
			httpx.HttpError(w, http.StatusNotFound, "Подписка не найдена")
			return
		}
		h.logger.Error("Ошибка Delete()", zap.Int("id", id), zap.Error(err))
		httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) listHandler(w http.ResponseWriter, r *http.Request) {
	var filter services.ListFilter
	if err := validateListSubscription(r.URL.Query(), &filter); err != nil {
		httpx.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	data, err := h.svc.List(r.Context(), filter)
	if err != nil {
		h.logger.Error("Ошибка List()", zap.Error(err))
		httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		return
	}

	resp := make([]subscriptionRes, 0, len(data))
	
	for _, dataItem := range data {
		respItem := subscriptionRes{
			ID:          dataItem.ID,
			ServiceName: dataItem.ServiceName,
			Price:       dataItem.Price,
			UserID:      dataItem.UserID,
			StartDate:   dataItem.StartDate.Local().Format("01-2006"),
		}
		if dataItem.EndDate != nil {
			respItem.EndDate = dataItem.EndDate.Local().Format("01-2006")
		}
		resp = append(resp, respItem)
	}

	if err := httpx.WriteJSON(w, http.StatusOK, resp); err != nil {
		switch {
		case errors.Is(err, httpx.ErrJSONMarshal):
			h.logger.Error("не удалось сериализовать JSON", zap.Error(err))
			httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		case errors.Is(err, httpx.ErrWriteBody):
			h.logger.Warn("клиент закрыл соединение, ответ не был отправлен", zap.Error(err))
		}
	}
}

func (h *Handler) totalCostHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	if err := validateTotalCost(query); err != nil {
		httpx.HttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	periodStart, _ := time.Parse("01-2006", strings.TrimSpace(query.Get("period_start")))
	periodEnd, _ := time.Parse("01-2006", strings.TrimSpace(query.Get("period_end")))
	
	filter := services.ListFilter{}
	if userID := strings.TrimSpace(query.Get("user_id")); userID != "" {
		filter.UserID, filter.HasUserID = userID, true
	}
	if serviceName := strings.TrimSpace(query.Get("service_name")); serviceName != "" {
		filter.ServiceName, filter.HasServiceName = serviceName, true
	}

	sum, err := h.svc.TotalCost(r.Context(), periodStart, periodEnd, filter)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrValidation):
			httpx.HttpError(w, http.StatusBadRequest, err.Error())
		default:
			h.logger.Error("Ошибка TotalCost()", zap.Error(err))
			httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		}
		return
	}

	if err := httpx.WriteJSON(w, http.StatusOK, map[string]int{"total_cost": sum}); err != nil {
		switch {
		case errors.Is(err, httpx.ErrJSONMarshal):
			h.logger.Error("не удалось сериализовать JSON", zap.Error(err))
			httpx.HttpError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
		case errors.Is(err, httpx.ErrWriteBody):
			h.logger.Warn("клиент закрыл соединение, ответ не был отправлен", zap.Error(err))
		}
	}
}
