package main

import (
	"strconv"
)

func prs(x k) (r k) { // `p"…"
	t, n := typ(x)
	if t != C || n == atom {
		panic("type")
	}
	if n == 0 {
		dec(x)
		return mk(N, atom)
	}
	p := p{p: 8 + x<<2, e: n + 8 + x<<2, lp: 7 + x<<2, ln: 1}
	r = mk(L, 1)
	m.k[2+r] = mk(S, atom)
	mys(8+m.k[2+r]<<2, 0) // ;→`
	for p.p < p.e {       // ex;ex;…
		y := p.ex(p.noun())
		if x == 0 {
			break
		}
		r = lcat(r, y)
		if !p.t(sSem) {
			break
		}
		p.p++
	}
	if p.p < p.e {
		p.xx() // unprocessed input
	}
	_, n = typ(r)
	if n == 2 {
		inc(m.k[3+r])
		dec(r)
		r = m.k[3+r] // r[1]
	}
	dec(x)
	return r
}

type p struct {
	p  k // current position, m.c[p.p:...]
	m  k // pos after matched token (token: m.c[p.p:p.m])
	e  k // pos after last byte available
	ln k // current line number
	lp k // m.c index of last newline
}

func (p *p) t(f func([]c) int) bool { // test for next token
	p.w()
	if p.p == p.e {
		return false
	}
	if n := f(m.c[p.p:p.e]); n > 0 {
		p.m = p.p + k(n)
	}
	return p.m > p.p
}
func (p *p) a(f func([]c) k) (r k) { // accept, parse and advance
	n := p.m - p.p
	p.p = p.m
	return f(m.c[p.p-n : p.p])
}
func (p *p) w() { // remove whitespace and count lines
	for {
		switch p.get() {
		case 0:
			return
		case ' ', '\t', '\r':
		case '\n':
			if p.p != p.lp+1 {
				p.lp = p.p - 1
				p.ln++
			}
			p.p--
			return
		case '/':
			if p.p == p.lp+2 || m.c[p.p-2] == ' ' || m.c[p.p-2] == '\t' {
				for {
					if c := p.get(); c == 0 {
						return
					} else if c == '\n' {
						p.p--
						break
					}
				}
			} else {
				p.p--
				return
			}
		default:
			p.p--
			return
		}
	}
}
func (p *p) get() c {
	if p.p == p.e {
		return 0
	}
	p.p++
	return m.c[p.p-1]
}
func (p *p) xx() {
	panic("parse: " + string(m.c[p.lp+1:p.p+1]) + " <-")
}

// Parsers
func (p *p) ex(x k) (r k) {
	switch {
	//case p.t(sNam): // TODO ... atNoun
	}
	return x
}
func (p *p) noun() (r k) {
	switch {
	case p.t(sHex):
		r = p.a(pHex)
		return p.idxr(r)
	case p.t(sNum):
		r = p.a(pNum)
		for p.t(sNum) {
			y := p.a(pNum)
			rt, yt, _, _ := typs(r, y)
			if rt < yt {
				r = to(r, yt)
			} else if yt < rt {
				y = to(y, rt)
			}
			r = cat(r, y)
		}
		return p.idxr(r)
	case p.t(sStr):
		r = p.a(pStr)
		return p.idxr(r)
	case p.t(sSym):
		r = p.a(pSym)
		for p.p != p.e && m.c[p.p] == '`' { // `a`b`c without whitespace
			if p.t(sSym) {
				r = cat(r, p.a(pSym))
			}
		}
		return p.idxr(enl(r))
	case p.t(sNam):
		return p.idxr(p.a(pNam))
	}
	return mk(N, atom)
}
func (p *p) idxr(x k) (r k) { return x } // TODO

func pHex(b []byte) (r k) { // 0x1234 `c|`C
	if n := k(len(b)); n == 3 { // allow short form 0x1
		r = mk(C, atom)
		m.c[8+r<<2] = xtoc(b[2])
	} else if n%2 != 0 {
		panic("parse hex")
	} else {
		n = (n - 2) / 2
		r, b = mk(C, n), b[2:]
		rc := 8 + r<<2
		for i := k(0); i < n; i++ {
			m.c[rc+i] = (xtoc(b[2*i]) << 4) | xtoc(b[2*i+1])
		}
		if n == 1 {
			m.k[r] = C<<28 | atom
		}
	}
	return r
}
func xtoc(x c) c {
	switch {
	case x < ':':
		return x - '0'
	case x < 'G':
		return 10 + x - 'A'
	default:
		return 10 + x - 'a'
	}
}
func pNum(b []byte) (r k) { // 0|1f|-2.3e+4|1i2: `i|`f|`z
	for i, c := range b {
		if c == 'i' {
			r = to(pNum(b[:i]), Z)
			y := to(pNum(b[i+1:]), F)
			m.f[3+r>>1] = m.f[1+y>>1]
			dec(y)
			return r
		}
	}
	f := 0
	if len(b) > 1 {
		if c := b[len(b)-1]; c == 'f' || c == '.' {
			b = b[:len(b)-1]
			f = 1
		} else if b[0] == '.' {
			b = b[1:]
			f = 2
		}
	}
	if x, err := strconv.Atoi(string(b)); err == nil { // TODO remove strconv
		if f == 0 {
			r = mk(I, atom)
			m.k[2+r] = k(i(x))
		} else {
			r = mk(F, atom)
			if f == 1 {
				m.f[1+r>>1] = float64(x)
			} else {
				m.f[1+r>>1] = 0.1 * float64(x)
			}
		}
		return r
	}
	if x, err := strconv.ParseFloat(string(b), 64); err == nil { // TODO remove strconv
		r = mk(F, atom)
		m.f[1+r>>1] = x
		return r
	}
	panic("parse number")
}
func pStr(b []byte) (r k) { // "a"|"a\nbc": `c|`C
	r = pQot(b)
	if _, n := typ(r); n == 1 {
		m.k[r] = C<<28 | atom
	}
	return r
}
func pNam(b []byte) (r k) { // name: `n
	r = mk(S, atom)
	mys(8+r<<2, btou(b))
	return r
}
func pSym(b []byte) (r k) { // `name|`"name": `n
	if len(b) == 1 {
		r = mk(S, atom)
		mys(8+r<<2, 0)
	} else if len(b) > 1 || b[1] != '"' {
		r = mk(S, atom)
		mys(8+r<<2, btou(b[1:]))
	} else {
		r = pQot(b[1:])
		_, n := typ(r)
		m.k[r] = C<<28 | atom
		rc := 8 + r<<2
		mys(rc, btou(m.c[rc:rc+n]))
	}
	return r
}
func pQot(b []byte) (r k) { // "a\nb": `C
	r = mk(C, k(len(b)-2))
	p := 8 + r<<2
	q := false
	for _, c := range b[1 : len(b)-1] {
		if c == '\\' && !q {
			q = true
		} else {
			if q {
				q = false
				switch c {
				case 'r':
					c = '\r'
				case 'n':
					c = '\n'
				case 't':
					c = '\t'
				}
			}
			m.c[p] = c
			p++
		}
	}
	return srk(r, C, k(len(b)-2), k(p-(8+r<<2)))
}

// Scanners return the length of the matched input or 0
func sHex(b []byte) (r int) {
	if !(len(b) > 1 && b[0] == '0' && b[1] == 'x') {
		return 0
	}
	for i, c := range b[2:] {
		if crHx(c) == false {
			return 2 + i
		}
	}
	return len(b)
}

func sNum(b []byte) (r int) {
	n := sFlt(b)
	if len(b) > n && b[n] == 'i' {
		n += 1 + sFlt(b[n+1:])
	}
	return n
}
func sFlt(b []byte) (r int) { // -0.12e-12|1f
	if len(b) > 1 && b[0] == '-' {
		r++
	}
	for i := r; i < len(b); i++ {
		if c := b[i]; cr09(c) {
			r++
		} else {
			if c == '.' {
				r += 1 + sDec(b[i+1:])
			}
			break
		}
	}
	if len(b) > r && b[r] == 'e' {
		r += 1 + sExp(b[r+1:])
	}
	if len(b) > r && b[r] == 'f' {
		r++
	}
	if r == 1 && b[0] == '-' {
		return 0
	}
	return r
}
func sDec(b []byte) (r int) {
	for _, c := range b {
		if !cr09(c) {
			break
		}
		r++
	}
	return r
}
func sExp(b []byte) (r int) {
	if len(b) > 0 && (b[0] == '+' || b[0] == '-') {
		r++
	}
	for i := r; i < len(b); i++ {
		if c := b[i]; !cr09(c) {
			break
		}
		r++
	}
	return r
}
func sNam(b []byte) (r int) {
	for i, c := range b {
		if cr09(c) || craZ(c) { // TODO: dot?
			if i == 0 && cr09(c) {
				return 0
			}
		} else {
			return i
		}
	}
	return len(b)
}
func sStr(b []byte) (r int) {
	if len(b) < 2 || b[0] != '"' {
		return 0
	}
	q := false
	for i, c := range b[1:] {
		if !q && c == '\\' {
			q = true
		} else {
			if q == false && c == '"' {
				return i + 2
			}
			q = false
		}
	}
	return 0
}
func sSym(b []byte) (r int) { // `alp012|`"any\"thing"|`a.b.c
	if b[0] != '`' {
		return 0
	}
	if len(b) > 2 && b[1] == '"' {
		return 1 + sStr(b[1:])
	}
	for i, c := range b[1:] {
		if !(cr09(c) || craZ(c) || c == '.') {
			return 1 + i
		}
	}
	return len(b)
}
func sSem(b []byte) int {
	if b[0] == ';' || b[0] == '\n' {
		return 1
	}
	return 0
}
func cr09(c c) bool { return c >= '0' && c <= '9' }
func craZ(c c) bool { return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') }
func crHx(c c) bool { return cr09(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F') }
func cOps(c c) bool {
	for _, b := range []byte("+-%*|&<>=~,^#_$?@.") { // TODO: store in ktree for kwac compat
		if c == b {
			return true
		}
	}
	return false
}
