package cmd

import (
	"context"
	"fmt"
	"isp_checker"
	"isp_checker/internal/check"
	"isp_checker/internal/utils"
	"net"
	"sync/atomic"

	"github.com/jniltinho/netprivate"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.org/x/sync/semaphore"
)

var ipCmd = &cobra.Command{
	Use:   "ip <ip-address>",
	Args:  cobra.ExactArgs(1),
	Short: "Check IP against available blacklists.",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		ip := net.ParseIP(args[0])
		if ip == nil {
			return fmt.Errorf("invalid IP address: %s", args[0])
		}

		if netprivate.Is(ip) {
			return fmt.Errorf("ip: %s is in the private range", ip)
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var sem = semaphore.NewWeighted(int64(conf.concurrency))
		ip := net.ParseIP(args[0])

		var items []check.Item
		var hosts []string
		hosts, err = isp_checker.GetDNSBLs(conf.dnsbl)
		for _, host := range hosts {
			items = append(items, check.Item{
				IP:        ip,
				Blacklist: fmt.Sprintf("%s.%s.", utils.ReverseIP(ip.String()), host),
				Host:      host,
			})
		}

		var blacklisted uint64

		if !conf.simple {
			log.Info().Int("queries", len(items)).Int("dsnbl", len(hosts)).Msg("processing")
		}

		for _, i := range items {
			if err = sem.Acquire(context.Background(), 1); err != nil {
				return err
			}
			go func(item check.Item) {
				if b, _ := processIpCheck(sem, item); b {
					atomic.AddUint64(&blacklisted, 1)
				}
			}(i)
		}

		if err = sem.Acquire(context.Background(), int64(conf.concurrency)); err != nil {
			return err
		}

		if conf.simple {
			if blacklisted > 0 {
				fmt.Print("listed:", conf.ipListed)
			} else {
				fmt.Print("unlisted:", ip.String())
			}
		}

		if !conf.simple {
			log.Info().Uint64("blacklisted", blacklisted).Int("queries", len(items)).Msg("Finished")
		}
		return err
	},
}

func init() {
	checkCmd.AddCommand(ipCmd)
}
