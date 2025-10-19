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

func (s *UserService) HandleCreateUser(request model.CreateUserRequest) (model.CreateUserResponse, error) {
	if request.Name == "" || request.Email == "" {
		return model.CreateUserResponse{}, fmt.Errorf("name and email are required")
	}

	user := &model.User{
		Name:  request.Name,
		Email: request.Email,
	}

	if err := database.SaveUser(user); err != nil {
		return model.CreateUserResponse{}, fmt.Errorf("user creation failed: %w", err)
	}

	response := model.CreateUserResponse{
		User: *user,
	}
	return response, nil
}

func (s *UserService) HandleUpdateUser(request model.UpdateUserRequest) (model.UpdateUserResponse, error) {
	if request.ID <= 0 {
		return model.UpdateUserResponse{}, fmt.Errorf("invalid user ID")
	}

	// Check if user exists
	existingUser, err := database.GetUserByID(request.ID)
	if err != nil {
		return model.UpdateUserResponse{}, fmt.Errorf("user not found: %w", err)
	}

	// Update user fields
	if request.Name != "" {
		existingUser.Name = request.Name
	}
	if request.Email != "" {
		existingUser.Email = request.Email
	}

	if err := database.SaveUser(existingUser); err != nil {
		return model.UpdateUserResponse{}, fmt.Errorf("user update failed: %w", err)
	}

	response := model.UpdateUserResponse{
		User: *existingUser,
	}
	return response, nil
}

func (s *UserService) HandleDeleteUser(request model.DeleteUserRequest) (model.DeleteUserResponse, error) {
	if request.ID <= 0 {
		return model.DeleteUserResponse{}, fmt.Errorf("invalid user ID")
	}

	// Check if user exists before deleting
	_, err := database.GetUserByID(request.ID)
	if err != nil {
		return model.DeleteUserResponse{
			Success: false,
			Message: fmt.Sprintf("User with ID %d not found", request.ID),
		}, nil
	}

	if err := database.DeleteUser(request.ID); err != nil {
		return model.DeleteUserResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to delete user: %v", err),
		}, nil
	}

	response := model.DeleteUserResponse{
		Success: true,
		Message: fmt.Sprintf("User with ID %d deleted successfully", request.ID),
	}
	return response, nil
}
