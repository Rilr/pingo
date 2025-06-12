package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"pingo/static"
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

type Ticket struct {
	ID         int    `json:"id"`
	Summary    string `json:"summary"`
	RecordType string `json:"recordType"`
	Board      struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Info struct {
			BoardHref string `json:"board_href"`
		} `json:"_info"`
	} `json:"board"`
	Status struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Sort int    `json:"Sort"`
		Info struct {
			StatusHref string `json:"status_href"`
		} `json:"_info"`
	} `json:"status"`
	Company struct {
		ID         int    `json:"id"`
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
		Info       struct {
			CompanyHref string `json:"company_href"`
			MobileGUID  string `json:"mobileGuid"`
		} `json:"_info"`
	} `json:"company"`
	Site struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Info struct {
			SiteHref   string `json:"site_href"`
			MobileGUID string `json:"mobileGuid"`
		} `json:"_info"`
	} `json:"site"`
	SiteName        string `json:"siteName"`
	AddressLine1    string `json:"addressLine1"`
	AddressLine2    string `json:"addressLine2"`
	City            string `json:"city"`
	StateIdentifier string `json:"stateIdentifier"`
	Zip             string `json:"zip"`
	Country         struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Info struct {
			CountryHref string `json:"country_href"`
		} `json:"_info"`
	} `json:"country"`
	ContactName string `json:"contactName"`
	Type        struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Info struct {
			TypeHref string `json:"type_href"`
		} `json:"_info"`
	} `json:"type"`
	SubType struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Info struct {
			SubTypeHref string `json:"subType_href"`
		} `json:"_info"`
	} `json:"subType"`
	Team struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Info struct {
			TeamHref string `json:"team_href"`
		} `json:"_info"`
	} `json:"team"`
	Priority struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Sort  int    `json:"sort"`
		Level string `json:"level"`
		Info  struct {
			PriorityHref string `json:"priority_href"`
			ImageHref    string `json:"image_href"`
		} `json:"_info"`
	} `json:"priority"`
	ServiceLocation struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Info struct {
			LocationHref string `json:"location_href"`
		} `json:"_info"`
	} `json:"serviceLocation"`
	Source struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Info struct {
			SourceHref string `json:"source_href"`
		} `json:"_info"`
	} `json:"source"`
	Severity                   string    `json:"severity"`
	Impact                     string    `json:"impact"`
	AllowAllClientsPortalView  bool      `json:"allowAllClientsPortalView"`
	CustomerUpdatedFlag        bool      `json:"customerUpdatedFlag"`
	AutomaticEmailContactFlag  bool      `json:"automaticEmailContactFlag"`
	AutomaticEmailResourceFlag bool      `json:"automaticEmailResourceFlag"`
	AutomaticEmailCcFlag       bool      `json:"automaticEmailCcFlag"`
	ClosedFlag                 bool      `json:"closedFlag"`
	Approved                   bool      `json:"approved"`
	EstimatedExpenseCost       float64   `json:"estimatedExpenseCost"`
	EstimatedExpenseRevenue    float64   `json:"estimatedExpenseRevenue"`
	EstimatedProductCost       float64   `json:"estimatedProductCost"`
	EstimatedProductRevenue    float64   `json:"estimatedProductRevenue"`
	EstimatedTimeCost          float64   `json:"estimatedTimeCost"`
	EstimatedTimeRevenue       float64   `json:"estimatedTimeRevenue"`
	BillingMethod              string    `json:"billingMethod"`
	SubBillingMethod           string    `json:"subBillingMethod"`
	DateResplan                time.Time `json:"dateResplan"`
	DateResponded              time.Time `json:"dateResponded"`
	ResolveMinutes             int       `json:"resolveMinutes"`
	ResPlanMinutes             int       `json:"resPlanMinutes"`
	RespondMinutes             int       `json:"respondMinutes"`
	IsInSLA                    bool      `json:"isInSla"`
	HasChildTicket             bool      `json:"hasChildTicket"`
	HasMergedChildTicketFlag   bool      `json:"hasMergedChildTicketFlag"`
	BillTime                   string    `json:"billTime"`
	BillExpenses               string    `json:"billExpenses"`
	BillProducts               string    `json:"billProducts"`
	Location                   struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Info struct {
			LocationHref string `json:"location_href"`
		} `json:"_info"`
	} `json:"location"`
	Department struct {
		ID         int    `json:"id"`
		Identifier string `json:"identifier"`
		Name       string `json:"name"`
		Info       struct {
			DepartmentHref string `json:"department_href"`
		} `json:"_info"`
	} `json:"department"`
	MobileGUID string `json:"mobileGuid"`
	SLA        struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Info struct {
			SLAHref string `json:"sla_href"`
		} `json:"_info"`
	} `json:"sla"`
	SLAStatus            string `json:"slaStatus"`
	RequestForChangeFlag bool   `json:"requestForChangeFlag"`
	Currency             struct {
		ID                      int    `json:"id"`
		Symbol                  string `json:"symbol"`
		CurrencyCode            string `json:"currencyCode"`
		DecimalSeparator        string `json:"decimalSeparator"`
		NumberOfDecimals        int    `json:"numberOfDecimals"`
		ThousandsSeparator      string `json:"thousandsSeparator"`
		NegativeParenthesesFlag bool   `json:"negativeParenthesesFlag"`
		DisplaySymbolFlag       bool   `json:"displaySymbolFlag"`
		CurrencyIdentifier      string `json:"currencyIdentifier"`
		DisplayIDFlag           bool   `json:"displayIdFlag"`
		RightAlign              bool   `json:"rightAlign"`
		Name                    string `json:"name"`
		Info                    struct {
			CurrencyHref string `json:"currency_href"`
		} `json:"_info"`
	} `json:"currency"`
	Info struct {
		LastUpdated         time.Time `json:"lastUpdated"`
		UpdatedBy           string    `json:"updatedBy"`
		DateEntered         time.Time `json:"dateEntered"`
		EnteredBy           string    `json:"enteredBy"`
		ActivitiesHref      string    `json:"activities_href"`
		ScheduleentriesHref string    `json:"scheduleentries_href"`
		DocumentsHref       string    `json:"documents_href"`
		ConfigurationsHref  string    `json:"configurations_href"`
		TasksHref           string    `json:"tasks_href"`
		NotesHref           string    `json:"notes_href"`
		ProductsHref        string    `json:"products_href"`
		TimeentriesHref     string    `json:"timeentries_href"`
		ExpenseEntriesHref  string    `json:"expenseEntries_href"`
	} `json:"_info"`
	EscalationStartDateUTC  time.Time `json:"escalationStartDateUTC"`
	EscalationLevel         int       `json:"escalationLevel"`
	MinutesBeforeWaiting    int       `json:"minutesBeforeWaiting"`
	RespondedSkippedMinutes int       `json:"respondedSkippedMinutes"`
	ResplanSkippedMinutes   int       `json:"resplanSkippedMinutes"`
	RespondedHours          float64   `json:"respondedHours"`
	RespondedBy             string    `json:"respondedBy"`
	ResplanHours            float64   `json:"resplanHours"`
	ResplanBy               string    `json:"resplanBy"`
	ResolutionHours         float64   `json:"resolutionHours"`
	MinutesWaiting          int       `json:"minutesWaiting"`
	CustomFields            []struct {
		ID               int    `json:"id"`
		Caption          string `json:"caption"`
		Type             string `json:"type"`
		EntryMethod      string `json:"entryMethod"`
		NumberOfDecimals int    `json:"numberOfDecimals"`
		ConnectWiseID    string `json:"connectWiseId"`
	} `json:"customFields"`
}

// ManageAuth generates a base64 encoded string for authentication
func ManageAuth() string {
	companyName := static.Manage.User
	publicKey := static.Manage.PubKey
	privateKey := static.Manage.PrvKey
	combinedStr := companyName + "+" + publicKey + ":" + privateKey
	base64Str := base64.StdEncoding.EncodeToString([]byte(combinedStr))
	return base64Str
}

// PostTicketPayload generates a JSON payload for creating a new service ticket
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

// Checks if the device address is reachable before attempting to SSH into it
func InitTtyToHost() {
	if !TestAddress(devAddr, 2, 1*time.Second, 10*time.Second) {
		AddtoLog(fmt.Sprintf("Device address %s is unresponsive before attempting to SSH", devAddr))
		os.Exit(3)
	} else {
		user := static.DeviceTty.User
		cred := static.DeviceTty.Cred
		AddtoLog(fmt.Sprintf("Attempting to Tunnel into: %s", devAddr))
		if err := sshIntoHost(devAddr, user, cred, "ipsec restart"); err != nil {
			AddtoLog(fmt.Sprintf("Failed to run command on device address %s: %v", devAddr, err))
			os.Exit(4)
		} else {
			AddtoLog(fmt.Sprintf("Command ran successfully on device address %s", devAddr))
			putTicketNote(0, "Tunnel was restarted successfully.")
			os.Exit(0)
		}
	}
}