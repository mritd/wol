// Copyright © 2018 mritd <mritd1234@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/mritd/wol/pkg/utils"
	"github.com/mritd/wol/pkg/wol"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var machine wol.Machine

var rootCmd = &cobra.Command{
	Use:   "wol",
	Short: "A simple WOL tool",
	Long: `
A simple WOL tool`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&machine.BroadcastInterface, "interface", "i", "", "Broadcast Interface")
	rootCmd.PersistentFlags().StringVarP(&machine.BroadcastIP, "ip", "b", "255.255.255.255", "Broadcast IP")
	rootCmd.PersistentFlags().IntVarP(&machine.Port, "port", "p", 7, "UDP Port")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wol.yaml)")
}

func initConfig() {

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		utils.CheckAndExit(err)
		cfgFile := home + string(filepath.Separator) + ".wol.yaml"
		if _, err = os.Stat(cfgFile); err != nil {
			os.Create(cfgFile)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".wol")
		viper.SetConfigType("yaml")

	}

	viper.AutomaticEnv()
	viper.ReadInConfig()

}
