package handlers

import (
	"CatsSocial/api/responses"
	"CatsSocial/db/functions"
	"CatsSocial/db/models"
	"CatsSocial/utils"
	"errors"
	"net/http"
	"net/mail"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Database *functions.User
}

func validateUser(req struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}) error {
	lenEmail := len(req.Email)
	lenPassword := len(req.Password)
	lenName := len(req.Name)

	if lenEmail == 0 || lenPassword == 0 || lenName == 0 {
		return errors.New("email, name, and password are required")
	}

	if !validate_email(req.Email) {
		return errors.New("email is not in a valid format")
	}

	if lenPassword < 5 || lenEmail < 5 {
		return errors.New("email and password length must be at least 5 characters")
	}

	if lenPassword > 15 {
		return errors.New("password length cannot exceed 15 characters")
	}

	if lenName > 50 {
		return errors.New("name length cannot exceed 50 characters")
	}

	return nil
}

func validate_email(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func validateLogin(req struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}) error {
	lenEmail := len(req.Email)
	lenPassword := len(req.Password)

	if lenEmail == 0 || lenPassword == 0 {
		return errors.New("email, and password are required")
	}

	if !validate_email(req.Email) {
		return errors.New("email is not in a valid format")
	}

	if lenPassword < 5 {
		return errors.New("email and password length must be at least 5 characters")
	}

	if lenPassword > 15 {
		return errors.New("password length cannot exceed 15 characters")
	}

	return nil
}

func (u *User) Register(ctx *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Email    string `json:"email"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Validate request body
	if err := validateUser(req); err != nil {
		status, response := responses.ErrorBadRequests(err.Error())
		return ctx.Status(status).JSON(response)
	}

	// Create user object
	usr := models.User{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}

	// Register user
	result, err := u.Database.Register(ctx.UserContext(), usr)
	if err != nil {
		if err.Error() == "EXISTING_EMAIL" {
			status, response := responses.ErrorConflict(err.Error())
			return ctx.Status(status).JSON(response)
		}

		status, response := responses.ErrorServers(err.Error())
		return ctx.Status(status).JSON(response)
	}

	// generate access token
	accessToken, err := utils.GenerateAccessToken(result.Email, result.Id)
	if err != nil {
		status, response := responses.ErrorServers(err.Error())
		return ctx.Status(status).JSON(response)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"email":       result.Email,
			"name":        result.Name,
			"accessToken": accessToken,
		},
	})
}

func (u *User) Login(ctx *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	if err := validateLogin(req); err != nil {
		status, response := responses.ErrorBadRequests(err.Error())
		return ctx.Status(status).JSON(response)
	}

	// login user
	result, err := u.Database.Login(ctx.UserContext(), req.Email, req.Password)
	if err != nil {
		if err.Error() == "USER_NOT_FOUND" {
			status, response := responses.ErrorNotFound(err.Error())
			return ctx.Status(status).JSON(response)
		}

		if err.Error() == "INVALID_PASSWORD" {
			status, response := responses.ErrorBadRequests(err.Error())
			return ctx.Status(status).JSON(response)
		}

		status, response := responses.ErrorServers(err.Error())
		return ctx.Status(status).JSON(response)
	}

	// generate access token
	accessToken, err := utils.GenerateAccessToken(result.Email, result.Id)
	if err != nil {
		status, response := responses.ErrorServers(err.Error())
		return ctx.Status(status).JSON(response)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged successfully",
		"data": fiber.Map{
			"email":       result.Email,
			"name":        result.Name,
			"accessToken": accessToken,
		},
	})
}
