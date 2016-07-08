// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/cloudfoundry-incubator/uaa-go-client"
	"github.com/cloudfoundry-incubator/uaa-go-client/schema"
)

type FakeClient struct {
	FetchTokenStub        func(forceUpdate bool) (*schema.Token, error)
	fetchTokenMutex       sync.RWMutex
	fetchTokenArgsForCall []struct {
		forceUpdate bool
	}
	fetchTokenReturns struct {
		result1 *schema.Token
		result2 error
	}
	FetchKeyStub        func() (string, error)
	fetchKeyMutex       sync.RWMutex
	fetchKeyArgsForCall []struct{}
	fetchKeyReturns     struct {
		result1 string
		result2 error
	}
	DecodeTokenStub        func(uaaToken string, desiredPermissions ...string) error
	decodeTokenMutex       sync.RWMutex
	decodeTokenArgsForCall []struct {
		uaaToken           string
		desiredPermissions []string
	}
	decodeTokenReturns struct {
		result1 error
	}
	RegisterOauthClientStub        func(*schema.OauthClient) (*schema.OauthClient, error)
	registerOauthClientMutex       sync.RWMutex
	registerOauthClientArgsForCall []struct {
		arg1 *schema.OauthClient
	}
	registerOauthClientReturns struct {
		result1 *schema.OauthClient
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) FetchToken(forceUpdate bool) (*schema.Token, error) {
	fake.fetchTokenMutex.Lock()
	fake.fetchTokenArgsForCall = append(fake.fetchTokenArgsForCall, struct {
		forceUpdate bool
	}{forceUpdate})
	fake.recordInvocation("FetchToken", []interface{}{forceUpdate})
	fake.fetchTokenMutex.Unlock()
	if fake.FetchTokenStub != nil {
		return fake.FetchTokenStub(forceUpdate)
	} else {
		return fake.fetchTokenReturns.result1, fake.fetchTokenReturns.result2
	}
}

func (fake *FakeClient) FetchTokenCallCount() int {
	fake.fetchTokenMutex.RLock()
	defer fake.fetchTokenMutex.RUnlock()
	return len(fake.fetchTokenArgsForCall)
}

func (fake *FakeClient) FetchTokenArgsForCall(i int) bool {
	fake.fetchTokenMutex.RLock()
	defer fake.fetchTokenMutex.RUnlock()
	return fake.fetchTokenArgsForCall[i].forceUpdate
}

func (fake *FakeClient) FetchTokenReturns(result1 *schema.Token, result2 error) {
	fake.FetchTokenStub = nil
	fake.fetchTokenReturns = struct {
		result1 *schema.Token
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) FetchKey() (string, error) {
	fake.fetchKeyMutex.Lock()
	fake.fetchKeyArgsForCall = append(fake.fetchKeyArgsForCall, struct{}{})
	fake.recordInvocation("FetchKey", []interface{}{})
	fake.fetchKeyMutex.Unlock()
	if fake.FetchKeyStub != nil {
		return fake.FetchKeyStub()
	} else {
		return fake.fetchKeyReturns.result1, fake.fetchKeyReturns.result2
	}
}

func (fake *FakeClient) FetchKeyCallCount() int {
	fake.fetchKeyMutex.RLock()
	defer fake.fetchKeyMutex.RUnlock()
	return len(fake.fetchKeyArgsForCall)
}

func (fake *FakeClient) FetchKeyReturns(result1 string, result2 error) {
	fake.FetchKeyStub = nil
	fake.fetchKeyReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DecodeToken(uaaToken string, desiredPermissions ...string) error {
	fake.decodeTokenMutex.Lock()
	fake.decodeTokenArgsForCall = append(fake.decodeTokenArgsForCall, struct {
		uaaToken           string
		desiredPermissions []string
	}{uaaToken, desiredPermissions})
	fake.recordInvocation("DecodeToken", []interface{}{uaaToken, desiredPermissions})
	fake.decodeTokenMutex.Unlock()
	if fake.DecodeTokenStub != nil {
		return fake.DecodeTokenStub(uaaToken, desiredPermissions...)
	} else {
		return fake.decodeTokenReturns.result1
	}
}

func (fake *FakeClient) DecodeTokenCallCount() int {
	fake.decodeTokenMutex.RLock()
	defer fake.decodeTokenMutex.RUnlock()
	return len(fake.decodeTokenArgsForCall)
}

func (fake *FakeClient) DecodeTokenArgsForCall(i int) (string, []string) {
	fake.decodeTokenMutex.RLock()
	defer fake.decodeTokenMutex.RUnlock()
	return fake.decodeTokenArgsForCall[i].uaaToken, fake.decodeTokenArgsForCall[i].desiredPermissions
}

func (fake *FakeClient) DecodeTokenReturns(result1 error) {
	fake.DecodeTokenStub = nil
	fake.decodeTokenReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) RegisterOauthClient(arg1 *schema.OauthClient) (*schema.OauthClient, error) {
	fake.registerOauthClientMutex.Lock()
	fake.registerOauthClientArgsForCall = append(fake.registerOauthClientArgsForCall, struct {
		arg1 *schema.OauthClient
	}{arg1})
	fake.recordInvocation("RegisterOauthClient", []interface{}{arg1})
	fake.registerOauthClientMutex.Unlock()
	if fake.RegisterOauthClientStub != nil {
		return fake.RegisterOauthClientStub(arg1)
	} else {
		return fake.registerOauthClientReturns.result1, fake.registerOauthClientReturns.result2
	}
}

func (fake *FakeClient) RegisterOauthClientCallCount() int {
	fake.registerOauthClientMutex.RLock()
	defer fake.registerOauthClientMutex.RUnlock()
	return len(fake.registerOauthClientArgsForCall)
}

func (fake *FakeClient) RegisterOauthClientArgsForCall(i int) *schema.OauthClient {
	fake.registerOauthClientMutex.RLock()
	defer fake.registerOauthClientMutex.RUnlock()
	return fake.registerOauthClientArgsForCall[i].arg1
}

func (fake *FakeClient) RegisterOauthClientReturns(result1 *schema.OauthClient, result2 error) {
	fake.RegisterOauthClientStub = nil
	fake.registerOauthClientReturns = struct {
		result1 *schema.OauthClient
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.fetchTokenMutex.RLock()
	defer fake.fetchTokenMutex.RUnlock()
	fake.fetchKeyMutex.RLock()
	defer fake.fetchKeyMutex.RUnlock()
	fake.decodeTokenMutex.RLock()
	defer fake.decodeTokenMutex.RUnlock()
	fake.registerOauthClientMutex.RLock()
	defer fake.registerOauthClientMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ uaa_go_client.Client = new(FakeClient)
