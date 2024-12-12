package handlers

type UserSaver interface {
	SaveNewUser(firstName, lastName string, phoneNumber int) error
}

//func NewMessage(log *slog.Logger, userSaver UserSaver)
