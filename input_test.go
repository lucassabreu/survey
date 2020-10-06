package survey

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
	"github.com/stretchr/testify/assert"
)

func init() {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true
}

func TestInputRender(t *testing.T) {

	suggestFn := func(string) (s []string) { return s }

	tests := []struct {
		title    string
		prompt   Input
		data     InputTemplateData
		expected string
	}{
		{
			"Test Input question output without default",
			Input{Message: "What is your favorite month:"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: ", defaultIcons().Question.Text),
		},
		{
			"Test Input question output with default",
			Input{Message: "What is your favorite month:", Default: "April"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: (April) ", defaultIcons().Question.Text),
		},
		{
			"Test Input answer output",
			Input{Message: "What is your favorite month:"},
			InputTemplateData{Answer: "October", ShowAnswer: true},
			fmt.Sprintf("%s What is your favorite month: October\n", defaultIcons().Question.Text),
		},
		{
			"Test Input question output without default but with help hidden",
			Input{Message: "What is your favorite month:", Help: "This is helpful"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Test Input question output with default and with help hidden",
			Input{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] (April) ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput)),
		},
		{
			"Test Input question output without default but with help shown",
			Input{Message: "What is your favorite month:", Help: "This is helpful"},
			InputTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: ", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
		{
			"Test Input question output with default and with help shown",
			Input{Message: "What is your favorite month:", Default: "April", Help: "This is helpful"},
			InputTemplateData{ShowHelp: true},
			fmt.Sprintf("%s This is helpful\n%s What is your favorite month: (April) ", defaultIcons().Help.Text, defaultIcons().Question.Text),
		},
		{
			"Test Input question output with completion",
			Input{Message: "What is your favorite month:", Suggest: suggestFn},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for suggestions] ", defaultIcons().Question.Text, string(defaultPromptConfig().SuggestInput)),
		},
		{
			"Test Input question output with suggestions and help hidden",
			Input{Message: "What is your favorite month:", Suggest: suggestFn, Help: "This is helpful"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] [%s for suggestions] ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput), string(defaultPromptConfig().SuggestInput)),
		},
		{
			"Test Input question output with suggestions and default and help hidden",
			Input{Message: "What is your favorite month:", Suggest: suggestFn, Help: "This is helpful", Default: "April"},
			InputTemplateData{},
			fmt.Sprintf("%s What is your favorite month: [%s for help] [%s for suggestions] (April) ", defaultIcons().Question.Text, string(defaultPromptConfig().HelpInput), string(defaultPromptConfig().SuggestInput)),
		},
		{
			"Test Input question output with suggestions shown",
			Input{Message: "What is your favorite month:", Suggest: suggestFn},
			InputTemplateData{
				PageEntries:   core.OptionAnswerList([]string{"January", "February", "March", "etc..."}),
				SelectedIndex: 1,
				Answer:        "February",
			},
			fmt.Sprintf(
				"%s What is your favorite month: February [Use arrows to navegate, enter to select, type to complement answer]\n"+
					"  January\n%s February\n  March\n  etc...\n",
				defaultIcons().Question.Text, defaultPromptConfig().Icons.SelectFocus.Text,
			),
		},
		{
			"Test Input question output with suggestion complemented",
			Input{Message: "What is your favorite month:", Suggest: suggestFn},
			InputTemplateData{
				Answer: "February and",
			},
			fmt.Sprintf(
				"%s What is your favorite month: [%s for suggestions] February and",
				defaultIcons().Question.Text, defaultPromptConfig().SuggestInput,
			),
		},
	}

	for _, test := range tests {
		r, w, err := os.Pipe()
		assert.Nil(t, err, test.title)

		test.prompt.WithStdio(terminal.Stdio{Out: w})
		test.data.Input = test.prompt

		// set the runtime config
		test.data.Config = defaultPromptConfig()

		err = test.prompt.Render(
			InputQuestionTemplate,
			test.data,
		)
		assert.Nil(t, err, test.title)

		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)

		assert.Contains(t, buf.String(), test.expected, test.title)
	}
}

func TestInputPrompt(t *testing.T) {
	tests := []PromptTest{
		{
			"Test Input prompt interaction",
			&Input{
				Message: "What is your name?",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird")
				c.ExpectEOF()
			},
			"Larry Bird",
		},
		{
			"Test Input prompt interaction with default",
			&Input{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("")
				c.ExpectEOF()
			},
			"Johnny Appleseed",
		},
		{
			"Test Input prompt interaction overriding default",
			&Input{
				Message: "What is your name?",
				Default: "Johnny Appleseed",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("Larry Bird")
				c.ExpectEOF()
			},
			"Larry Bird",
		},
		{
			"Test Input prompt interaction and prompt for help",
			&Input{
				Message: "What is your name?",
				Help:    "It might be Satoshi Nakamoto",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("?")
				c.ExpectString("It might be Satoshi Nakamoto")
				c.SendLine("Satoshi Nakamoto")
				c.ExpectEOF()
			},
			"Satoshi Nakamoto",
		},
		{
			// https://en.wikipedia.org/wiki/ANSI_escape_code
			// Device Status Report - Reports the cursor position (CPR) to the
			// application as (as though typed at the keyboard) ESC[n;mR, where n is the
			// row and m is the column.
			"Test Input prompt with R matching DSR",
			&Input{
				Message: "What is your name?",
			},
			func(c *expect.Console) {
				c.ExpectString("What is your name?")
				c.SendLine("R")
				c.ExpectEOF()
			},
			"R",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			RunPromptTest(t, test)
		})
	}
}
