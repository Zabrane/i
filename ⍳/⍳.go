// ⍳ interpret
package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ktye/i"
)

// args:
// 0: read from stdin, execute each line, continue on error
// filename: execute file, exit on error
// else: exec argv

type v = interface{}
type l = []v

func main() {
	a := setup()

	if len(os.Args) > 1 {
		if b, err := ioutil.ReadFile(os.Args[1]); err == nil {
			i.E(i.P(string(b)), a)
		} else {
			p(i.E(i.P(jon(" ", os.Args[1:]).(string)), a))
		}
		return
	}

	r := bufio.NewScanner(os.Stdin)
	for r.Scan() {
		p(run(r.Text(), a))
	}
}

func run(t string, a map[v]v) (r interface{}) {
	defer func() {
		if c := recover(); c != nil {
			for _, s := range strings.Split(string(debug.Stack()), "\n") {
				if strings.HasPrefix(s, "\t") {
					println(s[1:])
				}
			}
			r = c
		}
	}()
	pr := i.P(t)
	return i.E(pr, a)
}

func p(x v) {
	if x == nil {
		return
	}
	s, o := x.(string)
	if !o {
		s = fmt(x).(string)
	}
	println(s)
}

var fmt func(v) v
var jon func(v, v) v
var num func(v) v

func init() {
	a := make(map[v]v)
	i.E(l{}, a)
	fmt = a["$:"].(func(x v) v)
	jon = a["jon"].(func(x, y v) v)
	num = a["num"].(func(x v) v)
}
