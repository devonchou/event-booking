package test

import (
	"encoding/json"
	"event-booking-api/app/domain/dao"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *ApiTestSuite) TestAddEvent() {
	tests := []struct {
		name           string
		eventName      string
		description    string
		location       string
		eventTime      string
		token          string
		expectedStatus int
		expectedUserId int
	}{
		{"SuccessAddEvent", "Test Event 3", "This is a test event 3", "Tokyo", "2024-08-26T12:00:00Z", suite.user1Token, http.StatusCreated, 2},
		{"FailureMissingEventName", "", "This is a test event 3", "Tokyo", "2024-08-26T12:00:00Z", suite.user1Token, http.StatusBadRequest, 0},
		{"FailureMissingDescription", "Test Event 3", "", "Tokyo", "2024-08-26T12:00:00Z", suite.user1Token, http.StatusBadRequest, 0},
		{"FailureMissingLocation", "Test Event 3", "This is a test event 3", "", "2024-08-26T12:00:00Z", suite.user1Token, http.StatusBadRequest, 0},
		{"FailureMissingEventTime", "Test Event 3", "This is a test event 3", "Tokyo", "", suite.user1Token, http.StatusBadRequest, 0},
		{"FailureMissingToken", "Test Event 3", "This is a test event 3", "Tokyo", "2024-08-26T12:00:00Z", "", http.StatusUnauthorized, 0},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			payloads := fmt.Sprintf(`{"name": "%s", "description": "%s", "location": "%s", "event_time": "%s"}`, tt.eventName, tt.description, tt.location, tt.eventTime)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/api/events", strings.NewReader(payloads))
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusCreated {
				return
			}

			var actualEventTime time.Time
			err := suite.dbClient.
				QueryRow("SELECT event_time FROM events WHERE name = ? AND description = ? AND location = ? AND user_id = ?", tt.eventName, tt.description, tt.location, tt.expectedUserId).
				Scan(&actualEventTime)
			assert.NoError(suite.T(), err)

			expectedEventTime, err := time.Parse(time.RFC3339, tt.eventTime)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), expectedEventTime, actualEventTime)
		})
	}
}

func (suite *ApiTestSuite) TestGetAllEvent() {
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{"SuccessGetAllEvent", http.StatusOK},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/api/events", nil)
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var response struct {
				ResponseKey     string              `json:"response_key"`
				ResponseMessage string              `json:"response_message"`
				Data            []dao.EventResponse `json:"data"`
			}

			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), 2, len(response.Data))
		})
	}
}

func (suite *ApiTestSuite) TestGetEventById() {
	tests := []struct {
		name           string
		eventId        int
		expectedStatus int
	}{
		{"SuccessGetEvent", 2, http.StatusOK},
		{"FailureNotFound", 3, http.StatusNotFound},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/events/%v", tt.eventId), nil)
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var response struct {
				ResponseKey     string            `json:"response_key"`
				ResponseMessage string            `json:"response_message"`
				Data            dao.EventResponse `json:"data"`
			}

			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), tt.eventId, response.Data.ID)
		})
	}
}

func (suite *ApiTestSuite) TestUpdateEventById() {
	tests := []struct {
		name           string
		eventId        int
		eventName      string
		description    string
		location       string
		eventTime      string
		token          string
		expectedStatus int
	}{
		{"SuccessUpdateEvent", 1, "Updated Event", "This is updated event", "Updated Location", "2024-08-28T12:00:00Z", suite.user1Token, http.StatusOK},
		{"FailureMissingToken", 1, "Updated Event", "This is updated event", "Updated Location", "2024-08-28T12:00:00Z", "", http.StatusUnauthorized},
		{"FailureNotTheEventOwner", 1, "Updated Event", "This is updated event", "Updated Location", "2024-08-28T12:00:00Z", suite.user2Token, http.StatusUnauthorized},
		{"FailureEventNotFound", 4, "Updated Event", "This is updated event", "Updated Location", "2024-08-28T12:00:00Z", suite.user2Token, http.StatusNotFound},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			payloads := fmt.Sprintf(`{"name": "%s", "description": "%s", "location": "%s", "event_time": "%s"}`, tt.eventName, tt.description, tt.location, tt.eventTime)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", fmt.Sprintf("/api/events/%v", tt.eventId), strings.NewReader(payloads))
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var actualName, actualDescription, actualLocation string
			var actualEventTime time.Time
			err := suite.dbClient.
				QueryRow("SELECT name, description, location, event_time FROM events WHERE id = ?", tt.eventId).
				Scan(&actualName, &actualDescription, &actualLocation, &actualEventTime)
			assert.NoError(suite.T(), err)

			expectedEventTime, err := time.Parse(time.RFC3339, tt.eventTime)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), tt.eventName, actualName)
			assert.Equal(suite.T(), tt.description, actualDescription)
			assert.Equal(suite.T(), tt.location, actualLocation)
			assert.Equal(suite.T(), expectedEventTime, actualEventTime)
		})
	}
}

func (suite *ApiTestSuite) TestDeleteEventById() {
	tests := []struct {
		name           string
		eventId        int
		token          string
		expectedStatus int
	}{
		{"SuccessDeleteEvent", 1, suite.user1Token, http.StatusOK},
		{"FailureMissingToken", 2, "", http.StatusUnauthorized},
		{"FailureNotTheEventOwner", 2, suite.user1Token, http.StatusUnauthorized},
		{"FailureEventNotFound", 4, suite.user2Token, http.StatusNotFound},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/events/%v", tt.eventId), nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var count int
			err := suite.dbClient.QueryRow("SELECT COUNT(*) FROM events WHERE id = ? AND deleted_at IS NULL", tt.eventId).Scan(&count)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), 0, count)
		})
	}
}

func (suite *ApiTestSuite) TestRegisterUserForEvent() {
	tests := []struct {
		name           string
		eventId        int
		token          string
		expectedStatus int
		expectedUserId int
	}{
		{"SuccessRegister", 2, suite.user1Token, http.StatusCreated, 2},
		{"FailureMissingToken", 2, "", http.StatusUnauthorized, 0},
		{"FailureEventNotFound", 4, suite.user2Token, http.StatusNotFound, 0},
		{"FailureConflict", 1, suite.user2Token, http.StatusConflict, 0},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", fmt.Sprintf("/api/events/%v/register", tt.eventId), nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusCreated {
				return
			}

			var count int
			err := suite.dbClient.QueryRow("SELECT COUNT(*) FROM registers WHERE event_id = ? AND user_id = ?", tt.eventId, tt.expectedUserId).Scan(&count)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), 1, count)
		})
	}
}

func (suite *ApiTestSuite) TestUnregisterUserForEvent() {
	tests := []struct {
		name           string
		eventId        int
		token          string
		expectedStatus int
		expectedUserId int
	}{
		{"SuccessUnregister", 1, suite.user2Token, http.StatusOK, 3},
		{"FailureMissingToken", 1, "", http.StatusUnauthorized, 0},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/events/%v/register", tt.eventId), nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var count int
			err := suite.dbClient.QueryRow("SELECT COUNT(*) FROM registers WHERE event_id = ? AND user_id = ? AND deleted_at IS NULL", tt.eventId, tt.expectedUserId).Scan(&count)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), 0, count)
		})
	}
}

func (suite *ApiTestSuite) TestGetAttendeesEmailById() {
	tests := []struct {
		name           string
		eventId        int
		token          string
		expectedStatus int
	}{
		{"SuccessGetAttendeesEmail", 1, suite.user1Token, http.StatusOK},
		{"FailureMissingToken", 1, "", http.StatusUnauthorized},
		{"FailureNotTheEventOwner", 1, suite.user2Token, http.StatusUnauthorized},
		{"FailureEventNotFound", 4, suite.user2Token, http.StatusNotFound},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/api/events/%v/attendees", tt.eventId), nil)
			if tt.token != "" {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
			}
			suite.app.ServeHTTP(w, req)

			assert.Equal(suite.T(), tt.expectedStatus, w.Code)

			if tt.expectedStatus != http.StatusOK {
				return
			}

			var response struct {
				ResponseKey     string   `json:"response_key"`
				ResponseMessage string   `json:"response_message"`
				Data            []string `json:"data"`
			}

			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(suite.T(), err)

			assert.Equal(suite.T(), 2, len(response.Data))
		})
	}
}
