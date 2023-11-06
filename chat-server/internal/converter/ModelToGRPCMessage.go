package converter

import (
	"github.com/aywan/balun_miserv_s2/chat-server/internal/model"
	desc "github.com/aywan/balun_miserv_s2/chat-server/pkg/grpc/v1/chat_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ModelToGRPCMessage(m model.Message) *desc.Message {
	msg := &desc.Message{
		Id:        m.ID,
		CreatedAt: timestamppb.New(m.CreatedAt),
		Type:      desc.MessageType(m.MsgType),
		Text:      m.Text,
	}

	if m.UserID.Valid {
		id := m.UserID.Int64
		msg.UserId = &id
	}

	return msg
}
