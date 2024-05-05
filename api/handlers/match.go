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
	"github.com/gofiber/fiber/v2"
)

type (
	MatchHandler struct {
		Match        functions.Match
		CatDatabase  *functions.Cat
		UserDatabase *functions.User
	}

	MatchIssuer struct {
		Name      string `json:"name"`
		Email     string `json:"email"`
		CreatedAt string `json:"createdAt"`
	}

	MatchDetailResponse struct {
		Id             int               `json:"id"`
		Issuedby       MatchIssuer       `json:"issuedBy"`
		MatchCatDetail CatDetailResponse `json:"matchCatDetail"`
		UserCatDetail  CatDetailResponse `json:"userCatDetail"`
		Message        string            `json:"message"`
		CreatedAt      string            `json:"createdAt"`
	}

	GetMatchesResponse struct {
		Data []MatchDetailResponse `json:"data"`
	}
)

func (m *MatchHandler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, fiber.ErrUnauthorized):
		return fiber.ErrUnauthorized
	case errors.Is(err, fiber.ErrForbidden):
		return fiber.ErrForbidden
	case errors.Is(err, functions.ErrNoRow):
		status, response := responses.ErrorNotFound("not found")
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

func (m *MatchHandler) Create(c *fiber.Ctx) error {
	userIDClaim := c.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	var payload struct {
		matchCatId string `json:"matchCatId"`
		userCatId  string `json:"userCatId"`
		message    string `json:"message"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	var (
		lenMessage = len(payload.message)
		isValid    = func(length int) bool {
			return length >= 5 && length <= 120
		}
	)

	if !isValid(lenMessage) {
		return c.SendStatus(http.StatusBadRequest)
	}

	catID, err := strconv.Atoi(c.Params(payload.userCatId))
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	cat, err := m.CatDatabase.FindByIDUser(c.UserContext(), catID, userID)
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	matchCat, err := m.CatDatabase.FindByID(c.UserContext(), catID)
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	//400 if both matchCatId &userCatId already matched
	if cat.HasMatched {
		return c.SendStatus(http.StatusBadRequest)
	}
	if matchCat.HasMatched {
		return c.SendStatus(http.StatusBadRequest)
	}

	//400 if the cat’s gender is same
	if cat.Sex == matchCat.Sex {
		return c.SendStatus(http.StatusBadRequest)
	}

	//400 if matchCatId & userCatId is from the same owner
	if cat.UserID == matchCat.UserID {
		return c.SendStatus(http.StatusBadRequest)
	}

	matchCatIdInt, err := strconv.Atoi(payload.matchCatId)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	userCatIdInt, err := strconv.Atoi(payload.userCatId)
	if err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := m.Match.Create(c.UserContext(), models.Match{
		UserId:     userID,
		MatchCatId: matchCatIdInt,
		UserCatId:  userCatIdInt,
		Message:    payload.message,
	}); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusCreated)
}

func (m *MatchHandler) getIssuerDetail(c *fiber.Ctx, issuerId int) (MatchIssuer, error) {

	user, err := m.UserDatabase.GetUserById(c.UserContext(), strconv.Itoa(issuerId))
	if err != nil {
		return MatchIssuer{}, c.SendStatus(http.StatusInternalServerError)
	}

	return MatchIssuer{
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (m *MatchHandler) getCatDetail(c *fiber.Ctx, catId int) (CatDetailResponse, error) {

	cat, err := m.CatDatabase.FindByID(c.UserContext(), catId)
	if err != nil {
		return CatDetailResponse{}, c.SendStatus(http.StatusInternalServerError)
	}

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
	}, nil
}

func (m *MatchHandler) convertMatchModelToDetailResponse(c *fiber.Ctx, match models.Match) (MatchDetailResponse, error) {

	MatchIssuerDetail, err := m.getIssuerDetail(c, match.UserId)
	if err != nil {
		return MatchDetailResponse{}, err
	}

	MatchCatRes, err := m.getCatDetail(c, match.MatchCatId)
	if err != nil {
		return MatchDetailResponse{}, err
	}

	UserCatRes, err := m.getCatDetail(c, match.UserCatId)
	if err != nil {
		return MatchDetailResponse{}, err
	}

	return MatchDetailResponse{
		Id:             match.Id,
		Issuedby:       MatchIssuerDetail,
		MatchCatDetail: MatchCatRes,
		UserCatDetail:  UserCatRes,
		Message:        match.Message,
		CreatedAt:      match.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (m *MatchHandler) convertMatchesToGetMatchesResponse(
	c *fiber.Ctx,
	matches []models.Match,
) (GetMatchesResponse, error) {
	var result []MatchDetailResponse
	for _, match := range matches {

		MatchDetail, err := m.convertMatchModelToDetailResponse(c, match)
		if err != nil {
			return GetMatchesResponse{}, err
		}

		result = append(result, MatchDetail)
	}

	return GetMatchesResponse{
		Data: result,
	}, nil
}

func (m *MatchHandler) Get(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)

	matches, err := m.Match.GetRelatedMatches(c.UserContext(), userId)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	matchesResponse, err := m.convertMatchesToGetMatchesResponse(c, matches)
	if err != nil {
		return err
	}

	return c.JSON(map[string]interface{}{
		"message": "success",
		"data":    matchesResponse,
	})
}

func (m *MatchHandler) Approve(c *fiber.Ctx) error {
	var payload struct {
		matchId string `json:"matchId"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	match, err := m.Match.GetMatchById(c.UserContext(), payload.matchId)
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	if match.Status == "removed" || match.Status == "approved" {
		return m.handleError(c, fiber.ErrBadRequest)
	}

	err = m.Match.UpdateStatus(c.UserContext(), match, "approved")
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	userCat, err := m.CatDatabase.FindByID(c.UserContext(), match.UserCatId)
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	userCat.HasMatched = true

	err = m.CatDatabase.Update(c.UserContext(), userCat)
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	matchCat, err := m.CatDatabase.FindByID(c.UserContext(), match.MatchCatId)
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	matchCat.HasMatched = true

	err = m.CatDatabase.Update(c.UserContext(), matchCat)
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	userCatMatches, err := m.Match.GetRelatedCatMatches(c.UserContext(), match.UserCatId)
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	for _, match := range userCatMatches {

		err = m.Match.UpdateStatus(c.UserContext(), match, "removed")
		if err != nil {
			if err == functions.ErrNoRow {
				return m.handleError(c, fiber.ErrNotFound)
			}
			return m.handleError(c, err)
		}
	}

	matchCatMatches, err := m.Match.GetRelatedCatMatches(c.UserContext(), match.MatchCatId)
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	for _, match := range matchCatMatches {

		err = m.Match.UpdateStatus(c.UserContext(), match, "removed")
		if err != nil {
			if err == functions.ErrNoRow {
				return m.handleError(c, fiber.ErrNotFound)
			}
			return m.handleError(c, err)
		}
	}

	return c.SendStatus(http.StatusOK)
}

func (m *MatchHandler) Reject(c *fiber.Ctx) error {
	var payload struct {
		matchId string `json:"matchId"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	match, err := m.Match.GetMatchById(c.UserContext(), payload.matchId)
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	if match.Status == "removed" || match.Status == "approved" {
		return m.handleError(c, fiber.ErrBadRequest)
	}

	err = m.Match.UpdateStatus(c.UserContext(), match, "removed")
	if err != nil {
		if err == functions.ErrNoRow {
			return m.handleError(c, fiber.ErrNotFound)
		}
		return m.handleError(c, err)
	}

	return c.SendStatus(http.StatusOK)
}

func (m *MatchHandler) Delete(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	matchId := c.Params("matchId")

	err := m.Match.Delete(c.UserContext(), userId, matchId)
	if err != nil {
		if errors.Is(err, functions.ErrNoRow) {
			return c.SendStatus(http.StatusNotFound)
		}

		if errors.Is(err, functions.ErrUnauthorized) {
			return c.SendStatus(http.StatusUnauthorized)
		}

		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}
