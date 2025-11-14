package usecases

import (
	"context"
	"errors"
	"log/slog"
	"pr-review/internal/http/dto"
	"pr-review/internal/models"
	"pr-review/internal/repository/postgres"
	"pr-review/internal/utils"

	"github.com/google/uuid"
)

var (
	ErrTeamNotFound     = errors.New("team not found")
	ErrTeamAlredyExists = errors.New("team already exists")
)

func (uc *Usecases) GetTeam(ctx context.Context, name string) ([]*models.User, error) {
	const op = "usecases.GetTeam"
	log := uc.log.With(slog.String("op", op), slog.String("name", name))

	exists, err := uc.db.TeamExists(ctx, name)
	if err != nil {
		log.Error("error getting team existing", slog.String("error", err.Error()))
		return nil, err
	}
	if !exists {
		log.Warn("team not found")
		return nil, ErrTeamNotFound
	}

	members, err := uc.db.GetTeamMembers(ctx, name)
	if err != nil {
		log.Error("error getting team members", slog.String("error", err.Error()))
		return nil, err
	}

	log.Debug("successfully got team members")
	return members, nil
}

// CreateTeam создаёт команду по имени, добавляет в неё участников/изменяет данные существующих пользователей по id
// При изменении isActive у участника на false, все PR'ы снимаются с него и распределяются среди участников его команды
// Если участников не хватает до двух PR'ов, то выставляется флаг need_more_reviewers
// Для того, чтобы создание команды и добавление участников объединить в единую атомарную операцию, используется транзакция
func (uc *Usecases) CreateTeam(ctx context.Context, reqDTO *dto.AddTeamRequest) error {
	const op = "usecases.CreateTeam"
	log := uc.log.With(slog.String("op", op), slog.String("name", reqDTO.Name))

	teamId := uuid.NewString()
	team := &models.Team{
		Id:   teamId,
		Name: reqDTO.Name,
	}

	tx, err := uc.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			txErr := tx.Rollback(ctx)
			if txErr != nil {
				log.Error("error rolling back transaction", slog.String("error", txErr.Error()))
			}
		} else {
			txErr := tx.Commit(ctx)
			if txErr != nil {
				log.Error("error commiting transaction", slog.String("error", txErr.Error()))
			} else {
				log.Debug("successfully created team")
			}
		}
	}()

	// создаём команду
	err = uc.db.CreateTeam(ctx, tx, team)
	if err != nil {
		if errors.Is(err, postgres.ErrTeamAlredyExists) {
			log.Warn("team already exists")
			return ErrTeamAlredyExists
		}
		log.Error("error creating team", slog.String("error", err.Error()))
		return err
	}

	// добавляем/обновляем участников
	err = uc.db.AddOrUpdateTeamMembers(ctx, tx, teamId, reqDTO.Members)
	if err != nil {
		log.Error("error adding team members", slog.String("error", err.Error()))
		return err
	}

	// переассайниваем PR'ы с каждого неактивного участника
	for _, member := range reqDTO.Members {
		if member.IsActive {
			continue
		}

		prIds, err := uc.db.UnassignPRsFromUser(ctx, tx, member.Id)
		if err != nil {
			log.Error("error updating user team", slog.String("error", err.Error()))
			return err
		}

		for _, prId := range prIds {
			teammates, err := uc.db.GetTeammates(ctx, tx, prId, teamId)
			if err != nil {
				log.Error("error updating user team", slog.String("error", err.Error()))
				return err
			}

			// перемешаем teammates, чтобы назначать assignee в случайном порядке
			utils.Shuffle(teammates)

			newAssignerId, err := uc.db.AssignPRToUser(ctx, tx, prId, teammates)
			if err != nil {
				log.Error("error updating user team", slog.String("error", err.Error()))
				return err
			}

			if newAssignerId == "" { // новый аппрувер для PR'а не был найден, выставляем need_more_reviewers = true
				err = uc.db.SetNeedMoreReviewers(ctx, tx, prId)
				if err != nil {
					log.Error("error updating user team", slog.String("error", err.Error()))
					return err
				}
			}
		}
	}

	return nil
}
