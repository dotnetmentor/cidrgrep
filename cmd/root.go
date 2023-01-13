package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/dotnetmentor/cidrgrep/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type GlobalOptions struct {
	CIDR  string
	Debug bool
}

var (
	opt GlobalOptions
)

var rootCmd = &cobra.Command{
	Use:          "cidrgrep",
	Short:        "cidrgrep - CIDR grepping",
	Long:         `cidrgrep - like grep but for CIDR/IP matching`,
	SilenceUsage: true,
	Version:      fmt.Sprintf("%s (commit=%s)", version.Version, version.Commit),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 2 && args[0] == "completion" {
			fmt.Println("doing completion work...")
			switch args[1] {
			case "bash":
				cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		}

		if opt.CIDR == "" {
			return fmt.Errorf("CIDR flag required")
		}
		log.Output(os.Stderr)

		if !opt.Debug {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}

		_, cidr, err := net.ParseCIDR(opt.CIDR)
		if err != nil {
			return err
		}

		found := false
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			line := scanner.Text()
			ip := net.ParseIP(line)
			if ip == nil {
				os.Stderr.WriteString(fmt.Sprintf("\"%s\" is not a valid IP\n", line))
				continue
			}

			if cidr.Contains(ip) {
				found = true
				os.Stdout.WriteString(fmt.Sprintf("%s\n", ip))
				log.Debug().Msgf("The IP \"%s\" is within the CIDR \"%s\"", ip, cidr)
			} else {
				log.Debug().Msgf("The IP \"%s\" is NOT within the CIDR \"%s\"", ip, cidr)
			}
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		if !found {
			return fmt.Errorf("no matches found")
		}

		return nil
	},
}

func Execute() {
	rootCmd.SetOut(os.Stdout)
	rootCmd.SetErr(os.Stderr)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&opt.CIDR, "cidr", "c", "", "CIDR to match")
	rootCmd.Flags().BoolVarP(&opt.Debug, "debug", "d", false, "Run in debug mode")
}
