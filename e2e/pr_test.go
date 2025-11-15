package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"
	"pr-review/internal/usecases"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReassignPR проверяет переназначение PR с одного ревьювера на другого
func TestReassignPR(t *testing.T) {
	st := NewSuite()
	st.Start()
	t.Cleanup(func() {
		st.srv.Stop()
	})

	crTeamReq := &dto.AddTeamRequest{
		Name: "payments-reassign",
		Members: []*models.Member{
			{
				Id:       "86a832a4-a5d1-4e8c-9000-e5bdc206d9a3",
				Username: gofakeit.Name() + uuid.NewString(),
				IsActive: true,
			},
			{
				Id:       "86a832a4-a5d1-4e8c-9002-e5bdc206d951",
				Username: gofakeit.Name() + uuid.NewString(),
				IsActive: true,
			},
			{
				Id:       "86a83214-a5d1-4e8c-90a2-e5bdc206d951",
				Username: gofakeit.Name() + uuid.NewString(),
				IsActive: true,
			},
			{
				Id:       "81a83214-a5d1-4e8c-90a2-e5bdc206d951",
				Username: gofakeit.Name() + uuid.NewString(),
				IsActive: true,
			},
		},
	}
	_, code := createTeam(t, st, crTeamReq)
	require.Equal(t, 201, code)

	// Создаём PR
	response, code, prId, _ := createPR(t, st, "86a83214-a5d1-4e8c-90a2-e5bdc206d951")
	require.Equal(t, 201, code)
	require.Len(t, response.PR.Reviewers, 2)

	// Переназначаем PR с одного ревьювера на другого
	reassignReq := &dto.ReassignPRRequest{
		PullRequestID: prId,
		OldReviewerID: response.PR.Reviewers[0],
	}

	bytes, code := reassignPR(t, st, reassignReq)
	var res dto.ReassignPRResponse
	err := json.Unmarshal(bytes, &res)
	require.NoError(t, err)
	require.Equal(t, 200, code)
	require.Len(t, res.PR.Reviewers, 2)
	for _, reviewer := range res.PR.Reviewers {
		require.NotEqual(t, reassignReq.OldReviewerID, reviewer)
	}
}

// TestReassignPRNoCandidates проверяет переназначение PR с одного ревьювера на другого в случае, если отсутствуют кандидаты для назначения
func TestReassignPRNoCandidates(t *testing.T) {
	st := NewSuite()
	st.Start()
	t.Cleanup(func() {
		st.srv.Stop()
	})

	crTeamReq := &dto.AddTeamRequest{
		Name: gofakeit.Name() + uuid.NewString(),
		Members: []*models.Member{
			{
				Id:       "76a832a4-a5d1-4e8c-9000-e5bdc206d9a3",
				Username: gofakeit.Name() + uuid.NewString(),
				IsActive: true,
			},
			{
				Id:       "76a832a4-a5d1-4e8c-9002-e5bdc206d951",
				Username: gofakeit.Name() + uuid.NewString(),
				IsActive: true,
			},
			{
				Id:       "76a83214-a5d1-4e8c-90a2-e5bdc206d951",
				Username: gofakeit.Name() + uuid.NewString(),
				IsActive: true,
			},
		},
	}
	data, code := createTeam(t, st, crTeamReq)
	var crTeamRes dto.AddTeamResponse
	err := json.Unmarshal(data, &crTeamRes)
	require.NoError(t, err)
	require.Equal(t, 201, code)

	// Создаём PR
	response, code, prId, _ := createPR(t, st, "76a83214-a5d1-4e8c-90a2-e5bdc206d951")
	require.Equal(t, 201, code)
	require.Len(t, response.PR.Reviewers, 2)

	// Переназначаем PR с одного ревьювера на другого
	reassignReq := &dto.ReassignPRRequest{
		PullRequestID: prId,
		OldReviewerID: response.PR.Reviewers[0],
	}

	bytes, code := reassignPR(t, st, reassignReq)
	var res dto.ErrorResponse
	err = json.Unmarshal(bytes, &res)
	require.NoError(t, err)
	require.Equal(t, 409, code)
	require.Equal(t, dto.ErrCodeNoCandidates, res.Error.Code)
	require.Equal(t, usecases.ErrNoCandidatesToAssign.Error(), res.Error.Message)

	// проверка на то, что ревьювер с OldReviewerID не является ревьювером PR
	data, code = reassignPR(t, st, &dto.ReassignPRRequest{
		PullRequestID: prId,
		OldReviewerID: "76a83214-a5d1-4e8c-90a2-e5bdc206d951",
	})
	var res2 dto.ErrorResponse
	err = json.Unmarshal(data, &res2)
	require.NoError(t, err)
	require.Equal(t, 409, code)
	require.Equal(t, dto.ErrCodeUserNotReviewerOfPR, res2.Error.Code)
	require.Equal(t, usecases.ErrUserNotReviewerOfPR.Error(), res2.Error.Message)
}

// TestMergePR проверяет успешное смёрдживание PR после ревью одним человеком
func TestMergePR(t *testing.T) {
	st := NewSuite()
	st.Start()
	t.Cleanup(func() {
		st.srv.Stop()
	})

	crTeamReq := &dto.AddTeamRequest{
		Name: gofakeit.Name() + uuid.NewString(),
		Members: []*models.Member{
			{
				Id:       "96a832a4-a5d1-4e8c-1000-e5bdc206d9a3",
				Username: gofakeit.Name() + uuid.NewString(),
				IsActive: true,
			},
			{
				Id:       "96a832a4-a5d1-4e8c-1002-e5bdc206d951",
				Username: gofakeit.Name() + uuid.NewString(),
				IsActive: true,
			},
		},
	}
	_, code := createTeam(t, st, crTeamReq)
	require.Equal(t, 201, code)

	// Создаём PR
	response, code, prId, _ := createPR(t, st, "96a832a4-a5d1-4e8c-1002-e5bdc206d951")
	require.Equal(t, 201, code)
	require.Len(t, response.PR.Reviewers, 1)

	// Смёрдживаем PR
	mergeReq := &dto.MergePRRequest{
		PullRequestID: prId,
	}

	data, code := mergePR(t, st, mergeReq)
	require.Equal(t, 200, code)

	var res1 dto.MergePRResponse
	err := json.Unmarshal(data, &res1)
	require.NoError(t, err)
	require.Equal(t, 200, code)
	require.Equal(t, res1.PR.Status, models.StatusMerged)

	// проверка идемпотентности
	data, code = mergePR(t, st, mergeReq)
	require.Equal(t, 200, code)

	var res2 dto.MergePRResponse
	err = json.Unmarshal(data, &res2)
	require.NoError(t, err)
	require.Equal(t, 200, code)
	require.Equal(t, res2.PR.Status, models.StatusMerged)
	require.Equal(t, res1.PR.MergedAt, res2.PR.MergedAt)

	// проверка, что MERGED нельзя переназначить
	reassignReq := &dto.ReassignPRRequest{
		PullRequestID: prId,
		OldReviewerID: response.PR.Reviewers[0],
	}

	bytes, code := reassignPR(t, st, reassignReq)
	var res dto.ErrorResponse
	err = json.Unmarshal(bytes, &res)
	require.NoError(t, err)
	require.Equal(t, 409, code)
	require.Equal(t, dto.ErrCodeCannotReassignMergedPR, res.Error.Code)
	require.Equal(t, usecases.ErrPRMerged.Error(), res.Error.Message)
}

func TestStatistics(t *testing.T) {
	st := NewSuite()
	st.Start()
	t.Cleanup(func() {
		st.srv.Stop()
	})
	err := st.db.DeleteUsers(t.Context())
	assert.NoError(t, err)

	members := []*models.Member{
		{
			Id:       "16a832a4-a5d1-4e8c-9000-e5bdc206d9a3",
			Username: gofakeit.Name() + uuid.NewString(),
			IsActive: true,
		},
		{
			Id:       "26a832a4-a5d1-4e8c-9002-e5bdc206d951",
			Username: gofakeit.Name() + uuid.NewString(),
			IsActive: true,
		},
		{
			Id:       "36a83214-a5d1-4e8c-90a2-e5bdc206d951",
			Username: gofakeit.Name() + uuid.NewString(),
			IsActive: true,
		},
	}
	createTeam(t, st, &dto.AddTeamRequest{
		Name:    gofakeit.Name() + uuid.NewString(),
		Members: members,
	})

	for _, member := range members {
		for i := 0; i < 3; i++ {
			response, code, prId, _ := createPR(t, st, member.Id)
			require.Equal(t, 201, code)
			require.Len(t, response.PR.Reviewers, 2)

			mergeReq := &dto.MergePRRequest{
				PullRequestID: prId,
			}
			data, code := mergePR(t, st, mergeReq)
			require.Equal(t, 200, code)

			var res dto.MergePRResponse
			err := json.Unmarshal(data, &res)
			require.NoError(t, err)
			require.Equal(t, 200, code)
			require.Equal(t, res.PR.Status, models.StatusMerged)
		}
	}

	body, code := statistics(t, st)
	var resp dto.StatisticsResponse
	err = json.Unmarshal(body, &resp)
	require.NoError(t, err)
	require.Equal(t, 200, code)
	require.Equal(t, uint64(3), resp.Count)

	s := 0
	for _, v := range resp.Statistics {
		s += v
	}
	require.Equal(t, 9, s)
}

func createPR(t *testing.T, st *Suite, authorId string) (*dto.CreatePRResponse, int, string, string) {
	id := uuid.NewString()
	name := gofakeit.City()
	body := []byte(fmt.Sprintf(`{
		"pull_request_id": "%s",
		"pull_request_name": "%s",
		"author_id": "%s"
	}`, id, name, authorId))
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/pullRequest/create", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	st.srv.TestReq(req, recorder)
	var resp dto.CreatePRResponse
	err := json.Unmarshal(recorder.Body.Bytes(), &resp)
	require.NoError(t, err)

	return &resp, recorder.Result().StatusCode, id, name
}

func reassignPR(t *testing.T, st *Suite, reqBody *dto.ReassignPRRequest) ([]byte, int) {
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/pullRequest/reassign", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	st.srv.TestReq(req, recorder)

	var res dto.ReassignPRResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &res)
	require.NoError(t, err)

	return recorder.Body.Bytes(), recorder.Result().StatusCode
}

func mergePR(t *testing.T, st *Suite, reqBody *dto.MergePRRequest) ([]byte, int) {
	body, err := json.Marshal(reqBody)
	require.NoError(t, err)
	req := httptest.NewRequestWithContext(t.Context(), "POST", "/pullRequest/merge", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	st.srv.TestReq(req, recorder)

	return recorder.Body.Bytes(), recorder.Result().StatusCode
}

func statistics(t *testing.T, st *Suite) ([]byte, int) {
	req := httptest.NewRequestWithContext(t.Context(), "GET", "/pullRequest/statistics", nil)
	recorder := httptest.NewRecorder()

	st.srv.TestReq(req, recorder)

	return recorder.Body.Bytes(), recorder.Result().StatusCode
}
