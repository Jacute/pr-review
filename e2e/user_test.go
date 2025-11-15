package e2e

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
)

// TestUnassignPRAfterDeactivateUserByAddTeam проверяет, что при обновлении команды и деактивации пользователя,
// у всех открытых PR'ов этого пользователя ассайн удаляется.
func TestUnassignPRAfterDeactivateUserByAddTeam(t *testing.T) {
	st := NewSuite()
	st.Start()
	t.Cleanup(func() {
		st.srv.Stop()
	})

	// 1. Создаём команду
	_, code := createTeam(t, st, &dto.AddTeamRequest{
		Name: "payments1234",
		Members: []*models.Member{
			{
				Id:       "86a832a4-a5d1-4e8c-93a2-e5bdc206d9a3",
				Username: "Alice1",
				IsActive: true,
			},
			{
				Id:       "86a832a4-a5d1-4e8c-93a2-e5bdc206d951",
				Username: "Bo1b",
				IsActive: true,
			},
			{
				Id:       "86a83214-a5d1-4e8c-93a2-e5bdc206d951",
				Username: "Bo21b",
				IsActive: true,
			},
		},
	})
	require.Equal(t, 201, code)

	// 2. Создаём PR. Должно быть два аппрувера - Alice1 и Bo1b
	response, code, id, name := createPR(t, st)
	expectedRes2 := dto.CreatePRResponse{
		PR: &models.PullRequest{
			PullRequestShort: models.PullRequestShort{
				Id:       id,
				Title:    name,
				AuthorId: "86a83214-a5d1-4e8c-93a2-e5bdc206d951",
				Status:   models.StatusOpen,
			},
			Reviewers: []string{
				"86a832a4-a5d1-4e8c-93a2-e5bdc206d9a3",
				"86a832a4-a5d1-4e8c-93a2-e5bdc206d951",
			},
			MergedAt: nil,
		},
	}

	require.Equal(t, 201, code)
	require.ElementsMatch(t, expectedRes2.PR.Reviewers, response.PR.Reviewers)
	require.Equal(t, expectedRes2.PR.PullRequestShort, response.PR.PullRequestShort)
	require.Equal(t, expectedRes2.PR.MergedAt, response.PR.MergedAt)

	// 3. Проверяем, что у пользователя Bo1b есть assignee
	expectedRes3 := dto.GetReviewResponse{
		UserId: "86a832a4-a5d1-4e8c-93a2-e5bdc206d951",
		PullRequests: []*models.PullRequest{
			{
				PullRequestShort: models.PullRequestShort{
					Id:       id,
					Title:    name,
					AuthorId: "86a83214-a5d1-4e8c-93a2-e5bdc206d951",
					Status:   models.StatusOpen,
				},
				Reviewers: []string{
					"86a832a4-a5d1-4e8c-93a2-e5bdc206d9a3",
					"86a832a4-a5d1-4e8c-93a2-e5bdc206d951",
				},
				MergedAt: nil,
			},
		},
	}
	req := httptest.NewRequestWithContext(t.Context(), "GET", "/users/getReview?user_id=86a832a4-a5d1-4e8c-93a2-e5bdc206d951", nil)
	recorder := httptest.NewRecorder()
	st.srv.TestReq(req, recorder)
	require.Equal(t, 200, recorder.Result().StatusCode)
	var res3 dto.GetReviewResponse
	json.Unmarshal(recorder.Body.Bytes(), &res3)

	require.Len(t, res3.PullRequests, 1)
	require.ElementsMatch(t, expectedRes3.PullRequests[0].Reviewers, res3.PullRequests[0].Reviewers)
	require.Equal(t, expectedRes3.PullRequests[0].PullRequestShort, res3.PullRequests[0].PullRequestShort)
	require.Equal(t, expectedRes3.PullRequests[0].MergedAt, res3.PullRequests[0].MergedAt)
	require.Equal(t, res3.UserId, "86a832a4-a5d1-4e8c-93a2-e5bdc206d951")

	// 4. Переводим пользователя Bo1b в другую команду и деактивируем
	_, code = createTeam(t, st, &dto.AddTeamRequest{
		Name: "payments12345",
		Members: []*models.Member{
			{
				Id:       "86a832a4-a5d1-4e8c-93a2-e5bdc206d9a3",
				Username: "Alice1",
				IsActive: true,
			},
			{
				Id:       "86a832a4-a5d1-4e8c-93a2-e5bdc206d951",
				Username: "Bo1b",
				IsActive: false,
			},
		},
	})

	require.Equal(t, 201, code)

	// 5. Проверяем, что у пользователя нет PR'ов
	req = httptest.NewRequestWithContext(t.Context(), "GET", "/users/getReview?user_id=86a832a4-a5d1-4e8c-93a2-e5bdc206d951", nil)
	recorder = httptest.NewRecorder()
	st.srv.TestReq(req, recorder)
	require.Equal(t, 200, recorder.Result().StatusCode)
	var res4 dto.GetReviewResponse
	json.Unmarshal(recorder.Body.Bytes(), &res4)

	require.Equal(t, res4.UserId, "86a832a4-a5d1-4e8c-93a2-e5bdc206d951")
	require.Len(t, res4.PullRequests, 0)
}

func TestDeactivateActivateUser(t *testing.T) {
	st := NewSuite()
	st.Start()
	t.Cleanup(func() {
		st.srv.Stop()
	})

	crTeamReq := &dto.AddTeamRequest{
		Name: "team-deactivate-activate-user",
		Members: []*models.Member{
			{
				Id:       "81c8fcaf-e519-4250-8e03-b6b22274c332",
				Username: gofakeit.City(),
				IsActive: true,
			},
		},
	}
	_, code := createTeam(t, st, crTeamReq)
	require.Equal(t, 201, code)

	boolVar := false
	_, code = setIsActive(t, st, &dto.SetIsActiveRequest{
		UserId:   "81c8fcaf-e519-4250-8e03-b6b22274c332",
		IsActive: &boolVar,
	})
	require.Equal(t, 200, code)

	res, code := getTeam(t, st, crTeamReq.Name)
	require.Equal(t, 200, code)
	require.Len(t, res.Members, 1)
	require.Equal(t, res.Members[0].IsActive, false)

	boolVar = true
	_, code = setIsActive(t, st, &dto.SetIsActiveRequest{
		UserId:   "81c8fcaf-e519-4250-8e03-b6b22274c332",
		IsActive: &boolVar,
	})
	require.Equal(t, 200, code)

	res, code = getTeam(t, st, crTeamReq.Name)
	require.Equal(t, 200, code)
	require.Len(t, res.Members, 1)
	require.Equal(t, res.Members[0].IsActive, true)
}

func setIsActive(t *testing.T, st *Suite, reqBody *dto.SetIsActiveRequest) (*dto.SetIsActiveResponse, int) {
	body, _ := json.Marshal(reqBody)
	httpReq := httptest.NewRequestWithContext(t.Context(), "POST", "/users/setIsActive", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	st.srv.TestReq(httpReq, recorder)

	var res dto.SetIsActiveResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &res)
	require.NoError(t, err)

	return &res, recorder.Result().StatusCode
}
