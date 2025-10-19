package service

import (
	"fmt"

	"github.com/maasumiyaat/soap/database"
	"github.com/maasumiyaat/soap/model"
)

type UserService struct{}

func (s *UserService) HandleGetUserByID(request model.GetUserByIDRequest) (model.GetUserByIDResponse, error) {
	user, err := database.GetUserByID(request.ID)
	if err != nil {
		return model.GetUserByIDResponse{}, fmt.Errorf("user retrieval failed: %w", err)
	}

	response := model.GetUserByIDResponse{
		User: *user,
	}
	return response, nil
}
