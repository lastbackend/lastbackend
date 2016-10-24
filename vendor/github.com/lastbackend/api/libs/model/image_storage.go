package model

type ImageStorage struct {
	UUID     string
	UserID   string
	Username string
	Password string
	Email    string
	Host     string
	Main     bool
}

type ImageStorages []ImageStorage
