package handlers

type UserSaver interface {
	SaveNewUser(firstName, lastName string, phoneNumber int) error
}

//func NewNotification(log *slog.Logger, userSaver UserSaver)
