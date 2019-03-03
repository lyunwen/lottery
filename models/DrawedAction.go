package models

type DrawedAction struct {
	ID          int
	AwardID     int
	PeopleCount int
	BackMoney   int
	Memo        string
	//ToDo Or Done
	Status string
}
