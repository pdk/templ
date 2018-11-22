package functions

import (
	"log"
	"strings"
	"time"
)

// Stringer is a thing that has a String method.
type Stringer interface {
	String() string
}

// Now returns the current timestamp.
func Now() string {
	return time.Now().Format(time.RFC822)
}

// Prefix prefixes stuff on front.
func Prefix(data interface{}, prefix string) interface{} {

	switch s := data.(type) {
	case string:
		return prefix + s
	case []string:
		n := []string{}
		for _, i := range s {
			n = append(n, prefix+i)
		}
		return n
	case Stringer:
		return prefix + s.String()
	case []interface{}:
		n := []string{}
		for _, i := range s {
			switch t := i.(type) {
			case string:
				n = append(n, prefix+t)
			case Stringer:
				n = append(n, prefix+t.String())
			default:
				log.Fatalf("cannot apply Prefix to %s", i)
			}
		}
		return n
	default:
		log.Fatalf("cannot apply Prefix to %s", s)
	}

	// never gonna happen
	return nil
}

// Postfix postfixes stuff on the end.
func Postfix(data interface{}, postfix string) interface{} {

	switch s := data.(type) {
	case string:
		return s + postfix
	case []string:
		n := []string{}
		for _, i := range s {
			n = append(n, i+postfix)
		}
		return n
	case Stringer:
		return s.String() + postfix
	case []interface{}:
		n := []string{}
		for _, i := range s {
			switch t := i.(type) {
			case string:
				n = append(n, t+postfix)
			case Stringer:
				n = append(n, t.String()+postfix)
			default:
				log.Fatalf("cannot apply Postfix to %s", i)
			}
		}
		return n
	default:
		log.Fatalf("cannot apply Postfix to %s", s)
	}

	// never gonna happen
	return nil
}

// Join joins a list of strings with a joiner.
func Join(data []string, joiner string) string {

	return strings.Join(data, joiner)
}

// PrePostJoin will prefix and postfix each item, then join with joiner.
func PrePostJoin(data []interface{}, prefix, postfix, joiner string) string {

	n := []string{}
	for _, i := range data {
		switch s := i.(type) {
		case string:
			n = append(n, prefix+s+postfix)
		case Stringer:
			n = append(n, prefix+s.String()+postfix)
		default:
			log.Fatalf("cannot apply PrePostJoin to %s", s)
		}
	}

	return strings.Join(n, joiner)
}
