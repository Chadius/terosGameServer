package terosgameserver

import (
	"bytes"
	"context"
	terosgamerules "github.com/Chadius/terosGameRules"
	"github.com/Chadius/terosGameServer/rpc/github.com/chadius/teros_game_server"
)

// Server implements the RulesStrategy service
type Server struct {
	gameRules terosgamerules.RulesStrategy
}

// ReplayBattleScript uses the given data to process game rules
func (s *Server) ReplayBattleScript(cts context.Context, data *teros_game_server.DataStreams) (*teros_game_server.Results, error) {

	scriptDataByteStream := bytes.NewBuffer(data.GetScriptData())
	squaddieDataByteStream := bytes.NewBuffer(data.GetSquaddieData())
	powerDataByteStream := bytes.NewBuffer(data.GetPowerData())

	var outputGameResults bytes.Buffer

	transformErr := s.GetGameRules().ReplayBattleScript(scriptDataByteStream, squaddieDataByteStream, powerDataByteStream, &outputGameResults)
	outputImage := &teros_game_server.Results{TextData: outputGameResults.Bytes()}
	return outputImage, transformErr
}

func (s *Server) GetGameRules() terosgamerules.RulesStrategy {
	return s.gameRules
}

// NewServer returns a new Server object with the given gameRules.
//   Defaults to using the production GameRules if none is given.
func NewServer(gameRules terosgamerules.RulesStrategy) *Server {
	//var transformerToUse creatingsymmetry.TransformerStrategy
	//transformerToUse = &creatingsymmetry.FileTransformer{}
	//if transformer != nil {
	//	transformerToUse = transformer
	//}
	return &Server{
		gameRules: gameRules,
	}
}
