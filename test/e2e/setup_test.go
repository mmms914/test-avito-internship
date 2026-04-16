//go:build e2e
// +build e2e

package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	client     *http.Client
	baseURL    string
	adminToken string
	userToken  string
}

func (s *E2ETestSuite) checkStatus(resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		s.T().Logf("Expected status %d, got %d", expected, resp.StatusCode)
		s.T().Logf("Response body: %s", string(body))

		var errResp map[string]interface{}
		if json.Unmarshal(body, &errResp) == nil {
			if errDetail, ok := errResp["error"]; ok {
				s.T().Logf("Error details: %v", errDetail)
			}
		}
	}
}

func (s *E2ETestSuite) SetupSuite() {
	s.client = &http.Client{Timeout: 30 * time.Second}

	s.baseURL = os.Getenv("E2E_API_URL")
	if s.baseURL == "" {
		s.baseURL = "http://localhost:8081"
	}

	s.Require().Eventually(func() bool {
		resp, err := s.client.Get(s.baseURL + "/_info")
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}, 30*time.Second, 1*time.Second, "Service not ready")

	s.adminToken = s.getToken("admin")
	s.userToken = s.getToken("user")
}

func (s *E2ETestSuite) getToken(role string) string {
	reqBody := map[string]string{"role": role}
	body, _ := json.Marshal(reqBody)

	resp, err := s.client.Post(
		fmt.Sprintf("%s/dummyLogin", s.baseURL),
		"application/json",
		bytes.NewBuffer(body),
	)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.checkStatus(resp, http.StatusOK)
	s.Require().Equal(http.StatusOK, resp.StatusCode)

	var result struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	s.Require().NoError(err)
	s.Require().NotEmpty(result.Token)

	return result.Token
}

func (s *E2ETestSuite) doRequest(method, url string, token string, body interface{}) *http.Response {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBody)
	}

	var req *http.Request
	var err error
	if reqBody == nil {
		req, err = http.NewRequest(method, fmt.Sprintf("%s%s", s.baseURL, url), nil)
	} else {
		req, err = http.NewRequest(method, fmt.Sprintf("%s%s", s.baseURL, url), reqBody)
	}
	s.Require().NoError(err)

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.client.Do(req)
	s.Require().NoError(err)

	return resp
}

func (s *E2ETestSuite) TearDownSuite() {
	s.client.CloseIdleConnections()
}
