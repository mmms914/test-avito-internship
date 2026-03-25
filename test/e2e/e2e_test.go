//go:build e2e
// +build e2e

package e2e

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestE2E(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}
