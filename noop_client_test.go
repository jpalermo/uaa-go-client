package uaa_go_client_test

import (
	. "github.com/cloudfoundry-incubator/uaa-go-client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NoopUaaClient", func() {

	var client Client

	BeforeEach(func() {
		client = NewNoOpUaaClient()
	})

	Context("New", func() {
		It("returns a no-op token client", func() {
			Expect(client).NotTo(BeNil())
			Expect(client).To(BeAssignableToTypeOf(&NoOpUaaClient{}))
		})
	})

	Context("FetchToken", func() {
		It("returns an empty access token", func() {
			token, err := client.FetchToken(true)
			Expect(err).NotTo(HaveOccurred())
			Expect(token.AccessToken).To(BeEmpty())
		})

		It("returns an empty access token", func() {
			token, err := client.FetchToken(true)
			Expect(err).NotTo(HaveOccurred())
			Expect(token.AccessToken).To(BeEmpty())
		})
	})

	Context("FetchKey", func() {
		It("returns an empty token key", func() {
			key, err := client.FetchKey()
			Expect(err).NotTo(HaveOccurred())
			Expect(key).To(BeEmpty())
		})
	})

	Context("DecodeToken", func() {
		It("returns an empty decode", func() {
			decoded := client.DecodeToken("some token", "some perm")
			Expect(decoded).To(BeNil())
		})
	})

})
