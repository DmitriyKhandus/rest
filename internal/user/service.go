package user

import (
	"context"

	"github.com/DmitriyKhandus/rest-api/pkg/logging"
)

type service struct {
	storage Storage
	logger  *logging.Logger
}

func (s *service) Create(ctx context.Context, dto CreateUserDTO) (User, error) {
	return User{}, nil
}
