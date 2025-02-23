package pgxfilter_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPgxfilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pgxfilter Suite")
}
