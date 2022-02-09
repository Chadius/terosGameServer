package terosgameserver

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chadius/terosgamerules"
	"github.com/chadius/terosgameserver/rpc/github.com/chadius/teros_game_server"
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

	return s.runScriptToCreateResults(scriptDataByteStream, squaddieDataByteStream, powerDataByteStream)
}

func (s *Server) runScriptToCreateResults(scriptDataByteStream *bytes.Buffer, squaddieDataByteStream *bytes.Buffer, powerDataByteStream *bytes.Buffer) (*teros_game_server.Results, error) {
	var packagePanicErr error
	defer func() {
		if panicContext := recover(); panicContext != nil {
			packagePanicErr = fmt.Errorf("package panic: %v", panicContext)
		}
	}()
	var outputGameResults bytes.Buffer

	replayErr := s.GetGameRules().ReplayBattleScript(scriptDataByteStream, squaddieDataByteStream, powerDataByteStream, &outputGameResults)
	results := &teros_game_server.Results{TextData: outputGameResults.Bytes()}
	if packagePanicErr != nil {
		return results, packagePanicErr
	}
	return results, replayErr
}

func (s *Server) GetGameRules() terosgamerules.RulesStrategy {
	return s.gameRules
}

// NewServer returns a new Server object with the given gameRules.
//   Defaults to using the production GameRules if none is given.
func NewServer(gameRules terosgamerules.RulesStrategy) *Server {
	var gameRulesToUse terosgamerules.RulesStrategy
	gameRulesToUse = &terosgamerules.GameRules{}
	if gameRules != nil {
		gameRulesToUse = gameRules
	}
	return &Server{
		gameRules: gameRulesToUse,
	}
}
