package inbox

import "github.com/google/uuid"

type InboxMessage struct {
	EventID     uuid.UUID
	ProcessName string
}
