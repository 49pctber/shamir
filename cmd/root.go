/*
Copyright © 2024 49pctber

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
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shamir",
	Short: "An implementation of Shamir's secret sharing scheme.",
	Long: `This application is a practical example of a Shamir Secret Sharing scheme.
	
It allows you to divide a secret string S into n "shares".
Any k of those n shares can be used to reconstruct the original secret S.

Having k-1 or fewer shares will provide *no* information about the secret other than its length.`,
	Run: func(cmd *cobra.Command, args []string) {

		secretstring, err := cmd.Flags().GetString("secret")
		if err != nil || len(secretstring) > 0 {
			distributeCmd.Run(cmd, args)
			os.Exit(0)
		}

		var resp string
		fmt.Print("No command specified. Search current directory for shares in .txt files? (y/n): ")
		fmt.Scanf("%s", &resp)

		if resp[0] == 'y' || resp[0] == 'Y' {
			reconstructCmd.Run(cmd, args)
			os.Exit(0)
		} else {
			cmd.Root().Help()
			os.Exit(0)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("secret", "s", "", "the secret to share")
	rootCmd.PersistentFlags().StringP("directory", "d", ".", "input/output directory")
	rootCmd.PersistentFlags().IntP("nshares", "n", 0, "number of shares to produce")
	rootCmd.PersistentFlags().IntP("threshold", "k", 0, "the number of shares needed to reconstruct the secret")
	rootCmd.PersistentFlags().IntP("primitive", "p", 0x11d, "primitive polynomial to use when constructing Galois field")
}
