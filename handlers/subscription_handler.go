package handlers

import (
	"encoding/json"
	"go-rest-crud/model"
	"go-rest-crud/service"
	"log"
	"net/http"
)

type SubscriptionHandler struct {
	service *service.SubscriptionService
}

func NewSubscriptionHandler(service *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{service: service}
}

// CreateSubscription godoc
// @Summary      Create subscription
// @Description  Creates a new subscription. Checks for duplicate service_name, user_id, start_date combination.
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        subscription body model.CreateSubscriptionRequest true "Subscription data"
// @Success      201 {object} model.SubscriptionSwagger "Subscription created"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      409 {object} map[string]string "Subscription already exists"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router      /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] POST /subscriptions - Creating subscription")

	if r.ContentLength == 0 || r.Header.Get("Content-Type") != "application/json" {
		log.Printf("[HANDLER] Invalid JSON: Content-Type must be application/json")
		writeError(w, http.StatusUnsupportedMediaType, "Invalid JSON:Content-Type must be application/json")
		return
	}

	var input model.CreateSubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("[HANDLER] Invalid JSON: %v", err)
		writeError(w, http.StatusBadRequest, "invalid JSON format")
		return
	}

	subscription, err := h.service.CreateSubscription(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(subscription)
}

// GetSubscriptions godoc
// @Summary      Get subscriptions
// @Description  Returns list of subscriptions with optional filtering by service_name, user_id, start_date, end_date
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        service_name query string false "Filter by service name"
// @Param        user_id query string false "Filter by user UUID"
// @Param        start_date query string false "Filter by start date (MM-YYYY)"
// @Param        end_date query string false "Filter by end date (MM-YYYY)"
// @Success      200 {array} model.SubscriptionSwagger "List of subscriptions"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router      /subscriptions [get]
func (h *SubscriptionHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] GET /subscriptions - Fetching subscriptions")

	filter := model.FilterSubscriptionRequest{
		ServiceName: r.URL.Query().Get("service_name"),
		UserID:      r.URL.Query().Get("user_id"),
		StartDate:   r.URL.Query().Get("start_date"),
		EndDate:     r.URL.Query().Get("end_date"),
	}

	subscriptions, err := h.service.GetAllSubscriptions(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscriptions)
}

// GetSubscriptionsSumma godoc
// @Summary      Get summa of subscription prices
// @Description  Returns the sum of all price values for subscriptions matching the filters
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        service_name query string false "Filter by service name"
// @Param        user_id query string false "Filter by user UUID"
// @Param        start_date query string false "Filter by start date (MM-YYYY)"
// @Param        end_date query string false "Filter by end date (MM-YYYY)"
// @Success      200 {object} model.SummaResponse "Total price"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router      /subscriptions/summa [get]
func (h *SubscriptionHandler) GetSubscriptionsSumma(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HANDLER] GET /subscriptions/summa")

	filter := model.FilterSubscriptionRequest{
		ServiceName: r.URL.Query().Get("service_name"),
		UserID:      r.URL.Query().Get("user_id"),
		StartDate:   r.URL.Query().Get("start_date"),
		EndDate:     r.URL.Query().Get("end_date"),
	}

	summa, err := h.service.GetAllSubscriptionsSumma(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summa)
}

// GetSubscription godoc
// @Summary      Get subscription by ID
// @Description  Returns a single subscription by its UUID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path string true "Subscription UUID"
// @Success      200 {object} model.SubscriptionSwagger "Subscription found"
// @Failure      400 {object} map[string]string "Invalid UUID"
// @Failure      404 {object} map[string]string "Subscription not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router      /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("[HANDLER] GET /subscriptions/{id} - Get subscription: %s", id)

	subscription, err := h.service.GetSubscriptionByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// UpdateSubscription godoc
// @Summary      Update subscription
// @Description  Updates an existing subscription by its UUID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path string true "Subscription UUID"
// @Param        subscription body model.UpdateSubscriptionRequest true "Subscription update data"
// @Success      200 {object} model.SubscriptionSwagger "Subscription updated"
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      404 {object} map[string]string "Subscription not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router      /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("[HANDLER] PUT /subscriptions/{id} - Updating subscription: %s", id)

	if r.ContentLength == 0 || r.Header.Get("Content-Type") != "application/json" {
		log.Printf("[HANDLER] Invalid JSON: Content-Type must be application/json")
		writeError(w, http.StatusUnsupportedMediaType, "Invalid JSON:Content-Type must be application/json")
		return
	}

	var input model.UpdateSubscriptionRequest
	input.ID = id
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("[HANDLER] Invalid JSON: %v", err)
		writeError(w, http.StatusBadRequest, "invalid JSON format")
		return
	}

	subscription, err := h.service.UpdateSubscription(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subscription)
}

// DeleteSubscription godoc
// @Summary      Delete subscription
// @Description  Deletes a subscription by its UUID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id path string true "Subscription UUID"
// @Success      204 "No Content"
// @Failure      400 {object} map[string]string "Invalid UUID"
// @Failure      404 {object} map[string]string "Subscription not found"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router      /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	log.Printf("[HANDLER] PUT /subscriptions/{id} - Delete subscription: %s", id)

	err := h.service.DeleteSubscription(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
