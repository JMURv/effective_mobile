package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionRoutes(t *testing.T) {
	ts := NewTestEnv(t)

	client := &http.Client{}
	decode := func(r *http.Response, dst any) {
		defer r.Body.Close()
		_ = json.NewDecoder(r.Body).Decode(dst)
	}

	var createdID uuid.UUID
	testUserID := uuid.New()

	t.Run("CreateSubscription_Success", func(t *testing.T) {
		body := map[string]any{
			"service_name": "netflix",
			"price":        500,
			"user_id":      testUserID.String(),
			"start_date":   time.Now(),
			"end_date":     nil,
		}

		b, _ := json.Marshal(body)

		resp, err := client.Post(
			ts.Server.URL+"/subscriptions",
			"application/json",
			bytes.NewBuffer(b),
		)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("ListSubscriptions_ShouldReturnData", func(t *testing.T) {
		resp, err := client.Get(ts.Server.URL + "/subscriptions")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res []map[string]any
		decode(resp, &res)

		assert.NotEmpty(t, res)

		// save ID for later tests
		if len(res) > 0 {
			if idStr, ok := res[0]["id"].(string); ok {
				createdID, _ = uuid.Parse(idStr)
			}
		}
	})

	t.Run("GetByID_Success", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdID)

		resp, err := client.Get(ts.Server.URL + "/subscriptions/" + createdID.String())
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]any
		decode(resp, &res)

		assert.Equal(t, createdID.String(), res["id"])
	})

	t.Run("UpdateSubscription_Success", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdID)

		body := map[string]any{
			"service_name": "netflix-updated",
			"price":        999,
		}

		b, _ := json.Marshal(body)

		req, err := http.NewRequest(
			http.MethodPut,
			ts.Server.URL+"/subscriptions/"+createdID.String(),
			bytes.NewBuffer(b),
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("CalculateTotalCost_Success", func(t *testing.T) {
		reqBody := map[string]any{
			"service_name": "netflix-updated",
			"user_id":      testUserID.String(),
			"from":         time.Now().Add(-24 * time.Hour),
			"to":           time.Now().Add(24 * time.Hour),
		}

		b, _ := json.Marshal(reqBody)

		req, err := http.NewRequest(
			http.MethodPost,
			ts.Server.URL+"/subscriptions/total",
			bytes.NewBuffer(b),
		)
		require.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]int64
		decode(resp, &res)

		assert.Contains(t, res, "total_cost")
	})

	t.Run("DeleteSubscription_Success", func(t *testing.T) {
		require.NotEqual(t, uuid.Nil, createdID)

		req, err := http.NewRequest(
			http.MethodDelete,
			ts.Server.URL+"/subscriptions/"+createdID.String(),
			nil,
		)
		require.NoError(t, err)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})
}
