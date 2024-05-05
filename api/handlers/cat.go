package handlers

import (
	"CatsSocial/api/responses"
	"CatsSocial/db/functions"
	"CatsSocial/db/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofiber/fiber/v2"
)

type Race string

const (
	Persian           Race = "Persian"
	Maine_Coon        Race = "Maine Coon"
	Siamese           Race = "Siamese"
	Ragdoll           Race = "Ragdoll"
	Bengal            Race = "Bengal"
	Sphynx            Race = "Sphynx"
	British_Shorthair Race = "British Shorthair"
	Abyssinian        Race = "Abyssinian"
	Scottish_Fold     Race = "Scottish Fold"
	Birman            Race = "Birman"
)

type Sex string

const (
	Male      Sex = "male"
	FemaleSex     = "female"
)

type (
	Cat struct {
		Database     *functions.Cat
		UserDatabase *functions.User
	}

	CatPayload struct {
		Name        string   `json:"name"`
		Race        string   `json:"race"`
		Sex         string   `json:"sex"`
		AgeInMonth  int      `json:"ageInMonth"`
		Description string   `json:"description"`
		ImageUrls   []string `json:"imageUrl"`
	}

	QueryFilterGetCats struct {
		Id         string `json:"id"`
		Limit      int    `json:"limit"`
		Offset     int    `json:"offset"`
		Race       string `json:"Race"`
		Sex        string `json:"sex"`
		HasMatched bool   `json:"hasMatched"`
		AgeInMonth string `json:"ageInMonth"`
		Owned      bool   `json:"owned"`
		Search     string `json:"search"`
	}

	CatResponse struct {
		Id        string `json:"id"`
		CreatedAt string `json:"createdAt"`
	}

	CatDetailResponse struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		Race        string   `json:"race"`
		Sex         string   `json:"sex"`
		AgeInMonth  int      `json:"ageInMonth"`
		Description string   `json:"description"`
		ImageUrls   []string `json:"imageUrl"`
		HasMatched  bool     `json:"hasMatched"`
		CreatedAt   string   `json:"createdAt"`
	}

	Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	GetCatsResponse struct {
		Data []CatDetailResponse `json:"data"`
		//Meta Meta                `json:"meta"`
	}
)

func (app CatPayload) Validate() error {
	return validation.ValidateStruct(&app,
		// Name cannot be empty, and the length must be between 1 and 30.
		validation.Field(&app.Name, validation.Required, validation.Length(1, 30)),
		// Race cannot be empty, and should be in the Race enum.
		validation.Field(&app.Race, validation.Required, validation.In("Persian", "Maine Coon", "Siamese", "Ragdoll", "Bengal", "Sphynx", "British Shorthair", "Abyssinian", "Scottish Fold", "Birman")),
		// Sex cannot be empty and should be either "male" or "female".
		validation.Field(&app.Sex, validation.Required, validation.In("male", "female")),
		// Stock cannot be empty, and minimum value is 1 and maximum value is 120082
		validation.Field(&app.AgeInMonth, validation.NotNil, validation.Min(1), validation.Max(120082)),
		// Description cannot be empty, and the length must be between 1 and 200.
		validation.Field(&app.Description, validation.Required, validation.Length(1, 200)),
		// Tags cannot be empty, and should have at least 0 items.
		validation.Field(&app.ImageUrls, validation.Required, validation.Each(is.URL)),
	)
}

func (app QueryFilterGetCats) Validate() error {
	return validation.ValidateStruct(&app,
		// Limit should be greater than 0.
		validation.Field(&app.Limit, validation.Min(0)),
		// Offset cannot should be greater than 0.
		validation.Field(&app.Offset, validation.Min(0)),
		// Race should be in the Race enum.
		validation.Field(&app.Race, validation.In("Persian", "Maine Coon", "Siamese", "Ragdoll", "Bengal", "Sphynx", "British Shorthair", "Abyssinian", "Scottish Fold", "Birman")),
		// Sex should be either "male" or "female".
		validation.Field(&app.Sex, validation.Required, validation.In("male", "female")),
	)
}

func parseAgeInMonthQuery(query string) (string, int, error) {
	parts := strings.Split(query, "=")
	if len(parts) != 2 || parts[0] != "ageInMonth" {
		return "", 0, fmt.Errorf("invalid query format: %s", query)
	}

	operator := "="
	if strings.HasPrefix(parts[1], "<") {
		operator = "=<"
		parts[1] = parts[1][1:]
	} else if strings.HasPrefix(parts[1], ">") {
		operator = "=>"
		parts[1] = parts[1][1:]
	}

	value, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", 0, fmt.Errorf("invalid value in query: %s", query)
	}

	return operator, value, nil
}

func (p *Cat) convertCatEntityToResponse(cat models.Cat) CatResponse {
	return CatResponse{
		Id:        strconv.Itoa(cat.ID),
		CreatedAt: cat.CreatedAt.Format(time.RFC3339),
	}
}

func (p *Cat) convertCatEntityToDetailResponse(cat models.Cat) CatDetailResponse {
	return CatDetailResponse{
		Id:          strconv.Itoa(cat.ID),
		Name:        cat.Name,
		Race:        cat.Race,
		Sex:         cat.Sex,
		AgeInMonth:  cat.AgeInMonth,
		Description: cat.Description,
		ImageUrls:   cat.ImageUrls,
		HasMatched:  cat.HasMatched,
		CreatedAt:   cat.CreatedAt.Format(time.RFC3339),
	}
}

func (p *Cat) convertCatsToGetCatsResponse(
	cats []models.Cat,
	limit, offset, total int,
) GetCatsResponse {
	var result []CatDetailResponse
	for _, cat := range cats {
		result = append(result, p.convertCatEntityToDetailResponse(cat))
	}

	return GetCatsResponse{
		Data: result,
		// Meta: Meta{
		// 	Limit:  limit,
		// 	Offset: offset,
		// 	Total:  total,
		// },
	}
}

func (p *Cat) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, fiber.ErrUnauthorized):
		return fiber.ErrUnauthorized
	case errors.Is(err, fiber.ErrForbidden):
		return fiber.ErrForbidden
	case errors.Is(err, functions.ErrNoRow):
		status, response := responses.ErrorNotFound("no cat found")
		return c.Status(status).JSON(response)
	default:
		validationErrors, ok := err.(validation.Errors)
		if !ok {
			status, response := responses.ErrorServer(err.Error())
			return c.Status(status).JSON(response)
		}

		errMessages := []string{}
		for key, ve := range validationErrors {
			errMessages = append(errMessages, fmt.Sprintf(
				"field %s: %s",
				key,
				ve.Error()))
		}

		status, response := responses.ErrorBadRequests(strings.Join(errMessages, ""))
		return c.Status(status).JSON(response)
	}
}

func (p *Cat) GetCats(c *fiber.Ctx) error {
	var (
		userID   int
		err      error
		operator string
		value    int
	)

	var filter QueryFilterGetCats
	if err := c.QueryParser(&filter); err != nil {
		return p.handleError(c, errors.New(fmt.Sprintf("failed to parse query params: %v", err.Error())))
	}

	err = filter.Validate()
	if err != nil {
		return p.handleError(c, err)
	}

	if c.Locals("user_id") != nil {
		userIDClaim := c.Locals("user_id").(string)
		userID, err = strconv.Atoi(userIDClaim)
		if err != nil {
			return p.handleError(c, errors.New(fmt.Sprintf("failed parse user id: %v", err.Error())))
		}
	}

	if len(filter.AgeInMonth) != 0 {
		operator, value, err = parseAgeInMonthQuery(filter.AgeInMonth)
		if err != nil {
			return p.handleError(c, errors.New(fmt.Sprintf("failed parse user ageInMonthQuery: %v", err.Error())))
		}
	}

	filterDB := models.FilterGetCats{
		Id:                 filter.Id,
		Limit:              filter.Limit,
		Offset:             filter.Offset,
		Race:               filter.Race,
		Sex:                filter.Sex,
		HasMatched:         filter.HasMatched,
		AgeInMonthOperator: operator,
		AgeInMonthValue:    value,
		Owned:              filter.Owned,
		Search:             filter.Search,
	}

	cats, err := p.Database.FindAll(c.UserContext(), filterDB, userID)
	if err != nil {
		return p.handleError(c, err)
	}

	total, err := p.Database.Count(c.UserContext(), filterDB, userID)
	if err != nil {
		return p.handleError(c, err)
	}

	if filter.Limit == 0 {
		filter.Limit = 5
	}

	result := p.convertCatsToGetCatsResponse(cats, filter.Limit, filter.Offset, total)

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "success",
		"data":    result,
	})
}

func (p *Cat) AddCat(c *fiber.Ctx) error {
	userIDClaim := c.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return p.handleError(c, errors.New(fmt.Sprintf("failed parse user id: %v", err.Error())))
	}

	_, err = p.UserDatabase.GetUserById(c.UserContext(), userIDClaim)
	if err != nil {
		return p.handleError(c, fiber.ErrUnauthorized)
	}

	var payload CatPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	err = payload.Validate()
	if err != nil {
		return p.handleError(c, err)
	}

	cat, err := p.Database.Add(c.UserContext(), models.Cat{
		UserID:      userID,
		Name:        payload.Name,
		Race:        payload.Race,
		Sex:         payload.Sex,
		AgeInMonth:  payload.AgeInMonth,
		Description: payload.Description,
		ImageUrls:   payload.ImageUrls,
	})

	if err != nil {
		return p.handleError(c, err)
	}

	result := p.convertCatEntityToResponse(cat)

	return c.Status(http.StatusCreated).JSON(map[string]interface{}{
		"message": "success",
		"data":    result,
	})
}

func (p *Cat) UpdateCat(c *fiber.Ctx) error {
	userIDClaim := c.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return p.handleError(c, errors.New(fmt.Sprintf("failed parse user id: %v", err.Error())))
	}

	_, err = p.UserDatabase.GetUserById(c.UserContext(), userIDClaim)
	if err != nil {
		return p.handleError(c, fiber.ErrUnauthorized)
	}

	var payload CatPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	err = payload.Validate()
	if err != nil {
		return p.handleError(c, err)
	}

	catID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return p.handleError(c, errors.New("failed parse cat id"))
	}

	cat, err := p.Database.FindByIDUser(c.UserContext(), catID, userID)
	if err != nil {
		if err == functions.ErrNoRow {
			return p.handleError(c, fiber.ErrNotFound)
		}
		return p.handleError(c, err)
	}

	// 400 sex is edited when cat is already requested to match
	if cat.HasMatched && (len(payload.Sex) != 0) {
		return p.handleError(c, fiber.ErrBadRequest)
	}

	cat.Name = payload.Name
	cat.Race = payload.Race
	cat.Sex = payload.Sex
	cat.AgeInMonth = payload.AgeInMonth
	cat.Description = payload.Description
	cat.ImageUrls = payload.ImageUrls

	err = p.Database.Update(c.UserContext(), cat)
	if err != nil {
		return p.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(map[string]interface{}{})
}

func (p *Cat) DeleteCat(c *fiber.Ctx) error {
	userIDClaim := c.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return p.handleError(c, errors.New(fmt.Sprintf("failed parse user id: %v", err.Error())))
	}

	_, err = p.UserDatabase.GetUserById(c.UserContext(), userIDClaim)
	if err != nil {
		return p.handleError(c, fiber.ErrUnauthorized)
	}

	catID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return p.handleError(c, errors.New("failed parse cat id"))
	}

	_, err = p.Database.FindByIDUser(c.UserContext(), catID, userID)
	if err != nil {
		if err == functions.ErrNoRow {
			return p.handleError(c, fiber.ErrNotFound)
		}
		return p.handleError(c, err)
	}

	err = p.Database.DeleteByID(c.UserContext(), catID)
	if err != nil {
		return p.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(map[string]interface{}{})
}
