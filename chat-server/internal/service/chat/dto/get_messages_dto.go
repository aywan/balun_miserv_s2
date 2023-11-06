package dto

type GetMessagesDTO struct {
	ChatID          int64
	Limit           int64
	AfterMessageId  int64
	BeforeMessageId int64
}
