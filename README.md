<div align="center">

***埼玉 Saitama CLI***
</div>
The ultimate CLI tool for your coding training regimen. One punch is all you need to manage your problems.



Saitama is a simple, powerful, and fun command-line tool designed to help you track your coding problems. Inspired by the "One-Punch Man," it provides a no-nonsense, focused environment to manage your practice sessions directly from the terminal.

**✨ Features**
Interactive Problem Adding: No need for clunky flags. An interactive survey guides you through adding new problems.

Clean Problem Listing: View all your saved problems in a neat, colorful, and easy-to-read format.

Random Problem Picker: Ready for a challenge? The pick command will select 5 random problems for you to tackle.

Tag-Based Summary: Quickly see which topics you're focusing on with a summary of tags and their counts.

Colorful & Themed Interface: A CLI that's not just functional, but also fun to use.

**🚀 Installation**
For Saitama to be available globally (so you can run saitama from any directory), you need to have Go installed and your GOPATH configured.

**The simplest way to install is with go install:**

```
go install github.com/Thedrogon/saitama@latest
```

(Note: Replace your-username/saitama with the actual path to your repository once it's on GitHub/GitLab etc.)

This command will download, compile, and install the saitama binary in your Go bin directory. If this directory is in your system's PATH, you're ready to go!

📖 Usage (The Hero's Manual)
Here are the main techniques you'll need for your training.

1. Add a New Problem (saitama add)
Run the add command to start an interactive survey. It will ask for the Problem ID, Name, and Tags (comma-separated).
```
$ saitama add
? Enter the problem ID (e.g., LC1): › LC141
? Enter the problem name (e.g., 'Two Sum'): › Linked List Cycle
? Enter tags (comma separated, e.g., array,hashmap): › linkedlist, twopointers
```

👊 ONE PUNCH! Problem 'Linked List Cycle' added successfully!

2. List All Problems (saitama list)
View your entire arsenal of problems.

```
$ saitama list
--- Your Coding Problems ---
ID: LC1        Name: Two Sum                                    Tags: [array hashmap]
ID: LC20       Name: Valid Parentheses                          Tags: [stack string]
ID: LC141      Name: Linked List Cycle                          Tags: [linkedlist twopointers]
```

3. Pick Your Daily Challenge (saitama pick)
Let Saitama choose 5 random problems for your daily training session.
```
$ saitama pick
🚀 Here are your 5 random problems for today! 🚀
1. ID: LC20       Name: Valid Parentheses                          Tags: [stack string]
2. ID: LC1        Name: Two Sum                                    Tags: [array hashmap]
... (and 3 more)
```

4. View Tag Summary (saitama tags)
Get a high-level overview of your problem categories.
```
$ saitama tags
--- Problems by Tag ---
array                - 1 problem
hashmap              - 1 problem
stack                - 1 problem
string               - 1 problem
linkedlist           - 1 problem
twopointers          - 1 problem
```
5. Get Help (saitama wiki or saitama --help)
Displays the help menu with all available commands.

🔧 Development
Interested in contributing? Here’s how you can get the project running locally.

Clone the repository:
```
git clone https://github.com/Thedrogon/saitama.git]
cd saitama
```

Build the binary:
``
go build .
``

Run directly:
``
./saitama list
``

***📝 License***
This project is licensed under the MIT License. See the LICENSE file for details.
