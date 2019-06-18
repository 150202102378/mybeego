package utils

import uuid "github.com/satori/go.uuid"

//NewUUID: return new uuid
func NewUUID() uuid.UUID {
	u1, err := uuid.NewV1()
	if err != nil {
		u1 = NewUUID()
	}
	return u1
}
