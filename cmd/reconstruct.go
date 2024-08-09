package cmd

import (
	"encoding/base64"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	shamir "github.com/49pctber/shamir/internal"
	"github.com/spf13/cobra"
)

var reconstructCmd = &cobra.Command{
	Use:   "reconstruct [folder containing shares]",
	Short: "reconstruct secret(s) given a directory containing shares",
	Long:  `reconstruct secret(s) given a directory containing shares`,
	Run: func(cmd *cobra.Command, args []string) {

		r := regexp.MustCompile(`shamir-(\w+)-(\w+)-(\w+)-(.+)`)

		secretDict := make(map[string]map[int]shamir.Share, 0)
		primitivePolys := make(map[string]int, 0)

		dir, err := cmd.Flags().GetString("directory")
		if err != nil {
			fmt.Println("error specifying input directory")
			os.Exit(1)
		}

		dir, err = filepath.Abs(dir)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("\nSearching %s for .txt files...\n\n", dir)

		dottxtfound := false

		err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {

			if d.IsDir() {
				return nil
			}

			if !strings.HasSuffix(path, ".txt") {
				return nil
			}

			dottxtfound = true

			fmt.Printf("  Searching %s for shares...\n", path)

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			matches := r.FindAllSubmatch(data, -1)
			for _, match := range matches {
				fmt.Printf("    Found shamir-%s-%s-%s\n", match[1], match[2], match[3])
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
					secretDict[id] = make(map[int]shamir.Share, 0)
					primitivePolys[id] = int(primitivePoly)
				}

				secretDict[id][int(x)] = shamir.NewShare(x, y)
			}

			return nil
		})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if !dottxtfound {
			fmt.Println("No .txt files found. Exiting.\n")
			os.Exit(0)
		}

		fmt.Println("\nAttempting to reconstruct secrets from shares that were found...\n")

		for id, shares := range secretDict {
			sharesslice := make([]shamir.Share, 0)
			for _, share := range shares {
				sharesslice = append(sharesslice, share)
			}
			secret, err := shamir.RecoverSecret(primitivePolys[id], sharesslice)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("%s:\n%s\n\n", id, secret)
		}
	},
}

func init() {
	rootCmd.AddCommand(reconstructCmd)

	reconstructCmd.PersistentFlags().StringP("directory", "d", ".", "input directory")
}
