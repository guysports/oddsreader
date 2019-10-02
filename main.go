package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/guysports/oddsreader/pkg/cmd"
)

var cli struct {
	// Compare - Run the compare
	Compare cmd.Compare `cmd:"" help:"Retrieve latest odds from LiveOddsAPI site"`
	PM      cmd.PM      `cmd:"" help:"Retrieve latest odds from PM site"`
}

func main() {
	apikey := os.Getenv("ODDS_APIKEY")

	ctx := kong.Parse(&cli)
	err := ctx.Run(&cmd.Globals{
		APIKey:  apikey,
		Leagues: []string{"soccer_epl", "soccer_efl_champ", "soccer_spain_la_liga", "soccer_italy_serie_a", "soccer_germany_bundesliga"},
	})
	ctx.FatalIfErrorf(err)

}
