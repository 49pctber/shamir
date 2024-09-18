package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	shamir "github.com/49pctber/shamir"
	"github.com/spf13/cobra"
)

var distributeCmd = &cobra.Command{
	Use:   "distribute [secret to share]",
	Short: "Distrbute a secret S into n shares, where any k shares can reconstruct S.",
	Long:  `Distrbute a secret S into n shares, where any k shares can reconstruct S.`,
}

var distributeFileCmd = &cobra.Command{
	Use:   "file [filename]",
	Short: "distributes a file",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		invalid_command := false

		nshares, err := cmd.Flags().GetInt("nshares")
		if err != nil || nshares < 2 {
			fmt.Println("provide n >= 2")
			invalid_command = true
		}

		threshold, err := cmd.Flags().GetInt("threshold")
		if err != nil || threshold < 2 {
			fmt.Println("provide k >= 2")
			invalid_command = true
		}

		primitivePoly, err := cmd.Flags().GetInt("primitive")
		if err != nil {
			fmt.Println("error reading primitive polynomial")
			invalid_command = true
		}

		if invalid_command {
			os.Exit(1)
		}

		secretstring, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Println("error reading file")
			os.Exit(1)
		}

		fmt.Printf("Generating %d-of-%d secret sharing scheme...\n", threshold, nshares)

		s, err := shamir.NewShamirSecret(primitivePoly, nshares, threshold, []byte(secretstring))
		if err != nil {
			fmt.Printf("error distributing secret: %v\n", err)
			os.Exit(1)
		}

		dir, err := os.Getwd()
		if err != nil {
			fmt.Printf("error getting current directory: %v\n", dir)
			os.Exit(1)
		}

		for _, share := range s.GetShares() {

			fname := filepath.Clean(path.Join(dir, fmt.Sprintf("%s.txt", share.ShareLabel())))

			err := os.WriteFile(fname, []byte(share.String()), 0400)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("%s saved to %s\n", share.ShareLabel(), fname)
		}
	},
}

var distributeStringCmd = &cobra.Command{
	Use:   "string [secret]",
	Short: "distributes a string",
	Long:  `This is an example of sharing a secret string. The string is specified in the command, and then the shares are printed to the standard out.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		invalid_command := false

		nshares, err := cmd.Flags().GetInt("nshares")
		if err != nil || nshares < 2 {
			fmt.Println("provide n >= 2")
			invalid_command = true
		}

		threshold, err := cmd.Flags().GetInt("threshold")
		if err != nil || threshold < 2 {
			fmt.Println("provide k >= 2")
			invalid_command = true
		}

		primitivePoly, err := cmd.Flags().GetInt("primitive")
		if err != nil {
			fmt.Println("error reading primitive polynomial")
			invalid_command = true
		}

		if invalid_command {
			os.Exit(1)
		}

		secretstring := args[0]

		fmt.Printf("Generating %d-of-%d secret sharing scheme...\n", threshold, nshares)

		s, err := shamir.NewShamirSecret(primitivePoly, nshares, threshold, []byte(secretstring))
		if err != nil {
			fmt.Printf("error distributing secret: %v\n", err)
			os.Exit(1)
		}

		for i := range nshares {
			fmt.Printf("%s\n", s.ShareString(i))
		}
	},
}

func init() {
	rootCmd.AddCommand(distributeCmd)
	distributeCmd.PersistentFlags().IntP("nshares", "n", 0, "number of shares to produce")
	distributeCmd.PersistentFlags().IntP("threshold", "k", 0, "the number of shares needed to reconstruct the secret")
	distributeCmd.PersistentFlags().IntP("primitive", "p", 0x11d, "primitive polynomial to use when constructing Galois field")

	distributeCmd.AddCommand(distributeFileCmd)

	distributeCmd.AddCommand(distributeStringCmd)
}
