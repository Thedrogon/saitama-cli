// main.go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/AlecAivazis/survey/v2"
)

func main() {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	var rootCmd = &cobra.Command{
		Use:   "saitama",
		Short: color.CyanString("A CLI app to track your coding problems."),
		Long:  color.HiCyanString(`A colorful and simple CLI tool to manage, list, and randomly select coding problems to practice.`),
	}

	// Add commands to the root command
	rootCmd.AddCommand(addCmd(), listCmd(), pickCmd(), tagsCmd(), wikiCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// addCmd creates the "add" command
func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new coding problem via an interactive survey",
		Run: func(cmd *cobra.Command, args []string) {
			// This struct will hold the answers from the survey
			answers := struct {
				ID   string
				Name string
				Tags string // We'll ask for tags as a single string
			}{}

			// Define the interactive questions
			questions := []*survey.Question{
				{
					Name:      "id",
					Prompt:    &survey.Input{Message: "Enter the problem ID (e.g., LC1):"},
					Validate:  survey.Required,
				},
				{
					Name:      "name",
					Prompt:    &survey.Input{Message: "Enter the problem name (e.g., 'Two Sum'):"},
					Validate:  survey.Required,
				},
				{
					Name:     "tags",
					Prompt:   &survey.Input{Message: "Enter tags (comma separated, e.g., array,hashmap):"},
				},
			}

			// Ask the questions
			err := survey.Ask(questions, &answers)
			if err != nil {
				color.Red("An error occurred: %v", err)
				return
			}
            
            // Process tags from a comma-separated string to a slice
            tags := []string{}
            if answers.Tags != "" {
                tags = strings.Split(answers.Tags, ",")
                for i := range tags {
                    tags[i] = strings.TrimSpace(tags[i]) // Clean up whitespace
                }
            }


			problems, err := loadProblems()
			if err != nil {
				color.Red("Error loading problems: %v", err)
				return
			}
            
			newProblem := Problem{ID: answers.ID, Name: answers.Name, Tags: tags}
			problems = append(problems, newProblem)

			if err := saveProblems(problems); err != nil {
				color.Red("Error saving problem: %v", err)
				return
			}
			color.Green("👊 ONE PUNCH! Problem '%s' added successfully!", answers.Name)
		},
	}
	// Note: We no longer need flags for the add command!
	return cmd
}

// listCmd creates the "list" command
func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all saved coding problems",
		Run: func(cmd *cobra.Command, args []string) {
			problems, err := loadProblems()
			if err != nil {
				color.Red("Error loading problems: %v", err)
				return
			}
			if len(problems) == 0 {
				color.Yellow("No problems found. Add one with 'tracker add'.")
				return
			}

			color.Cyan("--- Your Coding Problems ---")
			for _, p := range problems {
				fmt.Printf("ID: %-10s Name: %-40s Tags: %v\n",
					color.HiYellowString(p.ID),
					color.WhiteString(p.Name),
					color.GreenString("%v", p.Tags),
				)
			}
		},
	}
}

// pickCmd creates the "pick" command for random selection
func pickCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "pick",
		Short: "Pick 5 random problems to solve",
		Run: func(cmd *cobra.Command, args []string) {
			problems, err := loadProblems()
			if err != nil {
				color.Red("Error loading problems: %v", err)
				return
			}
			if len(problems) < 5 {
				color.Red("Not enough problems to pick from. You need at least 5, but you have %d.", len(problems))
				return
			}

			// Shuffle the problems
			rand.Shuffle(len(problems), func(i, j int) {
				problems[i], problems[j] = problems[j], problems[i]
			})

			color.HiMagenta("🚀 Here are your 5 random problems for today! 🚀")
			for i := 0; i < 5; i++ {
				p := problems[i]
				fmt.Printf("%d. ID: %-10s Name: %-40s Tags: %v\n",
					i+1,
					color.HiYellowString(p.ID),
					color.WhiteString(p.Name),
					color.GreenString("%v", p.Tags),
				)
			}
		},
	}
}

// main.go

// tagsCmd creates the "tags" command
func tagsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tags",
		Short: "List all tags and their problem counts",
		Run: func(cmd *cobra.Command, args []string) {
			problems, err := loadProblems()
			if err != nil {
				color.Red("Error loading problems: %v", err)
				return
			}
			if len(problems) == 0 {
				color.Yellow("No problems found. Add one with 'saitama add'.")
				return
			}

			// Create a map to count occurrences of each tag
			tagCounts := make(map[string]int)
			for _, p := range problems {
				for _, tag := range p.Tags {
					tagCounts[tag]++
				}
			}

			color.Cyan("--- Problems by Tag ---")
			if len(tagCounts) == 0 {
				color.Yellow("No tags found.")
				return
			}

			for tag, count := range tagCounts {
				problemWord := "problem"
				if count > 1 {
					problemWord = "problems"
				}
				fmt.Printf("%-20s - %s %s\n",
					color.HiYellowString(tag),
					color.GreenString("%d", count),
					color.WhiteString(problemWord),
				)
			}
		},
	}
}

// main.go

// wikiCmd creates the "wiki" command
func wikiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "wiki",
		Short: "Display all available commands and operations",
		Run: func(cmd *cobra.Command, args []string) {
			// The .Root() method gets the top-level command ("saitama")
			// and .Help() displays its help message.
			cmd.Root().Help()
		},
	}
}