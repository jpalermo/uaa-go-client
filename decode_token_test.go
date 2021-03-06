package uaa_go_client_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/cloudfoundry-incubator/uaa-go-client"
	"github.com/cloudfoundry-incubator/uaa-go-client/config"
	"github.com/cloudfoundry-incubator/uaa-go-client/fakes"

	"github.com/dgrijalva/jwt-go"
	"code.cloudfoundry.org/clock/fakeclock"
	"code.cloudfoundry.org/lager"
	"code.cloudfoundry.org/lager/lagertest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("DecodeToken", func() {
	var (
		client            uaa_go_client.Client
		fakeSigningMethod *fakes.FakeSigningMethod
		// fakeUaaKeyFetcher *fakes.FakeUaaKeyFetcher
		signedKey      string
		UserPrivateKey string
		UAAPublicKey   string
		logger         lager.Logger

		token *jwt.Token
		err   error
	)

	verifyErrorType := func(err error, errorType uint32, message string) {
		validationError, ok := err.(*jwt.ValidationError)
		Expect(ok).To(BeTrue())
		Expect(validationError.Errors & errorType).To(Equal(errorType))
		Expect(err.Error()).To(Equal(message))
	}

	BeforeEach(func() {
		UserPrivateKey = "UserPrivateKey"
		UAAPublicKey = "UAAPublicKey"
		logger = lagertest.NewTestLogger("test")

		fakeSigningMethod = &fakes.FakeSigningMethod{}
		fakeSigningMethod.AlgStub = func() string {
			return "FAST"
		}
		fakeSigningMethod.SignStub = func(signingString string, key interface{}) (string, error) {
			signature := jwt.EncodeSegment([]byte(signingString + "SUPERFAST"))
			return signature, nil
		}
		fakeSigningMethod.VerifyStub = func(signingString, signature string, key interface{}) (err error) {
			if signature != jwt.EncodeSegment([]byte(signingString+"SUPERFAST")) {
				return errors.New("Signature is invalid")
			}

			return nil
		}

		jwt.RegisterSigningMethod("FAST", func() jwt.SigningMethod {
			return fakeSigningMethod
		})

		header := map[string]interface{}{
			"alg": "FAST",
		}

		alg := "FAST"
		signingMethod := jwt.GetSigningMethod(alg)
		token = jwt.New(signingMethod)
		token.Header = header

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

		client, err = uaa_go_client.NewClient(logger, cfg, clock)
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())

	})

	Describe("DecodeToken", func() {
		Context("when the token is valid", func() {
			BeforeEach(func() {
				claims := map[string]interface{}{
					"exp":   3404281214,
					"scope": []string{"route.advertise"},
				}
				token.Claims = claims

				signedKey, err = token.SignedString([]byte(UserPrivateKey))
				Expect(err).NotTo(HaveOccurred())

				server.AppendHandlers(
					getSuccessKeyFetchHandler(ValidPemPublicKey),
					getSuccessKeyFetchHandler(ValidPemPublicKey),
				)
			})

			It("caches the UAA public key", func() {
				err := client.DecodeToken("bearer "+signedKey, "route.advertise")
				Expect(err).NotTo(HaveOccurred())
				err = client.DecodeToken("bearer "+signedKey, "route.advertise")
				Expect(err).NotTo(HaveOccurred())

				Expect(len(server.ReceivedRequests())).To(Equal(1))
			})

			It("does not return an error", func() {
				err := client.DecodeToken("bearer "+signedKey, "route.advertise")
				Expect(err).NotTo(HaveOccurred())
			})

			It("does not return an error if the token type string is capitalized", func() {
				err := client.DecodeToken("Bearer "+signedKey, "route.advertise")
				Expect(err).NotTo(HaveOccurred())
			})

			It("does not return an error if the token type string is uppercase", func() {
				err := client.DecodeToken("BEARER "+signedKey, "route.advertise")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when a token is not valid", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					getSuccessKeyFetchHandler(ValidPemPublicKey),
				)
			})

			It("returns an error if the user token is not signed", func() {
				err = client.DecodeToken("bearer not-a-signed-token", "not a permission")
				Expect(err).To(HaveOccurred())
				verifyErrorType(err, jwt.ValidationErrorMalformed, "token contains an invalid number of segments")
				Expect(len(server.ReceivedRequests())).To(Equal(1))
			})

			It("returns an invalid token format when there is no token type", func() {
				err = client.DecodeToken("has-no-token-type", "not a permission")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid token format"))
				Expect(len(server.ReceivedRequests())).To(Equal(0))
			})

			It("returns an invalid token type when type is not bearer", func() {
				err = client.DecodeToken("basic some-auth", "not a permission")

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid token type: basic"))
				Expect(len(server.ReceivedRequests())).To(Equal(0))
			})
		})

		Context("when signature is invalid", func() {
			BeforeEach(func() {
				fakeSigningMethod.VerifyReturns(errors.New("invalid signature"))

				claims := map[string]interface{}{
					"exp":   3404281214,
					"scope": []string{"route.advertise"},
				}
				token.Claims = claims

				signedKey, err = token.SignedString([]byte(UserPrivateKey))
				Expect(err).NotTo(HaveOccurred())
				signedKey = "bearer " + signedKey
			})

			Context("uaa returns a verification key", func() {
				BeforeEach(func() {
					server.AppendHandlers(
						getSuccessKeyFetchHandler(ValidPemPublicKey),
						getSuccessKeyFetchHandler(ValidPemPublicKey),
					)
				})
				It("refreshes the key and returns an invalid signature error", func() {
					err := client.DecodeToken(signedKey, "route.advertise")
					Expect(err).To(HaveOccurred())

					Expect(len(server.ReceivedRequests())).To(Equal(2))
					verifyErrorType(err, jwt.ValidationErrorSignatureInvalid, "invalid signature")
				})
			})

			Context("when uaa returns an error", func() {
				BeforeEach(func() {
					server.AppendHandlers(
						getSuccessKeyFetchHandler(ValidPemPublicKey),
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET", TokenKeyEndpoint),
							ghttp.RespondWith(http.StatusGatewayTimeout, "booom"),
						),
					)
				})

				It("tries to refresh key and returns the uaa error", func() {
					err := client.DecodeToken(signedKey, "route.advertise")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("http-error-fetching-key"))
					Expect(len(server.ReceivedRequests())).To(Equal(2))
				})
			})
		})

		Context("when verification key needs to be refreshed to validate the signature", func() {
			BeforeEach(func() {
				fakeSigningMethod.VerifyStub = func(signingString string, signature string, key interface{}) error {
					switch k := key.(type) {
					case []byte:
						if string(k) == PemDecodedKey {
							return nil
						}
						return errors.New("invalid signature")
					default:
						return errors.New("invalid signature")
					}
				}
				claims := map[string]interface{}{
					"exp":   3404281214,
					"scope": []string{"route.advertise"},
				}
				token.Claims = claims

				signedKey, err = token.SignedString([]byte(UserPrivateKey))
				Expect(err).NotTo(HaveOccurred())
				signedKey = "bearer " + signedKey
			})

			Context("when a successful fetch happens", func() {
				BeforeEach(func() {
					server.AppendHandlers(
						getSuccessKeyFetchHandler(InvalidPemPublicKey),
						getSuccessKeyFetchHandler(ValidPemPublicKey),
					)
				})

				It("fetches new key and then validates the token", func() {
					err := client.DecodeToken(signedKey, "route.advertise")
					Expect(err).NotTo(HaveOccurred())
					Expect(len(server.ReceivedRequests())).To(Equal(2))
				})
			})

			Context("with multiple concurrent clients", func() {
				Context("when new key applies to all clients", func() {
					BeforeEach(func() {
						server.AppendHandlers(
							getSuccessKeyFetchHandler(ValidPemPublicKey),
							getSuccessKeyFetchHandler(ValidPemPublicKey),
						)
					})

					It("fetches new key and then validates the token", func() {
						wg := sync.WaitGroup{}
						for i := 0; i < 2; i++ {
							wg.Add(1)
							go func(wg *sync.WaitGroup) {
								defer GinkgoRecover()
								defer wg.Done()
								err := client.DecodeToken(signedKey, "route.advertise")
								Expect(err).NotTo(HaveOccurred())
							}(&wg)
						}
						wg.Wait()
						Expect(len(server.ReceivedRequests())).To(BeNumerically(">=", 1))
					})
				})

				Context("when new key applies to only one client and not others", func() {
					var (
						keyChannel      chan string
						expectErrorChan chan bool
					)

					BeforeEach(func() {
						keyChannel = make(chan string)
						expectErrorChan = make(chan bool)

						successHandler := func(w http.ResponseWriter, req *http.Request) {
							key := <-keyChannel
							w.Write([]byte(fmt.Sprintf("{\"alg\":\"alg\", \"value\": \"%s\" }", key)))
						}

						failureHandler := func(w http.ResponseWriter, req *http.Request) {
							w.WriteHeader(http.StatusInternalServerError)
							w.Write([]byte(""))
						}

						server.AppendHandlers(
							ghttp.CombineHandlers(
								failureHandler,
							),
							ghttp.CombineHandlers(
								successHandler,
							),
						)
					})

					AfterEach(func() {
						close(keyChannel)
						close(expectErrorChan)
					})

					It("fetches new key and validates the token", func() {
						wg := sync.WaitGroup{}
						for i := 0; i < 2; i++ {
							wg.Add(1)
							go func(wg *sync.WaitGroup) {
								defer GinkgoRecover()
								defer wg.Done()
								err := client.DecodeToken(signedKey, "route.advertise")
								select {
								case fail := <-expectErrorChan:
									if fail {
										Expect(err).To(HaveOccurred())
										Expect(err.Error()).To(Equal("http-error-fetching-key"))
									} else {
										Expect(err).NotTo(HaveOccurred())
									}
								}
							}(&wg)
						}
						// Error expected due to internal server error from UAA
						expectErrorChan <- true

						keyChannel <- ValidPemPublicKey

						// retrieved valid pem key from UAA, no error expected
						expectErrorChan <- false

						wg.Wait()
						Expect(len(server.ReceivedRequests())).To(Equal(2))
					})
				})
			})
		})

		Context("expired time", func() {
			BeforeEach(func() {
				claims := map[string]interface{}{
					"exp": time.Now().Unix() - 5,
				}
				token.Claims = claims

				signedKey, err = token.SignedString([]byte(UserPrivateKey))
				Expect(err).NotTo(HaveOccurred())

				signedKey = "bearer " + signedKey
				server.AppendHandlers(
					getSuccessKeyFetchHandler(ValidPemPublicKey),
				)
			})

			It("returns an error if the token is expired", func() {
				err = client.DecodeToken(signedKey, "route.advertise")
				Expect(err).To(HaveOccurred())
				verifyErrorType(err, jwt.ValidationErrorExpired, "token is expired")
			})
		})

		Context("permissions", func() {
			BeforeEach(func() {
				claims := map[string]interface{}{
					"exp":   time.Now().Unix() + 50000000,
					"scope": []string{"route.foo"},
				}
				token.Claims = claims

				signedKey, err = token.SignedString([]byte(UserPrivateKey))
				Expect(err).NotTo(HaveOccurred())

				signedKey = "bearer " + signedKey
				server.AppendHandlers(
					getSuccessKeyFetchHandler(ValidPemPublicKey),
				)
			})

			It("returns an error if the the user does not have requested permissions", func() {
				err = client.DecodeToken(signedKey, "route.my-permissions", "some.other.scope")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Token does not have 'route.my-permissions', 'some.other.scope' scope"))
			})
		})

	})
})
