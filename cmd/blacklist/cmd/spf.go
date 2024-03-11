package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/mileusna/spf"
	"github.com/spf13/cobra"
)

var spfCmd = &cobra.Command{
	Use:   "spf",
	Args:  cobra.OnlyValidArgs,
	Short: "Check SPF record.",
	Run: func(cmd *cobra.Command, args []string) {

		if conf.ip != "" && conf.sender != "" {
			checkSPF()
		} else {
			cmd.Help()
			os.Exit(0)
		}
	},
}

// isp-checker spf --ip=11.22.33.44 --sender=user@aol.com
func init() {
	spfCmd.PersistentFlags().StringVar(&conf.ip, "ip", "", "IP address to use for the SPF check")
	spfCmd.PersistentFlags().StringVar(&conf.sender, "sender", "", "Sender email to use for the SPF check")
	spfCmd.PersistentFlags().StringVar(&conf.dns, "dns", "1.1.1.1:53", "Default DNS server to use")
	rootCmd.AddCommand(spfCmd)
}

func checkSPF() {

	if conf.ip == "" && conf.sender == "" {
		fmt.Print("StartError")
		os.Exit(1)
	}

	// optional, set DNS server which will be used by resolver.
	// Default is Google's 8.8.8.8:53
	spf.DNSServer = conf.dns
	conf.domain = getDomain(conf.sender)

	if conf.domain == "" {
		fmt.Print("StartError")
		os.Exit(1)
	}

	ip := net.ParseIP(conf.ip)
	r := spf.CheckHost(ip, conf.domain, conf.sender, "")
	// returns spf check result
	// "PASS" / "FAIL" / "SOFTFAIL" / "NEUTRAL" / "NONE" / "TEMPERROR" / "PERMERROR"

	// if you only need to retrive SPF record as string from DNS
	//spfRecord, _ := spf.LookupSPF("domain.com")

	res := []string{"PASS", "FAIL", "SOFTFAIL", "NEUTRAL", "NONE", "TEMPERROR", "PERMERROR"}

	for _, v := range res {
		if r.String() == v {
			fmt.Print(strings.ToLower(v))
			return
		}
	}
	fmt.Print("StartError")

}

func getDomain(email string) string {
	at := strings.LastIndex(email, "@")
	if at >= 0 {
		domain := email[at+1:]
		return domain
	}
	return ""
}
