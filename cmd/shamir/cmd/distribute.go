package cmd

import (
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	"github.com/49pctber/shamir"
	"github.com/spf13/cobra"

	"github.com/skip2/go-qrcode"
)

//go:embed template/*
var templates embed.FS

func parseInput(cmd *cobra.Command) (int, int, int, bool, bool, bool, bool) {
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

	// check that primitivePoly is of primitive polynomial of degree 8
	validpolynomials := []int{
		0b100011101,
		0b100101011,
		0b101011111,
		0b101100011,
		0b101100101,
		0b101101001,
		0b111000011,
		0b111100111,
	}

	if !slices.Contains(validpolynomials, primitivePoly) {
		fmt.Println("not a primitive polynomial of degree 8")
		fmt.Println("choose one of the following: ", validpolynomials)
		invalid_command = true
	}

	qr, _ := cmd.Flags().GetBool("qr")
	card, _ := cmd.Flags().GetBool("card")
	file, _ := cmd.Flags().GetBool("file")
	print, _ := cmd.Flags().GetBool("print")

	if invalid_command {
		log.Fatal("invalid command")
	}

	return nshares, threshold, primitivePoly, qr, card, file, print
}

func generateSecret(secret []byte, primitivePoly, nshares, threshold int) *shamir.Shamir {
	s, err := shamir.NewShamirSecret(primitivePoly, nshares, threshold, secret)
	if err != nil {
		log.Fatalf("error distributing secret: %v\n", err)
	}

	return s
}

func distributePNGs(s *shamir.Shamir) error {
	for _, share := range s.GetShares() {
		fname, err := filepath.Abs(share.ShareLabel() + ".png")
		if err != nil {
			return err
		}

		err = qrcode.WriteFile(share.String(), qrcode.High, -10, fname)
		if err != nil {
			return err
		}

		fmt.Printf("%s: png saved to %s\n", share.ShareLabel(), fname)
	}

	return nil
}

func distributeCards(s *shamir.Shamir) error {
	for _, share := range s.GetShares() {
		fname, err := filepath.Abs(share.ShareLabel() + ".svg")
		if err != nil {
			return err
		}

		outfile, err := os.Create(fname)
		if err != nil {
			return err
		}

		template, err := template.ParseFS(templates, "template/card/card.tmpl")
		if err != nil {
			return err
		}

		q, err := qrcode.New(share.String(), qrcode.High)
		if err != nil {
			panic(err)
		}
		q.DisableBorder = true
		qrdata, err := q.PNG(-10)
		if err != nil {
			return err
		}

		err = template.Execute(outfile, struct {
			QrData string
			Label  string
		}{
			QrData: "data:image/png;base64," + base64.RawStdEncoding.EncodeToString(qrdata),
			Label:  strings.Join([]string{share.GetSecretId(), share.GetXString()}, "-"),
		})
		if err != nil {
			return err
		}

		fmt.Printf("%s: card saved to %s\n", share.ShareLabel(), fname)
	}

	return nil
}

func distributeFiles(s *shamir.Shamir) error {

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, share := range s.GetShares() {

		fname := filepath.Clean(path.Join(dir, fmt.Sprintf("%s.txt", share.ShareLabel())))

		err := os.WriteFile(fname, []byte(share.String()), 0400)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%s: text saved to %s\n", share.ShareLabel(), fname)
	}

	return nil
}

type ShareData struct {
	Label      string
	QrData     string
	Index      int
	TranslateX float64
	TranslateY float64
}

type TemplateData struct {
	SecretID string
	Shares   []ShareData
}

func distributePrintablePage(s *shamir.Shamir) error {

	if len(s.GetShares()) > 25 {
		return errors.New("too many shares to print on one page")
	}

	tmpl, err := template.ParseFS(templates, "template/page/*.tmpl")
	if err != nil {
		return err
	}

	fname := strings.Join([]string{shamir.SharePrefix, s.GetId(), "printable.svg"}, "-")
	outfile, err := os.Create(fname)
	if err != nil {
		return err
	}

	shares := s.GetShares()
	sharedata := make([]ShareData, len(shares))
	for i, share := range shares {

		qrraw, err := qrcode.Encode(share.String(), qrcode.High, -5)
		if err != nil {
			return err
		}

		sharedata[i] = ShareData{
			Index:      i,
			Label:      strings.Join([]string{share.GetSecretId(), share.GetXString()}, "-"),
			TranslateX: float64(i%5) * 1.5,
			TranslateY: float64(i/5) * 1.7,
			QrData:     "data:image/png;base64," + base64.RawStdEncoding.EncodeToString(qrraw),
		}

	}

	err = tmpl.ExecuteTemplate(outfile, "base.tmpl", TemplateData{
		SecretID: s.GetId(),
		Shares:   sharedata,
	})
	if err != nil {
		return err
	}

	fmt.Printf("printable page saved to %s\n", fname)
	return nil
}

func distribute(s *shamir.Shamir, qr, card, file, print bool) {
	if qr {
		err := distributePNGs(s)
		if err != nil {
			fmt.Printf("error producing cards: %v\n", err)
		}
	}

	if card {
		err := distributeCards(s)
		if err != nil {
			fmt.Printf("error producing cards: %v\n", err)
		}
	}

	if file {
		err := distributeFiles(s)
		if err != nil {
			fmt.Printf("error producing cards: %v\n", err)
		}
	}

	if print {
		err := distributePrintablePage(s)
		if err != nil {
			fmt.Printf("error producing printable SVG: %v\n", err)
		}
	}
}

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

		nshares, threshold, primitivePoly, qr, card, file, print := parseInput(cmd)

		secret, err := os.ReadFile(args[0])
		if err != nil {
			log.Fatalf("error reading file: %v\n", err)
		}

		s := generateSecret(secret, primitivePoly, nshares, threshold)
		fmt.Println(s)

		distribute(s, qr, card, file, print)
	},
}

var distributeStringCmd = &cobra.Command{
	Use:   "string [secret]",
	Short: "distributes a string",
	Long:  `This is an example of sharing a secret string. The string is specified in the command, and then the shares are printed to the standard out.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		nshares, threshold, primitivePoly, qr, card, file, print := parseInput(cmd)

		s := generateSecret([]byte(args[0]), primitivePoly, nshares, threshold)
		fmt.Println(s)

		distribute(s, qr, card, file, print)
	},
}

func init() {
	rootCmd.AddCommand(distributeCmd)
	distributeCmd.PersistentFlags().IntP("nshares", "n", 0, "number of shares to produce")
	distributeCmd.PersistentFlags().IntP("threshold", "k", 0, "the number of shares needed to reconstruct the secret")
	distributeCmd.PersistentFlags().IntP("primitive", "p", 0x11d, "primitive polynomial to use when constructing Galois field (must be of degree 8)")
	distributeCmd.PersistentFlags().Bool("qr", false, "create PNG QR codes for each share")
	distributeCmd.PersistentFlags().Bool("card", false, "create printable SVG cards for each share")
	distributeCmd.PersistentFlags().Bool("file", false, "save each share in a separate txt file")
	distributeCmd.PersistentFlags().Bool("print", false, "create a printable SVG file with QR codes for each share")

	distributeCmd.AddCommand(distributeFileCmd)

	distributeCmd.AddCommand(distributeStringCmd)
}
