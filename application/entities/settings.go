package entities

type Settings struct {
	Id int `json:"id"`

	// newStatusDelay = 5;
	NewStatusDelay int `json:"new_status_delay"` // what is variable type??
	// newStatusDelaySevere = 15;
	NewStatusDelaySevere int `json:"new_status_delay_severe"`
	// sentStatusDelay = 2;
	SentStatusDelay int `json:"sent_status_delay"`
	// sentStatusDelaySevere = 6;
	SentStatusDelaySevere int `json:"sent_status_delay_severe"`
	// receivedStatusDelay = 5;
	ReceivedStatusDelay int `json:"received_status_delay"`
	// receivedStatusDelaySevere = 15;
	ReceivedStatusDelaySevere int `json:"received_status_delay_severe"`
	// executingStatusDelay = 10;
	ExecutingStatusDelay int `json:"executing_status_delay"`
	// executingStatusDelaySevere = 20;
	ExecutingStatusDelaySevere int `json:"executing_status_delay_severe"`
	// finishedStatusDelay = 60;
	FinishedStatusDelay int `json:"finished_status_delay"`
	// finishedStatusDelaySevere = 120;
	FinishedStatusDelaySevere int `json:"finished_status_delay_severe"`
}
