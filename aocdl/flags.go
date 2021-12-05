package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strconv"
	"text/template"
)

func addFlags(config *configuration) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)

	ignored := new(bytes.Buffer)
	flags.SetOutput(ignored)

	var flagCfg configuration
	flags.StringVar(&flagCfg.SessionCookie, "session-cookie", "", "")
	flags.StringVar(&flagCfg.Output, "output", "", "")
	flags.IntVar(&flagCfg.Year, "year", 0, "")
	flags.IntVar(&flagCfg.Day, "day", 0, "")

	flags.BoolVar(&flagCfg.Force, "force", false, "")
	flags.BoolVar(&flagCfg.Wait, "wait", false, "")

	flags.StringVar(&flagCfg.StoryOut, "story-output", "", "")
	flags.StringVar(&flagCfg.TestOutput, "test-output", "", "")
	flags.StringVar(&flagCfg.TestTemplate, "test-template", "", "")
	flags.StringVar(&flagCfg.TestTemplateOutput, "test-template-output", "", "")
	flags.StringVar(&flagCfg.Template, "template", "", "")
	flags.StringVar(&flagCfg.TemplateOutput, "template-output", "", "")

	flagErr := flags.Parse(os.Args[1:])

	if flagErr == flag.ErrHelp {
		fmt.Println(titleAboutMessage)
		fmt.Println(usageMessage)
		fmt.Println(repositoryMessage)
		os.Exit(0)
	}

	if flagErr != nil {
		fmt.Fprintln(os.Stderr, flagErr)
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, usageMessage)
		os.Exit(1)
	}

	config.merge(&flagCfg)
}

func parseIntFlag(text string) (int, error) {
	if text == "" {
		return 0, nil
	}
	// Parse in base 10.
	value, err := strconv.ParseInt(text, 10, 0)
	return int(value), err
}

func renderOutput(config *configuration) error {
	tmpl, err := template.New("output").Parse(config.Output)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	data := make(map[string]int)
	data["Year"] = config.Year
	data["Day"] = config.Day

	err = tmpl.Execute(buf, data)
	if err != nil {
		return err
	}

	config.Output = buf.String()

	return nil
}
