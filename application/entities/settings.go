package entities

type Settings struct {
	Id int `json:"id"`

	// Use camelCase for json and PascalCase is used for go in this case because of the access modifiers.
	NewStatusDelay             int `json:"newStatusDelay"`
	NewStatusDelaySevere       int `json:"newStatusDelaySevere"`
	SentStatusDelay            int `json:"sentStatusDelay"`
	SentStatusDelaySevere      int `json:"sentStatusDelaySevere"`
	ReceivedStatusDelay        int `json:"receivedStatusDelay"`
	ReceivedStatusDelaySevere  int `json:"receivedStatusDelaySevere"`
	ExecutingStatusDelay       int `json:"executingStatusDelay"`
	ExecutingStatusDelaySevere int `json:"executingStatusDelaySevere"`
	FinishedStatusDelay        int `json:"finishedStatusDelay"`
	FinishedStatusDelaySevere  int `json:"finishedStatusDelaySevere"`
}

func NewSettings() *Settings {
	return &Settings{
		Id:                         0,
		NewStatusDelay:             5,
		NewStatusDelaySevere:       15,
		SentStatusDelay:            2,
		SentStatusDelaySevere:      6,
		ReceivedStatusDelay:        5,
		ReceivedStatusDelaySevere:  15,
		ExecutingStatusDelay:       10,
		ExecutingStatusDelaySevere: 20,
		FinishedStatusDelay:        60,
		FinishedStatusDelaySevere:  120,
	}
}
