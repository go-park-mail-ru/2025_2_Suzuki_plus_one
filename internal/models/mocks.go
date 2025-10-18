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
			// * This User is used in tests and openAPI. Keep it first
			ID:           "user1",
			Email:        "test@example.com",
			Username:     "testuser",
			PasswordHash: "Password123!",
		},
		{
			ID:           "user2",
			Email:        "alice@example.com",
			Username:     "alice",
			PasswordHash: "Password123!",
		},
		{
			ID:           "user3",
			Email:        "bob@example.com",
			Username:     "bob",
			PasswordHash: "Password123!",
		},
	}

	// Hash passwords for mock users
	for i, user := range users {
		// Make sure mock data is valid
		if utils.ValidateEmail(user.Email) != nil {
			panic("Invalid email in mock users: " + user.Email)
		}
		if utils.ValidatePassword(user.PasswordHash) != nil {
			panic("Invalid password in mock users: " + user.PasswordHash)
		}

		// Hash password
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
		User: UserAPI{
			ID:       "user1",
			Email:    "test@example.com",
			Username: "testuser",
		},
	}
}

func MockActors() []Actor {
	return []Actor{
		{
			ID:           "1",
			Name:         "Leonardo DiCaprio",
			DateOfBirth:  "1974-11-11",
			PlaceOfBirth: "Los Angeles, California, USA",
			Biography:    "Leonardo Wilhelm DiCaprio is an American actor and film producer. Known for his work in biopics and period films, he is the recipient of numerous accolades, including an Academy Award, a British Academy Film Award, and three Golden Globe Awards.",
			Preview:      "https://upload.wikimedia.org/wikipedia/commons/thumb/f/f6/Leonardo_Dicaprio_-_World_Premiere_%E2%80%98One_Battle_after_Another%E2%80%99.jpg/800px-Leonardo_Dicaprio_-_World_Premiere_%E2%80%98One_Battle_after_Another%E2%80%99.jpg",
		},
		{
			ID:           "2",
			Name:         "Tom Hanks",
			DateOfBirth:  "1956-07-09",
			PlaceOfBirth: "Concord, California, USA",
			Biography:    "Thomas Jeffrey Hanks is an American actor and filmmaker. Known for both his comedic and dramatic roles, he is one of the most popular and recognizable film stars worldwide, and is regarded as an American cultural icon.",
			Preview:      "https://upload.wikimedia.org/wikipedia/commons/thumb/3/39/TomHanksPrincEdw031223_%2811_of_41%29_%28cropped%29.jpg/800px-TomHanksPrincEdw031223_%2811_of_41%29_%28cropped%29.jpg",
		},
		{
			ID:           "3",
			Name:         "Meryl Streep",
			DateOfBirth:  "1949-06-22",
			PlaceOfBirth: "Summit, New Jersey, USA",
			Biography:    "Mary Louise Streep is an American actress. Often described as the best actress of her generation, Streep is particularly known for her versatility and accent adaptation. She has received numerous accolades throughout her career.",
			Preview:      "https://upload.wikimedia.org/wikipedia/commons/thumb/0/00/Meryl_Streep_interview_at_Festival_de_Cannes_2024_%28cropped_2%29.jpg/800px-Meryl_Streep_interview_at_Festival_de_Cannes_2024_%28cropped_2%29.jpg",
		},
		{
			ID:           "4",
			Name:         "Robert Downey Jr.",
			DateOfBirth:  "1965-04-04",
			PlaceOfBirth: "Manhattan, New York City, USA",
			Biography:    "Robert John Downey Jr. is an American actor. His career has been characterized by critical and popular success in his youth, followed by a period of substance abuse and legal troubles, before a resurgence of commercial success in middle age.",
			Preview:      "https://upload.wikimedia.org/wikipedia/commons/thumb/2/23/Robert_Downey_Jr._2014_Comic-Con.jpg/800px-Robert_Downey_Jr._2014_Comic-Con.jpg",
		},
		{
			ID:           "5",
			Name:         "Scarlett Johansson",
			DateOfBirth:  "1984-11-22",
			PlaceOfBirth: "Manhattan, New York City, USA",
			Biography:    "Scarlett Ingrid Johansson is an American actress. The world's highest-paid actress in 2018 and 2019, she has featured multiple times on the Forbes Celebrity 100 list. Her films have grossed over $14.3 billion worldwide.",
			Preview:      "https://upload.wikimedia.org/wikipedia/commons/thumb/a/ad/Scarlett_Johansson-8588.jpg/800px-Scarlett_Johansson-8588.jpg",
		},
		{
			ID:           "6",
			Name:         "Denzel Washington",
			DateOfBirth:  "1954-12-28",
			PlaceOfBirth: "Mount Vernon, New York, USA",
			Biography:    "Denzel Hayes Washington Jr. is an American actor, director, and producer. He has been described as an actor who reconfigured the concept of classic movie stardom. Throughout his career, Washington has received numerous accolades.",
			Preview:      "https://upload.wikimedia.org/wikipedia/commons/thumb/c/cc/Denzel_Washington_at_the_2025_Cannes_Film_Festival.jpg/800px-Denzel_Washington_at_the_2025_Cannes_Film_Festival.jpg",
		},
	}
}
