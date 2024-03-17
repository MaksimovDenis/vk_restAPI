package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	filmoteka "vk_restAPI"
	"vk_restAPI/package/service"
	mock_service "vk_restAPI/package/service/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type MoviesWithActors struct {
	Id          int    `json:"id" db:"id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	ReleaseDate string `json:"release_date" db:"release_date"`
	Rating      int    `json:"rating" db:"rating"`
	Actors      string `json:"actors" db:"actors"`
}

func TestHandler_handleCreateMovie(t *testing.T) {
	type mockBehavior func(s *mock_service.MockMoviesWithActors, input filmoteka.MoviesWithActors)

	testTable := []struct {
		name                string
		inputBody           string
		inputMovie          filmoteka.MoviesWithActors
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Only for administrator",
			inputBody:           `{"title":"Dune 2", "description":"New Film", "release_date":"2024-03-07", "rating":"9","actors":"Test"}`,
			mockBehavior:        func(s *mock_service.MockMoviesWithActors, input filmoteka.MoviesWithActors) {},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":"This function is only available to the administrator"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			movieService := mock_service.NewMockMoviesWithActors(c)
			testCase.mockBehavior(movieService, testCase.inputMovie)

			services := &service.Service{MoviesWithActors: movieService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/movies/create", handler.handleCreateMovie)

			req := httptest.NewRequest("POST", "/movies/create", bytes.NewBufferString(testCase.inputBody))

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}

func TestHandler_handleGetAllMovies(t *testing.T) {
	type mockBehaivior func(s *mock_service.MockMoviesWithActors)

	testTable := []struct {
		name                string
		mockBehavior        mockBehaivior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockMoviesWithActors) {
				movies := []filmoteka.MoviesWithActors{
					{
						Id:          1,
						Title:       "Dune 2",
						Description: "New film",
						ReleaseDate: "2024-03-07",
						Rating:      9,
						Actors:      "Zendeya",
					},
				}

				s.EXPECT().GetMovies().Return(movies, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"id":1,"title":"Dune 2","description":"New film","release_date":"2024-03-07","rating":9,"actors":"Zendeya"}]}`,
		},
		{
			name: "Empty list",
			mockBehavior: func(s *mock_service.MockMoviesWithActors) {
				movies := []filmoteka.MoviesWithActors{}
				s.EXPECT().GetMovies().Return(movies, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"error":"The list of movies is empty"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			movieService := mock_service.NewMockMoviesWithActors(c)
			testCase.mockBehavior(movieService)

			services := &service.Service{MoviesWithActors: movieService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/api/movies", handler.handleGetAllMovies)

			req := httptest.NewRequest("GET", "/api/movies", nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}

func TestHandler_handleGetAllMoviesSortedByTitle(t *testing.T) {
	type mockBehaivior func(s *mock_service.MockMoviesWithActors)

	testTable := []struct {
		name                string
		mockBehavior        mockBehaivior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockMoviesWithActors) {
				movies := []filmoteka.MoviesWithActors{
					{
						Id:          1,
						Title:       "Dune 2",
						Description: "New film",
						ReleaseDate: "2024-03-07",
						Rating:      9,
						Actors:      "Zendeya",
					},
					{
						Id:          2,
						Title:       "The Great Gatsby",
						Description: "Old film",
						ReleaseDate: "2014-03-07",
						Rating:      9,
						Actors:      "Leonardo DiCaprio",
					},
				}

				s.EXPECT().GetMoviesSortedByTitle().Return(movies, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"id":1,"title":"Dune 2","description":"New film","release_date":"2024-03-07","rating":9,"actors":"Zendeya"},{"id":2,"title":"The Great Gatsby","description":"Old film","release_date":"2014-03-07","rating":9,"actors":"Leonardo DiCaprio"}]}`,
		},
		{
			name: "Empty list",
			mockBehavior: func(s *mock_service.MockMoviesWithActors) {
				movies := []filmoteka.MoviesWithActors{}
				s.EXPECT().GetMoviesSortedByTitle().Return(movies, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"error":"The list of movies is empty"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			movieService := mock_service.NewMockMoviesWithActors(c)
			testCase.mockBehavior(movieService)

			services := &service.Service{MoviesWithActors: movieService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/api/movies/sort/title", handler.handleGetAllMoviesSortedByTitle)

			req := httptest.NewRequest("GET", "/api/movies/sort/title", nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}

func TestHandler_handleGetAllMoviesSortedByDate(t *testing.T) {
	type mockBehaivior func(s *mock_service.MockMoviesWithActors)

	testTable := []struct {
		name                string
		mockBehavior        mockBehaivior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "Empty list",
			mockBehavior: func(s *mock_service.MockMoviesWithActors) {
				movies := []filmoteka.MoviesWithActors{}
				s.EXPECT().GetMoviesSortedByDate().Return(movies, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"error":"The list of movies is empty"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			movieService := mock_service.NewMockMoviesWithActors(c)
			testCase.mockBehavior(movieService)

			services := &service.Service{MoviesWithActors: movieService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/api/movies/sort/date", handler.handleGetAllMoviesSortedByDate)

			req := httptest.NewRequest("GET", "/api/movies/sort/date", nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}

func TestHandler_handleGetMovieById(t *testing.T) {
	type mockBehaivior func(s *mock_service.MockMoviesWithActors, id int)

	testTable := []struct {
		name                string
		requestURL          string
		mockBehavior        mockBehaivior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:       "OK",
			requestURL: "/api/movies/1",
			mockBehavior: func(s *mock_service.MockMoviesWithActors, id int) {
				movies := filmoteka.MoviesWithActors{

					Id:          1,
					Title:       "Dune 2",
					Description: "New film",
					ReleaseDate: "2024-03-07",
					Rating:      9,
					Actors:      "Zendeya",
				}
				s.EXPECT().GetMovieById(id).Return(movies, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"title":"Dune 2","description":"New film","release_date":"2024-03-07","rating":9,"actors":"Zendeya"}`,
		},
		{
			name:       "Invalid ID parameter",
			requestURL: "/api/movies/",
			mockBehavior: func(s *mock_service.MockMoviesWithActors, id int) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"invalid id parameter"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			movieService := mock_service.NewMockMoviesWithActors(c)
			testCase.mockBehavior(movieService, 1)

			services := &service.Service{MoviesWithActors: movieService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/api/movies/", handler.handleGetMovieById)

			req := httptest.NewRequest("GET", testCase.requestURL, nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}

func TestHandler_handleSearchMovieByTitle(t *testing.T) {
	type mockBehaivior func(s *mock_service.MockMoviesWithActors, fragment string)

	testTable := []struct {
		name                string
		inputBody           string
		inputFragment       string
		mockBehavior        mockBehaivior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:          "OK",
			inputBody:     `{"fragment":"Du"}`,
			inputFragment: "Du",
			mockBehavior: func(s *mock_service.MockMoviesWithActors, fragment string) {
				movies := []filmoteka.MoviesWithActors{
					{
						Id:          1,
						Title:       "Dune 2",
						Description: "New film",
						ReleaseDate: "2024-03-07",
						Rating:      9,
						Actors:      "Zendeya",
					},
				}
				s.EXPECT().SearchMoviesByTitle(fragment).Return(movies, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"id":1,"title":"Dune 2","description":"New film","release_date":"2024-03-07","rating":9,"actors":"Zendeya"}]}`,
		},
		{
			name:          "Empty",
			inputBody:     `{"fragment":"Du"}`,
			inputFragment: "Du",
			mockBehavior: func(s *mock_service.MockMoviesWithActors, fragment string) {
				movies := []filmoteka.MoviesWithActors{}
				s.EXPECT().SearchMoviesByTitle(fragment).Return(movies, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"error":"The list of movies is empty"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			movieService := mock_service.NewMockMoviesWithActors(c)
			testCase.mockBehavior(movieService, testCase.inputFragment)

			services := &service.Service{MoviesWithActors: movieService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/api/searchbytitle", handler.handleSearchMoviesByTitle)

			req := httptest.NewRequest("POST", "/api/searchbytitle", bytes.NewBufferString(testCase.inputBody))

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}

func TestHandler_handleSearchMovieByActor(t *testing.T) {
	type mockBehaivior func(s *mock_service.MockMoviesWithActors, fragment string)

	testTable := []struct {
		name                string
		inputBody           string
		inputFragment       string
		mockBehavior        mockBehaivior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:          "OK",
			inputBody:     `{"fragment":"z"}`,
			inputFragment: "z",
			mockBehavior: func(s *mock_service.MockMoviesWithActors, fragment string) {
				movies := []filmoteka.MoviesWithActors{
					{
						Id:          1,
						Title:       "Dune 2",
						Description: "New film",
						ReleaseDate: "2024-03-07",
						Rating:      9,
						Actors:      "Zendeya",
					},
				}
				s.EXPECT().SearchMovieByActorName(fragment).Return(movies, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"id":1,"title":"Dune 2","description":"New film","release_date":"2024-03-07","rating":9,"actors":"Zendeya"}]}`,
		},
		{
			name:          "Empty",
			inputBody:     `{"fragment":"z"}`,
			inputFragment: "z",
			mockBehavior: func(s *mock_service.MockMoviesWithActors, fragment string) {
				movies := []filmoteka.MoviesWithActors{}
				s.EXPECT().SearchMovieByActorName(fragment).Return(movies, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"error":"The list of movies is empty"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			movieService := mock_service.NewMockMoviesWithActors(c)
			testCase.mockBehavior(movieService, testCase.inputFragment)

			services := &service.Service{MoviesWithActors: movieService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/api/searchbyactor", handler.handleSearchMoviesByActorName)

			req := httptest.NewRequest("POST", "/api/searchbyactor", bytes.NewBufferString(testCase.inputBody))

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}

func TestHandler_handleUpdateMovie(t *testing.T) {
	type mockBehavior func(s *mock_service.MockMovies, id int)

	testTable := []struct {
		name                string
		requestURL          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Only for administrator",
			requestURL:          "/movies/update/1",
			mockBehavior:        func(s *mock_service.MockMovies, id int) {},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":"This function is only available to the administrator"}`,
		},
		{
			name:                "Missing ID Parameter",
			requestURL:          "/movies/update",
			mockBehavior:        func(s *mock_service.MockMovies, id int) {},
			expectedStatusCode:  301,
			expectedRequestBody: ``,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			moviesService := mock_service.NewMockMovies(c)
			testCase.mockBehavior(moviesService, 1)

			services := &service.Service{Movies: moviesService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/movies/update/", handler.handleUpdateMovie)

			req := httptest.NewRequest("PUT", testCase.requestURL, nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}

func TestHandler_handleDeleteMovie(t *testing.T) {
	type mockBehavior func(s *mock_service.MockMovies, id int)

	testTable := []struct {
		name                string
		requestURL          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Only for administrator",
			requestURL:          "/movies/delete/1",
			mockBehavior:        func(s *mock_service.MockMovies, id int) {},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":"This function is only available to the administrator"}`,
		},
		{
			name:                "Missing ID Parameter",
			requestURL:          "/movies/delete",
			mockBehavior:        func(s *mock_service.MockMovies, id int) {},
			expectedStatusCode:  301,
			expectedRequestBody: ``,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			moviesService := mock_service.NewMockMovies(c)
			testCase.mockBehavior(moviesService, 1)

			services := &service.Service{Movies: moviesService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/movies/delete/", handler.handleDeleteMovie)

			req := httptest.NewRequest("DELETE", testCase.requestURL, nil)

			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}
