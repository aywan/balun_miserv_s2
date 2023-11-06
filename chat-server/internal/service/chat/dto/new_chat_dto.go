package dto

type NewChatDTO struct {
	OwnerID int64
	Name    string
	Users   []int64
}
