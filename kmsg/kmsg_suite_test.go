package kmsg_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKmsg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kmsg Suite")
}
