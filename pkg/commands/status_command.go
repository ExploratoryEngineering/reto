package commands

import (
	"fmt"

	"github.com/ExploratoryEngineering/reto/pkg/gitutil"
	"github.com/ExploratoryEngineering/reto/pkg/release"
)

type statusCommand struct {
	Verbose bool `kong:"short='V',help='Verbose output'"`
}

func okNotOK(v bool) string {
	if v {
		return "OK"
	}
	return "NOT OK"
}

func (c *statusCommand) Run(rc RunContext) error {
	ctx, err := release.GetContext()
	if err != nil {
		return err
	}

	changelogErr := release.ChangelogComplete(false)
	configErr := release.VerifyConfig(ctx.Config, false)

	fmt.Printf("Configuration:       %s\n", okNotOK(configErr == nil))
	fmt.Printf("Changelog            %s\n", okNotOK(changelogErr == nil))
	fmt.Printf("Version number:      %s\n", okNotOK(!ctx.Released))
	fmt.Printf("Uncommitted changes: %s\n", okNotOK(!gitutil.HasChanges(ctx.Config.SourceRoot)))
	fmt.Println()
	fmt.Printf("Active version:      %s\n", ctx.Version)
	fmt.Printf("Commit Hash:         %s\n", ctx.CommitHash)
	fmt.Printf("Name:                %s\n", ctx.Name)

	if rc.ReleaseCommands().Status.Verbose {
		fmt.Println()
		fmt.Println("Configuration:")
		fmt.Println("  Targets: ")
		for _, v := range ctx.Config.Targets {
			fmt.Printf("  - %s\n", v)
		}
		fmt.Println("  Files:")
		for _, v := range ctx.Config.Files {
			fmt.Printf("  - %s/%s\n", v.Name, v.Target)
		}
		fmt.Println()
	}
	return nil
}
