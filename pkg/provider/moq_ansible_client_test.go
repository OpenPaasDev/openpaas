// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package provider

import (
	"github.com/OpenPaasDev/openpaas/pkg/ansible"
	"sync"
)

// Ensure, that MockAnsibleClient does implement ansible.Client.
// If this is not the case, regenerate this file with moq.
var _ ansible.Client = &MockAnsibleClient{}

// MockAnsibleClient is a mock implementation of ansible.Client.
//
//	func TestSomethingThatUsesClient(t *testing.T) {
//
//		// make and configure a mocked ansible.Client
//		mockedClient := &MockAnsibleClient{
//			RunFunc: func(playbookFile string, varFile string) error {
//				panic("mock out the Run method")
//			},
//		}
//
//		// use mockedClient in code that requires ansible.Client
//		// and then make assertions.
//
//	}
type MockAnsibleClient struct {
	// RunFunc mocks the Run method.
	RunFunc func(playbookFile string, varFile string) error

	// calls tracks calls to the methods.
	calls struct {
		// Run holds details about calls to the Run method.
		Run []struct {
			// PlaybookFile is the playbookFile argument value.
			PlaybookFile string
			// VarFile is the varFile argument value.
			VarFile string
		}
	}
	lockRun sync.RWMutex
}

// Run calls RunFunc.
func (mock *MockAnsibleClient) Run(playbookFile string, varFile string) error {
	callInfo := struct {
		PlaybookFile string
		VarFile      string
	}{
		PlaybookFile: playbookFile,
		VarFile:      varFile,
	}
	mock.lockRun.Lock()
	mock.calls.Run = append(mock.calls.Run, callInfo)
	mock.lockRun.Unlock()
	if mock.RunFunc == nil {
		var (
			errOut error
		)
		return errOut
	}
	return mock.RunFunc(playbookFile, varFile)
}

// RunCalls gets all the calls that were made to Run.
// Check the length with:
//
//	len(mockedClient.RunCalls())
func (mock *MockAnsibleClient) RunCalls() []struct {
	PlaybookFile string
	VarFile      string
} {
	var calls []struct {
		PlaybookFile string
		VarFile      string
	}
	mock.lockRun.RLock()
	calls = mock.calls.Run
	mock.lockRun.RUnlock()
	return calls
}
