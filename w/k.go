package main

import (
	"strconv" // to be removed
	"unsafe"
)

type c = byte
type k = uint32
type i = int32
type f = float64
type z = complex128
type s = string

const (
	C, I, F, Z, S, L, D, N k = 1, 2, 3, 4, 5, 6, 7, 8
	atom                   k = 0x0fffffff
	NaI                    i = -2147483648
)

type (
	f1 func(k, k)
	f2 func(k, k, k)
)
type slice struct {
	p uintptr
	l int
	c int
}

//                 C  I  F   Z  S  L  D  0  1  2  3  4
var lns = [13]k{0, 1, 4, 8, 16, 8, 4, 8, 4, 4, 4, 4, 4}
var m struct { // linear memory (slices share underlying arrays)
	c []c
	k []k
	f []f
	z []z
}

var cpx = []func(k, k){nil, cpC, cpI, cpF, cpZ, cpF, cpL}      // copy
var swx = []func(k, k){nil, swC, swI, swF, swZ, swF, swI, swI} // swap
var nax = []func(k){nil, naC, naI, naF, naZ, naS}              // set missing/nan
var eqx = []func(k, k) bool{nil, eqC, eqI, eqF, eqZ, eqS, nil} // equal
var ltx = []func(k, k) bool{nil, ltC, ltI, ltF, ltZ, ltS}      // less than
var gtx = []func(k, k) bool{nil, gtC, gtI, gtF, gtZ, gtS}      // greater than
var stx = []func(k, k) k{nil, nil, stI, stF, stZ, stS}         // tostring (assumes 56 bytes space at dst)
var tox = []func(k, k){nil, func(r, x k) { m.k[r] = k(i(m.c[x])) }, func(r, x k) { m.f[r] = f(m.c[x]) }, func(r, x k) { m.z[r] = complex(f(m.c[x]), 0) }, func(r, x k) { m.c[r] = c(m.k[x]) }, nil, func(r, x k) { m.f[r] = f(i(m.k[x])) }, func(r, x k) { m.z[r] = complex(f(m.k[x]), 0) }, func(r, x k) { m.c[r] = c(m.f[x]) }, func(r, x k) { m.k[r] = k(i(f(m.f[x]))) }, nil, func(r, x k) { m.z[r] = complex(m.f[x], 0) }, func(r, x k) { m.c[r] = c(m.f[x<<1]) }, func(r, x k) { m.k[r] = k(i(m.f[x<<1])) }, func(r, x k) { m.f[r] = m.f[x<<1] }}

func cpC(dst, src k)  { m.c[dst] = m.c[src] }
func cpI(dst, src k)  { m.k[dst] = m.k[src] }
func cpF(dst, src k)  { m.f[dst] = m.f[src] }
func cpZ(dst, src k)  { m.z[dst] = m.z[src] }
func cpL(dst, src k)  { inc(m.k[src]); cpI(dst, src) }
func swC(dst, src k)  { m.c[dst], m.c[src] = m.c[src], m.c[dst] }
func swI(dst, src k)  { m.k[dst], m.k[src] = m.k[src], m.k[dst] }
func swF(dst, src k)  { m.f[dst], m.f[src] = m.f[src], m.f[dst] }
func swZ(dst, src k)  { m.z[dst], m.z[src] = m.z[src], m.z[dst] }
func naC(dst k)       { m.c[dst] = 32 }
func naI(dst k)       { m.k[dst] = 0x80000000 }
func naF(dst k)       { u := uint64(0x7FF8000000000001); m.f[dst] = *(*f)(unsafe.Pointer(&u)) }
func naZ(dst k)       { naF(dst << 1); naF(1 + dst<<1) }
func naS(dst k)       { mys(dst<<2, uint64(' ')<<(56)) }
func eqC(x, y k) bool { return m.c[x] == m.c[y] }
func eqI(x, y k) bool { return i(m.k[x]) == i(m.k[y]) }
func eqF(x, y k) bool { return m.f[x] == m.f[y] }
func eqZ(x, y k) bool { return eqF(x<<1, y<<1) && eqF(1+x<<1, 1+y<<1) }
func eqS(x, y k) bool { return m.k[x<<1] == m.k[y<<1] && m.k[1+x<<1] == m.k[1+y<<1] }
func ltC(x, y k) bool { return m.c[x] < m.c[y] }
func gtC(x, y k) bool { return m.c[x] > m.c[y] }
func ltI(x, y k) bool { return i(m.k[x]) < i(m.k[y]) }
func gtI(x, y k) bool { return i(m.k[x]) > i(m.k[y]) }
func ltF(x, y k) bool { return m.f[x] < m.f[y] }
func gtF(x, y k) bool { return m.f[x] > m.f[y] }
func ltZ(x, y k) bool { // real than imag
	if ltF(x<<1, y<<1) {
		return true
	} else if eqF(x<<1, y<<1) {
		return ltF(1+x<<1, 1+y<<1)
	}
	return false
}
func gtZ(x, y k) bool { // real than imag
	if gtF(x<<1, y<<1) {
		return true
	} else if eqF(x<<1, y<<1) {
		return gtF(1+x<<1, 1+y<<1)
	}
	return false
}
func ltS(x, y k) bool { return sym(x<<3) < sym(y<<3) }
func gtS(x, y k) bool { return sym(x<<3) > sym(y<<3) }
func stI(dst, x k) k { // TODO remove strconv
	s := strconv.Itoa(int(i(m.k[x])))
	n := k(len(s))
	copy(m.c[dst:dst+n], []byte(s))
	return n
}
func stF(dst, x k) k { // TODO remove strconv
	s := strconv.FormatFloat(m.f[x], 'g', 6, 64)
	n := k(len(s))
	copy(m.c[dst:dst+n], []byte(s))
	return n
}
func stZ(dst, x k) k {
	n := stF(dst, x<<1)
	m.c[dst+n] = 'i'
	return 1 + n + stF(dst+1+n, 1+x<<1)
}
func stS(dst, x k) k {
	u := sym(x << 3)
	for i := k(0); i < 8; i++ {
		if c := c(u >> (8 * (7 - i))); c == 0 {
			return i
		} else {
			m.c[dst+i] = c
		}
	}
	return 8
}
func ptr(x, t k) k { // convert k address to type dependend index of data section
	switch t {
	case C:
		return (2 + x) << 2
	case I, L:
		return 2 + x
	case F, S:
		return (2 + x) >> 1
	case Z:
		return (4 + x) >> 2
	}
	println(t)
	panic("type")
}
func mv(dst, src k) {
	t, n := typ(src)
	ln := k(1 << bk(t, n))
	rc := m.k[1+dst]
	dst, src = dst<<2, src<<2
	copy(m.c[dst:dst+ln], m.c[src:src+ln]) // copy bucket
	dst >>= 2
	m.k[dst] = t<<28 | n // restore header
	m.k[1+dst] = rc
}

func ini() { // start function
	m.f = make([]f, 1<<13)
	msl()
	m.k[2] = 16
	p := k(64)
	m.k[8] = p
	m.k[p] = 8
	for i := 9; i < 16; i++ {
		p *= 2
		m.k[i] = p
		m.k[p] = k(i)
	}
	m.k[0] = (I << 28) | 31
	for i, u := range []c(`:+-*%&|<>=!~,^#_$?@.01234'/\`) {
		m.c[136+i] = u
	}
	// TODO: K tree
	// TODO: size vector
}
func msl() { // update slice header after increasing m.f
	f := *(*slice)(unsafe.Pointer(&m.f))
	i := *(*slice)(unsafe.Pointer(&m.k))
	i.l, i.c, i.p = f.l*2, f.c*2, f.p
	m.k = *(*[]k)(unsafe.Pointer(&i))
	b := *(*slice)(unsafe.Pointer(&m.c))
	b.l, b.c, b.p = f.l*8, f.c*8, f.p
	m.c = *(*[]c)(unsafe.Pointer(&b))
	zz := *(*slice)(unsafe.Pointer(&m.z))
	zz.l, zz.c, zz.p = f.l/2, f.l/2, f.p
	m.z = *(*[]z)(unsafe.Pointer(&zz))
}
func grw() { // double memory
	s := m.k[2]
	if 1<<k(s) != len(m.c) {
		panic("grow")
	}
	m.k[s] = k(len(m.c)) >> 2
	m.f = append(m.f, make([]f, len(m.f))...)
	msl()
	m.k[2] = s + 1
	m.k[1<<(s-2)] = s // bucket type of new upper half
}
func bk(t, n k) k {
	sz := lns[t]
	if n != atom {
		sz *= n
	}
	if sz > 1<<31 {
		panic("size")
	}
	return buk(sz + 8) // complex values have an additional 8 byte padding after the header (does not change bucket type)
}
func mk(t, n k) k { // make type t of len n (-1:atom)
	bt := bk(t, n)
	fb, a := k(0), k(0)
	for i := bt; i < 31; i++ { // find next free bucket >= bt
		if k(i) >= m.k[2] {
			panic("grow")
			grw() // TODO: run a gc cycle (merge blocks) before growing?
		}
		if m.k[i] != 0 {
			fb, a = i, m.k[i]
			break
		}
	}
	m.k[fb] = m.k[1+a] // occupy
	if p := m.k[fb]; p > 0 && p < 40 {
		panic("illegal free pointer")
	}
	for i := fb - 1; i >= bt; i-- { // split large buckets
		u := a + 1<<(i-2) // free upper half
		m.k[1+u] = m.k[i]
		m.k[i] = u
		m.k[u] = i
	}
	m.k[a] = n | t<<28 // ok for atoms
	m.k[a+1] = 1       // refcount
	return a
}

func typ(a k) (k, k) { // type and length at addr
	return m.k[a] >> 28, m.k[a] & atom
}
func typs(x, y k) (xt, yt, xn, yn k) { xt, xn = typ(x); yt, yn = typ(y); return }
func inc(x k) k {
	t, n := typ(x)
	switch t {
	case L:
		if n == atom {
			panic("type")
		}
		for i := k(0); i < n; i++ {
			inc(m.k[2+x+i])
		}
	case D:
		if n != atom {
			panic("type")
		}
		inc(m.k[2+x])
		inc(m.k[3+x])
	}
	m.k[1+x]++
	return x
}
func use(x, t, n k) k {
	if m.k[1+x] == 1 && bk(typ(x)) == bk(t, n) {
		m.k[x] = t<<28 | n
		return x
	}
	return mk(t, n)
}
func use2(x, y, t, n k) k {
	if m.k[1+x] == 1 && bk(typ(x)) == bk(t, n) {
		m.k[x] = t<<28 | n
		return x
	} else if m.k[1+y] == 1 && bk(typ(y)) == bk(t, n) {
		m.k[y] = t<<28 | n
		return y
	}
	return mk(t, n)
}
func decret(x, r k) k {
	if r != x {
		dec(x)
	}
	return r
}
func decret2(x, y, r k) k {
	if r != x {
		dec(x)
	}
	if r != y {
		dec(y)
	}
	return r
}
func dec(x k) {
	if m.k[x]>>28 == 0 || m.k[1+x] == 0 {
		//xxd()
		panic("unref " + hxk(x))
	}
	t, n := typ(x)
	switch t {
	case L:
		if n == atom {
			panic("type")
		}
		for i := k(0); i < n; i++ {
			dec(m.k[2+x+i])
		}
	case D:
		if n != atom {
			panic("type")
		}
		dec(m.k[2+x])
		dec(m.k[3+x])
	}
	m.k[1+x]--
	if m.k[1+x] == 0 {
		free(x)
	}
}
func free(x k) {
	t, n := typ(x)
	bt := bk(t, n)
	m.k[x] = bt
	m.k[x+1] = m.k[bt]
	m.k[bt] = x
}
func srk(x, t, n, nn k) (r k) { // shrink bucket
	if m.k[x]>>28 != t {
		panic("type")
	}
	if bk(t, nn) < bk(t, n) { // alloc not split: prevent small object accumulation
		r = mk(t, nn)
		ln := nn * lns[t]
		if t == Z {
			ln += 8
		}
		rc, xc := 8+r<<2, 8+x<<2
		copy(m.c[rc:rc+ln], m.c[xc:xc+ln])
		dec(x)
		return r
	}
	m.k[x] = t<<28 | nn
	return x
}
func to(x, rt k) (r k) { // numeric conversions for types CIFZ
	if rt == 0 {
		return x
	}
	t, n := typ(x)
	if rt == t {
		return x
	}
	var g func(k, k)
	if t == S && rt == I { // symbol comparison to bool
		g = func(r, x k) {
			if m.f[x] == 0 {
				m.k[r] = 0
			} else {
				m.k[r] = 1
			}
		}
	} else {
		g = tox[4*(t-1)+rt-1]
	}
	r = mk(rt, n)
	if n == atom {
		n = 1
	}
	xp, rp := ptr(x, t), ptr(r, rt)
	for i := k(0); i < k(n); i++ {
		g(rp+i, xp+i)
	}
	dec(x)
	return r
}
func nm(x, rt k, fx []f1) (r k) { // numeric monad
	t, n := typ(x)
	min := C
	if fx[C] == nil {
		min = I
	}
	if fx[I] == nil {
		min = F
	} // TODO: Z only for ff == nil ?
	if min > t { // uptype x
		x, t = to(x, min), min
	}
	if t == Z && fx[Z] == nil { // e.g. real functions
		x, t = to(x, F), F
	}
	r = use(x, t, n)
	if n == atom {
		n = 1
	}
	switch t {
	case L:
		for j := k(0); j < k(n); j++ {
			m.k[r+2+j] = nm(inc(m.k[j+2+x]), rt, fx)
		}
	case D:
		if r != x {
			m.k[2+r] = inc(m.k[2+x])
		}
		m.k[3+r] = nm(m.k[3+x], rt, fx)
	case C, I, F, Z:
		rp, xp, f := ptr(r, t), ptr(x, t), fx[t]
		for i := k(0); i < k(n); i++ {
			f(rp+i, xp+i)
		}
	default:
		panic("type")
	}
	decret(x, r)
	if rt != 0 && t != rt {
		r = to(r, rt)
	}
	return r
}
func nd(x, y, rt k, fx []f2) (r k) { // numeric dyad
	xt, yt, xn, yn := typs(x, y)
	t, n, sc := xt, xn, k(0)
	if yt > t {
		t = yt
	}
	if t == C && fx[C] == nil {
		t = I
	}
	if (t == I && fx[I] == nil) || (t == Z && fx[Z] == nil) {
		t = F
	}
	if xt < L && yt < L {
		if xn == atom {
			n, sc = yn, 1
		} else if yn == atom {
			n, sc = xn, 2
		} else if xn != yn {
			panic("size")
		}
		if xt != t {
			x, xt = to(x, t), t
		}
		if yt != t {
			y, yt = to(y, t), t
		}
	}
	r = use2(x, y, t, n)
	if n == atom {
		n = 1
	}
	switch xt {
	case C, I, F, Z, S:
		if xt == S && k(len(fx)) < S {
			panic("type")
		}
		f, xp, yp, rp := fx[t], ptr(x, t), ptr(y, t), ptr(r, t)
		if sc == 1 {
			for i := k(0); i < n; i++ {
				f(rp+i, xp, yp+i)
			}
		} else if sc == 2 {
			for i := k(0); i < n; i++ {
				f(rp+i, xp+i, yp)
			}
		} else {
			for i := k(0); i < n; i++ {
				f(rp+i, xp+i, yp+i)
			}
		}
	//case L:
	//case D:
	default:
		panic("type")
	}
	decret2(x, y, r)
	if rt != 0 && t > rt {
		r = to(r, rt)
	}
	return r
}
func idx(x, t k) i { // int from a numeric scalar (trunc, ignore imag)
	switch t {
	case C:
		return i(m.c[8+x])
	case I:
		return i(m.k[2+x])
	case F:
		return i(m.f[1+x>>1])
	case Z:
		return i(m.f[2+x])
	}
	panic("type")
}
func explode(x k) (r k) { // explode an array (or atom) to a list of atoms
	t, n := typ(x)
	if t == L {
		return x
	} else if t > L {
		panic("type")
	}
	if n == atom {
		n = 1
	}
	cp, xp := cpx[t], ptr(x, t)
	r = mk(L, n)
	for i := k(0); i < n; i++ {
		rk := mk(t, atom)
		rp := ptr(rk, t)
		cp(rp, xp+i)
		m.k[2+r+i] = rk
	}
	dec(x)
	return r
}
func uf(x k) (r k) { // unify lists if possible
	xt, xn := typ(x)
	if xt != L {
		return x
	}
	ut := k(0)
	for j := k(0); j < xn; j++ {
		t, n := typ(m.k[2+x+j])
		switch {
		case t >= L || n != atom:
			return x
		case j == 0:
			ut = t
		case t != ut:
			return x
		}
	}
	r = mk(ut, xn)
	cp, rp := cpx[ut], ptr(r, ut)
	for i := k(0); i < xn; i++ {
		cp(rp+i, ptr(m.k[2+x+i], ut))
	}
	dec(x)
	return r
}

func idn(x k) (r k) { return x } // :x
func flp(x k) (r k) { // +x
	t, n := typ(x)
	if t > L || n == atom { // tables are not implemented
		panic("type")
	} else if t < L {
		return x
	}
	nr, tt := k(0), L
	for j := k(0); j < n; j++ {
		tj, nj := typ(m.k[2+x+j])
		switch {
		case j == 0:
			nr, tt = nj, tj
		case nj != nr:
			panic("size") // k7 does extends for atoms, and nan-fills short arrays
		case j > 0 && tt != L && tt != tj:
			tt = L
		}
	}
	if tt > L {
		panic("type")
	}
	if nr == atom {
		nr = 1
	}
	cp := cpx[tt]
	r = mk(L, nr)
	for i := k(0); i < nr; i++ {
		rr := mk(tt, n)
		m.k[2+r+i] = rr
	}
	if tt == L {
		for k := k(0); k < n; k++ {
			col := explode(inc(m.k[2+x+k]))
			for i := uint32(0); i < nr; i++ {
				m.k[2+k+m.k[2+r+i]] = inc(m.k[2+col+i])
			}
			dec(col)
		}
		for i := k(0); i < nr; i++ {
			m.k[2+r+i] = uf(m.k[2+r+i])
		}
	} else {
		for i := uint32(0); i < nr; i++ { // Rik = +Xki (cdn.mos.cms.futurecdn.net/XTZkbu7r5c4LZQ5SMzJDbV-970-80.jpg)
			for k := k(0); k < n; k++ {
				cp(k+ptr(m.k[2+r+i], tt), i+ptr(m.k[2+x+k], tt))
			}
		}
	}
	dec(x)
	return r
}
func neg(x k) k { // -x
	return nm(x, 0, []f1{nil, func(r, x k) { m.c[r] = -m.c[x] }, func(r, x k) { m.k[r] = k(-i(m.k[x])) }, func(r, x k) { m.f[r] = -m.f[x] }, func(r, x k) { m.z[r] = -m.z[x] }})
}
func fst(x k) (r k) { // *x
	t, n := typ(x)
	if t == D {
		inc(m.k[3+x])
		r = fst(m.k[3+x])
		dec(x)
		return r
	}
	if n == atom {
		return x
	} else if n == 0 {
		panic("nyi: fst empty") // what to return? missing value? panic?
	}
	if t == L {
		r = m.k[2+x]
		inc(r)
		dec(x)
		return r
	}
	r = mk(t, atom)
	switch t {
	case C:
		m.c[8+r<<2] = m.c[8+x<<2]
	case I:
		m.k[2+r] = m.k[2+x]
	case F, S:
		m.f[1+r>>1] = m.f[1+x>>1]
	case Z:
		m.z[1+r>>2] = m.z[1+x>>2]
	default:
		panic("nyi")
	}
	dec(x)
	return r
}
func inv(x k) k { // %x
	return nm(x, 0, []f1{nil, nil, func(r, x k) { m.f[r] = 1 / m.f[x] }, func(r, x k) { m.z[r] = 1 / m.z[x] }})
}
func wer(x k) (r k) { // &x
	t, n := typ(x)
	if t != I {
		panic("type")
	} else if n == atom {
		n = 1
	}
	nn := k(0)
	for j := k(0); j < n; j++ {
		if p := i(m.k[2+x+j]); p < 0 {
			panic("domain")
		} else {
			nn += k(p)
		}
	}
	r = mk(I, nn)
	jj := k(0)
	for j := k(0); j < n; j++ {
		for p := k(0); p < m.k[2+x+j]; p++ {
			m.k[2+r+jj] = j
			jj++
		}
	}
	dec(x)
	return r
}
func rev(x k) (r k) { // |x
	t, n := typ(x)
	if n == atom || n < 2 {
		if t == D {
			r = use(x, t, n)
			if r == x {
				m.k[r+2] = rev(m.k[x+2])
				m.k[r+3] = rev(m.k[x+3])
			} else {
				m.k[r+2] = rev(inc(m.k[x+2]))
				m.k[r+3] = rev(inc(m.k[x+3]))
			}
			return decret(x, r)
		}
		return x
	}
	r = use(x, t, n)
	if t < D {
		cp, m := cpx[t], n
		if t == L {
			cp = cpI
		}
		src, dst := ptr(x, t), ptr(r, t)
		if src == dst {
			cp, m = swx[t], n/2
		}
		for i := k(0); i < m; i++ {
			cp(dst+n-1-i, src+i)
		}
	} else {
		panic("nyi")
	}
	if t == L && x != r {
		for i := k(0); i < n; i++ {
			inc(m.k[2+i+r])
		}
	}
	return decret(x, r)
}
func asc(x k) (r k) { // <x
	t, n := typ(x)
	if n == atom || t >= L { // k7 also sorts lists of different numeric types
		panic("type")
	}
	a := mk(I, atom)
	m.k[2+a] = n
	r = til(a)

	lt, sw := ltx[t], swI
	src, ind, dst := ptr(x, t), 2+r, ptr(r, I)
	for i := k(1); i < n; i++ { // insertion sort, should be replaced
		for j := k(i); j > 0 && lt(src+m.k[ind+j], src+m.k[ind+j-1]); j-- {
			sw(dst+j, dst+(j-1))
		}
	}
	dec(x)
	return r
}
func dsc(x k) (r k) { return rev(asc(x)) } // >x
func grp(x k) (r k) { // =x
	t, n := typ(x)
	if n == atom {
		panic("value")
	}
	eq := eqx[t]
	r = mk(D, atom)
	u := unq(inc(x)) // TODO: keys are sorted in k7
	l, nu := mk(L, m.k[u]&atom), m.k[u]&atom
	m.k[2+r], m.k[3+r] = u, l
	up, xp := ptr(u, t), ptr(x, t)
	b := mk(C, n) // boolean
	bc := 8 + b<<2
	for j := k(0); j < nu; j++ { // over ?x
		nr := k(0)
		for jj := k(0); jj < n; jj++ { // over x
			m.c[bc+jj] = 0
			if eq(up+j, xp+jj) {
				m.c[bc+jj] = 1
				nr++
			}
		}
		lj, p := mk(I, nr), k(0)
		for jj := k(0); jj < n; jj++ { // over x
			if m.c[bc+jj] == 1 {
				m.k[2+lj+p] = jj
				p++
			}
		}
		m.k[2+l+j] = lj
	}
	dec(b)
	dec(x)
	return r
}
func til(x k) (r k) { // !x
	t, n := typ(x)
	if n != atom {
		panic("nyi !a")
	} else if t == D {
		r = inc(m.k[2+x])
		dec(x)
		return r
	} else if t > Z {
		panic("type")
	}
	if nn := idx(x, t); nn < 0 {
		dec(x)
		return eye(k(-nn))
	} else {
		dec(x)
		return jota(k(nn))
	}
}
func jota(n k) (r k) { // !n
	r = mk(I, n)
	for j := k(0); j < n; j++ {
		m.k[2+r+j] = j
	}
	return r
}
func eye(n k) (r k) { // !-n
	r = mk(L, n)
	for j := k(0); j < n; j++ {
		rj := mk(I, n)
		m.k[2+r+j] = rj
		for jj := k(0); jj < n; jj++ {
			if j == jj {
				m.k[2+rj+jj] = 1
			} else {
				m.k[2+rj+jj] = 0
			}
		}
	}
	return r
}
func not(x k) (r k) { // ~x
	return nm(x, I, []f1{nil, func(r, x k) {
		if m.c[x] == 0 {
			m.c[r] = 1
		} else {
			m.c[r] = 0
		}
	}, func(r, x k) {
		if m.k[x] == 0 {
			m.k[r] = 1
		} else {
			m.k[r] = 0
		}
	}, func(r, x k) {
		if m.f[x] == 0 {
			m.f[r] = 1
		} else {
			m.f[r] = 0
		}
	}, func(r, x k) {
		if m.z[x] == 0 {
			m.z[r] = 1
		} else {
			m.z[r] = 0
		}
	}})
}
func enl(x k) (r k) { // ,x (collaps uniform)
	t, n := typ(x)
	if t < L && n == atom {
		r = use(x, t, 1)
		if r == x {
			m.k[r] = t<<28 | 1
			return r
		}
		cp := cpx[t]
		src, dst := ptr(x, t), ptr(r, t)
		cp(dst, src)
		dec(x)
		return r
	}
	r = mk(L, 1)
	m.k[2+r] = x
	return r
}
func enlist(x k) (r k) { // dont unify
	r = mk(L, 1)
	m.k[2+r] = x
	return r
}
func srt(x k) (r k) { return atx(x, asc(inc(x))) } // ^x  TODO: replace with a sort implementation
func cnt(x k) (r k) { // #x
	t, n := typ(x)
	r = mk(I, atom)
	if t == D {
		_, n = typ(m.k[x+2])
	} else if n == atom {
		n = 1
	}
	m.k[2+r] = k(i(n))
	dec(x)
	return r
}
func flr(x k) k { // _x
	return nm(x, I, []f1{nil, func(r, x k) { m.c[r] = m.c[x] }, func(r, x k) { m.k[r] = m.k[x] }, func(r, x k) {
		y := f(i(m.f[x]))
		if m.f[x] < y {
			y -= 1.0
		}
		m.f[r] = y
	}, nil}) // TODO: k7 does not convert c to i
}
func str(x k) (r k) { // $x
	t, n := typ(x)
	if t == C {
		return x
	}
	if t < L {
		st, xp := stx[t], ptr(x, t)
		if n == atom {
			r = mk(C, 56)
			r = srk(r, C, 56, st(8+r<<2, xp))
		} else {
			r = mk(L, n)
			for i := k(0); i < n; i++ {
				y := mk(C, 56)
				m.k[2+r+i] = srk(y, C, 56, st(8+y<<2, xp+i))
			}
		}
	} else {
		switch t {
		case L:
			r = use(x, L, n)
			for i := k(0); i < n; i++ {
				if r != x {
					inc(m.k[2+i+x])
				}
				m.k[2+i+r] = str(m.k[2+i+x])
			}
		case D:
			r = mk(D, atom)
			m.k[2+r] = inc(m.k[2+x])
			m.k[3+r] = str(inc(m.k[3+x]))
		case N:
			r = mk(C, 0)
		case N + 1:
			r = mk(C, 2)
			m.c[8+r<<2] = m.c[136+m.k[2+x]]
			m.c[9+r<<2] = ':'
		case N + 2:
			if f := m.k[2+x]; f < 40 {
				r = mk(C, 1)
			} else {
				r = mk(C, 2)
				m.c[9+r<<2] = ':'
			}
			m.c[8+r<<2] = m.c[136+m.k[2+x]-20]
		default:
			panic("nyi")
		}
	}
	return decret(x, r)
}
func unq(x k) (r k) { // ?x
	t, n := typ(x)
	if n == atom {
		panic("nyi") // overloads, random numbers?
	} else if t == D { // what does ?d do?
		panic("type")
	} else if n < 2 {
		return x
	}
	r = mk(t, n)
	eq, cp := eqx[t], cpx[t]
	if t == L {
		eq = func(x, y k) bool { return match(m.k[x], m.k[y]) }
	}
	src, dst := ptr(x, t), ptr(r, t)
	nn := k(0)
	for i := k(0); i < n; i++ { // quadratic, should be improved
		u := true
		srci := src + i
		for j := k(0); j < nn; j++ {
			if eq(srci, dst+j) {
				u = false
				break
			}
		}
		if u {
			cp(dst+nn, srci)
			nn++
		}
	}
	dec(x)
	return srk(r, t, n, nn)
}
func tip(x k) (r k) { // @x
	r = mk(S, atom)
	m.k[2+r] = 0
	m.k[3+r] = 0
	t, n := typ(x)
	dec(x)
	tns := "_cifzn.a_1234" // TODO k7 compatibility, function types 1..4?
	s := tns[t]
	if n != atom && t < L {
		s -= 32
	}
	mys(8+r<<2, uint64(s)<<56)
	if s == '_' {
		mys(8+r<<2, 0)
	}
	return r
}
func evl(x k) (r k) { // .x
	t, n := typ(x)
	if t != L {
		return x
	}
	if n == 0 {
		panic("evl empty list?") // what TODO?
	}
	v := m.k[2+x]
	vt, _ := typ(v)
	switch vt {
	case S:
		if n == 1 {
			inc(v)
			dec(x)
			return v
		}
	case N + 1, N + 2:
		if n == 1 {
			return x
		} else if n+1 < vt-N {
			panic("nyi projection")
		}
		a := mk(L, n-1)
		for i := int(n - 2); i >= 0; i-- { // right to left
			m.k[2+k(i)+a] = evl(inc(m.k[3+k(i)+x]))
		}
		switch vt {
		case N + 1:
			f := table[m.k[2+v]].(func(k) k)
			r = f(inc(m.k[2+a]))
		case N + 2:
			f := table[m.k[2+v]].(func(k, k) k)
			r = f(inc(m.k[2+a]), inc(m.k[3+a]))
		default:
			panic("nyi")
		}
		dec(a)
		dec(x)
		return r
	}
	println("evl vt", vt)
	panic("nyi")
	return x
}
func kst(x k) (r k) { // `k@x
	t, n := typ(x)
	atm := n == atom
	if atm {
		n = 1
	}
	if n == 0 {
		r = use(x, C, 0)
		rc, rn := 8+r<<2, k(0)
		switch t { // these could also be in the k-tree
		case C:
			rn = putb(rc, rn, []c(`""`))
		case I:
			rn = putb(rc, rn, []c("!0"))
		case F:
			rn = putb(rc, rn, []c("!0.0"))
		case Z:
			rn = putb(rc, rn, []c("!0i0"))
		case S:
			rn = putb(rc, rn, []c("0#`"))
		case L:
			rn = putb(rc, rn, []c("()"))
		case D:
			rn = putb(rc, rn, []c("()!()"))
		default:
			panic("nyi")
		}
		m.k[r] = C<<28 | rn
		return decret(x, r)
	}
	switch t {
	case C: // ,"a" "a" "ab" "a\nb" ,0x01 0x010203
		r = mk(C, 2+2*n) // for both "a\nb" or 0x01234 or ,"\n"(short enough)
		// no need to shrink: 2*(10+n) is never <= 10+2*n
		rc, rn, xc := 8+r<<2, k(0), 8+x<<2
		if n == 1 && !atm {
			rn = putc(rc, rn, ',')
		}
		rn = putc(rc, rn, '"')
		hex := false
		for i := k(0); i < n; i++ {
			c := m.c[xc+i]
			if c < 32 || c > 126 || c == '"' {
				if c, o := qot(c); o {
					rn = putc(rc, rn, '\\')
					rn = putc(rc, rn, c)
				} else {
					hex = true
					break
				}
			} else {
				rn = putc(rc, rn, c)
			}
		}
		rn = putc(rc, rn, '"')
		if hex {
			rn = 0
			if n == 1 && !atm {
				rn = 1
			}
			rn = putc(rc, rn, '0')
			rn = putc(rc, rn, 'x')
			for i := k(0); i < n; i++ {
				c1, c2 := hxb(m.c[xc+i])
				rn = putc(rc, rn, c1)
				rn = putc(rc, rn, c2)
			}
		}
		m.k[r] = C<<28 | rn
	case I, F, Z:
		r = mk(C, 0)
		if n == 1 && !atm {
			m.c[8+r<<2] = ','
			m.k[r] = C<<28 | 1
		}
		rr := mk(C, 56)
		st, xp, rrc := stx[t], ptr(x, t), 8+rr<<2
		sp := mk(C, 1)
		m.c[8+sp<<2] = ' '
		for i := k(0); i < n; i++ {
			rn := st(rrc, xp+i)
			m.k[rr] = C<<28 | rn
			r = cat(r, inc(rr))
			if i < n-1 {
				r = cat(r, inc(sp))
			}
		}
		dec(sp)
		m.k[rr] = C<<28 | 56
		dec(rr)
		if t == F {
			_, n = typ(r)
			rc, dot := 8+r<<2, false
			for i := k(0); i < n; i++ {
				if m.c[rc+i] == '.' {
					dot = true
					break
				}
			}
			if !dot {
				f := mk(C, 1)
				m.c[8+f<<2] = 'f'
				r = cat(r, f)
			}
		}
	case S:
		if atm || n == 1 {
			rr := mk(C, 0)
			sn, rrc, rn, q := stS(8+rr<<2, ptr(x, S)), 8+rr<<2, k(1), false
			for i := k(0); i < sn; i++ {
				c := m.c[rrc+i]
				if !(cr09(c) || craZ(c) || c == '.') {
					q = true
				}
				if _, o := qot(c); o {
					rn++
				}
				rn++
			}
			if q {
				rn += 2
			}
			if !atm {
				rn++
			}
			r = mk(C, rn)
			rc, rn := 8+r<<2, k(0)
			if !atm {
				rn = putc(rc, rn, ',')
			}
			rn = putc(rc, rn, '`')
			if q {
				rn = putc(rc, rn, '"')
			}
			for i := k(0); i < sn; i++ {
				c, o := qot(m.c[rrc+i])
				if o {
					rn = putc(rc, rn, '\\')
				}
				rn = putc(rc, rn, c)
			}
			if q {
				rn = putc(rc, rn, '"')
			}
			dec(rr)
		} else {
			r = mk(C, 0)
			ix := mk(I, atom)
			for i := k(0); i < n; i++ {
				m.k[2+ix] = i
				r = cat(r, kst(atx(inc(x), inc(ix))))
			}
			dec(ix)
		}
	case L:
		r = mk(C, 1)
		rc := 8 + r<<2
		m.c[rc] = '('
		if n == 1 {
			m.c[rc] = ','
		}
		y := mk(C, 1)
		m.c[8+y<<2] = ';'
		ix := mk(I, atom)
		for i := k(0); i < n; i++ {
			m.k[2+ix] = i
			r = cat(r, kst(atx(inc(x), inc(ix))))
			if i < n-1 {
				r = cat(r, inc(y))
			}
		}
		dec(ix)
		if n != 1 {
			m.c[8+y<<2] = ')'
			r = cat(r, y)
		} else {
			dec(y)
		}
	case D:
		r = mk(C, 0)
		rr, encl := kst(inc(m.k[x+2])), false
		_, nr := typ(rr)
		for i := k(0); i < nr; i++ {
			if c := m.c[8+i+rr<<2]; c == '!' || c == '#' {
				encl = true
				break
			}
		}
		y := mk(C, 1)
		if encl {
			m.c[8+y<<2] = '('
			r = cat(r, inc(y))
			r = cat(r, rr)
			m.c[8+y<<2] = ')'
			r = cat(r, inc(y))
		} else {
			r = cat(r, rr)
		}
		m.c[8+y<<2] = '!'
		r = cat(r, y)
		r = cat(r, kst(inc(m.k[x+3])))
	case N, N + 1, N + 2:
		r = str(inc(x))
	default:
		println("kst t/n", t, n)
		panic("nyi")
	}
	dec(x)
	return r
}
func putc(rc, rn k, c c) k { // assumes enough space
	m.c[rc+rn] = c
	return rn + 1
}
func putb(rc, rn k, b []c) k {
	rc += rn
	copy(m.c[rc:rc+k(len(b))], b)
	return rn + k(len(b))
}
func qot(c c) (c, bool) {
	if c == '"' {
		return c, true
	} else if c == '\n' {
		return 'n', true
	} else if c == '\t' {
		return 't', true
	} else if c == '\r' {
		return 'r', true
	}
	return c, false
}
func sym(x k) uint64 { return *(*uint64)(unsafe.Pointer(&m.c[x])) }
func mys(x k, u uint64) k {
	var b [8]c
	b = *(*[8]c)(unsafe.Pointer(&u))
	copy(m.c[x:x+8], b[:])
	return x
}
func btou(b []c) uint64 {
	if len(b) > 8 {
		panic("size")
	}
	var u uint64
	for i := k(0); i < k(len(b)); i++ {
		u |= uint64(b[i]) << (8 * c(7-i))
	}
	return u
}

func add(x, y k) (r k) { // x+y
	return nd(x, y, 0, []f2{nil, func(r, x, y k) { m.c[r] = m.c[x] + m.c[y] }, func(r, x, y k) { m.k[r] = m.k[x] + m.k[y] }, func(r, x, y k) { m.f[r] = m.f[x] + m.f[y] }, func(r, x, y k) { m.z[r] = m.z[x] + m.z[y] }})
}
func sub(x, y k) (r k) { // x-y
	return nd(x, y, 0, []f2{nil, func(r, x, y k) { m.c[r] = m.c[x] - m.c[y] }, func(r, x, y k) { m.k[r] = m.k[x] - m.k[y] }, func(r, x, y k) { m.f[r] = m.f[x] - m.f[y] }, func(r, x, y k) { m.z[r] = m.z[x] - m.z[y] }})
}
func mul(x, y k) (r k) { // x*y
	return nd(x, y, 0, []f2{nil, func(r, x, y k) { m.c[r] = m.c[x] * m.c[y] }, func(r, x, y k) { m.k[r] = m.k[x] * m.k[y] }, func(r, x, y k) { m.f[r] = m.f[x] * m.f[y] }, func(r, x, y k) { m.z[r] = m.z[x] * m.z[y] }})
}
func div(x, y k) (r k) { // x%y
	return nd(x, y, 0, []f2{nil, nil, nil, func(r, x, y k) { m.f[r] = m.f[x] / m.f[y] }, func(r, x, y k) { m.z[r] = m.z[x] / m.z[y] }})
}
func min(x, y k) (r k) { // x&y
	return nd(x, y, 0, []f2{nil, func(r, x, y k) {
		if m.c[x] < m.c[y] {
			m.c[r] = m.c[x]
		} else {
			m.c[r] = m.c[y]
		}
	}, func(r, x, y k) {
		if i(m.k[x]) < i(m.k[y]) {
			m.k[r] = m.k[x]
		} else {
			m.k[r] = m.k[y]
		}
	}, func(r, x, y k) {
		if m.f[x] < m.f[y] {
			m.f[r] = m.f[x]
		} else {
			m.f[r] = m.f[y]
		}
	}, func(r, x, y k) {
		if ltZ(x, y) {
			m.z[r] = m.z[x]
		} else {
			m.z[r] = m.z[y]
		}
	}})
}
func max(x, y k) (r k) { // x|y
	return nd(x, y, 0, []f2{nil, func(r, x, y k) {
		if m.c[x] > m.c[y] {
			m.c[r] = m.c[x]
		} else {
			m.c[r] = m.c[y]
		}
	}, func(r, x, y k) {
		if i(m.k[x]) > i(m.k[y]) {
			m.k[r] = m.k[x]
		} else {
			m.k[r] = m.k[y]
		}
	}, func(r, x, y k) {
		if m.f[x] > m.f[y] {
			m.f[r] = m.f[x]
		} else {
			m.f[r] = m.f[y]
		}
	}, func(r, x, y k) {
		if gtZ(x, y) {
			m.z[r] = m.z[x]
		} else {
			m.z[r] = m.z[y]
		}
	}})
}
func cmp(fx []func(k, k) bool, t k) func(k, k, k) {
	f := fx[t]
	switch t {
	case C:
		return func(r, x, y k) {
			if f(x, y) {
				m.c[r] = 1
			} else {
				m.c[r] = 0
			}
		}
	case I:
		return func(r, x, y k) {
			if f(x, y) {
				m.k[r] = 1
			} else {
				m.k[r] = 0
			}
		}
	case F:
		return func(r, x, y k) {
			if f(x, y) {
				m.f[r] = 1
			} else {
				m.f[r] = 0
			}
		}
	case Z:
		return func(r, x, y k) {
			if f(x, y) {
				m.z[r] = 1
			} else {
				m.z[r] = 0
			}
		}
	case S:
		return func(r, x, y k) {
			if f(x, y) {
				m.f[r] = 1 // non-zero
			} else {
				m.f[r] = 0
			}
		}
	}
	panic("reach")
}
func les(x, y k) (r k) { // x<y
	return nd(x, y, I, []f2{nil, cmp(ltx, C), cmp(ltx, I), cmp(ltx, F), cmp(ltx, Z), cmp(ltx, S)})
}
func mor(x, y k) (r k) { // x>y
	return nd(x, y, I, []f2{nil, cmp(gtx, C), cmp(gtx, I), cmp(gtx, F), cmp(gtx, Z), cmp(ltx, S)})
}
func eql(x, y k) (r k) { // x=y
	return nd(x, y, I, []f2{nil, cmp(eqx, C), cmp(eqx, I), cmp(eqx, F), cmp(eqx, Z), cmp(eqx, S)})
}
func key(x, y k) (r k) { panic("nyi") } // x!y
func mch(x, y k) (r k) { // x~y
	r = mk(C, atom)
	m.k[2+r] = 0
	if match(x, y) {
		m.c[8+r<<2] = 1
	}
	dec(x)
	dec(y)
	return r
}
func cat(x, y k) (r k) { // x,y
	xt, yt, xn, yn := typs(x, y)
	switch {
	case xt < L && yt == xt:
		return ucat(x, y, xt, xn, yn)
	case xt > L || yt > L: // TODO: cat D?
		panic("type")
	case xt != L:
		x = explode(x)
		xt, xn = typ(x)
	case yt != L:
		y = explode(y)
		xt, yn = typ(y)
	}
	if m.k[xt+1] > 1 || bk(L, xn+yn) > bk(L, xn) {
		r = mk(L, xn+yn)
		for j := k(0); j < xn; j++ {
			m.k[2+r+j] = inc(m.k[2+x+j])
		}
		dec(x)
	} else {
		r = x
	}
	for j := k(0); j < yn; j++ {
		m.k[2+r+xn+j] = inc(m.k[2+y+j])
	}
	dec(y)
	return r
}
func ucat(x, y, t, xn, yn k) (r k) { // x, y same type < L
	if xn == atom {
		xn = 1
	}
	if yn == atom {
		yn = 1
	}

	cp, xp := cpx[t], ptr(x, t)
	if m.k[x+1] > 1 || bk(t, xn+yn) != bk(t, xn) {
		r = mk(t, xn+yn)
		rp := ptr(r, t)
		for i := k(0); i < xn; i++ {
			cp(rp+i, xp+i)
		}
	} else {
		r = inc(x)
		m.k[r] = t<<28 | (xn + yn)
	}
	rp, yp := xn+ptr(r, t), ptr(y, t)
	for i := k(0); i < yn; i++ {
		cp(rp+i, yp+i)
	}
	dec(x)
	dec(y)
	return r
}
func lcat(x, y k) (r k) { // append anything to a list; no unify
	_, nl := typ(x)
	if m.k[x+1] == 1 && bk(L, nl) == bk(L, nl+1) {
		m.k[2+x+nl] = y
		m.k[x] = L<<28 | (nl + 1)
		return x
	}
	r = mk(L, nl+1)
	for i := k(0); i < nl; i++ {
		m.k[2+i+r] = inc(m.k[2+i+x])
	}
	m.k[2+nl+r] = y
	dec(x)
	return r
}
func ept(x, y k) (r k) { panic("nyi") } // x^y
func tak(x, y k) (r k) { panic("nyi") } // x#y
func drp(x, y k) (r k) { panic("nyi") } // x_y
func cst(x, y k) (r k) { panic("nyi") } // x$y
func fnd(x, y k) (r k) { panic("nyi") } // x?y
func atx(x, y k) (r k) { // x@y
	xt, xn := typ(x)
	yt, yn := typ(y)
	if xn == atom {
		panic("type") // TODO overloads
	}
	switch {
	case xt < L && yt == I:
		cp, na, xp := cpx[xt], nax[xt], ptr(x, xt)
		r = mk(xt, yn)
		if yn == atom {
			yn = 1
		}
		rp, yp := ptr(r, xt), 2+y
		for i := k(0); i < yn; i++ {
			if i >= xn {
				// na(rc + sz*j)
				na(rp + i)
			} else {
				// cp(rc+sz*j, xc+sz*m.k[idx+j])
				cp(rp+i, xp+m.k[yp+i])
			}
		}
		dec(x)
		dec(y)
		return r
	case xt == L && yt == I:
		if yn == atom {
			r = inc(m.k[2+x+m.k[2+y]])
		} else {
			r = mk(L, yn)
			for i := k(0); i < yn; i++ {
				m.k[2+r+i] = inc(m.k[2+x+m.k[2+y+i]])
			}
			r = uf(r)
		}
		dec(x)
		dec(y)
		return r
	// case xt == L:
	//	missing element for a list is nax[type of first element]
	default:
		panic("nyi atx")
	}
}
func cal(x, y k) (r k) { panic("nyi") } // x.y

func match(x, y k) (rv bool) { // recursive match
	if x == y {
		return true
	}
	t, n := typ(x)
	tt, nn := typ(y)
	if tt != t || nn != n {
		return false
	}
	if n == atom {
		n = 1
	}
	switch t {
	case L:
		for j := k(0); j < n; j++ {
			if match(m.k[2+x+j], m.k[2+y+j]) == false {
				return false
			}
		}
		return true
	case D:
		if match(m.k[2+x], m.k[2+y]) == false || match(m.k[3+x], m.k[3+y]) == false {
			return false
		}
		return true
	default:
		eq := eqx[t]
		if eq == nil {
			panic("type")
		}
		x, y = ptr(x, t), ptr(y, t)
		for j := k(0); j < n; j++ {
			if eq(x+j, y+j) == false {
				return false
			}
		}
		return true
	}
	return false
}

func hxb(x c) (c, c) { h := "0123456780abcdef"; return h[x>>4], h[x&0x0F] }
func hxk(x k) s {
	b := []c{'0', 'x', '0', '0', '0', '0', '0', '0', '0', '0'}
	for j := k(0); j < 4; j++ {
		n := 8 * (3 - j)
		b[2+2*j], b[3+2*j] = hxb(c((x & (0xFF << n)) >> n))
	}
	return s(b)
}

func isnan(x f) bool { return x != x }

func buk(x uint32) (n k) { // from https://golang.org/src/math/bits/bits.go (Len32)
	x--
	if x >= 1<<16 {
		x >>= 16
		n = 16
	}
	if x >= 1<<8 {
		x >>= 8
		n += 8
	}
	n += k(l8t[x])
	if n < 4 {
		return 4
	}
	return n
}

var l8t = [256]c{
	0x00, 0x01, 0x02, 0x02, 0x03, 0x03, 0x03, 0x03, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x04, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x05, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06, 0x06,
	0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07, 0x07,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
	0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08, 0x08,
}

var table []interface{} // function table :+-*%&|<>=!~,^#_$?@.01234
func init() {
	// 0-19 monads, 20-29 dyads, 30-34 ioverbs
	for _, f := range []interface{}{
		idn, flp, neg, fst, inv, wer, rev, asc, dsc, grp, til, not, enl, srt, cnt, flr, str, unq, tip, evl,
		nil, add, sub, mul, div, min, max, les, mor, eql, key, mch, cat, ept, tak, drp, cst, fnd, atx, cal,
	} {
		table = append(table, f)
	}
}
