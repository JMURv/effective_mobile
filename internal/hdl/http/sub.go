package http

import (
	"net/http"

	"github.com/JMURv/golang-clean-template/internal/dto"
	"github.com/JMURv/golang-clean-template/internal/hdl/http/utils"
	_ "github.com/JMURv/golang-clean-template/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) RegisterSubscriptionRoutes() {
	h.Router.Post("/subscriptions", h.createSubscription)
	h.Router.Get("/subscriptions", h.listSubscriptions)
	h.Router.Get("/subscriptions/{id}", h.getSubscriptionByID)
	h.Router.Put("/subscriptions/{id}", h.updateSubscription)
	h.Router.Delete("/subscriptions/{id}", h.deleteSubscription)
	h.Router.Get("/subscriptions/total", h.calculateTotalCost)
}

func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	req := dto.CreateSubscriptionRequest{}
	if ok := utils.ParseAndValidate(w, r, &req); !ok {
		return
	}

	err := h.ctrl.Create(r.Context(), req)
	if err != nil {
		// TODO: not found, conflict, swagger, test
		utils.HandleError(w, err)
		return
	}

	utils.StatusResponse(w, http.StatusCreated)
}

func (h *Handler) getSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	res, err := h.ctrl.GetByID(r.Context(), id)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, res)
}

func (h *Handler) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	res, err := h.ctrl.List(r.Context())
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, res)
}

func (h *Handler) updateSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	req := dto.UpdateSubscriptionRequest{}
	if ok := utils.ParseAndValidate(w, r, &req); !ok {
		return
	}

	err = h.ctrl.Update(r.Context(), id, req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.StatusResponse(w, http.StatusOK)
}

func (h *Handler) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.ctrl.Delete(r.Context(), id)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.StatusResponse(w, http.StatusNoContent)
}

func (h *Handler) calculateTotalCost(w http.ResponseWriter, r *http.Request) {
	req := dto.CalculateTotalCostRequest{}
	if ok := utils.ParseAndValidate(w, r, &req); !ok {
		return
	}
	// TODO: custom validate if one date less than other
	res, err := h.ctrl.CalculateTotalCost(r.Context(), req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, res)
}
