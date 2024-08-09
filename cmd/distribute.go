package cmd

import (
	"fmt"
	"os"
	"path"

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
			fmt.Println("provide a secret to share using -S=\"[Your secret string here]\"")
			os.Exit(1)
		}

		nshares, err := cmd.Flags().GetInt("nshares")
		if err != nil || nshares < 2 {
			fmt.Println("provide n >= 2")
			os.Exit(1)
		}

		threshold, err := cmd.Flags().GetInt("threshold")
		if err != nil || threshold < 2 {
			fmt.Println("provide k >= 2")
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

		s, err := shamir.NewShamirSecret(primitivePoly, nshares, threshold, []byte(secretstring))
		if err != nil {
			fmt.Printf("error distributing secret: %v\n", err)
			os.Exit(1)
		}

		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.Mkdir(dir, os.ModeDir)
			if err != nil {
				panic(err)
			}
		}

		for i := range nshares {
			fname := path.Join(dir, fmt.Sprintf("%s.txt", s.ShareLabel(i)))
			os.WriteFile(fname, []byte(s.ShareString(i)), os.ModeAppend)
		}
	},
}

func init() {
	rootCmd.AddCommand(distributeCmd)

	distributeCmd.PersistentFlags().StringP("secret", "S", "", "the secret to share")
	distributeCmd.PersistentFlags().StringP("directory", "d", "shamir_shares", "output directory")
	distributeCmd.PersistentFlags().IntP("nshares", "n", 0, "number of shares to produce")
	distributeCmd.PersistentFlags().IntP("threshold", "k", 0, "the number of shares needed to reconstruct the secret")
	distributeCmd.PersistentFlags().IntP("primitive", "p", 0x11d, "primitive polynomial to use when constructing Galois field")
}
