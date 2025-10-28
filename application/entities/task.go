package entities

import "time"

type Task struct {
	Id          int        `json:"id"`
	Message     string     `json:"message"`
	Result      string     `json:"result"`
	Sender      int        `json:"sender"`
	Receivers   []int      `json:"receivers"`
	Status      int        `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	SentAt      *time.Time `json:"sent_at"`
	ReceivedAt  *time.Time `json:"received_at"`
	ExecutingAt *time.Time `json:"executing_at"`
	FinishedAt  *time.Time `json:"finished_at"`
	ClosedAt    *time.Time `json:"closed_at"`
}

func NewTask() *Task {
	now := time.Now()
	return &Task{
		Id:          0,
		Status:      STATUS_NEW,
		CreatedAt:   &now,
		SentAt:      nil,
		ReceivedAt:  nil,
		ExecutingAt: nil,
		FinishedAt:  nil,
		ClosedAt:    nil,
	}
}
