package entities

import (
	"errors"
	"strconv"
	"time"
)

type Task struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
	Result  string `json:"result"`
	Sender  int    `json:"sender"`

	Receivers   []int      `json:"receivers"`
	Status      int        `json:"status"`
	CreatedAt   *time.Time `json:"created_at"`
	SentAt      *time.Time `json:"sent_at"`
	ReceivedAt  *time.Time `json:"received_at"`
	ExecutingAt *time.Time `json:"executing_at"`
	FinishedAt  *time.Time `json:"finished_at"`
	ClosedAt    *time.Time `json:"closed_at"`
	Deleted     bool       `json:"deleted"`
	Version     int        `json:"version"`
	TaskGroup   int        `json:"task_group"`
	Children    []*Task    `json:"children"`
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
		Deleted:     false,
		Version:     1,
	}
}

func (task *Task) SetStatusDate(date time.Time) error {
	if task.Status == STATUS_NEW {
		task.CreatedAt = &date
		return nil
	}

	if task.Status == STATUS_SENT {
		task.SentAt = &date
		return nil
	}

	if task.Status == STATUS_RECEIVED {
		task.ReceivedAt = &date
		return nil
	}

	if task.Status == STATUS_EXECUTING {
		task.ExecutingAt = &date
		return nil
	}

	if task.Status == STATUS_FINISHED {
		task.FinishedAt = &date
		return nil
	}

	if task.Status == STATUS_CLOSED {
		task.ClosedAt = &date
		return nil
	}

	return errors.New("invalid status " + strconv.Itoa(task.Status))

}

// SetStatusDateIfNil sets status date in case it is not set yet
func (task *Task) SetStatusDateIfNil(date time.Time) error {
	if task.Status == STATUS_NEW {
		if task.CreatedAt == nil {
			task.CreatedAt = &date
		}
		return nil
	}

	if task.Status == STATUS_SENT {
		if task.SentAt == nil {
			task.SentAt = &date
		}
		return nil
	}

	if task.Status == STATUS_RECEIVED {
		if task.ReceivedAt == nil {
			task.ReceivedAt = &date
		}
		return nil
	}

	if task.Status == STATUS_EXECUTING {
		if task.ExecutingAt == nil {
			task.ExecutingAt = &date
		}
		return nil
	}

	if task.Status == STATUS_FINISHED {
		if task.FinishedAt == nil {
			task.FinishedAt = &date
		}
		return nil
	}

	if task.Status == STATUS_CLOSED {
		if task.ClosedAt == nil {
			task.ClosedAt = &date
		}
		return nil
	}

	return errors.New("invalid status " + strconv.Itoa(task.Status))
}

func (task *Task) GetStatusTime() *time.Time {
	if task.Status == STATUS_NEW {
		return task.CreatedAt
	}

	if task.Status == STATUS_SENT {
		return task.SentAt
	}

	if task.Status == STATUS_RECEIVED {
		return task.ReceivedAt
	}

	if task.Status == STATUS_EXECUTING {
		return task.ExecutingAt
	}

	if task.Status == STATUS_FINISHED {
		return task.FinishedAt
	}

	if task.Status == STATUS_CLOSED {
		return task.ClosedAt
	}

	return nil
}

// hasSenderOrReceiver checks if the given participantId is either the sender
// or the first receiver of the task
func (task *Task) HasSenderOrReceiver(participantId int) bool {
	if task.Sender == participantId {
		return true
	}
	//if task.Receivers == nil {
	//	return false
	//}
	if len(task.Receivers) > 0 && task.Receivers[0] == participantId {
		return true
	}
	return false
}

func (task *Task) HasSender(participantId int) bool {
	if task.Sender == participantId {
		return true
	}
	return false
}
