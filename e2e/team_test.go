package e2e

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pr-review/internal/http/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAddTeam проверяет создание команды с участниками
func TestAddTeam(t *testing.T) {
	st := NewSuite()
	st.Start()
	t.Cleanup(func() {
		st.srv.Stop()
	})

	testcases := []struct {
		name                string
		headers             map[string]string
		body                []byte
		teamName            string
		expectedStatus      int
		expectedResponse    string
		expectedGetResponse string
	}{
		{
			name:     "Valid Request",
			headers:  map[string]string{"Content-Type": "application/json"},
			teamName: "payments",
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
			expectedGetResponse: `{
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
				],
				"team_name": "payments"
			}`,
			expectedStatus: 201,
		},
		{
			name:     "enpty team",
			headers:  map[string]string{"Content-Type": "application/json"},
			teamName: "paymentsasd",
			body: []byte(`{
				"team_name": "paymentsasd",
				"members": []
			}`),
			expectedResponse: `{
			"team": {
				"team_name": "paymentsasd",
				"members": []
			}
			}`,
			expectedGetResponse: `{
				"team_name": "paymentsasd",
				"members": []
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

			req = httptest.NewRequestWithContext(t.Context(), "GET", "/team/get?team_name="+tc.teamName, nil)
			res = httptest.NewRecorder()
			st.srv.TestReq(req, res)
			assert.Equal(t, http.StatusOK, res.Result().StatusCode)
			require.JSONEq(t, tc.expectedGetResponse, res.Body.String())
		})
	}
}

// TestAddTeamWithExistingName проверяет, что при попытке создать команду с уже существующим именем
// возвращается ошибка с кодом 400 и соответствующим сообщением.
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

func createTeam(t *testing.T, st *Suite, reqBody *dto.AddTeamRequest) ([]byte, int) {
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/team/add", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	st.srv.TestReq(req, recorder)

	return recorder.Body.Bytes(), recorder.Result().StatusCode
}

func getTeam(t *testing.T, st *Suite, name string) (*dto.GetTeamResponse, int) {
	req := httptest.NewRequestWithContext(t.Context(), "GET", "/team/get?team_name="+name, nil)
	recorder := httptest.NewRecorder()

	st.srv.TestReq(req, recorder)

	var res dto.GetTeamResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &res)
	require.NoError(t, err)

	return &res, recorder.Result().StatusCode
}
