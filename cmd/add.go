// Copyright © 2016 Antoine GIRARD <antoine.girard@sapk.fr>
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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	valid "github.com/asaskevich/govalidator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func saveConfig() error {
	fmt.Println("Tring to save file ... ")
	f, err := os.Create(viper.ConfigFileUsed())
	if err != nil {
		return err
	}
	defer f.Close()

	cfg, err := json.MarshalIndent(viper.AllSettings(), "", "    ")
	if err != nil {
		return err
	}
	f.WriteString(string(cfg))
	fmt.Println("Done! ")
	return nil
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add fqdn.com [one/relative/file] [a/other/one]",
	Short: "Add a website to monitor",
	Long:  `Add a website to monitor and any specific relative url to check`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("FQDN not provided")
		}

		FQDN := args[0] //TODO check present of port OR check if IP
		if !valid.IsDNSName(FQDN) {
			return errors.New("FQDN '" + FQDN + "'not valid")
		}

		fmt.Println("Adding : " + FQDN) //+ strings.Join(args, " ")

		FQDNs := viper.GetStringMapStringSlice("fqdn")
		FQDNElm, ok := FQDNs[FQDN]
		if !ok {
			//FQDNs[FQDN]
			FQDNs[FQDN] = []string{"/"} //TODO parse sub url
		} else {
			fmt.Println("Getting old value : " + strings.Join(FQDNElm, ",")) //+ strings.Join(args, " ")
			FQDNs[FQDN] = []string{"/"}                                      //TODO parse sub url and merge
		}
		viper.Set("fqdn", FQDNs)

		return saveConfig()
	},
}

func init() {
	//addCmd.Flags().StringVarP(&fqdn, "fqdn", "f", "", "FQDN to monitor")
	RootCmd.AddCommand(addCmd)
}
