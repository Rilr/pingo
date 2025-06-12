package main

import (
	"encoding/json"
	"fmt"
)

type ContactRef struct {
	ID int `json:"id"`
}
type BoardRef struct {
	ID int `json:"id"`
}
type StatusRef struct {
	ID int `json:"id"`
}
type CompanyRef struct {
	ID int `json:"id"`
}
type TypeRef struct {
	ID int `json:"id"`
}
type SubTypeRef struct {
	ID int `json:"id"`
}
type ItemRef struct {
	ID int `json:"id"`
}
type PriorityRef struct {
	ID int `json:"id"`
}

type PostTicket struct {
	Summary    string      `json:"summary"`
	RecordType string      `json:"recordType"`
	Contact    ContactRef  `json:"contact"`
	Board      BoardRef    `json:"board"`
	Status     StatusRef   `json:"status"`
	Company    CompanyRef  `json:"company"`
	Type       TypeRef     `json:"type"`
	SubType    SubTypeRef  `json:"subType"`
	Item       ItemRef     `json:"item"`
	Priority   PriorityRef `json:"priority"`
}

func PostTicketPayload() []byte {
	var staticTicket = PostTicket{
		RecordType: "ServiceTicket",
		Contact:    ContactRef{ID: 1694}, // Bryan Pomares
		Board:      BoardRef{ID: 1},      // Help Desk
		Status:     StatusRef{ID: 579},   // Review by Dispatch
		Type:       TypeRef{ID: 193},     // Break/Fix
		SubType:    SubTypeRef{ID: 7},    // Network
		Item:       ItemRef{ID: 57},      // Failure
		Priority:   PriorityRef{ID: 6},   // Critical
	}

	payload := staticTicket
	payload.Summary = "SCRIPT TICKET - TCT VPN Tunnel Down"
	payload.Company = CompanyRef{ID: 19786} // TCT

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil
	}
	return jsonData
}