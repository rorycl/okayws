package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// getFileAsMarkdown converts the content in file f from markdown format
// to html
func getFileAsMarkdown(f string) ([]byte, error) {

	source, err := os.ReadFile(f)
	if err != nil {
		return source, err
	}

	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			// html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.NewLinkify(
				extension.WithLinkifyAllowedProtocols([][]byte{
					[]byte("http:"),
					[]byte("https:"),
				}),
			),
		),
	)

	var b bytes.Buffer
	err = markdown.Convert(source, &b)
	return b.Bytes(), err
}

// filePathAsURL returns a file path as a simplified URL
func filePathAsURL(f string) (string, error) {

	// trim suffix
	f = strings.TrimSuffix(f, filepath.Ext(f))

	// break into path components
	parts := strings.Split(f, string(os.PathSeparator))
	parts = func(ss []string) []string {
		r := []string{}
		for _, s := range ss {
			if s == "" {
				continue
			}
			r = append(r, s)
		}
		return r
	}(parts)

	okChars := make(map[string]bool)
	for _, p := range strings.Split("abcdefghijklmnopqrstuvwxyz-", "") {
		okChars[p] = true
	}
	dashChars := []string{" ", "_"}
	normParts := func(s string) (string, error) {
		output := ""
	LOOP:
		for _, i := range strings.Split(strings.ToLower(s), "") {
			for _, d := range dashChars {
				if i == d {
					output = output + "-"
					continue LOOP
				}
			}
			_, ok := okChars[i]
			if !ok {
				continue
			}
			output = output + i
		}
		output = strings.Trim(output, "-")
		if len(output) < 1 {
			return "", errors.New("path resolved to empty string")
		}
		return output, nil
	}
	for i, p := range parts {
		var err error
		parts[i], err = normParts(p)
		if err != nil {
			return "", err
		}
	}
	return string(os.PathSeparator) + strings.Join(parts, string(os.PathSeparator)), nil
}

type c struct {
	Content template.HTML
}

// fillTemplate parses the go template file at f with the contents of
// the map called dict
func fillTemplate(f string, cs c) (string, error) {
	var b bytes.Buffer
	t := template.Must(template.ParseFiles(f))
	err := t.Execute(&b, cs)
	return b.String(), err
}

func main() {

	exampleFile := "design/design.md"
	md, err := getFileAsMarkdown(exampleFile)
	if err != nil {
		fmt.Printf("markdown error: %v\n", err)
		os.Exit(1)
	}

	fp, err := filePathAsURL(exampleFile)
	if err != nil {
		fmt.Printf("filepathasurl error: %v\n", err)
		os.Exit(1)
	}

	ch := c{Content: template.HTML(string(md))}
	tplOutput, err := fillTemplate("templates/home.html", ch)
	if err != nil {
		fmt.Printf("template error: %v\n", err)
		os.Exit(1)
	}

	generalRouter := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, tplOutput)
	}

	r := mux.NewRouter()
	r.HandleFunc(fp, generalRouter)
	// log.Printf("url: fp %s", fp)

	hdl := handlers.RecoveryHandler()(handlers.LoggingHandler(os.Stdout, r))
	server := &http.Server{
		Addr:    "127.0.0.1:4001",
		Handler: hdl,
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
