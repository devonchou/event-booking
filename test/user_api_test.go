package test

import (
	"encoding/json"
	"event-booking-api/app/domain/dao"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func (suite *ApiTestSuite) TestAddUser() {
	tests := []struct {
		name           string
		email          string
		password       string
		roleId         int
		expectedStatus int
	}{
		{"SuccessAddAdmin", "admin2@example.com", "adminpass", 1, http.StatusCreated},
		{"SuccessAddUser", "user3@example.com", "userpass", 2, http.StatusCreated},
		{"FailureWrongEmailFormat", "wrongemail", "userpass", 2, http.StatusBadRequest},
		{"FailureMissingEmail", "", "userpass", 2, http.StatusBadRequest},
		{"FailureMissingPassword", "user3@example.com", "", 2, http.StatusBadRequest},
		{"FailureEmailExist", "user1@example.com", "userpass", 2, http.StatusConflict},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			payloads := fmt.Sprintf(`{"email": "%s", "password": "%s", "role_id": %v}`, tt.email, tt.password, tt.roleId)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/users", strings.NewReader(payloads))
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusCreated {
				return
			}

			var actualPasswordHash string
			err := suite.dbClient.QueryRow("SELECT password FROM users WHERE email = ? AND role_id = ?", tt.email, tt.roleId).Scan(&actualPasswordHash)
			assert.NoError(suite.T(), err)

			err = bcrypt.CompareHashAndPassword([]byte(actualPasswordHash), []byte(tt.password))
			assert.NoError(suite.T(), err)
		})
	}
}

func (suite *ApiTestSuite) TestGetAllUser() {
	tests := []struct {
		name           string
		token          string
		expectedStatus int
	}{
		{"SuccessGetAllUser", suite.adminToken, http.StatusOK},
		{"FailureMissingToken", "", http.StatusUnauthorized},
		{"FailureNotTheAdmin", suite.user1Token, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/users", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var response struct {
				ResponseKey     string             `json:"response_key"`
				ResponseMessage string             `json:"response_message"`
				Data            []dao.UserResponse `json:"data"`
			}

			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), 3, len(response.Data))
		})
	}
}

func (suite *ApiTestSuite) TestGetUserById() {
	tests := []struct {
		name           string
		userId         int
		token          string
		expectedStatus int
	}{
		{"SuccessGetUser", 2, suite.user1Token, http.StatusOK},
		{"FailureMissingToken", 2, "", http.StatusUnauthorized},
		{"FailureNotTheUser", 2, suite.user2Token, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/users/%v", tt.userId), nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var response struct {
				ResponseKey     string           `json:"response_key"`
				ResponseMessage string           `json:"response_message"`
				Data            dao.UserResponse `json:"data"`
			}

			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), tt.userId, response.Data.ID)
		})
	}
}

func (suite *ApiTestSuite) TestUpdateUserById() {
	tests := []struct {
		name             string
		userId           int
		email            string
		password         string
		token            string
		expectedStatus   int
		expectedEmail    string
		expectedPassword string
	}{
		{"SuccessUpdateUserEmail", 2, "update@example.com", "", suite.user1Token, http.StatusOK, "update@example.com", "userpass"},
		{"SuccessUpdateUserPassword", 3, "", "updatepass", suite.user2Token, http.StatusOK, "user2@example.com", "updatepass"},
		{"FailureWrongEmailFormat", 2, "wrongemail", "updatepass", suite.user1Token, http.StatusBadRequest, "", ""},
		{"FailureMissingToken", 2, "update@example.com", "updatepass", "", http.StatusUnauthorized, "", ""},
		{"FailureNotTheUser", 2, "update@example.com", "updatepass", suite.user2Token, http.StatusUnauthorized, "", ""},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			payloads := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, tt.email, tt.password)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/users/%v", tt.userId), strings.NewReader(payloads))
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var actualEmail, actualPasswordHash string
			err := suite.dbClient.QueryRow("SELECT email, password FROM users WHERE id = ?", tt.userId).Scan(&actualEmail, &actualPasswordHash)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), tt.expectedEmail, actualEmail)

			err = bcrypt.CompareHashAndPassword([]byte(actualPasswordHash), []byte(tt.expectedPassword))
			assert.NoError(suite.T(), err)
		})
	}
}

func (suite *ApiTestSuite) TestDeleteUserById() {
	tests := []struct {
		name           string
		userId         int
		token          string
		expectedStatus int
	}{
		{"SuccessDeleteUser", 2, suite.user1Token, http.StatusOK},
		{"FailureMissingToken", 2, "", http.StatusUnauthorized},
		{"FailureNotTheUser", 2, suite.user2Token, http.StatusUnauthorized},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/users/%v", tt.userId), nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var count int
			err := suite.dbClient.QueryRow("SELECT COUNT(*) FROM users WHERE id = ? AND deleted_at IS NULL", tt.userId).Scan(&count)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), 0, count)
		})
	}
}

func (suite *ApiTestSuite) TestLoginUser() {
	tests := []struct {
		name           string
		email          string
		password       string
		expectedStatus int
	}{
		{"SuccessValidCredentials", "user1@example.com", "userpass", http.StatusOK},
		{"FailureWrongPassword", "user1@example.com", "wrongpass", http.StatusUnauthorized},
		{"FailureEmailNotFound", "wrongemail@example.com", "userpass", http.StatusUnauthorized},
		{"FailureWrongEmailFormat", "wrongemail", "userpass", http.StatusBadRequest},
		{"FailureEmptyPassword", "user1@example.com", "", http.StatusBadRequest},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			loginCredentials := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, tt.email, tt.password)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/users/login", strings.NewReader(loginCredentials))
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var response struct {
				ResponseKey     string `json:"response_key"`
				ResponseMessage string `json:"response_message"`
				Data            string `json:"data"`
			}

			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.NotEmpty(suite.T(), response.Data)
		})
	}
}
