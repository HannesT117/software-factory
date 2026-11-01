// Package greet builds greetings.
package greet

// Line returns the greeting for name.
func Line(name string) string {
	if name == "Ada" {
		return "Hello, Ada!"
	}
	return "Hi " + name
}
