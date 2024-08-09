package cmd

import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	shamir "github.com/49pctber/shamir/internal"
	"github.com/spf13/cobra"
)

var reconstructCmd = &cobra.Command{
	Use:   "reconstruct [folder containing secrets]",
	Short: "reconstruct secret(s) given a directory",
	Long:  `reconstruct secret(s) given a directory`,
	Run: func(cmd *cobra.Command, args []string) {

		r := regexp.MustCompile(`^shamir-(\w+)-(\w+)-(\w+)-(.+)$`)

		secretDict := make(map[string][]shamir.Share, 0)
		primitivePolys := make(map[string]int, 0)

		dir, err := cmd.Flags().GetString("directory")
		if err != nil {
			fmt.Println("error specifying input directory")
			os.Exit(1)
		}

		err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {

			if d.IsDir() {
				return nil
			}

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			match := r.FindSubmatch(data)
			if len(match) > 0 {
				id := string(match[1])

				primitivePoly, err := strconv.ParseInt(string(match[2]), 16, 64)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				xdata, err := strconv.ParseInt(string(match[3]), 10, 64)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				ydata, err := base64.RawStdEncoding.DecodeString(string(match[4]))
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				x := shamir.GfElement(xdata)
				y := make([]shamir.GfElement, len(ydata))
				for i := range ydata {
					y[i] = shamir.GfElement(ydata[i])
				}

				if _, exists := secretDict[id]; !exists {
					secretDict[id] = make([]shamir.Share, 0)
					primitivePolys[id] = int(primitivePoly)
				}

				secretDict[id] = append(secretDict[id], shamir.NewShare(x, y))
			}

			return nil
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for id, shares := range secretDict {
			secret, err := shamir.RecoverSecret(primitivePolys[id], shares)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("%s: %s\n", id, secret)
		}

	},
}

func init() {
	rootCmd.AddCommand(reconstructCmd)

	reconstructCmd.PersistentFlags().StringP("directory", "d", "shamir_shares", "input directory")
}
