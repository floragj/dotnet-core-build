package fakes

import "sync"

type ProjectParser struct {
	ASPNetIsRequiredCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			Bool bool
		}
		Stub func(string) bool
	}
	NPMIsRequiredCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			Bool bool
		}
		Stub func(string) bool
	}
	NodeIsRequiredCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			Bool bool
		}
		Stub func(string) bool
	}
	ParseVersionCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Path string
		}
		Returns struct {
			Version string
			Err     error
		}
		Stub func(string) (string, error)
	}
}

func (f *ProjectParser) ASPNetIsRequired(param1 string) bool {
	f.ASPNetIsRequiredCall.Lock()
	defer f.ASPNetIsRequiredCall.Unlock()
	f.ASPNetIsRequiredCall.CallCount++
	f.ASPNetIsRequiredCall.Receives.Path = param1
	if f.ASPNetIsRequiredCall.Stub != nil {
		return f.ASPNetIsRequiredCall.Stub(param1)
	}
	return f.ASPNetIsRequiredCall.Returns.Bool
}
func (f *ProjectParser) NPMIsRequired(param1 string) bool {
	f.NPMIsRequiredCall.Lock()
	defer f.NPMIsRequiredCall.Unlock()
	f.NPMIsRequiredCall.CallCount++
	f.NPMIsRequiredCall.Receives.Path = param1
	if f.NPMIsRequiredCall.Stub != nil {
		return f.NPMIsRequiredCall.Stub(param1)
	}
	return f.NPMIsRequiredCall.Returns.Bool
}
func (f *ProjectParser) NodeIsRequired(param1 string) bool {
	f.NodeIsRequiredCall.Lock()
	defer f.NodeIsRequiredCall.Unlock()
	f.NodeIsRequiredCall.CallCount++
	f.NodeIsRequiredCall.Receives.Path = param1
	if f.NodeIsRequiredCall.Stub != nil {
		return f.NodeIsRequiredCall.Stub(param1)
	}
	return f.NodeIsRequiredCall.Returns.Bool
}
func (f *ProjectParser) ParseVersion(param1 string) (string, error) {
	f.ParseVersionCall.Lock()
	defer f.ParseVersionCall.Unlock()
	f.ParseVersionCall.CallCount++
	f.ParseVersionCall.Receives.Path = param1
	if f.ParseVersionCall.Stub != nil {
		return f.ParseVersionCall.Stub(param1)
	}
	return f.ParseVersionCall.Returns.Version, f.ParseVersionCall.Returns.Err
}
