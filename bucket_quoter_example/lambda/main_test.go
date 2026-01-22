package main

import (
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lambda integ tests", func() {
	BeforeEach(func() {
		os.Setenv("DRY_RUN", "false")
		os.Setenv("API_SUBSCRIPTION_KEY", "897d9f58-6b42-4ca7-8229-2e04056490b7")
		os.Setenv("LIMITER_HOST", "https://localhost:8443")
	})

	When("trigger lambda", func() {
		It("should be successful", func() {
			// run quoter on your dev machine before running this test or mock quoter
			resp, err := basicHandler()
			Ω(resp.StatusCode).Should(Equal(http.StatusOK))
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	AfterEach(func() {
		os.Unsetenv("DRY_RUN")
		os.Unsetenv("API_SUBSCRIPTION_KEY")
		os.Unsetenv("LIMITER_HOST")
	})
})
