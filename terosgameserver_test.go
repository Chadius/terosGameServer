package terosGameServer_test

import (
	"github.com/chadius/terosGameServer/terosgamerulesfakes"
	. "gopkg.in/check.v1"
	"io"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type UsePackageTestSuite struct{}

func (suite *UsePackageTestSuite) TestWhenFilesAreSupplied_ThenCallPackage(checker *C) {
	// Setup
	squaddieByteStream := []byte(`squaddies go here`)
	powerByteStream := []byte(`powers go here`)
	scriptByteStream := []byte(`script goes here`)
	expectedResponse := []byte(`rules responded`)

	// TODO turn the package terosGameRules into an Interface
	fakeGameRules := terosgamerulesfakes.FakeRulesStrategy{}
	fakeGameRules.ProcessAttackScriptStub = func(scriptFileHandle, squaddieFileHandle, powerFileHandle io.Reader, output io.Writer) error {
		output.Write(expectedResponse)
		return nil
	}

	// Act
	actualResponse := terosserver.ProcessAttackScript(scriptByteStream, squaddieByteStream, powerByteStream)

	// Assert
	actualScriptFileHandle, actualSquaddieFileHandle, actualPowerFileHandle, _ := fakeGameRules.ProcessAttackScriptArgsForCall(0)
	checker.Assert(actualScriptFileHandle, Equals, scriptByteStream)
	checker.Assert(actualSquaddieFileHandle, Equals, squaddieByteStream)
	checker.Assert(actualPowerFileHandle, Equals, powerByteStream)

	checker.Assert(200, Equals, actualResponse.ResponseCode)
	checker.Assert(expectedResponse, Equals, actualResponse.Data)
}
