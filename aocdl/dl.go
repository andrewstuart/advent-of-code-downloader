package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
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

	file, err := os.OpenFile(config.StoryOut, flags, 0660)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("file '%s' already exists; use '-force' to overwrite", config.StoryOut)
		}
		return err
	}
	defer file.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(io.MultiWriter(buf, file), resp.Body)
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	if config.TestOutput == "" && config.TestTemplate == "" {
		log.Printf("Not downloading test input for day %d. Config: %+v\n", config.Day, config)
		return nil
	}

	log.Printf("Downloading test input for day %d\n", config.Day)

	doc, err := goquery.NewDocumentFromReader(buf)
	if err != nil {
		return fmt.Errorf("error parsing story with goquery: %w", err)
	}

	preCode := doc.Find("pre code")
	test := preCode.First().Text()
	expected := preCode.Last().Text()

	// Write the test output to a file
	if config.TestOutput != "" {
		err := os.WriteFile(config.TestOutput, []byte(test), 0640)
		if err != nil {
			return fmt.Errorf("error writing test output: %w", err)
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
			"Test":   test,
			"Expect": expected,
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
		return fmt.Errorf("error creating request: %w", err)
	}

	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: config.SessionCookie,
	})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading input: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("error downloading input: %s", resp.Status)
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
		return fmt.Errorf("error opening output file: %w", err)
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	return nil
}
