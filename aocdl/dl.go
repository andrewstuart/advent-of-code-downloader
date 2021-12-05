package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"

	"github.com/PuerkitoBio/goquery"
)

func getStory(ctx context.Context, config *configuration) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d", config.Year, config.Day), nil)
	if err != nil {
		return err
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: config.SessionCookie,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	flags := os.O_WRONLY | os.O_CREATE
	if config.Force {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	file, err := os.OpenFile("a.html", flags, 0660)
	if os.IsExist(err) {
		file, err = os.OpenFile("b.html", flags, 0660)
	}
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file '%s' already exists; use '-force' to overwrite", config.Output)
		}
		return err
	}
	defer file.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(io.MultiWriter(buf, file), resp.Body)
	if err != nil {
		return err
	}

	if config.TestOutput == "" && config.TestTemplate == "" {
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return err
	}

	text := doc.Find("pre code").First().Text()

	if config.TestOutput != "" {
		err := ioutil.WriteFile(config.TestOutput, []byte(text), 0640)
		if err != nil {
			return err
		}
	}

	if config.TestTemplate != "" {
		tpl, err := template.ParseFiles(config.TestTemplate)
		if err != nil {
			return fmt.Errorf("error parsing template file: %w", err)
		}
		f, err := os.OpenFile(config.TestTemplateOutput, flags, 0640)
		if err != nil {
			return fmt.Errorf("error opening template output file: %w", err)
		}
		err = tpl.Execute(f, map[string]interface{}{
			"Config": config,
			"Test":   text,
		})
		if err != nil {
			return fmt.Errorf("error templating test output: %w", err)
		}
	}
	return nil
}

func download(ctx context.Context, config *configuration) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", config.Year, config.Day), nil)
	if err != nil {
		return err
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: config.SessionCookie,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	flags := os.O_WRONLY | os.O_CREATE
	if config.Force {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	file, err := os.OpenFile(config.Output, flags, 0660)
	if os.IsExist(err) {
		return fmt.Errorf("file '%s' already exists; use '-force' to overwrite", config.Output)
	} else if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
