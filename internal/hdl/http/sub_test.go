package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/JMURv/effective-mobile/internal/ctrl"
	"github.com/JMURv/effective-mobile/internal/dto"
	"github.com/JMURv/effective-mobile/internal/hdl"
	"github.com/JMURv/effective-mobile/internal/hdl/http/utils"
	"github.com/JMURv/effective-mobile/internal/hdl/validation"
	md "github.com/JMURv/effective-mobile/internal/models"
	"github.com/JMURv/effective-mobile/tests/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	validation.V = validator.New()
	os.Exit(m.Run())
}

func TestHandler_CreateSubscription(t *testing.T) {
	const uri = "/subscriptions"

	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := New(mctrl)

	testErr := errors.New("test error")

	tests := []struct {
		name       string
		payload    any
		status     int
		expect     func()
		assertions func(r *httptest.ResponseRecorder)
	}{
		{
			name:    "BadRequest_InvalidJSON",
			payload: "invalid-json",
			status:  http.StatusBadRequest,
			expect:  func() {},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Contains(t, res.Errors[0], "decode")
			},
		},
		{
			name:    "ValidationError",
			payload: map[string]any{"service_name": ""},
			status:  http.StatusBadRequest,
			expect:  func() {},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				err := json.NewDecoder(r.Body).Decode(&res)
				require.NoError(t, err)

				assert.Contains(t, res.Errors[0], "required")
			},
		},
		{
			name: "InternalError",
			payload: map[string]any{
				"service_name": "netflix",
				"price":        100,
				"user_id":      uuid.New().String(),
				"start_date":   time.Now(),
			},
			status: http.StatusInternalServerError,
			expect: func() {
				mctrl.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(testErr)
			},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Equal(t, hdl.ErrInternal.Error(), res.Errors[0])
			},
		},
		{
			name: "Success",
			payload: map[string]any{
				"service_name": "netflix",
				"price":        100,
				"user_id":      uuid.New().String(),
				"start_date":   time.Now(),
			},
			status: http.StatusCreated,
			expect: func() {
				mctrl.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			assertions: func(r *httptest.ResponseRecorder) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			var body bytes.Buffer
			if s, ok := tt.payload.(string); ok {
				body.WriteString(s)
			} else {
				_ = json.NewEncoder(&body).Encode(tt.payload)
			}

			req := httptest.NewRequest(http.MethodPost, uri, &body)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			h.createSubscription(w, req)

			assert.Equal(t, tt.status, w.Code)

			tt.assertions(w)
		})
	}
}

func TestHandler_GetSubscriptionByID(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := New(mctrl)

	id := uuid.New()
	testErr := errors.New("test error")

	tests := []struct {
		name       string
		id         string
		status     int
		expect     func()
		assertions func(*httptest.ResponseRecorder)
	}{
		{
			name:   "BadUUID",
			id:     "invalid",
			status: http.StatusBadRequest,
			expect: func() {},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.NotEmpty(t, res.Errors)
			},
		},
		{
			name:   "NotFound",
			id:     id.String(),
			status: http.StatusNotFound,
			expect: func() {
				mctrl.EXPECT().
					GetByID(gomock.Any(), id).
					Return(md.Subscription{}, ctrl.ErrNotFound)
			},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Equal(t, ctrl.ErrNotFound.Error(), res.Errors[0])
			},
		},
		{
			name:   "InternalError",
			id:     id.String(),
			status: http.StatusInternalServerError,
			expect: func() {
				mctrl.EXPECT().
					GetByID(gomock.Any(), id).
					Return(md.Subscription{}, testErr)
			},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Equal(t, hdl.ErrInternal.Error(), res.Errors[0])
			},
		},
		{
			name:   "Success",
			id:     id.String(),
			status: http.StatusOK,
			expect: func() {
				mctrl.EXPECT().
					GetByID(gomock.Any(), id).
					Return(md.Subscription{ID: id}, nil)
			},
			assertions: func(r *httptest.ResponseRecorder) {
				var res md.Subscription
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Equal(t, id, res.ID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			req := httptest.NewRequest(http.MethodGet, "/subscriptions/"+tt.id, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			h.getSubscriptionByID(w, req)

			assert.Equal(t, tt.status, w.Code)
			tt.assertions(w)
		})
	}
}

func TestHandler_ListSubscriptions(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := New(mctrl)

	testErr := errors.New("test error")

	tests := []struct {
		name       string
		status     int
		expect     func()
		assertions func(*httptest.ResponseRecorder)
	}{
		{
			name:   "InternalError",
			status: http.StatusInternalServerError,
			expect: func() {
				mctrl.EXPECT().
					List(gomock.Any()).
					Return(nil, testErr)
			},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Equal(t, hdl.ErrInternal.Error(), res.Errors[0])
			},
		},
		{
			name:   "Success",
			status: http.StatusOK,
			expect: func() {
				mctrl.EXPECT().
					List(gomock.Any()).
					Return([]md.Subscription{{}}, nil)
			},
			assertions: func(r *httptest.ResponseRecorder) {
				var res []md.Subscription
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Len(t, res, 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			req := httptest.NewRequest(http.MethodGet, "/subscriptions", nil)
			w := httptest.NewRecorder()

			h.listSubscriptions(w, req)

			assert.Equal(t, tt.status, w.Code)
			tt.assertions(w)
		})
	}
}

func TestHandler_UpdateSubscription(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := New(mctrl)

	id := uuid.New()
	testErr := errors.New("test error")

	tests := []struct {
		name       string
		id         string
		payload    any
		status     int
		expect     func()
		assertions func(*httptest.ResponseRecorder)
	}{
		{
			name:    "BadUUID",
			id:      "bad",
			payload: map[string]any{"service_name": "netflix"},
			status:  http.StatusBadRequest,
			expect:  func() {},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.NotEmpty(t, res.Errors)
			},
		},
		{
			name:    "ValidationError",
			id:      id.String(),
			payload: map[string]any{"price": -1},
			status:  http.StatusBadRequest,
			expect:  func() {},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.NotEmpty(t, res.Errors)
			},
		},
		{
			name:    "InternalError",
			id:      id.String(),
			payload: map[string]any{"service_name": "netflix"},
			status:  http.StatusInternalServerError,
			expect: func() {
				mctrl.EXPECT().
					Update(gomock.Any(), id, gomock.Any()).
					Return(testErr)
			},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Equal(t, hdl.ErrInternal.Error(), res.Errors[0])
			},
		},
		{
			name:    "Success",
			id:      id.String(),
			payload: map[string]any{"service_name": "netflix"},
			status:  http.StatusOK,
			expect: func() {
				mctrl.EXPECT().
					Update(gomock.Any(), id, gomock.Any()).
					Return(nil)
			},
			assertions: func(r *httptest.ResponseRecorder) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			var body bytes.Buffer
			_ = json.NewEncoder(&body).Encode(tt.payload)

			req := httptest.NewRequest(http.MethodPut, "/subscriptions/"+tt.id, &body)
			req.Header.Set("Content-Type", "application/json")

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			h.updateSubscription(w, req)

			assert.Equal(t, tt.status, w.Code)
			tt.assertions(w)
		})
	}
}

func TestHandler_DeleteSubscription(t *testing.T) {
	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := New(mctrl)

	id := uuid.New()
	testErr := errors.New("test error")

	tests := []struct {
		name   string
		id     string
		status int
		expect func()
	}{
		{
			name:   "BadUUID",
			id:     "bad",
			status: http.StatusBadRequest,
			expect: func() {},
		},
		{
			name:   "InternalError",
			id:     id.String(),
			status: http.StatusInternalServerError,
			expect: func() {
				mctrl.EXPECT().
					Delete(gomock.Any(), id).
					Return(testErr)
			},
		},
		{
			name:   "Success",
			id:     id.String(),
			status: http.StatusNoContent,
			expect: func() {
				mctrl.EXPECT().
					Delete(gomock.Any(), id).
					Return(nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			req := httptest.NewRequest(http.MethodDelete, "/subscriptions/"+tt.id, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()
			h.deleteSubscription(w, req)

			assert.Equal(t, tt.status, w.Code)
		})
	}
}

func TestHandler_CalculateTotalCost(t *testing.T) {
	const uri = "/subscriptions/total"

	mock := gomock.NewController(t)
	defer mock.Finish()

	mctrl := mocks.NewMockAppCtrl(mock)
	h := New(mctrl)

	testErr := errors.New("test error")

	now := time.Now()

	tests := []struct {
		name       string
		payload    any
		status     int
		expect     func()
		assertions func(*httptest.ResponseRecorder)
	}{
		{
			name:    "BadJSON",
			payload: "invalid-json",
			status:  http.StatusBadRequest,
			expect:  func() {},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Contains(t, res.Errors[0], "decode")
			},
		},
		{
			name: "ValidationError",
			payload: map[string]any{
				"service_name": "",
				"user_id":      "",
				"from":         now,
				"to":           now,
			},
			status: http.StatusBadRequest,
			expect: func() {},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.NotEmpty(t, res.Errors)
			},
		},
		{
			name: "InternalError",
			payload: map[string]any{
				"service_name": "netflix",
				"user_id":      uuid.New().String(),
				"from":         now.Add(-time.Hour),
				"to":           now,
			},
			status: http.StatusInternalServerError,
			expect: func() {
				mctrl.EXPECT().
					CalculateTotalCost(gomock.Any(), gomock.Any()).
					Return(dto.CalculateTotalCostResponse{}, testErr)
			},
			assertions: func(r *httptest.ResponseRecorder) {
				var res utils.ErrorsResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Equal(t, hdl.ErrInternal.Error(), res.Errors[0])
			},
		},
		{
			name: "Success",
			payload: map[string]any{
				"service_name": "netflix",
				"user_id":      uuid.New().String(),
				"from":         now.Add(-time.Hour),
				"to":           now,
			},
			status: http.StatusOK,
			expect: func() {
				mctrl.EXPECT().
					CalculateTotalCost(gomock.Any(), gomock.Any()).
					Return(dto.CalculateTotalCostResponse{
						TotalCost: 1500,
					}, nil)
			},
			assertions: func(r *httptest.ResponseRecorder) {
				var res dto.CalculateTotalCostResponse
				_ = json.NewDecoder(r.Body).Decode(&res)
				assert.Equal(t, int64(1500), res.TotalCost)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expect()

			var body bytes.Buffer
			if s, ok := tt.payload.(string); ok {
				body.WriteString(s)
			} else {
				_ = json.NewEncoder(&body).Encode(tt.payload)
			}

			req := httptest.NewRequest(http.MethodPost, uri, &body)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			h.calculateTotalCost(w, req)

			assert.Equal(t, tt.status, w.Code)
			tt.assertions(w)
		})
	}
}
