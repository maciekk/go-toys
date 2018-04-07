// Work in progress, a much earlier attempt at Go (and I think parsing of
// Steam data...)

package main

import (
	"bufio"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const sample_input = `
status {
    id : 123
    colour : GREEN
    data {
	a: 1
	b: 2
	c {
	    xyzzy : true
	}
    }
}
status {
    id : 456
    colour : RED
    data {
	a: x
	b: y
	c {
	    xyzzy : false
	}
    }
}
`

type bcl_value struct {
	is_plain bool
	// use this if is_plain == true
	plain_value string
	// else use this
	m map[string]bcl_value
}

// Assumes  the initial '{' already read.
func parse_clause(scanner *bufio.Scanner) map[string]bcl_value {
	re_plain := regexp.MustCompile(`\s*(\w+)\s*:\s*(.*)$`)
	re_multi := regexp.MustCompile(`\s*(\w+)\s*{$`)
	re_end_clause := regexp.MustCompile(`\s*}$`)
	m := make(map[string]bcl_value)
	for {
		scanner.Scan()
		l := scanner.Text()
		r := re_plain.FindStringSubmatch(l)
		if r != nil {
			fmt.Println("We got a plain setting:", r[1])
			m[r[1]] = bcl_value{is_plain: true, plain_value: r[2]}
			continue
		}
		r = re_multi.FindStringSubmatch(l)
		if r != nil {
			fmt.Println("We got a clause setting:", r[1])
			m[r[1]] = bcl_value{is_plain: false, m: parse_clause(scanner)}
			continue
		}
		if re_end_clause.FindString(l) != "" {
			// We hit the end of current clause.
			return m
		}
		fmt.Println("Unknown line format:", l)
	}
	panic("Ran out of lines while parsing a clause!")
}

func value(v bcl_value) string {
	if v.is_plain {
		return v.plain_value
	} else {
		// TODO(maciejk): provide actual values
		return "[MAP]"
	}
}

func main() {
	scanner := bufio.NewScanner(strings.NewReader(sample_input))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		l := scanner.Text()
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if l != "status {" {
			fmt.Println(l)
			panic("Expect start of 'status' clause")
		}
		m := parse_clause(scanner)
		fmt.Printf("STATUS: id:%s, colour:%s, data.c:%s\n",
			value(m["id"]), value(m["colour"]), value(m["data"].m["c"]))
		fmt.Println("Keys:", reflect.ValueOf(m).MapKeys())
	}
}
