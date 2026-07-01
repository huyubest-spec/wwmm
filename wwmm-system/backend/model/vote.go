package model

import "time"

type Vote struct {
	VoteID    int       `json:"voteId"`
	UserID    int       `json:"userId"`
	PhotoID   int       `json:"photoId"`
	TxHash    string    `json:"txHash"`
	CreatedAt time.Time `json:"createdAt"`
}
