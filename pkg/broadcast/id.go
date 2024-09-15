package broadcast

import "github.com/google/uuid"

type ReceiverID string

func NewReceiverID() ReceiverID {
	id, err := uuid.NewV7()
	if err != nil {
		return ReceiverID(uuid.New().String())
	}
	return ReceiverID(id.String())
}
