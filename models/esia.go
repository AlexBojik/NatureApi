package models

type Info struct {
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	ETag       string `json:"eTag"`
	Snils      string `json:"snils"`
}

type Elements struct {
	Elements []string `json:"elements"`
}

type Contact struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Address struct {
	ZipCode    string `json:"zipCode"`
	AddressStr string `json:"addressStr"`
	Building   string `json:"building"`
	Frame      string `json:"frame"`
	House      string `json:"house"`
	Flat       string `json:"flat"`
	Type       string `json:"type"`
}

type Documents struct {
	Type      string `json:"type"`
	Series    string `json:"series"`
	Number    string `json:"number"`
	IssuedBy  string `json:"issuedBy"`
	IssueId   string `json:"issueId"`
	IssueDate string `json:"issueDate"`
}
