package models

import (
	"time"
)

type (
	Cat struct {
		Id          int       `json:"id"`
		UserId      int       `json:"user_id"`
		Name        string    `json:"name"`
		Race        string    `json:"race"`
		Sex         string    `json:"sex"`
		AgeInMonth  int       `json:"ageInMonth"`
		Description string    `json:"description"`
		ImageUrls   []string  `json:"imageUrls"`
		HasMatched  bool      `json:"hasMatched"`
		CreatedAt   time.Time `json:"createdAt"`
	}

	FilterGetCats struct {
		Id                 string `json:"id"`
		Limit              int    `json:"limit"`
		Offset             int    `json:"offset"`
		Race               string `json:"Race"`
		Sex                string `json:"sex"`
		HasMatched         bool   `json:"hasMatched"`
		AgeInMonthOperator string `json:"ageInMonthOperator"`
		AgeInMonthValue    int    `json:"ageInMonthValue"`
		Owned              bool   `json:"owned"`
		Search             string `json:"search"`
	}
)
