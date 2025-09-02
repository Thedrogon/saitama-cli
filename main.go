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
	//rand.Seed(time.Now().UnixNano())
	banner := `
 ██████  █████  ██ ████████  █████  ███    ███  █████  
██      ██   ██ ██    ██    ██   ██ ████  ████ ██   ██ 
███████ ███████ ██    ██    ███████ ██ ████ ██ ███████ 
     ██ ██   ██ ██    ██    ██   ██ ██  ██  ██ ██   ██ 
███████ ██   ██ ██    ██    ██   ██ ██      ██ ██   ██ 
                                                       
        Your Coding Problem Training Partner 🥊
`       	

	var rootCmd = &cobra.Command{
		Use:   "saitama",
		Short: color.HiCyanString("A CLI app to track your coding problems."),
		Long: color.HiCyanString(banner) + "\n" +
			color.WhiteString("A powerful CLI tool to manage, organize, and randomly select coding problems.\n") +
			color.YellowString("Train like a hero! 💪"),
		Example: `  saitama add           # Add a new problem interactively
  saitama list          # List all problems
  saitama pick          # Get 5 random problems
  saitama search dp     # Search problems by tag
  saitama stats         # View problem statistics`,
	}

	// Add commands to the root command
	rootCmd.AddCommand(addCmd(), listCmd(), pickCmd(), tagsCmd(), wikiCmd())

	if err := rootCmd.Execute(); err != nil {
		
		os.Exit(1)
	}
}

// addCmd creates the "add" command
// addCmd creates the "add" command with improved UX
func addCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new coding problem interactively",
		Long:  color.HiGreenString("🔥 ONE PUNCH ADD! ") + "Add a new coding problem with an interactive questionnaire.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println()
			color.HiMagenta("═══════════════════════════════════════")
			color.HiMagenta("        🥊 ADD NEW PROBLEM 🥊         ")
			color.HiMagenta("═══════════════════════════════════════")
			fmt.Println()

			existingProblems, err := loadProblems()
			if err != nil {
				color.Red("❌ Error loading existing problems: %v", err)
				return
			}

			answers := struct {
				ID   string
				Name string
				Tags string
			}{}

			questions := []*survey.Question{
				{
					Name:   "id",
					Prompt: &survey.Input{Message: "🆔 Problem ID (e.g., LC1, CF123):"},
					Validate: survey.ComposeValidators(survey.Required, func(ans interface{}) error {
						id := ans.(string)
						if _, index := findProblemByID(existingProblems, strings.ToUpper(id)); index != -1 {
							return fmt.Errorf("ID '%s' already exists", id)
						}
						return nil
					}),
				},
				{
					Name:   "name",
					Prompt: &survey.Input{Message: "📝 Problem Name:"},
					Validate: survey.Required,
				},
				{
					Name:   "tags",
					Prompt: &survey.Input{Message: "🏷️  Tags (comma-separated):", Help: "e.g., array,hashmap,easy"},
				},
			}

			err = survey.Ask(questions, &answers)
			if err != nil {
				if err == survey.ErrInterrupt {
					color.Yellow("👋 Add operation cancelled.")
					return
				}
				color.Red("❌ Error during survey: %v", err)
				return
			}

			// Process tags
			var tags []string
			if answers.Tags != "" {
				tagList := strings.Split(answers.Tags, ",")
				for _, tag := range tagList {
					cleaned := strings.TrimSpace(strings.ToLower(tag))
					if cleaned != "" {
						tags = append(tags, cleaned)
					}
				}
			}

			// Create and save the problem
			newProblem := Problem{
				ID:        strings.ToUpper(answers.ID),
				Name:      answers.Name,
				Tags:      tags,
				DateAdded: time.Now(), // Set the date added
			}

			problems := append(existingProblems, newProblem)

			if err := saveProblems(problems); err != nil {
				color.Red("❌ Error saving problem: %v", err)
				return
			}

			fmt.Println()
			color.HiGreen("🎉 ONE PUNCH SUCCESS! 🎉")
			color.Green("✅ Problem '%s' added successfully!", answers.Name)
			color.Cyan("🆔 ID: %s", newProblem.ID)
			if len(tags) > 0 {
				color.Yellow("🏷️  Tags: %s", strings.Join(tags, ", "))
			}
			fmt.Println()
		},
	}
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