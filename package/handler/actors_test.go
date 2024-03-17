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

func TestHandler_handleCreateActor(t *testing.T) {
	type mockBehavior func(s *mock_service.MockActors, input filmoteka.Actors)

	testTable := []struct {
		name                string
		inputBody           string
		inputActor          filmoteka.Actors
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Only for administrator",
			inputBody:           `{"first_name":"John", "last_name":"Doe", "gender":"male", "date_of_birth":"1990-01-01"}`,
			mockBehavior:        func(s *mock_service.MockActors, input filmoteka.Actors) {},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":"This function is only available to the administrator"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			actorService := mock_service.NewMockActors(c)
			testCase.mockBehavior(actorService, testCase.inputActor)

			services := &service.Service{Actors: actorService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/actors/create", handler.handleCreateActor)

			req := httptest.NewRequest("POST", "/actors/create", bytes.NewBufferString(testCase.inputBody))

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

func TestHandler_handleGetAllActors(t *testing.T) {
	type mockBehaivior func(s *mock_service.MockActorsWithMovies)

	testTable := []struct {
		name                string
		mockBehavior        mockBehaivior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK",
			mockBehavior: func(s *mock_service.MockActorsWithMovies) {
				actor := []filmoteka.ActorsWithMovies{
					{
						Id:          1,
						FirstName:   "Jhon",
						LastName:    "Doe",
						Gender:      "male",
						DateOfBirth: "1970-01-01",
						Movies:      "Terminator",
					},
				}

				s.EXPECT().GetActors().Return(actor, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"data":[{"id":1,"first_name":"Jhon","last_name":"Doe","gender":"male","date_of_birth":"1970-01-01","movies":"Terminator"}]}`,
		},
		{
			name: "Empty List",
			mockBehavior: func(s *mock_service.MockActorsWithMovies) {
				actor := []filmoteka.ActorsWithMovies{}

				s.EXPECT().GetActors().Return(actor, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"error":"The list of actors is empty"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			actorService := mock_service.NewMockActorsWithMovies(c)
			testCase.mockBehavior(actorService)

			services := &service.Service{ActorsWithMovies: actorService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/api/actors", handler.handleGetAllActors)

			req := httptest.NewRequest("GET", "/api/actors", nil)

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

func TestHandler_handleGetActorById(t *testing.T) {
	type mockBehaivior func(s *mock_service.MockActorsWithMovies, id int)

	testTable := []struct {
		name                string
		requestURL          string
		mockBehavior        mockBehaivior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:       "OK",
			requestURL: "/api/actors/1",
			mockBehavior: func(s *mock_service.MockActorsWithMovies, id int) {
				actor := filmoteka.ActorsWithMovies{
					Id:          1,
					FirstName:   "Jhon",
					LastName:    "Doe",
					Gender:      "male",
					DateOfBirth: "1970-01-01",
					Movies:      "Terminator",
				}

				s.EXPECT().GetActorById(id).Return(actor, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1,"first_name":"Jhon","last_name":"Doe","gender":"male","date_of_birth":"1970-01-01","movies":"Terminator"}`,
		},
		{
			name:       "Invalid ID parameter",
			requestURL: "/api/actors/",
			mockBehavior: func(s *mock_service.MockActorsWithMovies, id int) {
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

			actorService := mock_service.NewMockActorsWithMovies(c)
			testCase.mockBehavior(actorService, 1)

			services := &service.Service{ActorsWithMovies: actorService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/api/actors/", handler.handleGetActorById)

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

func TestHandler_handleUpdateActor(t *testing.T) {
	type mockBehavior func(s *mock_service.MockActors, id int)

	testTable := []struct {
		name                string
		requestURL          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Only for administrator",
			requestURL:          "/actors/update/1",
			mockBehavior:        func(s *mock_service.MockActors, id int) {},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":"This function is only available to the administrator"}`,
		},
		{
			name:                "Missing ID Parameter",
			requestURL:          "/actors/update",
			mockBehavior:        func(s *mock_service.MockActors, id int) {},
			expectedStatusCode:  301,
			expectedRequestBody: ``,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			actorService := mock_service.NewMockActors(c)
			testCase.mockBehavior(actorService, 1)

			services := &service.Service{Actors: actorService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/actors/update/", handler.handleUpdateActor)

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

func TestHandler_handleDeleteActor(t *testing.T) {
	type mockBehavior func(s *mock_service.MockActors, id int)

	testTable := []struct {
		name                string
		requestURL          string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:                "Only for administrator",
			requestURL:          "/actors/delete/1",
			mockBehavior:        func(s *mock_service.MockActors, id int) {},
			expectedStatusCode:  403,
			expectedRequestBody: `{"error":"This function is only available to the administrator"}`,
		},
		{
			name:                "Missing ID Parameter",
			requestURL:          "/actors/delete",
			mockBehavior:        func(s *mock_service.MockActors, id int) {},
			expectedStatusCode:  301,
			expectedRequestBody: ``,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			//Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			actorService := mock_service.NewMockActors(c)
			testCase.mockBehavior(actorService, 1)

			services := &service.Service{Actors: actorService}
			handler := NewHandler(services)

			mux := http.NewServeMux()
			mux.HandleFunc("/actors/delete/", handler.handleDeleteActor)

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
