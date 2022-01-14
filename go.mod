module github.com/chadius/terosgameserver

go 1.15

replace github.com/chadius/terosgamerules v0.0.0-20220114204417-1dca36259870 => ../terosGameRules

require (
	github.com/chadius/terosgamerules v0.0.0-20220114204417-1dca36259870
	github.com/maxbrunsfeld/counterfeiter/v6 v6.4.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/twitchtv/twirp v8.1.1+incompatible // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c
)
