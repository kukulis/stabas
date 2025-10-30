package entities

import (
	"errors"
	"strconv"
)

const STATUS_NEW = 1
const STATUS_SENT = 2
const STATUS_RECEIVED = 3
const STATUS_EXECUTING = 4
const STATUS_FINISHED = 5
const STATUS_CLOSED = 6

var ALL_STATUSES = []int{STATUS_NEW, STATUS_SENT, STATUS_RECEIVED, STATUS_EXECUTING, STATUS_FINISHED, STATUS_CLOSED}

type Status struct{}

func ValidateStatus(status int) error {
	contains := false
	for _, st := range ALL_STATUSES {
		if st == status {
			contains = true
			break
		}
	}

	if !contains {
		return errors.New("Invalid status " + strconv.Itoa(status))
	}

	return nil
}
