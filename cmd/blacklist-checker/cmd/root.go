package cmd

import (
	"fmt"
	"isp_checker"
	"math/rand/v2"

	"isp_checker/internal/utils"

	"github.com/spf13/cobra"
)

func dsnblFileExists() error {
	if conf.dnsbl != "" {
		exists, err := utils.FileExists(conf.dnsbl)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("file %s does not exist", conf.dnsbl)
		}
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   isp_checker.Name,
	Short: "Check if your IP/CIDR BlackList --> DNSBL/SPF.",
	Long:  `A Simple tool check DNSBL and SPF.`,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&conf.dnsbl, "dsnbl", "", "DNSBL file to use, if empty it uses the internal list, should be a list of DNSBL to use, each one on a new line")
	rootCmd.PersistentFlags().BoolVarP(&conf.simple, "simple", "s", false, "Simple output")

}

type config struct {
	nameservers []string
	concurrency int
	dnsbl       string
	simple      bool
	ip          string
	sender      string
	domain      string
	dns         string
	ipListed    string
}

func (c config) Nameserver() string {
	if len(c.nameservers) > 1 {
		//rand.N(time.Now().UnixNano())
		return c.nameservers[rand.N(len(c.nameservers))]
	}
	return c.nameservers[0]
}

var conf config

func Execute() error {
	return rootCmd.Execute()
}
