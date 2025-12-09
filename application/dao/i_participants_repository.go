package dao

import "darbelis.eu/stabas/entities"

type IParticipantsRepository interface {
	GetParticipants() []*entities.Participant
	FindParticipant(id int) (*entities.Participant, error)
	AddParticipant(participant *entities.Participant) (*entities.Participant, error)
	RemoveParticipant(id int) error
	UpdateParticipant(participant *entities.Participant) error
	FindParticipantByName(name string) *entities.Participant
	FindParticipantByToken(token string) *entities.Participant
	UpdateParticipantToken(id int, token string) error
	UpdateParticipantPassword(id int, token string) error
}
