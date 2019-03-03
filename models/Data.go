package models



type Data struct {
	Version          string
	Count            Count
	Awards           []Award
	Actions          []DrawedAction
	Users            []User
	DrawedRecords    []DrawedRecord
	BackMoneyRecords []BackMoneyRecord
}
