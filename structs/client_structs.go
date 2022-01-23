package model

import "time"

type Client struct {
	Id        uint64 `json:"clientId"`
	Firstname string `json:"firstName"`
	Lastname  string `json:"lastName"`
	Username  string `json:"username"`
	Active    bool   `json:"active"`
}

type Profile struct {
	Id                uint64      `json:"profileId"`
	DateOfBirth       time.Time   `json:"dateOfBirth"`
	ProfilePicture    string      `json:"profilePicture"` // base64
	Contact           ContactInfo `json:"contact"`
	PhysicalAddress   Address     `json:"physicalAddress"`
	MailingAddress    Address     `json:"mailingAddress"`
	AdditionalDetails string      `json:"additionalDetails"`
	ClientId          uint64      `json:"clientId"`
}

type ContactInfo struct {
	PrimaryEmail         string      `json:"primaryEmail"`
	SecondaryEmail       string      `json:"secondaryEmail"`
	PrimaryPhoneNumber   PhoneNumber `json:"primaryPhoneNumber"`
	SecondaryPhoneNumber PhoneNumber `json:"secondaryPhoneNumber"`
}

type PhoneNumber struct {
	AreaCode          string `json:"areaCode"`
	Number            string `json:"number"`
	ContactByText     bool   `json:"contactByText"`
	ContactByVoice    bool   `json:"contactByVoice"`
	ContactByWhatsApp bool   `json:"contactByWhatsApp"`
}

type Address struct {
	AddressLineOne string `json:"addressLineOne"`
	AddressLineTwo string `json:"addressLineTwo"`
	City           string `json:"city"`
	State          string `json:"state"`
	Country        string `json:"country"`
	ClientId       uint64 `json:"clientId"`
}
