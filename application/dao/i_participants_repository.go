package dao

import "darbelis.eu/stabas/entities"

type IParticipantsRepository interface {
	GetParticipants() []*entities.Participant
	FindParticipant(id int) (*entities.Participant, error)
	AddParticipant(participant *entities.Participant) (int, error)
	RemoveParticipant(id int) error
	UpdateParticipant(participant *entities.Participant) error
}
