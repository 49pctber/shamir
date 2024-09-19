package cmd

import (
	"embed"
	"encoding/base64"
	"errors"
	"os"
	"strings"
	"text/template"

	"github.com/49pctber/shamir"
	"github.com/skip2/go-qrcode"
)

//go:embed template/*
var templates embed.FS

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

func savePrintableSVG(s *shamir.Shamir) error {

	if len(s.GetShares()) > 25 {
		return errors.New("too many shares to print on one page")
	}

	tmpl, err := template.ParseFS(templates, "template/*.tmpl")
	if err != nil {
		return err
	}

	outfile, err := os.Create(shamir.SharePrefix + "-" + s.GetId() + ".svg")
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

	return nil
}
