package terosgameserver_test

import (
	"bytes"
	"github.com/Chadius/terosGameServer/rulesstrategyfakes"
	"github.com/chadius/terosGameServer/rpc/github.com/chadius/teros_game_server"
	"github.com/chadius/terosGameServer/terosgameserver"
	"github.com/chadius/terosgamerules"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"net/http/httptest"
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
	fakeGameRules     *terosgamerules.RulesStrategy
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

func (suite *ServerUsesPackageSuite) fakeGameRulesWithResponse(expectedResponse []byte) *creatingsymmetryfakes.FakeTransformerStrategy {
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

//
//func (suite *ServerUsesPackageSuite) TestWhenClientMakesRequest_ResponseIsValid() {
//	// Act
//	suite.server.ServeHTTP(suite.responseRecorder, suite.request)
//
//	// Assert
//	response := suite.responseRecorder.Result()
//
//	require := require.New(suite.T())
//	require.Equal(200, response.StatusCode, "Status code is wrong")
//}
