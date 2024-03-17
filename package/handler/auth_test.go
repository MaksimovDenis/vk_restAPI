package handler

import (
	"bytes"
	"errors"
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

func TestHandler_signUp(t *testing.T) {
	type mockBehaivior func(s *mock_service.MockAuthorization, user filmoteka.User)

	testTable := []struct {
		name                string
		inputBody           string
		inputUser           filmoteka.User
		mockBehaivior       mockBehaivior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test", "password":"test","is_admin":true}`,
			inputUser: filmoteka.User{
				Username: "test",
				Password: "test",
				Is_admin: true,
			},
			mockBehaivior: func(s *mock_service.MockAuthorization, user filmoteka.User) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"id":1}`,
		},
		{
			name:      "Empty Fields",
			inputBody: `{"username":"test"}`,
			mockBehaivior: func(s *mock_service.MockAuthorization, user filmoteka.User) {
			},
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Username and password are required"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"username":"test", "password":"test","is_admin":true}`,
			inputUser: filmoteka.User{
				Username: "test",
				Password: "test",
				Is_admin: true,
			},
			mockBehaivior: func(s *mock_service.MockAuthorization, user filmoteka.User) {
				s.EXPECT().CreateUser(user).Return(1, errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"service failure"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			testCase.mockBehaivior(auth, testCase.inputUser)

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			//Test server
			mux := http.NewServeMux()
			mux.HandleFunc("/auth/sign-up", handler.handleSignUp)

			//Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-up",
				bytes.NewBufferString(testCase.inputBody))

			//Perform Request
			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}

func TestHandler_handleSignIn(t *testing.T) {
	type mockBehavior func(s *mock_service.MockAuthorization, username, password string)

	testTable := []struct {
		name                string
		inputBody           string
		username            string
		password            string
		mockBehavior        mockBehavior
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username":"test", "password":"test"}`,
			username:  "test",
			password:  "test",
			mockBehavior: func(s *mock_service.MockAuthorization, username, password string) {
				s.EXPECT().GenerateToken(username, password).Return("testtoken", nil)
			},
			expectedStatusCode:  200,
			expectedRequestBody: `{"token":"testtoken"}`,
		},
		{
			name:                "Empty Fields",
			inputBody:           `{"password":"test"}`,
			expectedStatusCode:  400,
			expectedRequestBody: `{"error":"Username and password are required"}`,
		},
		{
			name:      "Service Failure",
			inputBody: `{"username":"test", "password":"test"}`,
			username:  "test",
			password:  "test",
			mockBehavior: func(s *mock_service.MockAuthorization, username, password string) {
				s.EXPECT().GenerateToken(username, password).Return("", errors.New("service failure"))
			},
			expectedStatusCode:  500,
			expectedRequestBody: `{"error":"service failure"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Init Deps
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_service.NewMockAuthorization(c)
			if testCase.mockBehavior != nil {
				testCase.mockBehavior(auth, testCase.username, testCase.password)
			}

			services := &service.Service{Authorization: auth}
			handler := NewHandler(services)

			//Test server
			mux := http.NewServeMux()
			mux.HandleFunc("/auth/sign-in", handler.handleSignIn)

			//Test request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/auth/sign-in",
				bytes.NewBufferString(testCase.inputBody))

			//Perform Request
			mux.ServeHTTP(w, req)

			actual := strings.TrimSpace(w.Body.String())
			expected := strings.TrimSpace(testCase.expectedRequestBody)

			//Asserts
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, expected, actual)

		})
	}
}
