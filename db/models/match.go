package models

import "time"

type Match struct {
	Id          int       `json:"matchId"`
	UserId      int       `json:"userId"`
	MatchUserId int       `json:"matchUserId"`
	MatchCatId  int       `json:"matchCatId"`
	UserCatId   int       `json:"userCatId"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}
