// i interpret
package i

import (
	"math"
	"math/cmplx"
	"reflect"
)

func P(s s) v {
	return prs(s)
}
func E(l v, a kt) v {
	if a == nil {
		a = make(kt)
	}
	if len(a) == 0 {
		kinit(a)
	}
	return eva(a, l)
}

type (
	i  = int
	f  = float64
	fv = []f
	z  = complex128
	zv = []z
	s  = string
	sv = []s
	v  = interface{}
	l  = []v
	d  = dict
	kt = map[v]v
)
type rV = reflect.Value
type rT = reflect.Type

func rval(x v) rV { return reflect.ValueOf(x) }
func rtyp(x v) rT { return reflect.TypeOf(x) }

var rTb = rtyp(true)
var rTf = rtyp(0.0)
var rTz = rtyp(complex(0, 0))
var rTs = rtyp("")

func cpy(x v) v {
	switch t := x.(type) {
	case f:
		return x
	case z:
		return x
	case s:
		return x
	case l:
		r := make(l, len(t))
		for i := range r {
			r[i] = cpy(t[i])
		}
		return r
	case d:
		r := t
		r.k = cpy(t.k).(l)
		r.v = cpy(t.v).(l)
	}
	v := rval(x)
	switch v.Kind() {
	case reflect.Slice:
		n := v.Len()
		r := reflect.MakeSlice(v.Type(), n, n)
		for i := 0; i < n; i++ {
			y := cpy(v.Index(i).Interface())
			r.Index(i).Set(rval(y))
		}
		return r.Interface()
	case reflect.Chan, reflect.Interface, reflect.Ptr, reflect.UnsafePointer:
		// TODO: allow pointer, with or without deep copy, depending on type?
		e("type") // TODO: these types should be returned verbatim
	case reflect.Map, reflect.Struct:
		e("assert") // already converted to dict
	}
	return x
}

func e(s string) interface{} { panic(s); return nil }

func ln(v interface{}) int {
	r := rval(v)
	if r.Kind() == reflect.Slice {
		return r.Len()
	}
	return -1
}

func lz(l interface{}) interface{} {
	return reflect.Zero(rtyp(l).Elem()).Interface()
}

/*
func mk(l interface{}, n int) interface{} {
	switch l.(type) {
	case float64:
		return make([]float64, n)
	case complex128:
		return make([]complex128, n)
	case string:
		return make([]string, n)
	}
	return reflect.MakeSlice(reflect.SliceOf(rtyp(l)), n, n).Interface()
}
*/

type dict struct {
	k, v          l
	f, s, u, p, g bool // flipped, sorted, uniq, parted, grouped
	t             reflect.Type
}

func md(x interface{}) (d, bool) {
	var d d
	switch m := x.(type) {
	case map[v]v:
		off := 0
		if h, ok := m["_"]; ok {
			hdr := h.(dict)
			d.f, d.s, d.u, d.p, d.g = hdr.f, hdr.s, hdr.u, hdr.p, hdr.g
			d.k = cpy(hdr.k).(l)
			off = 1
		}
		n := len(m) - off
		if off == 0 {
			d.k = make(l, n)
			d.v = make(l, n)
			i := 0
			for k, v := range m {
				d.k[i] = cpy(k)
				d.v[i] = cpy(v)
				i++
			}
		} else {
			d.v = make([]interface{}, n)
			for i, k := range d.k {
				d.v[i] = cpy(m[k])
			}
		}
		return d, true
	}
	// convert from user supplied map or struct

	v := rval(x)
	d.t = v.Type()
	if kind := v.Kind(); kind == reflect.Map || kind == reflect.Struct {
		n := 0
		if kind == reflect.Map {
			n = v.Len()
		} else {
			n = v.NumField()
		}
		d.k = make(l, n)
		d.v = make(l, n)
		if kind == reflect.Map {
			keys := v.MapKeys()
			for i, k := range keys {
				d.k[i] = cpy(k.Interface())
				d.v[i] = cpy(v.MapIndex(k).Interface())
			}
		} else {
			t := v.Type()
			for i := 0; i < n; i++ {
				d.k[i] = t.Field(i).Name
				d.v[i] = cpy(v.Field(i).Interface())
			}
		}
		return d, true
	}
	return d, false
}
func (d dict) mp() interface{} {
	if d.t == nil {
		r := make(map[v]v)
		r["_"] = dict{d.k, nil, d.f, d.s, d.u, d.p, d.g, nil}
		for i, k := range d.k {
			r[k] = d.v[i]
		}
		return r
	}

	// convert back to original map or struct type.
	v := reflect.New(d.t)
	v = v.Elem()
	if v.Kind() == reflect.Map {
		v = reflect.MakeMap(d.t)
		keytype := v.Type().Key()
		valtype := v.Type().Elem()
		for i, k := range d.k {
			rk := rval(k)
			if t := rk.Type(); t != keytype {
				rk = rk.Convert(t)
			}
			rv := rval(d.v[i])
			if t := rv.Type(); t != valtype {
				rv = rv.Convert(t)
			}
			v.SetMapIndex(rk, rv)
		}
		return v.Interface()
	} else if v.Kind() == reflect.Struct {
		for i, k := range d.k {
			f := v.FieldByName(rval(k).String())
			if f.Kind() == reflect.Slice {
				w := rval(d.v[i])
				if w.IsValid() == false {
					continue
				}
				sv := reflect.MakeSlice(f.Type(), w.Len(), w.Len())
				reflect.Copy(sv, w)
			} // TODO: make other types, that need it.
			f.Set(rval(d.v[i]))
		}
		return v.Interface()
	}
	return e("type")
}
func (d dict) at(key v) (int, v) {
	for i, k := range d.k {
		if k == key {
			return i, d.v[i]
		}
	}
	return -1, nil
}

// function krange(x, f) { var r=[]; for(var z=0;z<x;z++) { r.push(f(z)); } return k(3,r); }
func krange(n int, f func(int) v) l {
	l := make(l, n)
	for i := 0; i < n; i++ {
		l[i] = f(i)
	}
	return l
}

// function kmap (x, f) { return k(3, l(x).v.map(f)); }
func kmap(x v, f func(v, int) v) v {
	n := ln(x)
	if n < 0 {
		e("type")
	}
	for i := 0; i < n; i++ {
		set(x, i, f(at(x, i), i))
	}
	return x
}

// function kzip (x, y, f) { return kmap(sl(x,y), function(z, i) { return f(z, y.v[i]); }); }
func kzip(x, y v, f func(v, v) v) v {
	return kmap(sl, func(v v, i int) v {
		return f(v, at(y, i))
	})
}

func sl(x, y v) v {
	if ln(x) != ln(y) {
		e("len")
	}
	return x
}

func impl(v v, t reflect.Type) reflect.Method {
	if rtyp(v).Implements(t) {
		return t.Elem().Method(0)
	}
	return reflect.Method{}
}

func idx(v v) int { return int(re(v)) }

func re(v v) float64 {
	switch w := v.(type) {
	case float64:
		return w
	case bool:
		if w {
			return 1
		}
		return 0
	case int:
		return float64(w)
	case complex128:
		if cmplx.IsNaN(w) {
			return math.NaN()
		}
		return real(w)
	}
	r := rval(v)
	if k := r.Kind(); k == reflect.Bool {
		if r.Bool() {
			return 1
		}
		return 0
	} else if k < reflect.Uint {
		return float64(r.Uint())
	}
	return float64(r.Int()) // panics
}

func at(l v, i int) interface{} {
	switch v := l.(type) {
	case []interface{}:
		return v[i]
	case []float64:
		return v[i]
	case []complex128:
		return v[i]
	}
	v := rval(l)
	return v.Index(i).Interface()
}

func set(l interface{}, i int, v interface{}) {
	switch t := l.(type) {
	case []interface{}:
		t[i] = v
		return
	case []float64:
		t[i] = v.(float64)
		return
	case []complex128:
		t[i] = v.(complex128)
		return
	}
	rval(l).Index(i).Set(rval(v))
}

func kinit(a kt) {}
