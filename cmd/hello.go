package cmd

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type HelloOptions struct {
	Message string
}

var (
	helloOpt HelloOptions
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Short description",
	Long:  `Long description`,
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info().Str("ts", time.Now().Format(time.RFC3339)).Msgf("Hello %s", helloOpt.Message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(helloCmd)
	helloCmd.Flags().StringVarP(&helloOpt.Message, "message", "m", "world", "message to print (eg. world)")
}
