package api

type TasksFilter struct {
	SortByTime       bool `json:"sort_by_time"`
	SortByStatusTime bool `json:"sort_by_status_time"`
	SortByStatus     bool `json:"sort_by_status"`
	SortByName       bool `json:"sort_by_name"`

	FilterBySender      int    `json:"sort_by_sender"`
	FilterByReceiver    int    `json:"sort_by_receiver"`
	FilterByMessagePart string `json:"sort_by_message_part"`
}
