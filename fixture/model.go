package fixture

import "github.com/google/uuid"

type Model interface {
	GetID() uuid.UUID
}
