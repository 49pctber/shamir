package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	shamir "github.com/49pctber/shamir/internal"
	"github.com/spf13/cobra"
)

var distributeCmd = &cobra.Command{
	Use:   "distribute [secret to share]",
	Short: "Distrbute a secret S into n shares, where any k shares can reconstruct S.",
	Long:  `Distrbute a secret S into n shares, where any k shares can reconstruct S.`,
	Run: func(cmd *cobra.Command, args []string) {

		secretstring, err := cmd.Flags().GetString("secret")
		if err != nil || len(secretstring) < 1 {
			cmd.Help()
			fmt.Println("\nprovide a secret to share using -s=\"<secret string>\"")
			os.Exit(1)
		}

		nshares, err := cmd.Flags().GetInt("nshares")
		if err != nil || nshares < 2 {
			cmd.Help()
			fmt.Println("\nprovide n >= 2")
			os.Exit(1)
		}

		threshold, err := cmd.Flags().GetInt("threshold")
		if err != nil || threshold < 2 {
			cmd.Help()
			fmt.Println("\nprovide k >= 2")
			os.Exit(1)
		}

		primitivePoly, err := cmd.Flags().GetInt("primitive")
		if err != nil {
			fmt.Println("error reading primitive polynomial")
			os.Exit(1)
		}

		dir, err := cmd.Flags().GetString("directory")
		if err != nil {
			fmt.Println("error specifying output directory")
			os.Exit(1)
		}

		fmt.Printf("Generating %d-of-%d secret sharing scheme...\n", threshold, nshares)

		s, err := shamir.NewShamirSecret(primitivePoly, nshares, threshold, []byte(secretstring))
		if err != nil {
			fmt.Printf("error distributing secret: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, 0755)
			if err != nil {
				panic(err)
			}
		}

		for i := range nshares {

			fname := filepath.Clean(path.Join(dir, fmt.Sprintf("%s.txt", s.ShareLabel(i))))
			sharestring := s.ShareString(i)

			err := os.WriteFile(fname, []byte(sharestring), 0400)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("%s saved to %s\n", sharestring, fname)

		}
	},
}

func init() {
	rootCmd.AddCommand(distributeCmd)
}
