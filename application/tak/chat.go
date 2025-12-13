package tak

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateChatCoTXML creates a chat message in CoT XML format
func CreateChatCoTXML(senderUID, senderCallsign, messageText, chatroom string) (string, error) {
	now := time.Now().UTC()
	stale := now.Add(60 * time.Second)
	messageID := uuid.New().String()

	event := Event{
		Version: "2.0",
		UID:     fmt.Sprintf("GeoChat.%s", senderUID),
		Type:    "b-t-f", // broadcast-to-friend (chat message)
		Time:    now.Format("2006-01-02T15:04:05.000Z"),
		Start:   now.Format("2006-01-02T15:04:05.000Z"),
		Stale:   stale.Format("2006-01-02T15:04:05.000Z"),
		How:     "h-g-i-g-o", // human-generated, internally computed, geochat, other
		Point: Point{
			Lat: "0.0",
			Lon: "0.0",
			Hae: "0.0",
			Ce:  "9999999.0",
			Le:  "9999999.0",
		},
		Detail: Detail{
			Chat: Chat{
				ID:             chatroom,
				Chatroom:       chatroom,
				GroupOwner:     "false",
				MessageID:      messageID,
				Parent:         "RootContactGroup",
				SenderCallsign: senderCallsign,
				ChatGroup: ChatGroup{
					UID0: senderUID,
					UID1: chatroom,
					ID:   chatroom,
				},
			},
			Link: Link{
				UID:            senderUID,
				ProductionTime: now.Format("2006-01-02T15:04:05.000Z"),
				Type:           "a-f-G-E-V-A", // friend, ground, equipment, vehicle, android
				ParentCallsign: senderCallsign,
				Relation:       "p-p", // peer-to-peer
			},
			Remarks: Remarks{
				Source: fmt.Sprintf("BAO.F.ATAK.%s", senderCallsign),
				To:     chatroom,
				Time:   now.Format("2006-01-02T15:04:05.000Z"),
				Text:   messageText,
			},
		},
	}

	output, err := xml.MarshalIndent(event, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal XML: %w", err)
	}

	return xml.Header + string(output), nil
}
