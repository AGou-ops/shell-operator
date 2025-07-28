package main

import (
	"fmt"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/flant/shell-operator/pkg/app"
	"github.com/flant/shell-operator/pkg/debug"
	"github.com/flant/shell-operator/pkg/filter/jq"
)

func main() {
	kpApp := kingpin.New(app.AppName, fmt.Sprintf("%s %s: %s", app.AppName, app.Version, app.AppDescription))

	// override usage template to reveal additional commands with information about start command
	kpApp.UsageTemplate(app.OperatorUsageTemplate(app.AppName))

	// print version
	kpApp.Command("version", "Show version.").Action(func(_ *kingpin.ParseContext) error {
		fmt.Printf("%s %s\n", app.AppName, app.Version)
		fl := jq.NewFilter()
		fmt.Println(fl.FilterInfo())
		return nil
	})

	// start main loop
	startCmd := kpApp.Command("start", "Start shell-operator.").
		Default().
		Action(start)
	app.DefineStartCommandFlags(kpApp, startCmd)

	debug.DefineDebugCommands(kpApp)
	debug.DefineDebugCommandsSelf(kpApp)

	kingpin.MustParse(kpApp.Parse(os.Args[1:]))
}
