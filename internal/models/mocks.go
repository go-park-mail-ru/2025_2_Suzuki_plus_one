package models

import "github.com/go-park-mail-ru/2025_2_Suzuki_plus_one/internal/utils"

// Mock data for movies
func MockMovies() []Movie {
	return []Movie{
		{
			ID:      "1",
			Title:   "Inception",
			Year:    2010,
			Genres:  []string{"Sci-Fi", "Thriller"},
			Preview: "https://upload.wikimedia.org/wikipedia/en/2/2e/Inception_%282010%29_theatrical_poster.jpg",
		},
		{
			ID:      "2",
			Title:   "Interstellar",
			Year:    2011,
			Genres:  []string{"Sci-Fi"},
			Preview: "https://upload.wikimedia.org/wikipedia/en/b/bc/Interstellar_film_poster.jpg",
		},
	}
}

// Mock credentials for multiple users
func MockUsers() []UserDB {
	users := []UserDB{
		{
			ID:           "user1",
			Email:        "test@example.com",
			Username:     "testuser",
			PasswordHash: "hashedpassword123",
		},
		{
			ID:           "user2",
			Email:        "alice@example.com",
			Username:     "alice",
			PasswordHash: "hashedpassword456",
		},
		{
			ID:           "user3",
			Email:        "bob@example.com",
			Username:     "bob",
			PasswordHash: "hashedpassword789",
		},
	}

	// Hash passwords for mock users
	for i, user := range users {
		hashedPassword, err := utils.HashPassword(user.PasswordHash)
		if err == nil {
			users[i].PasswordHash = hashedPassword
		}
	}

	return users
}

// Mock data for sign-up request
func MockSignUpRequest() SignUpRequest {
	return SignUpRequest{
		Email:     "test@example.com",
		Password:  "password123",
		Username:  "testuser",
		FirstName: "Test",
		LastName:  "User",
	}
}

// Mock data for sign-in request
func MockSignInRequest() SignInRequest {
	return SignInRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
}

// Mock data for authentication response
func MockAuthResponse() SignInResponse {
	return SignInResponse{
		Token: "exampletoken123",
		User: UserAPI{
			ID:       "user1",
			Email:    "test@example.com",
			Username: "testuser",
		},
	}
}
