// main.go
package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func main() {
	// rand.Seed is deprecated and no longer needed in modern Go.
	// The math/rand package is automatically seeded.

	// ASCII Art Banner
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
	rootCmd.AddCommand(
		addCmd(),
		listCmd(),
		pickCmd(),
		tagsCmd(),
		searchCmd(),
		deleteCmd(),
		editCmd(),
		statsCmd(),
		importCmd(),
		exportCmd(),
		wikiCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		// Cobra already prints the error, so we just exit
		os.Exit(1)
	}
}

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

			// FIX: The correct way to handle survey errors/interrupts is to check for err != nil.
			err = survey.Ask(questions, &answers)
			if err != nil {
				color.Yellow("👋 Add operation cancelled.")
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

// Enhanced list command with better formatting
func listCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all saved coding problems",
		Long:  "Display all your coding problems in a beautiful table format",
		Run: func(cmd *cobra.Command, args []string) {
			problems, err := loadProblems()
			if err != nil {
				color.Red("❌ Error loading problems: %v", err)
				return
			}
			if len(problems) == 0 {
				color.Yellow("📝 No problems found yet!")
				color.Cyan("💡 Add your first problem with: saitama add")
				return
			}

			fmt.Println()
			color.HiCyan("═══════════════════════════════════════════════════════════════════════════════")
			color.HiCyan("                            🗂️  YOUR CODING ARSENAL 🗂️                        ")
			color.HiCyan("═══════════════════════════════════════════════════════════════════════════════")
			fmt.Println()

			fmt.Printf("%-15s %-50s %-30s\n", color.HiYellowString("🆔 ID"), color.HiWhiteString("📝 NAME"), color.HiGreenString("🏷️ TAGS"))
			color.HiBlack("---------------------------------------------------------------------------------------------------")

			for i, p := range problems {
				tagStr := "none"
				if len(p.Tags) > 0 {
					tagStr = strings.Join(p.Tags, ", ")
				}

				if i%2 == 0 {
					fmt.Printf("%-15s %-50s %-30s\n", color.CyanString(p.ID), color.WhiteString(p.Name), color.GreenString(tagStr))
				} else {
					fmt.Printf("%-15s %-50s %-30s\n", color.HiCyanString(p.ID), color.HiWhiteString(p.Name), color.HiGreenString(tagStr))
				}
			}

			fmt.Println()
			color.HiBlack("---------------------------------------------------------------------------------------------------")
			color.Magenta("📊 Total: %d problems", len(problems))
			fmt.Println()
		},
	}
	return cmd
}

// Enhanced pick command
func pickCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pick [number]",
		Short: "Pick random problems to solve",
		Long:  "Get a random selection of problems for your training session",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			problems, err := loadProblems()
			if err != nil {
				color.Red("❌ Error loading problems: %v", err)
				return
			}

			count := 5
			if len(args) > 0 {
				if c, err := strconv.Atoi(args[0]); err == nil && c > 0 {
					count = c
				}
			}

			if len(problems) == 0 {
				color.Yellow("📝 No problems found!")
				color.Cyan("💡 Add some problems first with: saitama add")
				return
			}

			if len(problems) < count {
				color.Yellow("⚠️  Not enough problems! You have %d, but requested %d", len(problems), count)
				color.Cyan("💡 Showing all %d problems instead:", len(problems))
				count = len(problems)
			}

			rand.Shuffle(len(problems), func(i, j int) { problems[i], problems[j] = problems[j], problems[i] })

			fmt.Println()
			color.HiMagenta("═══════════════════════════════════════════════════════════════")
			color.HiMagenta("           🎯 TODAY'S TRAINING SELECTION! 🎯                 ")
			color.HiMagenta("═══════════════════════════════════════════════════════════════")
			fmt.Println()

			for i := 0; i < count; i++ {
				p := problems[i]
				tagStr := "No tags"
				if len(p.Tags) > 0 {
					tagStr = strings.Join(p.Tags, " • ")
				}
				color.HiYellow("🥊 %d. %s", i+1, p.ID)
				color.White("   📝 %s", p.Name)
				color.Green("   🏷️  %s", tagStr)
				fmt.Println()
			}
			color.HiGreen("💪 Good luck with your training! ONE PUNCH! 🥊")
			fmt.Println()
		},
	}
	return cmd
}

// New search command
func searchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search problems by name or tag",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			problems, err := loadProblems()
			if err != nil {
				color.Red("❌ Error loading problems: %v", err)
				return
			}

			query := strings.ToLower(args[0])
			var matches []Problem

			for _, p := range problems {
				// Check name
				if strings.Contains(strings.ToLower(p.Name), query) {
					matches = append(matches, p)
					continue
				}
				// Check tags
				for _, tag := range p.Tags {
					if strings.Contains(strings.ToLower(tag), query) {
						matches = append(matches, p)
						break
					}
				}
			}

			if len(matches) == 0 {
				color.Yellow("🔍 No problems found matching: '%s'", query)
				return
			}

			fmt.Println()
			color.HiCyan("🔍 Found %d problems matching '%s':", len(matches), query)
			fmt.Println()

			for i, p := range matches {
				tagStr := strings.Join(p.Tags, ", ")
				color.Yellow("%d. %s - %s", i+1, p.ID, p.Name)
				color.Green("   Tags: %s", tagStr)
				fmt.Println()
			}
		},
	}
}

// New delete command - REFACTORED to use findProblemByID
func deleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a problem by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			problems, err := loadProblems()
			if err != nil {
				color.Red("❌ Error loading problems: %v", err)
				return
			}

			targetID := strings.ToUpper(args[0])
			problem, index := findProblemByID(problems, targetID)

			if index == -1 {
				color.Red("❌ Problem with ID '%s' not found", targetID)
				return
			}

			confirm := false
			prompt := &survey.Confirm{
				Message: fmt.Sprintf("Delete problem '%s - %s'?", problem.ID, problem.Name),
			}
			err = survey.AskOne(prompt, &confirm)
			if err != nil {
				if err == survey.ErrInterrupt {
					color.Yellow("👋 Delete operation cancelled.")
					return
				}
				color.Red("❌ Error during confirmation: %v", err)
				return
			}

			if !confirm {
				color.Yellow("❌ Deletion cancelled")
				return
			}

			// Efficiently delete element from slice
			newProblems := append(problems[:index], problems[index+1:]...)

			if err := saveProblems(newProblems); err != nil {
				color.Red("❌ Error saving: %v", err)
				return
			}

			color.Green("✅ Problem '%s' deleted successfully!", problem.ID)
		},
	}
}

// New edit command - REFACTORED to use findProblemByID
func editCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit a problem by ID",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			problems, err := loadProblems()
			if err != nil {
				color.Red("❌ Error loading problems: %v", err)
				return
			}

			targetID := strings.ToUpper(args[0])
			problem, index := findProblemByID(problems, targetID)

			if index == -1 {
				color.Red("❌ Problem with ID '%s' not found", targetID)
				return
			}

			answers := struct {
				Name string
				Tags string
			}{}

			questions := []*survey.Question{
				{
					Name:   "name",
					Prompt: &survey.Input{Message: "📝 New name:", Default: problem.Name},
				},
				{
					Name:   "tags",
					Prompt: &survey.Input{Message: "🏷️  New tags:", Default: strings.Join(problem.Tags, ", ")},
				},
			}

			err = survey.Ask(questions, &answers)
			if err != nil {
				if err == survey.ErrInterrupt {
					color.Yellow("👋 Edit operation cancelled.")
					return
				}
				color.Red("❌ Error during survey: %v", err)
				return
			}

			// Update name
			problems[index].Name = answers.Name

			// Process and update tags
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
			problems[index].Tags = tags

			if err := saveProblems(problems); err != nil {
				color.Red("❌ Error saving: %v", err)
				return
			}
			color.Green("✅ Problem '%s' updated successfully!", problem.ID)
		},
	}
}

// Enhanced tags command
func tagsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "tags",
		Short: "List all tags with problem counts",
		Run: func(cmd *cobra.Command, args []string) {
			problems, err := loadProblems()
			if err != nil {
				color.Red("❌ Error loading problems: %v", err)
				return
			}
			if len(problems) == 0 {
				color.Yellow("📝 No problems found!")
				return
			}

			tagCounts := make(map[string]int)
			for _, p := range problems {
				for _, tag := range p.Tags {
					tagCounts[tag]++
				}
			}

			fmt.Println()
			color.HiCyan("═══════════════════════════════════")
			color.HiCyan("        🏷️  TAG ANALYTICS 🏷️         ")
			color.HiCyan("═══════════════════════════════════")
			fmt.Println()

			if len(tagCounts) == 0 {
				color.Yellow("🏷️  No tags found")
				return
			}

			for tag, count := range tagCounts {
				fmt.Printf("%-20s %s\n", color.HiYellowString("🏷️  "+tag), color.GreenString("(%d problems)", count))
			}
			fmt.Println()
		},
	}
}

// New stats command
func statsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show detailed statistics",
		Run: func(cmd *cobra.Command, args []string) {
			problems, err := loadProblems()
			if err != nil {
				color.Red("❌ Error loading problems: %v", err)
				return
			}
			if len(problems) == 0 {
				color.Yellow("📝 No problems found!")
				return
			}

			tagCounts := make(map[string]int)
			totalTags := 0
			for _, p := range problems {
				for _, tag := range p.Tags {
					tagCounts[tag]++
					totalTags++
				}
			}

			fmt.Println()
			color.HiMagenta("═══════════════════════════════════════")
			color.HiMagenta("         📊 SAITAMA STATISTICS 📊        ")
			color.HiMagenta("═══════════════════════════════════════")
			fmt.Println()

			color.HiYellow("🗂️  Total Problems: %d", len(problems))
			color.HiYellow("🏷️  Unique Tags: %d", len(tagCounts))
			if len(problems) > 0 {
				color.HiYellow("📈 Average Tags per Problem: %.1f", float64(totalTags)/float64(len(problems)))
			}
			fmt.Println()
		},
	}
}

// New import command - NOW FUNCTIONAL
func importCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import <file>",
		Short: "Import problems from a JSON backup file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filePath := args[0]

			// Safety check
			confirm := false
			prompt := &survey.Confirm{Message: "This will merge imported problems with your current list. Continue?"}
			if err := survey.AskOne(prompt, &confirm); err != nil || !confirm {
				color.Yellow("Import cancelled.")
				return
			}

			importedProblems, err := importProblems(filePath)
			if err != nil {
				color.Red("❌ Error importing problems: %v", err)
				return
			}

			currentProblems, err := loadProblems()
			if err != nil {
				color.Red("❌ Error loading current problems: %v", err)
				return
			}

			// Merge logic (skip duplicates based on ID)
			existingIDs := make(map[string]bool)
			for _, p := range currentProblems {
				existingIDs[p.ID] = true
			}

			var mergedProblems []Problem
			mergedCount := 0
			for _, p := range importedProblems {
				if !existingIDs[p.ID] {
					mergedProblems = append(mergedProblems, p)
					mergedCount++
				}
			}

			finalProblems := append(currentProblems, mergedProblems...)

			if err := saveProblems(finalProblems); err != nil {
				color.Red("❌ Error saving merged list: %v", err)
				return
			}
			color.Green("✅ Successfully imported %d new problems from %s!", mergedCount, filePath)
		},
	}
}

// New export command - NOW FUNCTIONAL
func exportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "export <file>",
		Short: "Export all problems to a JSON file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			filePath := args[0]
			problems, err := loadProblems()
			if err != nil {
				color.Red("❌ Error loading problems for export: %v", err)
				return
			}

			if err := exportProblems(problems, filePath); err != nil {
				color.Red("❌ Error exporting problems: %v", err)
				return
			}
			color.Green("✅ Successfully exported %d problems to %s!", len(problems), filePath)
		},
	}
}

// Enhanced wiki command
func wikiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "wiki",
		Short: "Show all available commands",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Root().Help(); err != nil {
				color.Red("❌ Could not display help information.")
			}
		},
	}
}

