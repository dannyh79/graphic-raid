package classroom_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClassroom(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Classroom Suite")
}
