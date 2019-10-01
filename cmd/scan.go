/*
Copyright © 2019 Lars Mogren <lars@thenatureofsoftware.se>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/TheNatureOfSoftware/k3pi/pkg"
	cmd2 "github.com/TheNatureOfSoftware/k3pi/pkg/cmd"
	"github.com/TheNatureOfSoftware/k3pi/pkg/misc"
	"github.com/TheNatureOfSoftware/k3pi/pkg/ssh"
	"github.com/kubernetes-sigs/yaml"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

// scanCmd represents the list command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scans for members of the Raspberries",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		scanRequest := &cmd2.ScanRequest{
			Cidr:              viper.GetString("cidr"),
			HostnameSubString: viper.GetString("substr"),
			SSHSettings:       sshSettings(),
			UserCredentials:   credentials(viper.GetStringSlice("basic-auth")),
		}
		cmdOpFactory := &pkg.CmdOperatorFactory{Create: ssh.NewCmdOperator}
		nodes, err := cmd2.ScanForRaspberries(scanRequest, misc.NewHostScanner(), cmdOpFactory)
		if err != nil {
			fmt.Errorf("failed to scan for Raspberries: %d", err)
		}
		y, err := yaml.Marshal(nodes)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}
		fmt.Println(string(y))
	},
}

// Splits slice of <username>:<password> and returns a map
func credentials(basicAuths []string) map[string]string {
	c := make(map[string]string)
	for _, v := range basicAuths {
		parts := strings.Split(v, ":")
		if len(parts) == 2 {
			c[parts[0]] = parts[1]
		}
	}
	return c
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().String("cidr", "192.168.1.0/24", "CIDR to scan for members")
	scanCmd.Flags().String("substr", "", "Substring that should be part of hostname")
	scanCmd.Flags().StringSlice("basic-auth", []string{}, "Username and password separated with ':' for authentication")
	_ = viper.BindPFlag("cidr", scanCmd.Flags().Lookup("cidr"))
	_ = viper.BindPFlag("substr", scanCmd.Flags().Lookup("substr"))
	_ = viper.BindPFlag("basic-auth", scanCmd.Flags().Lookup("basic-auth"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func sshSettings() *ssh.Settings {
	return &ssh.Settings{KeyPath: viper.GetString("ssh-key"),
		User: viper.GetString("user"), Port: viper.GetString("ssh-port")}
}
