package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
)

func suggestFiles(toComplete string) []string {
	files, _ := filepath.Glob(toComplete + "*")
	return files
}

// the questions to ask
var q = []*survey.Question{
	{
		Name: "file",
		Prompt: &survey.Input{
			Message: "Which file should be read?",
			Suggest: suggestFiles,
			Help:    "Any file; do not need to exist yet",
		},
		Validate: survey.ComposeValidators(
			survey.Required,
			func(file interface{}) error {
				if file == "?" {
					return errors.New("? is not a valid file name")
				}

				return nil
			},
		),
	},
}

func main() {
	answers := struct {
		File string
	}{}

	// ask the question
	err := survey.Ask(q, &answers)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// print the answers
	fmt.Printf("File chosen %s.\n", answers.File)
}
