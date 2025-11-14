package e2e

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddTeam(t *testing.T) {
	st := NewSuite()
	st.Start()
	t.Cleanup(func() {
		st.srv.Stop()
	})

	testcases := []struct {
		name             string
		headers          map[string]string
		body             []byte
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:    "Valid Request",
			headers: map[string]string{"Content-Type": "application/json"},
			body: []byte(`{
				"team_name": "payments",
				"members": [
					{
						"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d9ad",
						"username": "Alice",
						"is_active": true
					},
					{
						"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d9a1",
						"username": "Bob",
						"is_active": true
					}
				]
			}`),
			expectedResponse: `{
			"team": {
				"team_name": "payments",
				"members": [
					{
						"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d9ad",
						"username": "Alice",
						"is_active": true
					},
					{
						"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d9a1",
						"username": "Bob",
						"is_active": true
					}
				]
			}
			}`,
			expectedStatus: 201,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequestWithContext(t.Context(), "POST", "/team/add", bytes.NewBuffer(tc.body))
			for k, v := range tc.headers {
				req.Header.Set(k, v)
			}
			res := httptest.NewRecorder()

			st.srv.TestReq(req, res)

			assert.Equal(t, tc.expectedStatus, res.Result().StatusCode)
			require.JSONEq(t, tc.expectedResponse, res.Body.String())
		})
	}
}

func TestAddTeamWithExistingName(t *testing.T) {
	st := NewSuite()
	st.Start()
	t.Cleanup(func() {
		st.srv.Stop()
	})

	body := []byte(`{
		"team_name": "payments123",
		"members": [
			{
				"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d9a3",
				"username": "Alice1",
				"is_active": true
			},
			{
				"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d951",
				"username": "Bo1b",
				"is_active": true
			}
		]
	}`)
	resBody := []byte(`{
	"team": {
		"team_name": "payments123",
		"members": [
			{
				"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d9a3",
				"username": "Alice1",
				"is_active": true
			},
			{
				"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d951",
				"username": "Bo1b",
				"is_active": true
			}
		]
	}
	}`)
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	st.srv.TestReq(req, res)

	assert.Equal(t, 201, res.Result().StatusCode)
	require.JSONEq(t, string(resBody), res.Body.String())

	req = httptest.NewRequestWithContext(t.Context(), "POST", "/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()
	st.srv.TestReq(req, res)
	assert.Equal(t, 400, res.Result().StatusCode)
	require.JSONEq(t, `{
	"error": {
		"code": "TEAM_EXISTS",
		"message": "team_name already exists"
	}
	}`, res.Body.String())
}

func TestUnassignPRAfterUpdateUserByAddTeam(t *testing.T) {
	st := NewSuite()
	st.Start()
	t.Cleanup(func() {
		st.srv.Stop()
	})

	// 1. Создаём команду
	body := []byte(`{
		"team_name": "payments123",
		"members": [
			{
				"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d9a3",
				"username": "Alice1",
				"is_active": true
			},
			{
				"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d951",
				"username": "Bo1b",
				"is_active": true
			},
			{
				"user_id": "86a83214-a5d1-4e8c-93a2-e5bdc206d951",
				"username": "Bo21b",
				"is_active": true
			}
		]
	}`)
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	st.srv.TestReq(req, res)

	require.Equal(t, 201, res.Result().StatusCode)

	// 2. Создаём PR
	body = []byte(`{
		"pull_request_id": "5df8139e-b302-4851-93fd-4086091347a6",
		"pull_request_name": "aboba",
		"author_id": "86a83214-a5d1-4e8c-93a2-e5bdc206d951"
	}`)
	req = httptest.NewRequestWithContext(t.Context(), "POST", "/pullRequest/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()

	st.srv.TestReq(req, res)

	require.Equal(t, 201, res.Result().StatusCode)

	body = []byte(`{
		"team_name": "payments1234",
		"members": [
			{
				"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d9a3",
				"username": "Alice1",
				"is_active": true
			},
			{
				"user_id": "86a832a4-a5d1-4e8c-93a2-e5bdc206d951",
				"username": "Bo1b",
				"is_active": true
			}
		]
	}`)
	req = httptest.NewRequestWithContext(t.Context(), "POST", "/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res = httptest.NewRecorder()

	st.srv.TestReq(req, res)

	require.Equal(t, 201, res.Result().StatusCode)
}
