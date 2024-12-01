package ratelimiter

import (
	"context"
	"errors"
	"net"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/felipeivanaga/go-expert-rate-limiter/internal/pkg/mocks"
	"github.com/felipeivanaga/go-expert-rate-limiter/internal/pkg/ratelimiter/strategies"
	"github.com/stretchr/testify/suite"
)

type RateLimiterTestSuite struct {
	suite.Suite
	StrategyMock        *mocks.RedisLimiterStrategyMock
	MaxRequestsPerIP    int
	MaxRequestsPerToken int
	TimeWindowMillis    int

	Limiter *RateLimiter
}

func TestRateLimiter(t *testing.T) {
	suite.Run(t, new(RateLimiterTestSuite))
}

func (s *RateLimiterTestSuite) SetupTest() {
	strategyMock := new(mocks.RedisLimiterStrategyMock)

	ipMaxReqs := 10
	tokenMaxReqs := 100
	timeWindowMillis := 1000

	s.StrategyMock = strategyMock
	s.MaxRequestsPerIP = ipMaxReqs
	s.MaxRequestsPerToken = tokenMaxReqs
	s.TimeWindowMillis = timeWindowMillis

	s.Limiter = NewRateLimiter(strategyMock, ipMaxReqs, tokenMaxReqs, timeWindowMillis)
}

func (s *RateLimiterTestSuite) clearMocks() {
	s.StrategyMock.ExpectedCalls = nil
}

func (s *RateLimiterTestSuite) TestRateLimiterByIP() {
	s.Run("Should allow request", func() {
		defer s.clearMocks()

		ctx := context.Background()
		r := httptest.NewRequest("GET", "/", nil)

		mockReq := strategies.RateLimiterRequest{
			Key:      net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]).String(),
			Limit:    int64(s.MaxRequestsPerIP),
			Duration: time.Duration(s.TimeWindowMillis) * time.Millisecond,
		}

		mockRes := strategies.RateLimiterResult{
			Result:    strategies.Allow,
			Limit:     int64(s.MaxRequestsPerIP),
			Total:     1,
			Remaining: int64(s.MaxRequestsPerIP) - 1,
			ExpiresAt: time.Now().Add(time.Duration(s.TimeWindowMillis) * time.Millisecond),
		}

		s.StrategyMock.On("Check", ctx, &mockReq).Return(&mockRes, nil)

		result, err := s.Limiter.Check(ctx, r)

		s.Nil(err)
		s.Equal(mockRes, *result)
		s.StrategyMock.AssertExpectations(s.T())
	})

	s.Run("Should deny request", func() {
		defer s.clearMocks()

		ctx := context.Background()
		r := httptest.NewRequest("GET", "/", nil)

		mockReq := strategies.RateLimiterRequest{
			Key:      net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]).String(),
			Limit:    int64(s.MaxRequestsPerIP),
			Duration: time.Duration(s.TimeWindowMillis) * time.Millisecond,
		}

		mockRes := strategies.RateLimiterResult{
			Result:    strategies.Deny,
			Limit:     int64(s.MaxRequestsPerIP),
			Total:     int64(s.MaxRequestsPerIP),
			Remaining: 0,
			ExpiresAt: time.Now().Add(time.Duration(s.TimeWindowMillis) * time.Millisecond),
		}

		s.StrategyMock.On("Check", ctx, &mockReq).Return(&mockRes, nil)

		result, err := s.Limiter.Check(ctx, r)

		s.Nil(err)
		s.Equal(mockRes, *result)
		s.StrategyMock.AssertExpectations(s.T())
	})

	s.Run("Should return error", func() {
		defer s.clearMocks()

		ctx := context.Background()
		r := httptest.NewRequest("GET", "/", nil)

		mockReq := strategies.RateLimiterRequest{
			Key:      net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]).String(),
			Limit:    int64(s.MaxRequestsPerIP),
			Duration: time.Duration(s.TimeWindowMillis) * time.Millisecond,
		}

		mockErr := errors.New("any-error")

		s.StrategyMock.On("Check", ctx, &mockReq).Return(nil, mockErr)

		result, err := s.Limiter.Check(ctx, r)

		s.Error(err)
		s.Nil(result)
		s.Equal(mockErr, err)
		s.StrategyMock.AssertExpectations(s.T())
	})
}

func (s *RateLimiterTestSuite) TestRateLimiterByToken() {
	apiKey := "any-api-key"

	s.Run("Should allow request", func() {
		defer s.clearMocks()

		ctx := context.Background()

		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("API_KEY", apiKey)

		mockReq := strategies.RateLimiterRequest{
			Key:      apiKey,
			Limit:    int64(s.MaxRequestsPerToken),
			Duration: time.Duration(s.TimeWindowMillis) * time.Millisecond,
		}

		mockRes := strategies.RateLimiterResult{
			Result:    strategies.Allow,
			Limit:     int64(s.MaxRequestsPerToken),
			Total:     1,
			Remaining: int64(s.MaxRequestsPerToken) - 1,
			ExpiresAt: time.Now().Add(time.Duration(s.TimeWindowMillis) * time.Millisecond),
		}

		s.StrategyMock.On("Check", ctx, &mockReq).Return(&mockRes, nil)

		result, err := s.Limiter.Check(ctx, r)

		s.Nil(err)
		s.Equal(mockRes, *result)
		s.StrategyMock.AssertExpectations(s.T())
	})

	s.Run("Should deny request", func() {
		defer s.clearMocks()

		ctx := context.Background()

		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("API_KEY", apiKey)

		mockReq := strategies.RateLimiterRequest{
			Key:      apiKey,
			Limit:    int64(s.MaxRequestsPerToken),
			Duration: time.Duration(s.TimeWindowMillis) * time.Millisecond,
		}

		mockRes := strategies.RateLimiterResult{
			Result:    strategies.Deny,
			Limit:     int64(s.MaxRequestsPerToken),
			Total:     int64(s.MaxRequestsPerToken),
			Remaining: 0,
			ExpiresAt: time.Now().Add(time.Duration(s.TimeWindowMillis) * time.Millisecond),
		}

		s.StrategyMock.On("Check", ctx, &mockReq).Return(&mockRes, nil)

		result, err := s.Limiter.Check(ctx, r)

		s.Nil(err)
		s.Equal(mockRes, *result)
		s.StrategyMock.AssertExpectations(s.T())
	})

	s.Run("Should return error", func() {
		defer s.clearMocks()

		ctx := context.Background()

		r := httptest.NewRequest("GET", "/", nil)
		r.Header.Set("API_KEY", apiKey)

		mockReq := strategies.RateLimiterRequest{
			Key:      apiKey,
			Limit:    int64(s.MaxRequestsPerToken),
			Duration: time.Duration(s.TimeWindowMillis) * time.Millisecond,
		}

		mockErr := errors.New("any-error")

		s.StrategyMock.On("Check", ctx, &mockReq).Return(nil, mockErr)

		result, err := s.Limiter.Check(ctx, r)

		s.Error(err)
		s.Nil(result)
		s.Equal(mockErr, err)
		s.StrategyMock.AssertExpectations(s.T())
	})
}
