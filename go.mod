module github.com/erikgeiser/promptkit

go 1.16

require (
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/charmbracelet/bubbles v0.8.0
	github.com/charmbracelet/bubbletea v0.14.1
	github.com/charmbracelet/lipgloss v0.3.0
	github.com/muesli/reflow v0.3.0
	github.com/muesli/termenv v0.9.0
	golang.org/x/sys v0.0.0-20210816183151-1e6c022a8912 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
)

replace github.com/charmbracelet/bubbletea v0.14.1 => github.com/erikgeiser/bubbletea v0.14.2-0.20210831114530-c8560c7d318e
