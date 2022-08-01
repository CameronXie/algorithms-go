package tree

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

type Node interface {
	fmt.Stringer

	Left() Node
	Right() Node
}

func Print(node Node, w io.StringWriter) error {
	if isNil(node) {
		return errors.New("empty")
	}

	return printTree(node, w, "", "", true)
}

func printTree(node Node, w io.StringWriter, indent, position string, isLast bool) error {
	cornerSymbol := "|"
	if isLast {
		cornerSymbol = "`"
	}

	if position != "" {
		position = fmt.Sprintf("%v---%v: ", cornerSymbol, position)
	}

	if _, err := w.WriteString(fmt.Sprintf("%v%v%v\n", indent, position, node)); err != nil {
		return err
	}

	if position != "" {
		if !isLast {
			indent += "|   "
		} else {
			indent += "    "
		}
	}

	if left := node.Left(); !isNil(left) {
		if err := printTree(left, w, indent, "L", isNil(node.Right())); err != nil {
			return err
		}
	}

	if right := node.Right(); !isNil(right) {
		if err := printTree(right, w, indent, "R", true); err != nil {
			return err
		}
	}

	return nil
}

func isNil(data any) bool {
	return data == nil || reflect.ValueOf(data).IsNil()
}
