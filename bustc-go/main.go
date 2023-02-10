package main

import (
	"fmt"
	"os"
)

func main() {
	r := NewLexer(os.Stdin)
	for r.Next() {
		t := r.Token()
		println(fmt.Sprintf("%d:%d", t.Position.Line, t.Position.Column), t.Kind.String(), t.Text)
	}
	if err := r.Err(); err != nil {
		println(err.Error())
	}
}
