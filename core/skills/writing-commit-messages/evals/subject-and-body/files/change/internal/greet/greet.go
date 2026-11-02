// Package greet builds greetings.
package greet

// Line returns the greeting for name in the language lang.
func Line(name, lang string) string {
	if lang == "de" {
		return "Hallo, " + name + "!"
	}
	return "Hello, " + name + "!"
}
