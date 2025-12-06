package dao

import (
	"darbelis.eu/stabas/entities"
	"darbelis.eu/stabas/util"
	"errors"
	"fmt"
)

type IParticipantsRepository interface {
	GetParticipants() []*entities.Participant
	FindParticipant(id int) (*entities.Participant, error)
	AddParticipant(participant *entities.Participant) (int, error)
	RemoveParticipant(id int) error
	UpdateParticipant(participant *entities.Participant) error
}

type ParticipantsRepository struct {
	participants []*entities.Participant
}

func NewParticipantsRepository(initialParticipants []*entities.Participant) *ParticipantsRepository {
	return &ParticipantsRepository{
		participants: initialParticipants,
	}
}

func (rep *ParticipantsRepository) GetParticipants() []*entities.Participant {
	return util.ArrayFilter(rep.participants, func(participant *entities.Participant) bool {
		return !participant.Deleted
	})
}

func (rep *ParticipantsRepository) FindParticipant(id int) (*entities.Participant, error) {
	for _, participant := range rep.participants {
		if participant.Id == id {
			return participant, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("participant with id %d not found", id))
}

func (rep *ParticipantsRepository) getMaxId() int {
	return util.ArrayReduce(rep.participants, 0, func(maxId int, participant *entities.Participant) int {
		return util.MaxInt(maxId, participant.Id)
	})
}

func (rep *ParticipantsRepository) AddParticipant(participant *entities.Participant) (int, error) {
	if participant.Id == 0 {
		participant.Id = rep.getMaxId() + 1
	}

	rep.participants = append(rep.participants, participant)

	return participant.Id, nil
}

func (rep *ParticipantsRepository) RemoveParticipant(id int) error {
	rep.participants = util.ArrayFilter(rep.participants, func(participant *entities.Participant) bool { return participant.Id != id })

	return nil
}

func (rep *ParticipantsRepository) UpdateParticipant(participant *entities.Participant) error {
	existingParticipant, err := rep.FindParticipant(participant.Id)

	if err != nil {
		return err
	}

	existingParticipant.Name = participant.Name

	return nil
}
