package main

import (
	"github.com/dgnorton/dmapi"
	htmlTemplate "html/template"
	"io"
	"io/ioutil"
	textTemplate "text/template"
)

func fprintText(wr io.Writer, e *dmapi.Entries, templateFile string) error {
	bytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	txt, err := textTemplate.New("txt").Parse(string(bytes))
	if err != nil {
		return err
	}

	err = txt.ExecuteTemplate(wr, "txt", e)

	return err
}

func fprintHTML(wr io.Writer, e *dmapi.Entries, templateFile string) error {
	bytes, err := ioutil.ReadFile(templateFile)
	if err != nil {
		return err
	}

	html, err := htmlTemplate.New("html").Parse(string(bytes))
	if err != nil {
		return err
	}

	err = html.ExecuteTemplate(wr, "html", e)

	return err
}
