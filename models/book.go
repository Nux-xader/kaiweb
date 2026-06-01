package models

type KAIWebPassengerData struct {
	PassengerTitle  string `json:"penumpang_title[]"`
	PassengerIdType string `json:"penumpang_tandapengenal[]"`
	PassengerType   string `json:"penumpang_type[]"`
	PassengerName   string `json:"penumpang_nama[]"`
	PassengerId     string `json:"penumpang_notandapengenal[]"`
}
