package middlewares

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/felipeivanaga/go-expert-rate-limiter/internal/pkg/mocks"
	"github.com/felipeivanaga/go-expert-rate-limiter/internal/pkg/ratelimiter/strategies"
	"github.com/felipeivanaga/go-expert-rate-limiter/internal/pkg/responsehandler"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RateLimiterMiddlewareTestSuite struct {
	suite.Suite
	ResponseHandler *responsehandler.WebResponseHandler
	Limiter         *mocks.RateLimiterMock
	Router          *chi.Mux

	Middleware RateLimiterMiddlewareInterface
}

func (s *RateLimiterMiddlewareTestSuite) SetupTest() {
	responseHandler := &responsehandler.WebResponseHandler{}
	limiterMock := new(mocks.RateLimiterMock)
	router := chi.NewRouter()

	s.ResponseHandler = responseHandler
	s.Limiter = limiterMock
	s.Router = router

	s.Middleware = NewRateLimiterMiddleware(responseHandler, limiterMock)
}

func (s *RateLimiterMiddlewareTestSuite) clearMocks() {
	s.Limiter.ExpectedCalls = nil
}

func TestRateLimiterMiddleware(t *testing.T) {
	suite.Run(t, new(RateLimiterMiddlewareTestSuite))
}

func (s *RateLimiterMiddlewareTestSuite) TestHandle() {
	later := time.Now().Add(10 * time.Second)

	s.Router.Use(s.Middleware.Handle)
	registerTestRoute(s.Router)

	s.Run("Should allow request when rate limiter allows", func() {
		defer s.clearMocks()

		server := httptest.NewServer(s.Router)
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, nil)

		s.Limiter.On("Check", mock.Anything, mock.Anything).Return(&strategies.RateLimiterResult{
			Result:    strategies.Allow,
			Limit:     10,
			Total:     1,
			Remaining: 9,
			ExpiresAt: later,
		}, nil)

		resp, _ := http.DefaultClient.Do(req)

		s.Equal(http.StatusOK, resp.StatusCode)
		s.Equal("10", resp.Header.Get("X-RateLimit-Limit"))
		s.Equal("9", resp.Header.Get("X-RateLimit-Remaining"))
		s.Equal(strconv.FormatInt(later.Unix(), 10), resp.Header.Get("X-RateLimit-Reset"))
	})

	s.Run("Should deny request when rate limiter denies", func() {
		defer s.clearMocks()

		server := httptest.NewServer(s.Router)
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, nil)

		s.Limiter.On("Check", mock.Anything, mock.Anything).Return(&strategies.RateLimiterResult{
			Result:    strategies.Deny,
			Limit:     10,
			Total:     10,
			Remaining: 0,
			ExpiresAt: later,
		}, nil)

		resp, _ := http.DefaultClient.Do(req)

		s.Equal(http.StatusTooManyRequests, resp.StatusCode)
		s.Equal("10", resp.Header.Get("X-RateLimit-Limit"))
		s.Equal("0", resp.Header.Get("X-RateLimit-Remaining"))
		s.Equal(strconv.FormatInt(later.Unix(), 10), resp.Header.Get("X-RateLimit-Reset"))
	})

	s.Run("Should respond with internal server error when rate limiter returns error", func() {
		defer s.clearMocks()

		server := httptest.NewServer(s.Router)
		defer server.Close()

		req, _ := http.NewRequest("GET", server.URL, nil)

		s.Limiter.On("Check", mock.Anything, mock.Anything).Return(nil, errors.New("any-error"))

		resp, _ := http.DefaultClient.Do(req)
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)

		s.Equal(http.StatusInternalServerError, resp.StatusCode)
		s.Contains(string(body), "error while checking rate limit")
	})
}

func registerTestRoute(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Hello World!"}`))
	})
}
