package pgxmql_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPgxMQL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PgxMQL Suite")
}
