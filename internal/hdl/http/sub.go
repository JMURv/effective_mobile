package http

import (
	"net/http"

	"github.com/JMURv/effective-mobile/internal/dto"
	"github.com/JMURv/effective-mobile/internal/hdl/http/utils"
	_ "github.com/JMURv/effective-mobile/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *Handler) RegisterSubscriptionRoutes() {
	h.Router.Post("/subscriptions", h.createSubscription)
	h.Router.Get("/subscriptions", h.listSubscriptions)
	h.Router.Get("/subscriptions/{id}", h.getSubscriptionByID)
	h.Router.Put("/subscriptions/{id}", h.updateSubscription)
	h.Router.Delete("/subscriptions/{id}", h.deleteSubscription)
	h.Router.Post("/subscriptions/total", h.calculateTotalCost)
}

// createSubscription godoc
//
//	@Summary		Create subscription
//	@Description	Create a new subscription
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			request	body	dto.CreateSubscriptionRequest	true	"Subscription payload"
//	@Success		201
//	@Failure		400	{object}	utils.ErrorsResponse
//	@Failure		500	{object}	utils.ErrorsResponse
//	@Router			/subscriptions [post]
func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	req := dto.CreateSubscriptionRequest{}
	if ok := utils.ParseAndValidate(w, r, &req); !ok {
		return
	}

	err := h.ctrl.Create(r.Context(), req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.StatusResponse(w, http.StatusCreated)
}

// getSubscriptionByID godoc
//
//	@Summary		Get subscription by ID
//	@Description	Get a single subscription by UUID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Subscription ID"
//	@Success		200	{object}	models.Subscription
//	@Failure		400	{object}	utils.ErrorsResponse
//	@Failure		404	{object}	utils.ErrorsResponse
//	@Failure		500	{object}	utils.ErrorsResponse
//	@Router			/subscriptions/{id} [get]
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

// listSubscriptions godoc
//
//	@Summary		List subscriptions
//	@Description	Get all subscriptions
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.Subscription
//	@Failure		500	{object}	utils.ErrorsResponse
//	@Router			/subscriptions [get]
func (h *Handler) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	res, err := h.ctrl.List(r.Context())
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, res)
}

// updateSubscription godoc
//
//	@Summary		Update subscription
//	@Description	Update subscription by ID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string							true	"Subscription ID"
//	@Param			request	body	dto.UpdateSubscriptionRequest	true	"Update payload"
//	@Success		200
//	@Failure		400	{object}	utils.ErrorsResponse
//	@Failure		404	{object}	utils.ErrorsResponse
//	@Failure		500	{object}	utils.ErrorsResponse
//	@Router			/subscriptions/{id} [put]
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

// deleteSubscription godoc
//
//	@Summary		Delete subscription
//	@Description	Delete subscription by ID
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Subscription ID"
//	@Success		204
//	@Failure		400	{object}	utils.ErrorsResponse
//	@Failure		404	{object}	utils.ErrorsResponse
//	@Failure		500	{object}	utils.ErrorsResponse
//	@Router			/subscriptions/{id} [delete]
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

// calculateTotalCost godoc
//
//	@Summary		Calculate total subscription cost
//	@Description	Calculate total cost for subscriptions within a period
//	@Tags			subscriptions
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CalculateTotalCostRequest	true	"Calculation payload"
//	@Success		200		{object}	dto.CalculateTotalCostResponse
//	@Failure		400		{object}	utils.ErrorsResponse
//	@Failure		500		{object}	utils.ErrorsResponse
//	@Router			/subscriptions/total [post]
func (h *Handler) calculateTotalCost(w http.ResponseWriter, r *http.Request) {
	req := dto.CalculateTotalCostRequest{}
	if ok := utils.ParseAndValidate(w, r, &req); !ok {
		return
	}

	res, err := h.ctrl.CalculateTotalCost(r.Context(), req)
	if err != nil {
		utils.HandleError(w, err)
		return
	}

	utils.SuccessResponse(w, http.StatusOK, res)
}
