package model

import "time"

type Photo struct {
	PhotoID        int       `json:"photoId"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	ImageURL       string    `json:"imageUrl"`
	ImageHash      string    `json:"imageHash"`
	FileSize       int       `json:"fileSize"`
	PhotographerID int       `json:"photographerId"`
	Category       string    `json:"category"`
	ShootLocation  string    `json:"shootLocation"`
	ShootTime      string    `json:"shootTime"`
	CameraInfo     string    `json:"cameraInfo"`
	Status         int       `json:"status"`
	AuditComment   string    `json:"auditComment"`
	VoteCount      int       `json:"voteCount"`
	ViewCount      int       `json:"viewCount"`
	IsOnChain      int       `json:"isOnChain"`
	ChainTxHash    string    `json:"chainTxHash"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type PhotoFull struct {
	Photo
	PhotographerName     string `json:"photographerName"`
	PhotographerAvatar   string `json:"photographerAvatar"`
	HasVoted             bool   `json:"hasVoted"`
}

const (
	PhotoStatusPending  = 0
	PhotoStatusApproved = 1
	PhotoStatusRejected = 2
)
