package models

type LuckyUser struct {
	UserID     int
	RtxName    string
	CName      string
	Level      int
	BackMoney  int
	Drawer     string
	AwardType  AwardType
	AwardValue int
	Memo       string
}

type AwardType int

const (
	Running    AwardType = 1 // value --> 0
	Stopped    AwardType = 2 // value --> 1
	Rebooting  AwardType = 3 // value --> 2
	Terminated AwardType = 4 // value --> 3
)
