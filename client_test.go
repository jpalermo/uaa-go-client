package uaa_go_client_test

import (
	"net/url"
	"strings"
	"time"

	"github.com/cf-routing/uaa-go-client"
	"github.com/cf-routing/uaa-go-client/config"
	"github.com/pivotal-golang/clock/fakeclock"
	"github.com/pivotal-golang/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("UAA Client", func() {
	const (
		DefaultMaxNumberOfRetries   = 3
		DefaultRetryInterval        = 15 * time.Second
		DefaultExpirationBufferTime = 30
	)

	BeforeEach(func() {
		forceUpdate = false
		cfg = &config.Config{
			MaxNumberOfRetries:    DefaultMaxNumberOfRetries,
			RetryInterval:         DefaultRetryInterval,
			ExpirationBufferInSec: DefaultExpirationBufferTime,
		}
		server = ghttp.NewServer()

		url, err := url.Parse(server.URL())
		Expect(err).ToNot(HaveOccurred())

		addr := strings.Split(url.Host, ":")

		cfg.UaaEndpoint = "http://" + addr[0] + ":" + addr[1]
		Expect(err).ToNot(HaveOccurred())

		cfg.ClientName = "client-name"
		cfg.ClientSecret = "client-secret"
		clock = fakeclock.NewFakeClock(time.Now())
		logger = lagertest.NewTestLogger("test")
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("uaa_go_client.NewClient", func() {
		Context("when all values are valid", func() {
			It("returns a token fetcher instance", func() {
				client, err := uaa_go_client.NewClient(logger, cfg, clock)
				Expect(err).NotTo(HaveOccurred())
				Expect(client).NotTo(BeNil())
			})
		})

		Context("when values are invalid", func() {
			Context("when oauth config is nil", func() {
				It("returns error", func() {
					client, err := uaa_go_client.NewClient(logger, nil, clock)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Configuration cannot be nil"))
					Expect(client).To(BeNil())
				})
			})

			Context("when oauth config client id is empty", func() {
				It("returns error", func() {
					config := &config.Config{
						UaaEndpoint:  "http://some.url:80",
						ClientName:   "",
						ClientSecret: "client-secret",
					}
					client, err := uaa_go_client.NewClient(logger, config, clock)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("OAuth Client ID cannot be empty"))
					Expect(client).To(BeNil())
				})
			})

			Context("when oauth config client secret is empty", func() {
				It("returns error", func() {
					config := &config.Config{
						UaaEndpoint:  "http://some.url:80",
						ClientName:   "client-name",
						ClientSecret: "",
					}
					client, err := uaa_go_client.NewClient(logger, config, clock)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("OAuth Client Secret cannot be empty"))
					Expect(client).To(BeNil())
				})
			})

			Context("when oauth config tokenendpoint is empty", func() {
				It("returns error", func() {
					config := &config.Config{
						UaaEndpoint:  "",
						ClientName:   "client-name",
						ClientSecret: "client-secret",
					}
					client, err := uaa_go_client.NewClient(logger, config, clock)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("UAA endpoint cannot be empty"))
					Expect(client).To(BeNil())
				})
			})

			Context("when token fetcher config's max number of retries is zero", func() {
				It("creates the client", func() {
					config := &config.Config{
						UaaEndpoint:           "http://some.url:80",
						MaxNumberOfRetries:    0,
						RetryInterval:         2 * time.Second,
						ExpirationBufferInSec: 30,
						ClientName:            "client-name",
						ClientSecret:          "client-secret",
					}
					client, err := uaa_go_client.NewClient(logger, config, clock)
					Expect(err).NotTo(HaveOccurred())
					Expect(client).NotTo(BeNil())
				})
			})

			Context("when token fetcher config's expiration buffer time is negative", func() {
				It("sets the expiration buffer time to the default value", func() {
					config := &config.Config{
						MaxNumberOfRetries:    3,
						RetryInterval:         2 * time.Second,
						ExpirationBufferInSec: -1,
						UaaEndpoint:           "http://some.url:80",
						ClientName:            "client-name",
						ClientSecret:          "client-secret",
					}
					client, err := uaa_go_client.NewClient(logger, config, clock)
					Expect(err).NotTo(HaveOccurred())
					Expect(client).NotTo(BeNil())
				})
			})
		})
	})
})