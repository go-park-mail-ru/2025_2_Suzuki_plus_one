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
		{
			ID:      "3",
			Title:   "Pulp Fiction",
			Year:    1994,
			Genres:  []string{"Crime", "Drama"},
			Preview: "https://upload.wikimedia.org/wikipedia/en/3/3b/Pulp_Fiction_%281994%29_poster.jpg",
		},
		{
			ID:      "4",
			Title:   "The Shawshank Redemption",
			Year:    1994,
			Genres:  []string{"Drama"},
			Preview: "https://upload.wikimedia.org/wikipedia/en/8/81/ShawshankRedemptionMoviePoster.jpg",
		},
		{
			ID:      "5",
			Title:   "Forrest Gump",
			Year:    1994,
			Genres:  []string{"Drama", "Romance"},
			Preview: "https://upload.wikimedia.org/wikipedia/en/6/67/Forrest_Gump_poster.jpg",
		},
		{
			ID:      "6",
			Title:   "Fight Club",
			Year:    1999,
			Genres:  []string{"Drama"},
			Preview: "https://upload.wikimedia.org/wikipedia/en/f/fc/Fight_Club_poster.jpg",
		},
		{
			ID:      "7",
			Title:   "The Social Network",
			Year:    2010,
			Genres:  []string{"Biography", "Drama"},
			Preview: "https://upload.wikimedia.org/wikipedia/en/8/8c/The_Social_Network_film_poster.png",
		},
		{
			ID:      "8",
			Title:   "Back to the Future",
			Year:    1985,
			Genres:  []string{"Adventure", "Comedy", "Sci-Fi"},
			Preview: "https://upload.wikimedia.org/wikipedia/en/d/d2/Back_to_the_Future.jpg",
		},
		{
			ID:      "9",
			Title:   "The Silence of the Lambs",
			Year:    1991,
			Genres:  []string{"Crime", "Drama", "Thriller"},
			Preview: "https://upload.wikimedia.org/wikipedia/en/8/86/The_Silence_of_the_Lambs_poster.jpg",
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
		hashedPassword, err := utils.HashPasswordBcrypt(user.PasswordHash)
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
