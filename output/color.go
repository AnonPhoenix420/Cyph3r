
package output

import "fmt"

/*
ANSI color definitions
*/
const (
	// Core colors
	Blue       = "\033[34m"
	Red        = "\033[31m"
	Reset      = "\033[0m"

	// Feature-specific themes
	YellowBold = "\033[1;33m" // GeoIP results
	PinkBold   = "\033[1;35m" // Port scan results

	// Banner aesthetic
	MetallicGreen = "\033[1;38;5;82m" // bold + 256-color green
)

/*
Text helpers
*/
func BlueText(s string) string {
	return fmt.Sprintf("%s%s%s", Blue, s, Reset)
}

func RedText(s string) string {
	return fmt.Sprintf("%s%s%s", Red, s, Reset)
}

func YellowBoldText(s string) string {
	return fmt.Sprintf("%s%s%s", YellowBold, s, Reset)
}

func PinkBoldText(s string) string {
	return fmt.Sprintf("%s%s%s", PinkBold, s, Reset)
}

/*
Status helpers used by main.go
*/
func Up(msg string) {
	fmt.Printf("%s[UP]%s %s\n", Blue, Reset, msg)
}

func Down(msg string) {
	fmt.Printf("%s[DOWN]%s %s\n", Red, Reset, msg)
}

/*
GeoIP helpers
*/
func GeoLabel(label string) string {
	return YellowBoldText(label)
}

/*
Port scan helpers
*/
func PortLabel(label string) string {
	return PinkBoldText(label)
}

/*
Tool banner
*/
func Banner() {
	fmt.Println(MetallicGreen + `
  _____    __     __   ______    _    _    _____   ______ 
 / ____|  \ \   / /  | ___ \  | |  | |  |___  |  | ___ \
| |       \ \_/ /   | |_/ /  | |__| |    / /   | |_/ / 
| |        \   /    |  __/   |  __  |  |_ \   |  _  \ 
| |____     | |     | |      | |  | |  ___) |  | | \ \
 \_____|    |_|     \_|      |_|  |_|  |____/   \_|  \_|
` + Reset)

	fmt.Println("        CYPH3R — Network Diagnostics Utility")
	fmt.Println("     ⚠ Educational & Professional Use Only ⚠")
	fmt.Println()
}
