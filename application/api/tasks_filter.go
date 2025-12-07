package api

type TasksFilter struct {
	// TODO refactor by storing Sort as a string
	SortByTime       bool `json:"sort_by_time"`
	SortByStatusTime bool `json:"sort_by_status_time"`
	SortByStatus     bool `json:"sort_by_status"`
	SortByName       bool `json:"sort_by_name"`

	FilterBySender      int    `json:"filter_by_sender"`
	FilterByReceiver    int    `json:"filter_by_receiver"`
	FilterByMessagePart string `json:"filter_by_message_part"`
}
