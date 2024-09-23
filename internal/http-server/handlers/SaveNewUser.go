package handlers

import "log/slog"

type UserSaver interface {
	SaveNewUser(firstName, lastName string, phoneNumber int) error
}

func New(log *slog.Logger, userSaver UserSaver)
