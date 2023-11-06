package chat

const tableChat = "\"chat\""
const (
	colChatId        = "id"
	colChatCreatedAt = "created_at"
	colChatUpdatedAt = "updated_at"
	colChatDeletedAt = "deleted_at"
	colChatOwnerId   = "owner_id"
	colChatName      = "name"
)

const tableChatUser = "\"chat_user\""
const (
	colChatUserChatId        = "chat_id"
	colChatUserUserId        = "user_id"
	colChatUserLastMessageId = "last_message_id"
)

const tableMessage = "\"message\""
const (
	colMessageId        = "id"
	colMessageCreatedAt = "created_at"
	colMessageChatId    = "chat_id"
	colMessageUserId    = "user_id"
	colMessageType      = "type"
	colMessageText      = "text"
)

const tableChatMessage = "\"chat_message\""
const (
	colChatMessageChatId        = "chat_id"
	colChatMessageLastMessageId = "last_message_id"
)
