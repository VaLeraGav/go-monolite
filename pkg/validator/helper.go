package validator

import (
	"errors"

	"github.com/google/uuid"
)

func ParseUUID(uuidStr string) (uuid.UUID, error) {
	if uuidStr == "" {
		return uuid.UUID{}, errors.New("UUID не указан")
	}

	uuidParse, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, errors.New("некорректный формат UUID")
	}

	return uuidParse, nil
}
