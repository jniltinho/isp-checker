package cmd

import (
	"net/netip"

	"isp_checker/internal/check"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check available blacklists.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = dsnblFileExists(); err != nil {
			return err
		}
		for _, n := range conf.nameservers {
			_, err = netip.ParseAddrPort(n)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	checkCmd.PersistentFlags().StringSliceVarP(&conf.nameservers, "nameservers", "n", []string{"1.1.1.1:53"}, "Nameservers to use, a random one is used for each request")
	checkCmd.PersistentFlags().IntVarP(&conf.concurrency, "concurrency", "c", 25, "How many requests to process at once")
	rootCmd.AddCommand(checkCmd)
}

func processIpCheck(sem *semaphore.Weighted, item check.Item) (blacklisted bool, err error) {
	var responses []string
	blacklisted, responses, err = check.Check(
		sem,
		item,
		conf.Nameserver(),
	)
	if !conf.simple {
		if blacklisted {
			log.Warn().Str("dnsbl", item.Host).Bool("blacklisted", blacklisted).Strs("responses", responses).IPAddr("ip", item.IP).Send()
		} else {
			log.Trace().Str("dnsbl", item.Host).Bool("blacklisted", blacklisted).Strs("responses", responses).IPAddr("ip", item.IP).Send()
		}
	}

	if len(responses) > 0 {
		conf.ipListed = responses[0]
	}

	return blacklisted, err
}
