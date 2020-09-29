package fakes

import "sync"

type RootManager struct {
	SetupCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			ExistingRoot string
			SdkLocation  string
		}
		Returns struct {
			Root string
			Err  error
		}
		Stub func(string, string) (string, error)
	}
}

func (f *RootManager) Setup(param1 string, param2 string) (string, error) {
	f.SetupCall.Lock()
	defer f.SetupCall.Unlock()
	f.SetupCall.CallCount++
	f.SetupCall.Receives.ExistingRoot = param1
	f.SetupCall.Receives.SdkLocation = param2
	if f.SetupCall.Stub != nil {
		return f.SetupCall.Stub(param1, param2)
	}
	return f.SetupCall.Returns.Root, f.SetupCall.Returns.Err
}
