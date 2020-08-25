package apf

import "github.com/spf13/cobra"

type Cli struct {
	rootCmd cobra.Command
}

func NewCli(use string) *Cli {
	return &Cli{
		cobra.Command{
			Use: use,
		},
	}
}
