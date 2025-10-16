package db

import (
	"errors"
	"go.uber.org/zap"

	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/models"
	"github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/utils"
)

// Create a new user locking a db
//
// Returns:
//
//	*models.UserDB: pointer to the user (nil if creation failed)
//	error: error if any occurred during creation (nil if successful)
//	bool: true if user was created, false otherwise (already exists)
func (db *DataBase) CreateUser(email, password string, logger *zap.Logger) (*models.UserDB, error, bool) {
	// Validate email and password
	if err := utils.ValidateEmail(email); err != nil {
		return nil, err, false
	}
	if err := utils.ValidatePassword(password); err != nil {
		return nil, err, false
	}

	hashedPassword, err := utils.HashPasswordBcrypt(password)
	if err != nil {
		return nil, err, false
	}

	result := make(chan *models.UserDB)

	// Check if user already exists
	if user := db.FindUserByEmail(email); user != nil {
		logger.Error("User with this email already exists", zap.String("email", email))
		return user, errors.New("user with this email already exists"), false
	}

	go func() {
		defer close(result)
		db.mu.Lock()
		defer db.mu.Unlock()

		id := db.getNewID()
		// Simulate creating a user
		user := models.UserDB{
			ID:           id,
			Username:     "user" + id,
			Email:        email,
			PasswordHash: hashedPassword,
		}
		db.users = append(db.users, user)
		logger.Debug("User created",
			zap.String("id", user.ID),
			zap.String("username", user.Username),
			zap.String("email", user.Email))

		result <- &user
	}()

	return <-result, nil, true
}
