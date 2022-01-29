package terosgameserver_test

import (
	"bytes"
	"github.com/chadius/terosgamerules"
	"github.com/chadius/terosgameserver/rpc/github.com/chadius/teros_game_server"
	"github.com/chadius/terosgameserver/rulesstrategyfakes"
	"github.com/chadius/terosgameserver/terosgameserver"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServerUsesPackageSuite(t *testing.T) {
	suite.Run(t, new(ServerUsesPackageSuite))
}

type ServerUsesPackageSuite struct {
	suite.Suite
	request           *http.Request
	responseRecorder  *httptest.ResponseRecorder
	server            teros_game_server.TwirpServer
	fakeGameRules     *rulesstrategyfakes.FakeRulesStrategy
	scriptData        []byte
	squaddieData      []byte
	powerData         []byte
	gameRulesResponse []byte
}

func (suite *ServerUsesPackageSuite) SetupTest() {
	suite.scriptData = []byte(`images go here`)
	suite.squaddieData = []byte(`formula goes here`)
	suite.powerData = []byte(`outputSettings go here`)
	requestBody := suite.getDataStream()
	suite.request = suite.generateProtobufRequest(requestBody)
	suite.responseRecorder = httptest.NewRecorder()

	suite.gameRulesResponse = []byte(`rules responded`)
	suite.fakeGameRules = suite.fakeGameRulesWithResponse(suite.gameRulesResponse)
	suite.server = suite.getServer()
}

func (suite *ServerUsesPackageSuite) fakeGameRulesWithResponse(expectedResponse []byte) *rulesstrategyfakes.FakeRulesStrategy {
	fakeTransformerStrategy := rulesstrategyfakes.FakeRulesStrategy{}
	fakeTransformerStrategy.ReplayBattleScriptStub = func(scriptFileHandle, squaddieFileHandle, powerFileHandle io.Reader, output io.Writer) error {
		output.Write(expectedResponse)
		return nil
	}
	return &fakeTransformerStrategy
}

func (suite *ServerUsesPackageSuite) getServer() teros_game_server.TwirpServer {
	server := terosgameserver.NewServer(suite.fakeGameRules)
	twirpServer := teros_game_server.NewTerosGameServerServer(server)
	return twirpServer
}

func (suite *ServerUsesPackageSuite) generateProtobufRequest(requestBody *bytes.Buffer) *http.Request {
	testRequest, newRequestErr := http.NewRequest(
		http.MethodPost,
		"/twirp/chadius.terosGameServer.TerosGameServer/ReplayBattleScript",
		requestBody,
	)
	require := require.New(suite.T())
	require.Nil(newRequestErr)
	testRequest.Header.Set("Content-Type", "application/protobuf")
	return testRequest
}

func (suite *ServerUsesPackageSuite) getDataStream() *bytes.Buffer {
	dataStream := &teros_game_server.DataStreams{
		ScriptData:   suite.scriptData,
		SquaddieData: suite.squaddieData,
		PowerData:    suite.powerData,
	}

	protobuf, protobufErr := proto.Marshal(dataStream)
	requestBody := bytes.NewBuffer(protobuf)

	require := require.New(suite.T())
	require.Nil(protobufErr)

	return requestBody
}

func (suite *ServerUsesPackageSuite) TestWhenClientMakesRequest_ResponseIsValid() {
	// Act
	suite.server.ServeHTTP(suite.responseRecorder, suite.request)

	// Assert
	response := suite.responseRecorder.Result()

	require := require.New(suite.T())
	require.Equal(200, response.StatusCode, "Status code is wrong")
}

func (suite *ServerUsesPackageSuite) TestWhenClientMakesRequest_PackageIsCalledWithInputData() {
	// Act
	suite.server.ServeHTTP(suite.responseRecorder, suite.request)

	// Assert
	response := suite.responseRecorder.Result()

	require := require.New(suite.T())
	require.Equal(200, response.StatusCode, "Status code is wrong")

	suite.requireFakePackageWasCalledWithExpectedData(require)
}

func (suite *ServerUsesPackageSuite) requireFakePackageWasCalledWithExpectedData(require *require.Assertions) {
	require.Equal(1, suite.fakeGameRules.ReplayBattleScriptCallCount())

	actualScriptDataByteStream, actualSquaddieDataByteStream, actualPowerDataByteStream, _ := suite.fakeGameRules.ReplayBattleScriptArgsForCall(0)

	actualScriptData, scriptReadErr := ioutil.ReadAll(actualScriptDataByteStream)
	require.Nil(scriptReadErr, "Error while reading script data from mock object")
	require.Equal(0, bytes.Compare(suite.scriptData, actualScriptData), "script given to mock object is different")

	actualSquaddieData, squaddieReadErr := ioutil.ReadAll(actualSquaddieDataByteStream)
	require.Nil(squaddieReadErr, "Error while reading squaddie data from mock object")
	require.Equal(0, bytes.Compare(suite.squaddieData, actualSquaddieData), "squaddie data to mock object is different")

	actualPowerData, powerReadErr := ioutil.ReadAll(actualPowerDataByteStream)
	require.Nil(powerReadErr, "Error while reading power data from mock object")
	require.Equal(0, bytes.Compare(suite.powerData, actualPowerData), "power data given to mock object is different")
}

func (suite *ServerUsesPackageSuite) TestWhenClientMakesRequest_ResponseIsUnmarshalled() {
	// Act
	suite.server.ServeHTTP(suite.responseRecorder, suite.request)

	// Assert
	require := require.New(suite.T())
	suite.requireResponseDataMatches(require)
}

func (suite *ServerUsesPackageSuite) requireResponseDataMatches(require *require.Assertions) {
	output := &teros_game_server.Results{}
	unmarshalErr := proto.Unmarshal(suite.responseRecorder.Body.Bytes(), output)
	require.Nil(unmarshalErr, "Error while unmarshalling response body")
	require.Equal(suite.gameRulesResponse, output.TextData, "output message received from mock object is different")
}

type InjectGameRulesSuite struct {
	suite.Suite
}

func TestInjectTransformerSuite(t *testing.T) {
	suite.Run(t, new(InjectGameRulesSuite))
}

func (suite *InjectGameRulesSuite) TestDefaultsToProductionGameRulesPackage() {
	// Setup
	productionGameRules := &terosgamerules.GameRules{}

	// Act
	server := terosgameserver.NewServer(nil)

	// Assert
	require := require.New(suite.T())
	require.Equal(
		reflect.TypeOf(server.GetGameRules()),
		reflect.TypeOf(productionGameRules),
	)
}

func (suite *InjectGameRulesSuite) TestUsesInjectedGameRules() {
	// Setup
	fakeTransformer := &rulesstrategyfakes.FakeRulesStrategy{}

	// Act
	server := terosgameserver.NewServer(fakeTransformer)

	// Assert
	require := require.New(suite.T())
	require.Equal(
		reflect.TypeOf(server.GetGameRules()),
		reflect.TypeOf(fakeTransformer),
	)
}
