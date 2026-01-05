package output

import "fmt"

const (
	Blue  = "\033[34m"
	Red   = "\033[31m"
	Reset = "\033[0m"
)

func BlueText(s string) string {
	return fmt.Sprintf("%s%s%s", Blue, s, Reset)
}

func RedText(s string) string {
	return fmt.Sprintf("%s%s%s", Red, s, Reset)
}

// Status helpers used by main.go
func Up(msg string) {
	fmt.Printf("%s[UP]%s %s\n", Blue, Reset, msg)
}

func Down(msg string) {
	fmt.Printf("%s[DOWN]%s %s\n", Red, Reset, msg)
}

// Tool banner
func Banner() {
	fmt.Println(Blue + `
 ██████╗██╗   ██╗██████╗ ██╗  ██╗██████╗ ██████╗ 
██╔════╝╚██╗ ██╔╝██╔══██╗██║  ██║██╔══██╗██╔══██╗
██║      ╚████╔╝ ██████╔╝███████║██████╔╝██████╔╝
██║       ╚██╔╝  ██╔═══╝ ██╔══██║██╔═══╝ ██╔══██╗
╚██████╗   ██║   ██║     ██║  ██║██║     ██║  ██║
 ╚═════╝   ╚═╝   ╚═╝     ╚═╝  ╚═╝╚═╝     ╚═╝  ╚═╝
` + Reset)

	fmt.Println("        CYPH3R — Network Diagnostics Utility")
	fmt.Println("     ⚠ Educational & Professional Use Only ⚠")
	fmt.Println()
}
