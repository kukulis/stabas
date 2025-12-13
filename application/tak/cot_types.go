package tak

import "encoding/xml"

// CoT XML Structure for Chat Messages
type Event struct {
	XMLName xml.Name `xml:"event"`
	Version string   `xml:"version,attr"`
	UID     string   `xml:"uid,attr"`
	Type    string   `xml:"type,attr"`
	Time    string   `xml:"time,attr"`
	Start   string   `xml:"start,attr"`
	Stale   string   `xml:"stale,attr"`
	How     string   `xml:"how,attr"`
	Point   Point    `xml:"point"`
	Detail  Detail   `xml:"detail"`
}

type Point struct {
	Lat string `xml:"lat,attr"`
	Lon string `xml:"lon,attr"`
	Hae string `xml:"hae,attr"`
	Ce  string `xml:"ce,attr"`
	Le  string `xml:"le,attr"`
}

type Detail struct {
	Chat    Chat    `xml:"__chat"`
	Link    Link    `xml:"link"`
	Remarks Remarks `xml:"remarks"`
}

type Chat struct {
	ID             string    `xml:"id,attr"`
	Chatroom       string    `xml:"chatroom,attr"`
	GroupOwner     string    `xml:"groupOwner,attr"`
	MessageID      string    `xml:"messageId,attr"`
	Parent         string    `xml:"parent,attr"`
	SenderCallsign string    `xml:"senderCallsign,attr"`
	ChatGroup      ChatGroup `xml:"chatgrp"`
}

type ChatGroup struct {
	UID0 string `xml:"uid0,attr"`
	UID1 string `xml:"uid1,attr"`
	ID   string `xml:"id,attr"`
}

type Link struct {
	UID            string `xml:"uid,attr"`
	ProductionTime string `xml:"production_time,attr"`
	Type           string `xml:"type,attr"`
	ParentCallsign string `xml:"parent_callsign,attr"`
	Relation       string `xml:"relation,attr"`
}

type Remarks struct {
	Source string `xml:"source,attr"`
	To     string `xml:"to,attr"`
	Time   string `xml:"time,attr"`
	Text   string `xml:",chardata"`
}
