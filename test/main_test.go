package test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestApiTestSuite(t *testing.T) {
	suite.Run(t, new(ApiTestSuite))
}
