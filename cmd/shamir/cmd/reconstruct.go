package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	shamir "github.com/49pctber/shamir"
	"github.com/spf13/cobra"
)

var reconstructCmd = &cobra.Command{
	Use:   "reconstruct",
	Short: "reconstruct secret",
	Long:  `reconstruct secret based on strings or files`,
}

var reconstructFileCmd = &cobra.Command{
	Use:   "file",
	Short: "searches the directory for shares in files prefixed with shamir",
	Long:  `searches the directory for shares in files prefixed with shamir`,
	Run: func(cmd *cobra.Command, args []string) {

		shares := make([]shamir.Share, 0)

		dir, err := cmd.Flags().GetString("directory")
		if err != nil {
			log.Fatalf("error specifying input directory: %v\n", err)
		}

		dir, err = filepath.Abs(dir)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Searching %s for files prefixed with %s...\n", dir, shamir.SharePrefix)

		err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {

			if err != nil {
				fmt.Printf("Error accessing %s: %v\n", path, err)
				return err
			}

			// ignore subdirectories
			if d.IsDir() && path != dir {
				return filepath.SkipDir
			} else if path == dir {
				return nil
			}

			if !strings.HasPrefix(filepath.Base(path), shamir.SharePrefix) {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			new_shares, err := shamir.NewSharesFromString(string(data))
			if err != nil {
				return err
			}

			shares = append(shares, new_shares...)

			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		if len(shares) == 0 {
			fmt.Println("No shares found. Exiting.")
			return
		} else {
			fmt.Printf("Found %d shares.\n", len(shares))
		}

		secretDict := make(map[string][]shamir.Share, 0)

		for _, share := range shares {
			if _, ok := secretDict[share.GetSecretId()]; !ok {
				secretDict[share.GetSecretId()] = make([]shamir.Share, 0)
			}
			secretDict[share.GetSecretId()] = append(secretDict[share.GetSecretId()], share)
		}

		fmt.Println("Attempting to reconstruct secrets from shares that were found...")

		for id, shares := range secretDict {
			secret, err := shamir.RecoverSecret(shares)
			if err != nil {
				log.Fatal(err)
			}

			fname := "secret-" + id
			abs, err := filepath.Abs(fname)
			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile(fname, secret, 0700)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Secret %s saved to %s\n", id, abs)
		}
	},
}

var reconstructStringCmd = &cobra.Command{
	Use:   "string [shares...]",
	Short: "reconstruct secret given a sequences of shares",
	Long:  `reconstruct secret given a sequences of shares`,
	Run: func(cmd *cobra.Command, args []string) {

		shares := make([]shamir.Share, 0)

		for _, arg := range args {
			new_shares, err := shamir.NewSharesFromString(arg)
			if err != nil {
				log.Fatal(err)
			}

			shares = append(shares, new_shares...)
		}

		if len(shares) == 0 {
			fmt.Println("No valid shares specified. Exiting.")
			return
		}

		secretDict := make(map[string][]shamir.Share, 0)

		for _, share := range shares {
			fmt.Printf("Found %s\n", share.String())
			if _, ok := secretDict[share.GetSecretId()]; !ok {
				secretDict[share.GetSecretId()] = make([]shamir.Share, 0)
			}
			secretDict[share.GetSecretId()] = append(secretDict[share.GetSecretId()], share)
		}

		fmt.Println("Attempting to reconstruct secrets from shares that were found...")

		for id, shares := range secretDict {
			secret, err := shamir.RecoverSecret(shares)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s:\n%s\n", id, secret)
		}
	},
}

func init() {
	rootCmd.AddCommand(reconstructCmd)

	reconstructCmd.AddCommand(reconstructFileCmd)
	reconstructFileCmd.PersistentFlags().StringP("directory", "d", "", "directory to search and save results")

	reconstructCmd.AddCommand(reconstructStringCmd)
}
