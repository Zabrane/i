package main

// go test
//  lots of output, last line should be ok
// go test -short
//  tests also https://raw.githubusercontent.com/kparc/ref/master/src/md/index.md
//  short:~short but it's go's non-default option

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"testing"
)

type (
	l  = []interface{}
	d  = [2]interface{}
	iv = []int
	sv = []string
)

func TestIni(t *testing.T) {
	t.Skip()
	ini()
	st := Stats()
	if st.UsedBlocks() != 1 {
		t.Fatal()
	}
	//mk(1, 9000)
	//pfl()
	//xxd()
}

func TestK(t *testing.T) {
	// t.Skip()
	testCases := []struct {
		x, r s
	}{
		//TODO b64@ b64?
		{"(!0) bin 3", "-1"}, // k7: 0?
		// parse error: {`{{z+y*x}/[0;x;y]}`, ""},
		//{"`csv?(\"1,2,3\";\"4,5,6\")", "(1 4;2 5;3 6)"},
		//{"`csv?(\"abc|def\";\"1.2|3\";\"4.5|6\")", "+`abc`def!(1.2 4.5;3 6)"},
		//{"`csv?(\"0.7|8\";\"1.2|3\";\"4.5|6\")", "+`abc`def!(0.7 1.2 4.5;8 3 6)"},
		//{`"fi"0:("abc|def";"1.2|3";"4.5|")`, "+`abc`def!(1.2 4.5;3 0N)"},
		//{`"fz ns"0:("a,b,c,d,e,f,g";"1,2,30,ign,abc,ABC,ign";"3,4,40,ign,def,DEF,ign")`, "+`a`b`e`f!(1 3f;2a30 4a40;`abc`def;(\"ABC\";\"DEF\"))"},
		// {"t:+`a`b!(1 2;3 4);t~t[]", "1"},
		//{"x:3 3#!9;x[1 2;0]", "0 1 2"}, // TODO matrix indexing
		//{"x:3 3#!9;x[;0 2]", "0 2"},
		//{"x:3 3#!9;x[1 2;1 2]:0;x", "(0 1 2;3 0 0;6 0 0)"},
		//{"x:3 3#!9;x[1 2;1 2]:2 2#1+!4", "(0 1 2;3 1 2;6 3 4)"},
		//{"x:3 3#!9;x[1 2;1 2]*:10", "(0 1 2;3 40 50;6 70 80)"},
		//{"t:+`a`b!(1 2;3 4);t[1 0;`b`a]", "(4 2;3 1)"},
		{"1", "1"},
		{"(0;0n;0N;0w;-0w)", "(0;0n;0N;0w;-0w)"},
		{"(0;0n;0N;0w;-0w)=(0;0n;0N;0w;-0w)", "1 1 1 1 1"},
		{"0.00123", "0.00123"},
		{"0.001234567890", "0.001234568"},
		{"12.3456789 1.23456789 0.123456789 0.0123456789 0.00123456789", "12.34568 1.234568 0.1234568 0.01234568 0.001234568"},
		{"- 12.3456789 1.23456789 0.123456789 0.0123456789 0.00123456789", "-12.34568 -1.234568 -0.1234568 -0.01234568 -0.001234568"},
		{"1.23456789e7 1.23456789e50 1.23456789e150", "1.234568e7 1.234568e50 1.234568e150"},
		{"1e2 1e20 1e25 1e42 1e50 1e100", "100 1e20 1e25 1e42 1e50 1e100"},
		{"1e-2 1e-20 1e-25 1e-42 1e-50 1e-100", "0.01 1e-20 1e-25 1e-42 1e-50 1e-100"},
		{".1 .12 .123 .1234", "0.1 0.12 0.123 0.1234"},
		{"-.1 -.12 -.123 -.1234", "-0.1 -0.12 -0.123 -0.1234"},
		{".95 -.985", "0.95 -0.985"},
		{"`p \"1+ \\\\x\"", "(+;1;(8::;`x))"}, // k7: (+;1;(;`x))
		{"`a", "`a"},
		{"`a`b", "`a`b"},
		{"*1 2 3", "1"},
		{"%4f", "0.25"},
		{"%4 5", "0.25 0.2"},
		{"%0i1", "1a270"},
		{"1 2,4 5", "1 2 4 5"},
		{"1 2+2 3", "3 5"},
		{"3+1 2 3", "4 5 6"},
		{"1 2 3+4", "5 6 7"},
		{"1 2 3+4f", "5 6 7f"},
		{"3@4", "3"},        // k7 class error
		{"3@!4", "3 3 3 3"}, // k7 class error
		{"(1 2;3 4)+2 3", "(3 4;6 7)"},
		{"2 3+(1 2;3 4)", "(3 4;6 7)"},
		{"(1 2;3 4)+2", "(3 4;5 6)"},
		{"2+(1 2;3 4)", "(3 4;5 6)"},
		{"3+0x02", "5"},
		{"(1 2;3 4)+5", "(6 7;8 9)"},
		{"2+(3f;,4 5)", "(5f;,6 7)"},
		{"1 2 3<4f", "1 1 1"},
		{"1 2>3 4f", "0 0"},
		{"`a`b`c=`c`b`a", "0 1 0"},
		{"`a`b!2 3", "`a`b!2 3"},
		{"`a!4", "(,`a)!,4"},
		{"`a!1 2 3", "(,`a)!,1 2 3"},
		{"`a!!0", "(,`a)!,!0"},
		{"`a`b`c!!3", "`a`b`c!0 1 2"},
		{"`a`b`c!(1;2;3)", "`a`b`c!1 2 3"},
		{"`a`b`c!(1;`a;3)", "`a`b`c!(1;`a;3)"},
		{"(`a`b!1 2),`c!3", "`a`b`c!1 2 3"},
		{"d:`a`b!(1 2;3 4);d+d", "`a`b!(2 4;6 8)"},
		{"t:+`a`b!(1 2;3 4);t+t", "+`a`b!(2 4;6 8)"},
		// {"(`a`b!(1 2;3 4))+`a`c!(5 6;7 8)", "`a`b`c!(6 8;3 4;7 8)"}, // nyi `a`b!..+`a`c!..
		{"(1;2+3;4)", "1 5 4"},
		{"1;2", "2"},
		{"1 2 3[0 2]", "1 3"},
		{"(1;(`a;3);4)[1;0]", "`a"},
		{"x:(1f;2;3);x -1 0 1 2 3", "(0n;1f;2;3;0n)"},
		{"x:(1;2f;3);x -1 0 1 2 3", "(0N;1;2f;3;0N)"},
		{"x:(`a;2f;`c);x -1 0 1 2 3", "(`;`a;2f;`c;`)"},
		{"x:((2;3f);2;3);x -1 0 1 2 3", `(" ";(2;3f);2;3;" ")`},
		{"x:1 2 3;x -1 0 1 2 3 4 5", "0N 1 2 3 0N 0N 0N"},
		{"(`a`b!1 2)[`b]", "2"},
		{"`a`b`c^`b", "`a`c"},
		{"`a`b`c^,`b", "`a`c"},
		{`"",!0`, `""`},
		{"!0x03", "0x000102"},
		{"!2 3", "(0 0 0 1 1 1;0 1 2 0 1 2)"},
		{"2#!3", "0 1"},
		{"2#!30", "0 1"},
		{"5#!3", "0 1 2 0 1"},
		{"-2#!3", "1 2"},
		{"-5#!3", "1 2 0 1 2"},
		{"0#!3", "!0"},
		{"2#!0", "0N 0N"},
		{"-3#!0", "0N 0N 0N"},
		// {"3#()", `("";"";"")`}, // or "   "? TODO check k7
		{"(!0)#1 2 3", "1"},
		{"2 3#!5", "(0 1 2;3 4 0)"},
		{"3 -2#!5", "(3 4;0 1;2 3)"},
		{"3 2 3#!5", "((0 1 2;3 4 0);(1 2 3;4 0 1);(2 3 4;0 1 2))"},
		{"2 3#(`a;2)", "((`a;2;`a);(2;`a;2))"},
		{"2 3#(`a;(1;`b))", "((`a;(1;`b);`a);((1;`b);`a;(1;`b)))"},
		{"-3#(1;(`a`b;4.5))", "((`a`b;4.5);1;(`a`b;4.5))"},
		{"`a#`a`c!1 2", "(,`a)!,1"},
		// {"`a`b`c#`a`c!1 2", "`a`b`c!1 0N 2"}, index error instead of na(k7)
		{"`b_`a`b`c!1 2 3", "`a`c!1 3"},
		{"`c`b_`a`b`c!1 2 3", "(,`a)!,1"},
		{"3 5_!10", "(3 4;5 6 7 8 9)"},
		{"0 3 5_!10", "(0 1 2;3 4;5 6 7 8 9)"},
		{"0 3 3 5_!6", "(0 1 2;!0;3 4;,5)"},
		{"0 3_!3", "(0 1 2;!0)"},
		{"3 3_!3", "(!0;!0)"},
		{"2 4_(1;`a;3f;\"x\";8)", `((3f;"x");,8)`},
		{"2 4_(1;2;`a;`b;8)", "(`a`b;,8)"},
		{`<(2f;1)`, "1 0"},
		{`<((1;2f;3f);(1;2f))`, "1 0"},
		{`<((1;2f;(2;"alpha"));(1;2f;(2;"alph")))`, "0 1"},
		{`<,""`, ",0"},
		{"=1 2 2 3 2", "1 2 3!(,0;1 2 4;,3)"},
		{"=\"alphabeta\"", "\"abehlpt\"!(0 4 8;,5;,6;,3;,1;,2;,7)"},
		{"=(1;2f;2f;3;3;1)", "(1;3;2f)!(0 5;3 4;1 2)"},
		{`=""`, `""!""`},
		{`=,""`, `(,"")!,,0`}, // (!0)!0#,!0
		{"`i$5", "5"},
		{"`c$5", "0x05"},
		{"`f$5 6", "5 6f"},
		{"`z$5", "5a"},
		{"`i$`xyz", "1"},
		{"`i$`", "0"},
		{"`n$\"alpha\"", "`alpha"},
		{"`n$\"x\"", "`x"},
		{"`$\"xyz\"", "`xyz"},
		{"`c$()", `""`},
		{"`c$(3 4;5 6)", "(0x0304;0x0506)"},
		{"`c$`a`b!3 4", "`a`b!0x0304"},
		{"`c$+`a`b!(3 4;5 6)", "+`a`b!(0x0304;0x0506)"},
		{"`i$0%0", "0N"},
		{"`f$0N", "0n"},
		{"`f$3f", "3f"}, // k7 does not allow *f
		{"`z$0N", "0na0n"},
		{"`i$\"123\"", "123"},
		{"`i$\"\"", "0N"},
		{"`f$\"\"", "0n"},
		{"`f$\"1.23\"", "1.23"},
		{"`z$\"\"", "0na0n"},
		{"`z$\"1a30\"", "1a30"},
		{`0$"abc"`, `""`},
		{`2$"abc"`, `"ab"`},
		{`-2$"abc"`, `"bc"`},
		{`4$"abc"`, `"abc "`},
		{`-5$"abc"`, `"  abc"`},
		{`3$"abc"`, `"abc"`},
		{"4$1 2 3", "1 2 3 0N"},
		{"-4$1 2 3", "0N 1 2 3"},
		{"()!()", "()!()"},
		{"(!0)!()", "(!0)!()"},
		{"(!0f)!()", "(!0f)!()"},
		{"(!0i)!()", "(!0a)!()"},
		{"(0#`)!()", "(0#`)!()"},
		{"+(!0)!()", "+(!0)!()"},
		{`#"\\"`, "1"},
		{`"\\"`, `"\\"`},
		{"2 3 1 4 0?3", "1"},
		{"2 3 1 4 0?3 1 0", "1 2 4"},
		{"4 2 1f?5 6f", "3 3"}, // TODO: or 0N?
		{"2+", "2+"},
		{"2~3", "0"},
		{"(*3.0)~3f", "1"},
		{"+[;3]", "+[;3]"},
		{"(2+).,3", "5"},
		{"(1;(2;3f)).(1 1)", "3f"},
		{"(1;(2;3f)).,1", "(2;3f)"},
		{"(1;(2;3f))@1", "(2;3f)"},
		{"+(1;3 4)", "(1 3;1 4)"},
		{"+(1 2 3;4 5 6;7 8 9)", "(1 4 7;2 5 8;3 6 9)"},
		{"+(1 2 3;4 5 6)", "(1 4;2 5;3 6)"},
		{"++(1 2 3;4 5 6)", "(1 2 3;4 5 6)"},
		{"+(,1;2 3;4 5 6)", "(1 2 4;0N 3 5;0N 0N 6)"},
		{"+(1 2;3 4 5;,6)", "(1 3 6;2 4 0N;0N 5 0N)"},
		{`+((1;2f);(3f;4))`, `((1;3f);(2f;4))`},
		{`+((1;2f);(3;4f))`, `(1 3;2 4f)`},
		{`+(("abc";"def");,"ghi")`, `(("abc";"ghi");("def";" "))`},
		{"{1+x}", "{1+x}"},
		{"1 3 4 8 10 bin 1 -3 3 5 9", "0 -1 1 2 3"},
		{"1 3 4 8 10f bin 3.5", "1"},
		{"`abc`def`ghi bin `d`e`j", "0 1 2"},
		{`(1;2f;3f) bin 2.5`, "1"},
		{`(1;2f;3f) bin (0;1;2.0;2.1;3f;3.1)`, "-1 0 1 1 2 2"},
		{"x:1;y:2;x+y", "3"},
		{"x:1;x", "1"},
		{"x:3;x:4", "4"},
		{"x:3", "3"},
		{"::x:3", "3"},
		{"x:3;-x", "-3"},
		{"x:1 2;x[0]:3 4;x", "(3 4;2)"},
		{".d:3;.d:1 2;.d", "1 2"},
		{`.d:"";.d`, `""`},
		{"g:{x+y};g[3;4]", "7"},
		{"x:1", "1"},
		{"e:3;f:4;e+f", "7"},
		{"{2*x}3", "6"},
		{"+", "+"},
		{"'", "'"},
		{"~`a``b`c", "0 1 0 0"},
		{"~(;)", "1 1"},
		{"@(;)[0]", "`"},
		{"#(;)[0]", "1"},
		{"~(;)[0]", "1"},
		{"~(+)", "0"},
		{"~!0", "!0"},
		{"`p\"1+2 3\"", "(+;1;2 3)"},
		{"+/(0x01;2;3f)", "6f"},
		{"+/1 2 3", "6"},
		{"*/4 5 6", "120"},
		{"+/,3", "3"},
		{"+/3", "3"},
		{"+/1 2 3f", "6f"},
		{"+/1 2 3i", "6a"},
		{"+/(1;(2f;3))", "(3f;4)"},
		{"4+/1 2 3f", "10f"},
		{"3+/2 3 4", "12"},
		{"3+/(0x02;3f;4)", "12f"},
		{"(,3)+/2 3 4", ",12"},
		{"3 4 5+/2 3", "8 9 10"},
		{"1 2 3+/4 5", "10 11 12"},
		{`-\4 8 9`, "4 -4 -13"},
		{`-\(2 3;5 8;2 1)`, `(2 3;-3 -5;-5 -6)`},
		{`3-\5 9 3`, "-2 -11 -14"},
		{`(,3)-\5 9 3`, "(,-2;,-11;,-14)"},
		{`2 3-\5 9 3`, "(-3 -2;-12 -11;-15 -14)"},
		{`2 3-\(5f;9;0x03)`, "(-3 -2f;-12 -11f;-15 -14f)"},
		{`3 2 1-\,3`, ",0 -1 -2"},
		{`2 3-\5 9 3`, "(-3 -2;-12 -11;-15 -14)"},
		{`+\4 8 9`, "4 12 21"},
		{`+\(4f;0x08;9)`, "4 12 21f"},
		{"+/`a`b`c!1 2 3", "6"},
		{"+\\`a`b`c!1 2 3", "`a`b`c!1 3 6"},
		{"%:'4 5", "0.25 0.2"},
		{"%:'`a`b!4 8", "`a`b!0.25 0.125"},
		{"12%'3 4", "4 3f"},
		{"12 15%'3", "4 5f"},
		{"1 2 3{2*x+y}'4 5 6", "10 14 18"},
		{"-':8 2 5", "8 -6 3"},
		{"-':`a`b`c!3 8 1", "`a`b`c!3 5 -7"},
		{"1 2 3+/:7 8", "(8 9 10;9 10 11)"},
		{`1 2 3+\:7 8`, "(8 9;9 10;10 11)"},
		{`{4>x}{x+1}/1`, "4"},
		{`{4>x}{x+1}\1`, "1 2 3 4"},
		{`{4>#x}{x,"k"}/"o"`, `"okkk"`},
		{"3(2,)/1", "2 2 2 1"},
		{`3(2,)\1`, "(1;2 1;2 2 1;2 2 2 1)"},
		{`3(-:)\1`, "1 -1 1 -1"},
		{"(1_)\\!3", "(0 1 2;1 2;,2;!0)"},
		{"(1_)/!3", "!0"},
		{`|+\!3`, `3 1 0`},
		{`"-"\:"ab-cd--ef-gh-"`, `("ab";"cd";"";"ef";"gh";"")`},
		{"`\\:\"ab\ncd\n\ne\n\n\"", `("ab";"cd";"";,"e";"")`}, // 1 trailing nl removed
		{`"x"/:(,"a";,"b";,"c")`, `"axbxc"`},
		{"`/:(\"aa\";\"bb\")", `"aa\nbb\n"`},
		{`2\:0`, "!0"},
		{`2 3\:0`, "0 0"},
		{`3 4\:30`, "1 2"},
		{`5 3 4\:30`, "2 1 2"},
		{`2\:234`, "1 1 1 0 1 0 1 0"},
		{`8\:234`, "3 5 2"},
		{`16\:234`, "14 10"},
		{`3/:4 3 2 1`, "142"},
		{`16/:14 10`, "234"},
		{`1 2 3/:3 2 1`, "25"},
		{`10 10 10\:3`, "0 0 3"},
		{`10 8 6\:215 345 7`, "(4 7 0;3 1 1;5 3 1)"},
		{`24 60 60/:2 23 12`, "8592"}, // 2h23m12s → 8592s (decode)
		{`24 60 60\:8592`, "2 23 12"}, // 8592s → 2h23m12s (encode)
		//{`1 -3 4\:30`, "0 -2 2"}, // k7: 0 -2 2
		{"(1;(2 3))[1][0]", "2"},
		{"$[0;2;3]", "3"},
		{"$[0;2;1;4]", "4"},
		{"$[0;1;0;2;3;4;5]", "4"},
		{"$[0;1;2]", "2"},
		{"$[0x00;1;2]", "2"},
		{"$[`;1;2]", "2"},
		{"$[`a;1;2]", "1"},
		{"$[\"\";1;2]", "2"},
		{"$[0 0;1;2]", "1"},
		{"x:`a`b!(1;2 3);x[`b]", "2 3"},
		{"x:`a`b!(1;2 3);x[`b;0]", "2"},
		{"x:`a`b!(1;2 3);x.b", "2 3"},
		{"x:`a`b!(1;2 3);x.b 1", "3"},
		{"a:`a`b!(1;`b`c!(2;3f));a.b.c", "3f"},
		{"()", "()"},
		{`"."\:""`, "()"},
		{"*()", `""`},
		{"*!0", "0N"},
		{"*!0f", "0n"},
		{`*""`, `" "`},
		{"*0#`", "`"},
		{"(a;b):3 4;a+b", "7"},
		{"(a;b):1;a+b", "2"},
		{"x:8;x-:", "-8"},
		{"x:8;x-:;x", "8"},
		{"x:`a`b!1 2;x[`b]:3;x", "`a`b!1 3"},
		{"y:`b;x:`a`b!1 2;x[y]:3;x", "`a`b!1 3"},
		{"x:`a`b!1 2;x.b:3;x", "`a`b!1 3"},
		{"x:`a`b!1 2;x[`c]:3;!x", "`a`b`c"},
		{"x:`a`b!1 2;x[`c`d]:3 4;x", "`a`b`c`d!1 2 3 4"},
		{"x:`a!(1;2f);x.b:0;x", "`a`b!((1;2f);0)"},
		{"x:!3;x[0]:4;x", "4 1 2"},
		{"x:!3;x[0 2]:4;x", "4 1 4"},
		{"x:!3;x[0 2]:-1 4;x", "-1 1 4"},
		{"x:(1;2f;`a);x[0]:3;x", "(3;2f;`a)"},
		{"x:(1;2f;`a);x[0 2]:-(1 2);x", "(-1;2f;-2)"},
		{"x:!4;@[`x;1 2;*;10];x", "0 10 20 3"},
		{"x:(1;2f;0x03);@[`x;0 1;*;10];x", "(10;20f;0x03)"},
		{"x:(1;2f;`a);x[0 2]:3 4f;x", "3 2 4f"},
		{"x:2f;x*:3;x", "6f"},
		{"(x;y):2 3f;(x;y)*:3;x+y", "15f"},
		{"x:!5;x[2 3]*:2;x", "0 1 4 6 4"},
		{"x:(1;2 3 4);x[1]:4 5;x", "(1;4 5)"},
		{"x:1 2;@[`x;0;3];x", "3 2"},
		{"x:1 2;@[`x;0;3]", "`x"},
		{"x:`a`b!3 4;x[!x]", "3 4"},
		{"@[1 2;0;3]", "3 2"},
		{"(3 3#!9)[1;0 2]", "3 5"},
		{"x:3 3#!9;x[1;2]:0;x", "(0 1 2;3 4 0;6 7 8)"},
		{"@[(1;2 3 4);1;5 5]", "(1;5 5)"},
		{"@[`a`b!1 2;`a;3]", "`a`b!3 2"},
		{"(4 3#1+!12)", "(1 2 3;4 5 6;7 8 9;10 11 12)"},
		{"(4 3#1+!12)[1]", "4 5 6"},
		{"(4 3#1+!12)[;,1]", "(,2;,5;,8;,11)"},
		{"(4 3#1+!12)[1;]", "4 5 6"},
		{"(4 3#1+!12)[,1;]", ",4 5 6"},
		{"(4 3#1+!12)[0 2;1]", "2 8"},
		{"(4 3#1+!12)[0 2;1 2 3]", "(2 3 0N;8 9 0N)"},
		{"(4 3#1+!12f)[-1 0 2;-1 1 2 3]", "(0n 0n 0n 0n;0n 2 3 0n;0n 8 9 0n)"},
		{"(1;2 3)[!2;!2]", "(1 1;2 3)"}, // k7: class error
		{"(,1;3 4)[;]", "(1 0N;3 4)"},
		{"x:`a`b!(1;`c`d!(2;3f));y:x;x.b.d:4;x.b.d+y.b.d", "7f"},
		{"x:3 3#!9;x[1;0 2]+:1 2;x", "(0 1 2;4 4 7;6 7 8)"},
		{"x:3 3#!9;.[`x;1 2;0];x", "(0 1 2;3 4 0;6 7 8)"},
		{"x:3 3#!9;.[`x;1 2;+;2];x", "(0 1 2;3 4 7;6 7 8)"},
		{"x:`a`b!(1;2 3);x.b[0]:4;x", "`a`b!(1;4 3)"},
		{"x:`a`b!(1;`c`d!1 2);x.b.c:4;x", "`a`b!(1;`c`d!4 2)"},
		{"x:`a`b!(1;`c`d!1 2);x.b[`c]:4;x", "`a`b!(1;`c`d!4 2)"},
		{"a:1;b:{;a:3;x+a}4;a+b", "8"},
		{"a:1;b:{;a::3;x+a}4;a+b", "10"},
		{"`p{x+y}", "((+;`x;`y);\"{x+y}\";`x`y;,`f)"},
		{"`p{1+x;a:y}", "((`;(+;1;`x);(:;`a;`y));\"{1+x;a:y}\";`x`y;`a`f)"}, // k7: .{...}
		{"{$[x<3;2+x;x]}1", "3"},
		{"{$[x<3;f x+1;x]}1", "3"},
		{"f:{2*x};n:3;f n+1", "8"},
		{"+`a`b!(1 2;3 4)", "+`a`b!(1 2;3 4)"},
		{"++`a`b!(1 2;3 4)", "`a`b!(1 2;3 4)"},
		{"@`a`b!(1 2;3 4)", "`a"},
		{"@+`a`b!(1 2;3 4)", "`A"},
		{"t:+`a`b!(1 2;3 4);t`b", "3 4"},
		{"t:+`a`b!(1 2;3 4);t 1", "`a`b!2 4"},
		{"t:+`a`b!(1 2;3 4);t 1 0", "+`a`b!(2 1;4 3)"},
		{"t:+`a`b!(1 2;3 4);t[1;`b`a]", "4 2"},
		{"t:+`a`b!(1 2;3 4);1+t", "+`a`b!(2 3;4 5)"},
		{"t:+`a`b!(1 2;3 4);t+1", "+`a`b!(2 3;4 5)"},
		{`."1+2"`, "3"},
		{`.(+;1;2)`, "3"},
		{`. 1 2!3 4`, "3 4"},
		{"2/(!10)", "0 0 1 1 2 2 3 3 4 4"},
		{`2 3 4\3`, "1 0 3"},
		{"3/6.2", "2"},
		{`3\0 1 2 3 4 5`, "0 1 2 0 1 2"},
		{`3\- 0 1 2 3 4 5`, "0 2 1 0 2 1"},
		{`8\10`, "2"},
		{`8\-10`, "6"}, // not -2
		{`.5p 0 .5p 0~1p\.5p 1p 1.5p 2p`, "1"},
		{`x:1 2 3;y:4 5 6;p:,/x,/:\:y;p`, "(1 4;1 5;1 6;2 4;2 5;2 6;3 4;3 5;3 6)"},
		{"(-)", "-"},
		{"(-).(1 2)", "-1"}, // TODO (-).1 2
		{"(+)1 2", "1 2+"},
		{"3 4 5  1 2", "4 5"},
		{`{x+1}3 4 5`, "4 5 6"},
		{`{x-1}3 4 5`, "2 3 4"},
		{`{x[0]-1}3 4 5`, "2"},
		{`h:{(y+(-/(*x)**x;2f**/*x);1+*|x)};g:{{{(4f>+/(*x)**x)&255>*|x}h[;y]/x}[(x;0);x]};*|g 0 0.65`, "22"},
		{"++", "+:+"},
		{"2*+", "2*+"},
		{"(.5*+).(2 4)", "3f"},
		{"(-*:)2 3", "-2"},
		{"-+/-[5 3 1;1 2 3]", "-3"},
		{"ina:3;ina", "3"},
		{`{3>*x}{x+1}\1`, "1 2 3"},
		{`{3>*x}{x+1}\1 0`, "(1 0;2 1;3 2)"},
		{"(!10f)@(3 5;6 8;(1;2 3))", "(3 5f;6 8f;(1f;2 3f))"},
		{"{x>3}#!7", "4 5 6"},
		{"{x>3}_!7", "0 1 2 3"},
		{"{x>3}#`a`b`c!2 5 6", "`b`c!5 6"},
		{"{x>3}_`a`b`c!2 5 6", "(,`a)!,2"},
		{"5'12 13 18 20", "10 10 15 20"},
		{"5'12.34", "10"},
		{"rand 2", "0.1833065 0.9024364"},
		{"rand -1", ",0.0856132"},
		{"rand -3", "0.0856132 -1.800962 1.698196"},
		{"rand 2i", "1.623531a273.0228 2.475345a136.6822"},
		{"rand \"a\"+`c$!26", `"e"`},
		{"3 rand 2f", "0.366613 1.804873 1.027283"},
		{"3 rand 10", "1 9 5"},
		{"4 rand \"a\"+`c$!26", `"exng"`},
		{"-5 rand !6", "3 2 4 0 1"},
		{"2p", "6.283185"},
		{".1p", "0.3141593"},
		{"0i1 1a90", "1a90 1a90"},
		{"sqrt 2 4", "1.414214 2"},
		{"sin 0 .5p", "0 1f"},
		{"cos 0 1p", "1 -1f"},
		{"abs -3 2 0", "3 2 0"},
		{"abs -3 2 0f", "3 2 0f"},
		{"abs 3i4 6", "5 6f"},
		{"2 abs (-3 4;3i4)", "(9 16;25f)"},
		{"exp 1 2", "2.718282 7.389056"},
		{"log exp 1 2", "1 2f"},
		{"2 3 exp 3", "8 27f"},
		{"2 log 8", "3f"},
		{"{x+y}/3", "3"},
		{`{x+y}/'("";!0;!0f)`, `(" ";0N;0n)`},
		{`+/'("";!0;!0f)`, `(" ";0;0f)`},
		{`-/'("";!0;!0f)`, `(" ";0;0f)`},
		{`*/'("";!0;!0f)`, `(" ";1;1f)`},
		{`%/'("";!0;!0f)`, `(" ";0N;0n)`},
		{`&/'("";!0;!0f)`, `(" ";1;1f)`},
		{`|/'("";!0;!0f)`, `(" ";0;0f)`},
		{`#/'("";!0;!0f)`, `(" ";0N;0n)`},
		{"2#8", "8 8"},
		{"0 1 2#'9", "(!0;,9;9 9)"},
		{"{*/x#2}'!8", "1 2 4 8 16 32 64 128"},
		{`7(2*)\1`, "1 2 4 8 16 32 64 128"},
		{`_.1+2 exp !8`, "1 2 4 8 16 32 64 128"},
		{`1,*\7#2`, "1 2 4 8 16 32 64 128"},
		{`*\1_&1 1 7`, "1 2 4 8 16 32 64 128"},
		{`*\&0 1 7`, "1 2 4 8 16 32 64 128"},
		{"norm 3 4", "5f"},
		{"norm 3i4", "5f"},
		{"norm 0x0304", "5f"},
		{"2 norm 3 4", "25f"},
		{"real 1i3 -2i5", "1 -2f"},
		{"imag 2a270", "-2f"},
		{"imag (3i 1a90;4i5 0i3)", "(0 1f;5 3f)"},
		{"conj 0i1 1a30 2a270", "1a270 1a330 2a90"},
		{"(180%1p)*phase 1a30 2a270", "30 -90f"},
		{"3i4 5i12 8i15~3 5 8 cmplx 4 12 15", "1"},
		{"0 cmplx 1 2 3", "1a90 2a90 3a90"},
		{"1 2 3f cmplx 0", "1 2 3a"},
		{"expi (1p%180)*30 60 90f", "1a30 1a60 1a90"},
		{"expi ((1p%180)*30 60 90f;0 1p)", "(1a30 1a60 1a90;1 1a180)"},
		{"1 2 3 expi (1p%180)*30 60 90f", "1a30 2a60 3a90"},
		{"3.2 expi (1p%180)*30 60 90f", "3.2a30 3.2a60 3.2a90"},
		{"1 2 3 expi .5p", "1a90 2a90 3a90"},
		{"+1 2 3", "1 2 3"},
		{"+,1 2 3", "(,1;,2;,3)"},
		{"1 2 3/4 5 6", "32"},
		{"(,1 2 3f)/4 5 6f", ",32f"},
		{"1 2 3a/+,4 5 6a", ",32a"},
		{"(,1 2 3f)/+,4 5 6f", ",,32f"},
		{"1 2f/(4 5 6f;7 8 9f)", "18 21 24f"},
		{"(1 2 3;4 5 6)/7 8 9", "50 122"},
		{"(+,1 2)/,3 4 5", "(3 4 5;6 8 10)"},
		{"(1 2;3 4;5 6)/(7 8 9 10;11 12 13 14)", "(29 32 35 38;65 72 79 86;101 112 123 134)"},
		{`(1 2f;-3 2f;-3 0f)\1 2 3`, "-0.6470588 0.4264706"},
		{`a:4 4#rand 16i;b:rand 4i;x:a\b;1e-14>|/abs b-a/x`, "1"},
		{"!-3", "(1 0 0;0 1 0;0 0 1)"},
		{"=3", "(1 0 0;0 1 0;0 0 1)"},
		{"!-3f", "(1 0 0f;0 1 0f;0 0 1f)"},
		{"=3a", "(1 0 0a;0 1 0a;0 0 1a)"},
		{"diag !3f", "(0 0 0f;0 1 0f;0 0 2f)"},
		{"diag !3i", "(0 0 0a;0 1 0a;0 0 2a)"},
		{"diag diag 5 3 1", "5 3 1"},
		{"cond (12 -51 4f;6 167 -68f;-4 24 -41f)", "22.4"},
		{"cond 2*(12 -51 4f;6 167 -68f;-4 24 -41f)", "22.4"},
		{"cond (1 0 0f;0 1 0f;0 0 1f)", "1f"},
		{"cond (1 1f;0 1f)", "4f"},
		{`A:(2i8 10i2 5i7;9i8 10i9 8i1;3i3 7i4 4i2;3i10 8i4 4i2;5i4 5i10 1i1);B:(8i7 2i6;9i5 7i6;10i2 8i4;9i5 3i6;1i8 6i3);A\B`, "(0.02775684a225.7687 0.591145a15.53089 0.7323478a346.7655;0.4919331a29.44904 0.2971654a306.5227 0.4829134a18.47038)"},
		{"+/!4024", "8094276"},
		{"+/!4024f", "8094276f"},
		{"+/0 cmplx!4024", "8094276a90"},
		{"avg (1 2 3;0x010203;1 2 3f;1 2 3i;0i1 0i2 0i3;15)", "(2f;2f;2f;2a;2a90;15)"},
		{"3 avg 3 1 0 0 2 3 4 5", "3 2 1.333333 0.3333333 0.6666667 1.666667 3 4"},
		{"3 avg 3 1 0 0 2 3 4 5i", "3 2 1.333333 0.3333333 0.6666667 1.666667 3 4a"},
		{"imag 3 avg 0 cmplx 3 1 0 0 2 3 4 5", "3 2 1.333333 0.3333333 0.6666667 1.666667 3 4"},
		{"0 avg 3 1 0 0 2 3 4 5", "3 2 1.333333 1 1.2 1.5 1.857143 2.25"},
		{"100 avg 3 1 0 0 2 3 4 5", "3 2 1.333333 1 1.2 1.5 1.857143 2.25"},
		{"imag 0 avg 0 cmplx 3 1 0 0 2 3 4 5", "3 2 1.333333 1 1.2 1.5 1.857143 2.25"},
		{"0.7 avg 3 1 0 0 2 3 4 5", "3 1.6 0.3 0 1.4 2.7 3.7 4.7"},
		{"0.7 avg 3 1 0 0 2 3 4 5i", "3 1.6 0.3 0 1.4 2.7 3.7 4.7a"},
		{"med (4 9 5 1 0;4 5 1 0)", "4 4"},
		{"0 med 1 2 3 -1 -3 5 0", "1 2 2 2 1 2 1"},
		{"0 med`abc`de`fgh`ac`def`fhg`fab", "`abc`de`de`de`de`def`def"},
		{"dev (1 2 3;2 4 4 4 5 5 5 7 9f)", "1 2f"},
		{"dev 1a30*(5*rand -1000)cmplx(3*rand -1000)", "4.963708a30.52689 3.059594a120.5269"}, // 5@30 3@120
		{"var (1 2 3;2 4 4 4 5 5 5 7 9f)", "1 4f"},
		{"1 9 8 4 4 2 var 0 2 5 2 8 4", "10.26667 7.9 1.6"},
		{".5 med 4 5 1 0", "2.5"},
		{".95 med rand -1000", "1.697773"},      // (sqrt 2)*erfinv(-1+2*p) = 1.6448536269514722
		{"-.95 med 3+2*rand -1000", "6.331052"}, // 6.289707
		{"-0.95 med rand 1000i", "2.462811"},    // sqrt(-2*log(1-p)) = 2.448 (binormal 95%)
		{`"aaabaaabaaa" find "aa"`, "(0 2;4 2;8 2)"},
		{`"aabcde" find "abc"`, ",1 3"},
		{`"abc" find ,"d"`, "()"},
		{"`m@`abc`d!(1.23;4)", `("abc:1.23";"d  :4")`},
		{"`m@+`abc`d!(1 2;3 4)", `("abc d";"--- -";"1   3";"2   4")`},
		{"`m@+`abc`d!(1 2;3.2 2p)", `("abc d       ";"--- --------";"1   3.2     ";"2   6.283185")`},
		{"`m@((1.2 3);(4 5f))", `("1.2 3";"4   5")`},
		{"0+`hex?\"ffFF\"", "255 255"},
		{"`hex?`hex \"abc\"", `"abc"`},
		{"`hex?\"0x1234\"", "0x1234"},
		{"`hex \"abc\"", `"616263"`},
		{"`hex \"a\"", `"61"`},
		{"`csv@(!3;4+!3f;`a`b`c)", `("0,4,a";"1,5,b";"2,6,c")`},
		{"`csv@+`I`F`S!(!3;.1 .2 .3;`a`b`c)", `("I,F,S";"0,0.1,a";"1,0.2,b";"2,0.3,c")`},
		//{"`csv@(!3;.1a30 .2a40 .3a50;`a`b`c)", `("0,0.1,30,a";"1,0.2,40,b";"2,0.3,50,c")`},
		//{"`csv@+`I`Z`S!(!3;.1a30 .2a40 .3a50;`a`b`c)", `("I,Z,Z,S";"0,0.1,30,a";"1,0.2,40,b";"2,0.3,50,c")`},
		{"(\"ii\";\",\")0:(\"1,2,3\";\"4,5,6\")", "(1 4;2 5)"},
		{"(\"if\";\"|\")0:(\"1|2|3\";\"||6\")", "(1 0N;2 0n)"},
		{"1 rot`a`b`c", "`b`c`a"},
		{"-1 rot`a`b`c", "`c`a`b"},
		{"4 rot`a`b`c", "`b`c`a"},
		{"-4 rot`a`b`c", "`c`a`b"},
		{"prm 3", "(0 1 2;1 0 2;2 0 1;0 2 1;1 2 0;2 1 0)"},
		{"prm 1+!3f", "(1 2 3f;2 1 3f;3 1 2f;1 3 2f;2 3 1f;3 2 1f)"},
		{`prm "abc"`, `("abc";"bac";"cab";"acb";"bca";"cba")`},
		{"2 5 7 bin 3", "0"},
		{"2 5 7 bin !9", "-1 -1 0 0 0 1 1 2 2"},
	}
	for _, occ := range []bool{true, false} {
		for _, tc := range testCases {
			ini()
			fmt.Printf("%s → %s\n", tc.x, tc.r)
			x := prs(K([]byte(tc.x)))
			if occ {
				inc(x)
			}
			y := kst(evl(x))
			r := string(G(y).([]c))
			if r != tc.r {
				t.Fatalf("expected %s got %s\n", tc.r, r)
			}
			if occ {
				dec(x)
			}
			dec(y)
			clear()
			check(t)
		}
	}
}

func TestParse(t *testing.T) {
	// t.Skip()
	testCases := []struct {
		x, r, t s
	}{
		{"", "", "`"},
		{"`a", ",`a", "`N"},
		{"`a`b", ",`a`b", "`."},
		{"`a / b", ",`a", "`N"},
		{".x", "`.x", "`n"},
		{"{1+x}", "{1+x}", "`1"},
		{"{x+y}", "{x+y}", "`2"},
		{"{1+(2;`a;y)}", "{1+(2;`a;y)}", "`2"},
		{"/alpha\n`a", "(`;;,`a)", "`."},
		{"/alpha\n`a /beta\n/gamma", "(`;;,`a;)", "`."},
		{"0x01", "0x01", "`c"},
		{"0xF", "0x0f", "`c"},
		{"0x", `""`, "`C"},
		{"0x1234", "0x1234", "`C"},
		{`"a"`, `"a"`, "`c"},
		{`"a\t\n\r\"xyz"`, `"a\t\n\r\"xyz"`, "`C"},
		{"10", "10", "`i"},
		{"10 20", "10 20", "`I"},
		{"-3", "-3", "`i"},
		{".1", "0.1", "`f"},
		{"2.", "2f", "`f"},
		{"-1.23e-005", "-1.23e-5", "`f"},
		{"1 2f", "1 2f", "`F"},
		{"1.23 3", "1.23 3", "`F"},
		{"1 2 3. ", "1 2 3f", "`F"},
		{"2i-3", "3.605551a303.6901", "`z"},
		{"-2.0e+012i-3.6", "2e12a180", "`z"},
		{"1 2 3. 2i-3", "1 2 3 3.605551a303.6901", "`Z"},
		{"+", "+", "`2"},
		{"+:", "+:", "`1"},
		{"0:", "0:", "`2"},
		{"0:1", "(0::;1)", "`."},
		{"x-:", "(-:;`x)", "`."},
		{"'", "'", "`2"},
		{"1+2;", "(`;(+;1;2);)", "`."},
		{"(/)", ",/", "`."},
		{"2/3", "((/;2);3)", "`."},
		{"(/:)", ",/:", "`."},
		{"in", "in", "`2"},
		{"within", "within", "`2"},
		{"bin", "bin", "`2"},
		{"like", "like", "`2"},
		{"1 2@0", "(1 2;0)", "`."},
		{"+/", "(/;+)", "`."},
		{"+/1 2 3", "((/;+);1 2 3)", "`."},
		{"1+/3 4 5", "((/;+);1;3 4 5)", "`."},
		{"*1 2 3", "(*:;1 2 3)", "`."},
		{"1+(2;`a;3.5)", "(+;1;(;2;,`a;3.5))", "`."},
		{"1;2\n3", "(`;1;2;3)", "`."},
		{"1 2 3[0 2]", "(1 2 3;0 2)", "`."},
		{"`a`b`c[2]", "(,`a`b`c;2)", "`."},
		{"x[0][1]", "((`x;0);1)", "`."},
		{"(1;(2;3);4)[1;1]", "((;1;(;2;3);4);1;1)", "`."},
		{"(1;(`a;3);4)[1;0]", "((;1;(;,`a;3);4);1;0)", "`."},
		{"x[1;2;3]", "(`x;1;2;3)", "`."},
		{"(`a`b!1 2)", "(!;,`a`b;1 2)", "`."},
		{"(`a`b!1 2)[`b]", "((!;,`a`b;1 2);,`b)", "`."},
		{"2+", "(+;2)", "`."},
		{"{x+y}[2;3]", "({x+y};2;3)", "`."},
		{"g:{x+y};g[3;4]", "(`;(:;`g;{x+y});(`g;3;4))", "`."},
		{"-':8 2 5", "((':;-);8 2 5)", "`."},
		{`|+\!3`, `(|:;((\;+);(!:;3)))`, "`."},
		{"{-x}+3", `(+;{-x};3)`, "`."},
		{"$[0;1;0;2;3]", "($;0;1;0;2;3)", "`."},
		{"(-).(1 2)", "(.;,-;1 2)", "`."},
		{"x+1", "(+;`x;1)", "`."},
		{"x-1", "(-;`x;1)", "`."},
		{"x -1", "(`x;-1)", "`."},
		{"3.5-1", "(-;3.5;1)", "`."},
		{"3.5 -1", "3.5 -1", "`F"},
		{"x[1]-2", "(-;(`x;1);2)", "`."},
		{"-*3", "(-:;(*:;3))", "`."},
		{"+-", "(.;+:;-)", "`."},
		{"+-:", "(.;+:;-:)", "`."},
		{"+-*", "(.;+:;(.;-:;*))", "`."},
		{"3*+", "(.;(*;3);+)", "`."},
		{"3+2+", "(.;(+;3);(+;2))", "`."},
		{". {$[x>3;2;3]}5", "(.:;({$[x>3;2;3]};5))", "`."},
	}
	for i, occ := range []bool{true, false} {
		for j, tc := range testCases {
			for _, at := range []bool{false, true} {
				ini()
				x := mkb([]c(tc.x))
				if occ {
					inc(x)
				}
				y := prs(x)
				exp := tc.r
				if at {
					fmt.Printf("@`p(%q) ~ %v\n", tc.x, tc.t)
					y = kst(tip(y))
					exp = tc.t
				} else {
					fmt.Printf("`p(%q) ~ %v\n", tc.x, tc.r)
					y = kst(y)
				}
				r := string(G(y).([]c))
				if !reflect.DeepEqual(r, exp) {
					t.Fatalf("[%d/%d]: expected: %v got %v\n", j, i, exp, r)
				}
				dec(y)
				if occ {
					dec(x)
				}
				clear()
				check(t)
			}
		}
	}
}
func TestNumMonad(t *testing.T) {
	// t.Skip()
	ini()
	xv := []interface{}{c(3), []c{3, 5}, -5, iv{3, -9}, 3.2, []f{-3.5, 2.9, 0}, 2 - 4i, []z{4 - 2i, 3 + 4i}}
	testCases := []struct {
		f func(k) k
		s s
		r []interface{}
	}{
		{neg, "-", l{c(253), []c{253, 251}, 5, iv{-3, 9}, -3.2, []f{3.5, -2.9, -0}, -2 + 4i, []z{-4 + 2i, -3 - 4i}}},
		{fst, "*", l{c(3), c(3), -5, 3, 3.2, -3.5, 2 - 4i, 4 - 2i}},
		{rev, "|", l{c(3), []c{5, 3}, -5, iv{-9, 3}, 3.2, []f{0, 2.9, -3.5}, 2 - 4i, []z{3 + 4i, 4 - 2i}}},
		{not, "~", l{0, iv{0, 0}, 0, iv{0, 0}, 0, iv{0, 0, 1}, 0, iv{0, 0}}},
		{enl, ",", l{[]c{3}, l{[]c{3, 5}}, iv{-5}, l{iv{3, -9}}, []f{3.2}, l{[]f{-3.5, 2.9, 0}}, []z{2 - 4i}, l{[]z{4 - 2i, 3 + 4i}}}},
		{cnt, "#", l{1, 2, 1, 2, 1, 3, 1, 2}},
		{str, "$", l{c(3), []c{3, 5}, []c("-5"), l{[]c("3"), []c("-9")}, []c("3.2"), l{[]c("-3.5"), []c("2.9"), []c("0")}, []c("4.472136a296.5651"), l{[]c("4.472136a333.4349"), []c("5a53.1301")}}},
		{kst, "`k@", l{[]c("0x03"), []c("0x0305"), []c("-5"), []c("3 -9"), []c("3.2"), []c("-3.5 2.9 0"), []c("4.472136a296.5651"), []c("4.472136a333.4349 5a53.1301")}},
		{tip, "@", l{"c", "C", "i", "I", "f", "F", "z", "Z"}},
		{evl, ".", xv},
	}
	for _, occ := range []bool{true, false} {
		for j, tc := range testCases {
			for i := range xv {
				// fmt.Println("TC", xv[i])
				x := K(xv[i])
				if x == 0 {
					t.Fatalf("cannot import go type %T", xv[i])
				}
				if occ {
					inc(x)
				}
				y := tc.f(x)
				if occ {
					dec(x)
				}
				r := G(y)
				fmt.Printf("%s(%v) = %v\n", tc.s, xv[i], r)
				if !reflect.DeepEqual(r, tc.r[i]) {
					t.Fatalf("[%d/%d]: expected: %v got %v (@%d)\n", j, i, tc.r[i], r, y)
				}
				dec(y)
				clear()
				check(t)
			}
		}
	}
}
func TestMonad(t *testing.T) {
	// t.Skip()
	ini()
	testCases := []struct {
		f    func(k) k
		s    s
		x, r interface{}
	}{
		// flp needs ini
		//{flp, "+", l{iv{1, 2}, iv{3, 4}}, l{iv{1, 3}, iv{2, 4}}},
		//{flp, "+", l{iv{1, 2, 3}, []f{3, 4, 5}}, l{l{1, 3.0}, l{2, 4.0}, l{3, 5.0}}},
		//{flp, "+", l{l{1, 2.2, l{3, 4}}, l{"a", []c{'x'}, 5 + 2i}}, l{l{1, "a"}, l{2.2, []c{'x'}}, l{l{3, 4}, 5 + 2i}}},
		//{flp, "+", l{l{1, 2.2}, l{3, 4.4}}, l{iv{1, 3}, []f{2.2, 4.4}}},
		{til, "!", 3, iv{0, 1, 2}},
		{til, "!", -1, l{iv{1}}},
		{til, "!", -3, l{iv{1, 0, 0}, iv{0, 1, 0}, iv{0, 0, 1}}},
		{til, "!", d{sv{"a", "b"}, iv{1, 2}}, sv{"a", "b"}},
		// TODO !a(odometer)
		{fst, "*", l{3, 4, 5}, 3},
		{fst, "*", "alpha", "alpha"},
		{fst, "*", l{"alpha"}, "alpha"},
		{fst, "*", d{l{"x", "y"}, l{iv{5, 3}, 4}}, iv{5, 3}},
		{fst, "*", d{sv{"x", "y"}, iv{7, 2}}, 7},
		{inv, "%", 4, 0.25},
		{inv, "%", []f{4.0, 8.0}, []f{0.25, 0.125}},
		{str, "$", "", []c("")},
		{str, "$", "a", []c("a")},
		{str, "$", sv{"", "a", "bb", "a\t\r\nb"}, l{[]c(""), []c("a"), []c("bb"), []c("a\t\r\nb")}},
		{str, "$", l{1, c(3), l{4, 5.0}}, l{[]c("1"), c(3), l{[]c("4"), []c("5")}}},
		{str, "$", d{sv{"x", "y"}, iv{1, 2}}, d{sv{"x", "y"}, l{[]c("1"), []c("2")}}},
		{kst, "`k", iv{1, 2, 3}, []c("1 2 3")},
		{kst, "`k", l{1, 2, l{4, 5}}, []c("(1;2;(4;5))")},
		{kst, "`k", d{l{5, 5.5}, iv{1, 2}}, []c("(5;5.5)!1 2")},
		{rev, "|", l{}, l{}},
		{rev, "|", l{iv{3}}, l{iv{3}}},
		{rev, "|", l{1, 2}, l{2, 1}},
		{rev, "|", l{1, l{3, 4}}, l{l{3, 4}, 1}},
		{rev, "|", d{iv{1, 2}, iv{3, 4}}, d{iv{2, 1}, iv{4, 3}}},
		{rev, "|", d{sv{"alpha", "beta"}, l{3, iv{3, 5}}}, d{sv{"beta", "alpha"}, l{iv{3, 5}, 3}}},
		{wer, "&", iv{0, 0, 1, 1, 0, 1}, iv{2, 3, 5}},
		{wer, "&", iv{}, iv{}},
		{wer, "&", iv{2}, iv{0, 0}},
		{wer, "&", iv{1, 2, 3}, iv{0, 1, 1, 2, 2, 2}},
		{asc, "<", iv{1, 2, 3, 4}, iv{0, 1, 2, 3}},
		{asc, "<", iv{1, 4, 3, 2}, iv{0, 3, 2, 1}},
		{asc, "<", iv{4, 2, 3, 4}, iv{1, 2, 0, 3}},
		{asc, "<", []f{4, 1, 2}, iv{1, 2, 0}},
		{asc, "<", []c{6, 4, 2, 1}, iv{3, 2, 1, 0}},
		{asc, "<", []z{4, 1, 2}, iv{1, 2, 0}},
		{asc, "<", []z{0, 1 + 1i, 1, 2}, iv{0, 2, 1, 3}},
		{asc, "<", sv{"b", "ab", "a", "aa"}, iv{2, 3, 1, 0}},
		{dsc, ">", iv{1, 4, 3, 2}, iv{1, 2, 3, 0}},
		//{grp, "=", []c{'c', 'b', 'a', 'c', 'a', 'b', 'c'}, d{[]c{'c', 'b', 'a'}, l{iv{0, 3, 6}, iv{1, 5}, iv{2, 4}}}},
		{grp, "=", iv{1, 2, 3}, d{iv{1, 2, 3}, l{iv{0}, iv{1}, iv{2}}}},
		{grp, "=", iv{3, 3, 1, 3, 2, 1}, d{iv{1, 2, 3}, l{iv{2, 5}, iv{4}, iv{0, 1, 3}}}},
		{grp, "=", []f{5.5, 1, 3, 3, 2}, d{[]f{1, 2, 3, 5.5}, l{iv{1}, iv{4}, iv{2, 3}, iv{0}}}},
		{grp, "=", []z{3, 3 + 1i, 3, 3 + 1i, 3 + 1i}, d{[]z{3, 3 + 1i}, l{iv{0, 2}, iv{1, 3, 4}}}},
		{grp, "=", sv{"alpha", "beta", "alpha", "gamma", "alpha", "beta"}, d{sv{"alpha", "beta", "gamma"}, l{iv{0, 2, 4}, iv{1, 5}, iv{3}}}},
		{enl, ",", "alpha", sv{"alpha"}},
		{enl, ",", l{1, 2, l{3, 4.5}}, l{l{1, 2, l{3, 4.5}}}},
		{enl, ",", d{iv{3, 4}, sv{"x", "y"}}, l{d{iv{3, 4}, sv{"x", "y"}}}},
		{srt, "^", iv{3, 1, 2, 3, 5}, iv{1, 2, 3, 3, 5}},
		{cnt, "#", "alpha", 1},
		{cnt, "#", l{}, 0},
		{cnt, "#", l{1, 2, l{3, 4}}, 3},
		{cnt, "#", d{iv{3, 4}, sv{"x", "y"}}, 2},
		{tip, "@", l{}, "."},
		{tip, "@", d{iv{1, 2}, iv{3, 4}}, "a"},
		{evl, ".", l{uint16(2), l{uint16(2), 3}}, 3},
		{evl, ".", l{uint16(2), l{uint16(6), iv{3, 4}}}, iv{-4, -3}},
		{evl, ".", l{uint16(2), iv{3, 4}}, iv{-3, -4}},
		{unq, "?", []c{1, 2, 43, 2}, []c{1, 2, 43}},
		{unq, "?", iv{1, 2, 3, 2}, iv{1, 2, 3}},
		{unq, "?", []f{5, 0, 0, 0, 8, 0, 0, 0, 5, 0, 0, 5}, []f{5, 0, 8}},
		{unq, "?", []z{0, 4i, 5i, 4i, 0, 3}, []z{0, 4i, 5i, 3}},
		{unq, "?", l{1, 2, 3, 1}, l{1, 2, 3}},
		{unq, "?", l{1i, l{2, sv{"a"}}, l{3, "b"}, l{2, sv{"a"}}, 1i}, l{1i, l{2, sv{"a"}}, l{3, "b"}}},
	}
	for _, occ := range []bool{true, false} {
		for j, tc := range testCases {
			//fmt.Println("TC", i, j, tc.s, tc.x, "occ", occ)
			x := K(tc.x)
			_ = Stats().UsedBlocks()
			if x == 0 {
				t.Fatalf("cannot import go type %T", tc.x)
			}
			if occ {
				inc(x)
			}
			y := tc.f(x)
			fpck("1")
			if occ {
				dec(x)
			}
			r := G(y)
			fmt.Printf("%s[%v] = %v\n", tc.s, tc.x, r)
			if !reflect.DeepEqual(r, tc.r) {
				t.Fatalf("monad[%d]: expected: %v got %v (@%d)\n", j, tc.r, r, y)
			}
			dec(y)
			clear()
			check(t)
		}
	}
}
func TestDyad(t *testing.T) {
	// t.Skip()
	testCases := []struct {
		f       func(k, k) k
		s       s
		x, y, r interface{}
	}{
		{atx, "@", iv{2, 1, 3, 5}, iv{2, 0, 1}, iv{3, 2, 1}},
		{ept, "^", iv{1, 5, 3, 3, 2}, iv{6, 3}, iv{1, 5, 2}},
		{drp, "_", 2, iv{1, 2, 3, 4}, iv{3, 4}},
		{drp, "_", -2, iv{1, 2, 3, 4}, iv{1, 2}},
		{drp, "_", 2, l{1, 2.0, "a"}, sv{"a"}},
		{drp, "_", -2, l{1, 2.0, "a"}, iv{1}},
		{fnd, "?", iv{1, 2}, iv{2}, iv{1}},
		{fnd, "?", []c{5, 4, 3, 2}, c(2), 3},
	}
	for _, occ := range []bool{true, false} {
		for j, tc := range testCases {
			// fmt.Println("TC", i, j, tc.s, tc.x, "occ", occ)
			x := K(tc.x)
			y := K(tc.y)
			_ = Stats().UsedBlocks()
			if x == 0 || y == 0 {
				t.Fatalf("cannot import go type %T", tc.x)
			}
			if occ {
				inc(x)
				inc(y)
			}
			z := tc.f(x, y)
			fpck("1")
			if occ {
				dec(x)
				dec(y)
			}
			r := G(z)
			fmt.Printf("%s[%v, %v] = %v\n", tc.s, tc.x, tc.y, r)
			if !reflect.DeepEqual(r, tc.r) {
				t.Fatalf("dyad[%d]: expected: %v got %v\n", j, tc.r, r)
			}
			dec(z)
			clear()
			check(t)
		}
	}
}
func TestTo(t *testing.T) {
	// t.Skip()
	ini()
	testCases := []struct {
		x, r interface{}
		t    k
	}{
		{c(1), 1, I},
		{c(1), 1.0, F},
		{[]c{1, 2, 3, 4}, iv{1, 2, 3, 4}, I},
		{1, c(1), C},
		{5, 5.0, F},
		{1, 1 + 0i, Z},
		{iv{1, 2, 3, 4}, []z{1, 2, 3, 4}, Z},
		{2.3, c(2), C},
		{2.3, 2, I},
		{2.3, 2.3 + 0i, Z},
		{-2.3 + 4.5i, -2.3, F},
		{2.3 - 4.5i, c(2), C},
		{complex(math.NaN(), 0), -2147483648, I},
		{complex(0, math.NaN()), -2147483648, I},
	}
	for _, occ := range []bool{true, false} {
		for _, tc := range testCases {
			x := K(tc.x)
			if occ {
				inc(x)
			}
			y := to(x, tc.t)
			if occ {
				dec(x)
			}
			r := G(y)
			fmt.Printf("to(%v,%d) = %v\n", tc.x, tc.t, tc.r)
			if !reflect.DeepEqual(r, tc.r) {
				t.Fatalf("expected: %v got %v %T %T\n", tc.r, r, tc.r, r)
			}
			dec(y)
			clear()
			check(t)
		}
	}
}
func TestStr(t *testing.T) {
	// t.Skip()
	ini()
	for _, x := range []s{"a", "b", "aa", "bb", "alpha", "betagamm"} {
		n := len(x)
		if n > 8 {
			n = 8
		}
		if r := G(K(x)); r != x[:n] {
			t.Fatalf("expected %s got %s\n", x, r)
		}
	}
	if u := sym(8 + K("abcdefgh")<<2); u != 0x6162636465666768 {
		t.Fatalf("%x\n", u)
	}
}
func TestRef(t *testing.T) { // generate readme.md
	k, e := ioutil.ReadFile("k.go")
	if e != nil {
		t.Fatal(e)
	}
	fl := make(map[s]i)
	for i, b := range bytes.Split(k, []c("\n")) {
		if bytes.HasPrefix(b, []c("func ")) && b[8] == '(' {
			fl[s(b[5:8])] = int32(1 + i)
		}
	}
	w, e := os.Create("readme.md")
	if e != nil {
		t.Fatal(e)
	}
	defer w.Close()
	r := []c(ref)
	w.Write([]c("<pre>\n"))
	for i := 1; i < len(r)-2; i++ {
		if !craZ(r[i-1]) && craZ(r[i]) && craZ(r[i+1]) && craZ(r[i+2]) {
			key := r[i : i+3]
			if l, o := fl[s(key)]; o {
				fmt.Fprintf(w, `<a href="../../blob/master/k.go#L%d">%s</a>`, l, s(key))
			} else {
				w.Write(r[i : i+3])
			}
			i += 2
		} else {
			w.Write(r[i : i+1])
		}
	}
	w.Write([]c(`

k7+
 complex(type z): 1i2, 1a30, cmplx, expi, real, imag, phase, conj, rand 3i(binormal)
 matrix: x/y(mul), A\B(solve), A\0(qr), A\1(inv), diag A, diag v, norm, cond
 stat: x med (pct/erf/cum), dev z (principal axis), x var, var z (cov), x avg (cum/win/exp)
k7-
 32bit, time/duration, :expr, ksql, crypto, network, multithread
 
Type/memory system
32-bit system, buddy allocater with 8 byte headers
types (cifzsla01234) byte8, int32, float64, complex128, symbol64, list32, dict64, funcs
space for 15 types
8 byte header:
  4 bits type p>>28
 28 bits vector size, atom: p&0x0fffffff == 0x0fffffff
 free block: type is 0:     p&0xf0000000 == 0
 bucket type is stored only in free blocks at p (uint32 value)
 32 bits (p+1) are refcount for used blocks or pointer to next free

Initial memory (64kB)
 p[0]        block header
 p[1]        rng state
 p[2]        total allocated memory log2 (initial 64k, max 4G) uint32
 p[3]        points to a dict of built-ins S(name)!L(fcodes)
 followed by free list:
 p[4..31]    points to free block of bucket size n = 4..31
 byte[136…168] symbols :+-*%&|<>=!~,^#_$?@.0123456789'/\
 byte[169…181] type names _cifzn.a_1234
 p[47]       0x2f src pointer
 p[48]       0x30 points to k tree keys (^S)
 p[49]       0x31 points to k tree values (L)
 byte[?]     type size vector: 0,1,4,8,16,8,4,0,0,0,0,0,0
             A01234 need only a single block but may have length>0

Function codes
  0-19 monadic primitives :+-*%&|<>=!~,^#_$?@.
 20-29 monadic ioverbs 0123456789
 30-32 monadic operator functions '/\
 33-38 monadic derived functions f' f/ f\ f': f/: f\:
 39-79 monadic builtins
 80-159 dyadic versions

Functions have type N+1…N+4 (valence)
 basic functions and builtins: atoms
  x+2 is the function code
 lambda functions: marked with length 0
  x+2 string form C
  x+3 (arg list;parse tree), variables are always x,y,z
 projection: length 1(over lambda) 2(over basic/builtins)
  x+2 function code or pointer to lambda function
  x+3 full argument list with holes (N)
 composition: length 3, type N+1 or N+2
  x+2, x+3: point to verbs
 derived verbs, e.g. evaluating (/;+) have type N+1 with code > 256
  x+2 derived function code<<8
  x+3 points to the function operand
 call will adjust the valence if a derived function has two arguments

</pre>`))
}
func TestRefcard(t *testing.T) { // go test -short (it's the long test, but short=false is the default)
	if !testing.Short() {
		return
	}
	var b []byte
	if req, e := http.Get(`https://raw.githubusercontent.com/kparc/ref/master/src/md/index.md`); e != nil {
		t.Fatal(e)
	} else {
		defer req.Body.Close()
		if b, e = ioutil.ReadAll(req.Body); e != nil {
			t.Fatal(e)
		}
	}
	ini()
	x := spl(mkb(b), mkc('\n'))
	kx(mks(".tref"), x) // TODO: implementation that works for k7 and i
}
func check(t *testing.T) {
	// Number of used blocks after an expression should be:
	// 1(block 0) + 3(built-in dict,k,v) + 2(k-tree k,v)
	// vars := m.k[m.k[kkey]] & atom
	if u := Stats().UsedBlocks(); u != 6 {
		xxd()
		t.Fatalf("leak: %d", u)
	}
	fpck("")
}
func pfl() {
	for i := 4; i < 32; i++ {
		println(i, strconv.FormatUint(uint64(m.k[i]), 16), strconv.FormatUint(uint64(m.k[i]<<2), 16))
	}
}
func xdec(x k) {
	if m.k[x]>>28 != 0 {
		dec(x)
	}
}
func xxd() { // memory dump
	h := k(0)
	for i := k(0); i < k(len(m.k)); i += 4 {
		a, b, c, d := m.k[i+0], m.k[i+1], m.k[i+2], m.k[i+3]
		if a == 0 && b == 0 && c == 0 && d == 0 {
			continue
		}
		fmt.Printf("0x%04x %08x %08x %08x %08x", i, a, b, c, d)
		if i == h {
			tp := m.k[i] >> 28
			if tp == 0 {
				fmt.Printf("  %d", m.k[i])
				h += 1 << (m.k[i] - 2)
				nf := m.k[i+1]
				if nf > 0 && nf < 64 {
					fmt.Printf(" illegal fp")
				} else if nf > 0 && m.k[nf]>>28 != 0 {
					fmt.Printf(" fp is not free")
				}
			} else {
				atoms := "?cifzsla01234"
				vects := "?CIFZSLA01234"
				tp, n := typ(i)
				bt := bk(tp, n)
				if n == atom {
					fmt.Printf(" %c%d +%d", atoms[tp], bt, b)
				} else {
					fmt.Printf(" %c%d #%d +%d", vects[tp], bt, n, b)
				}
				h += 1 << (bt - 2)
			}
		}
		fmt.Println()
	}
}
func fpck(s s) { // check free pointers
	for i := 4; i < 32; i++ {
		nf := m.k[i]
		if nf > 0 && (nf < 64 || m.k[nf]>>28 != 0) {
			xxd()
			panic("fpck " + s + " bad pointer in free-list: @" + strconv.Itoa(int(i)))
		}
	}
	h := k(0)
	for i := k(0); i < k(len(m.k)); i += 4 {
		if i == h {
			tp := m.k[i] >> 28
			if tp == 0 {
				h += 1 << (m.k[i] - 2)
				nf := m.k[i+1]
				if nf > 0 && (nf < 64 || m.k[nf]>>28 != 0) {
					xxd()
					panic("fpck " + s + " illegal free-pointer")
				}
			} else {
				tp, n := typ(i)
				bt := bk(tp, n)
				h += 1 << (bt - 2)
			}
		}
	}
}

type Bucket struct {
	Type       uint32
	Used, Free uint32 // num blocks
	Net        uint32
}
type MemStats map[uint32]Bucket

func (b Bucket) Overhead() uint32 {
	return b.Used*uint32(1<<b.Type) - b.Net
}
func (s MemStats) UsedBlocks() (t uint32) {
	for _, b := range s {
		t += b.Used
	}
	return t
}
func Stats() MemStats {
	st := make(MemStats)
	a := uint32(0)
	o := uint32(0)
	for a < 1<<(m.k[2]-2) {
		tp := m.k[a] >> 28
		if tp == 0 {
			t := m.k[a]
			if t < 4 || t > 31 {
				xxd()
				fmt.Printf("free block at %x with bt %d\n", a, t)
				panic("size")
			}
			b := st[t]
			b.Type = t
			b.Free++
			st[t] = b
			o = 1 << (t - 2)
		} else {
			tt, n := typ(a)
			t := bk(tt, n)
			if t < 4 || t > 31 {
				println(a, t)
				panic("size")
			}
			b := st[t]
			b.Type = t
			b.Used++
			if n == atom {
				n = 1
			}
			b.Net += n * lns[tp]
			st[t] = b
			o = 1 << (t - 2)
		}
		a += o
	}
	return st
}

// type conversions between go and k (used by k_test.go)

func K(x interface{}) k { // convert go value to k type, returns 0 on error
	if x == nil {
		return mk(N, atom)
	}
	kstr := func(dst k, s string) {
		mys(dst, btou([]c(s)))
	}
	var r k
	switch a := x.(type) {
	case bool:
		r = mkc(0)
		if a {
			m.c[8+r<<2] = 1
		}
	case byte:
		r = mkc(a)
	case int:
		r = mki(k(a))
	case uint16: // function index
		if a < 20 {
			r = mk(N+1, atom)
		} else {
			r = mk(N+2, atom)
		}
		m.k[2+r] = k(a)
	case float64:
		r = mk(F, atom)
		m.f[1+r>>1] = a
	case complex128:
		r = mk(Z, atom)
		m.z[1+r>>2] = a
	case string:
		r = mks(a)
	case []bool:
		buf := make([]byte, len(a))
		for i, v := range a {
			if v {
				buf[i] = 1
			}
		}
		return K(buf)
	case []byte:
		r = mk(C, k(len(a)))
		for i, v := range a {
			m.c[8+i+int(r<<2)] = v
		}
	case []int:
		r = mk(I, k(len(a)))
		for i, v := range a {
			m.k[2+i+int(r)] = k(v)
		}
	case []float64:
		r = mk(F, k(len(a)))
		for i, v := range a {
			m.f[1+i+int(r>>1)] = v
		}
	case []complex128:
		r = mk(Z, k(len(a)))
		for i, v := range a {
			m.z[1+i+int(r>>2)] = v
		}
	case []string:
		r = mk(S, k(len(a)))
		for i := range a {
			kstr(8+8*k(i)+r<<2, a[i])
		}
	case []interface{}:
		if len(a) == 1 { // collapse list of atom to single element vector
			rr := K(a[0])
			t, n := typ(rr)
			if n == atom { // TODO: allow ,d?
				r = rr
				m.k[r] = t<<28 | 1
				return r
			} else {
				dec(rr)
			}
		}
		r = mk(L, k(len(a)))
		for i, v := range a {
			u := K(v)
			m.k[2+i+int(r)] = u
		}
	case [2]interface{}:
		key := K(a[0])
		val := K(a[1])
		_, nk := typ(key)
		_, nv := typ(val)
		if nk != nv {
			return 0
		}
		r = mk(A, atom)
		m.k[2+r] = key
		m.k[3+r] = val
	}
	return r
}
func G(x k) interface{} { // convert k value to go type (returns nil on error)
	t, n := typ(x)
	str := func(xp k) s {
		r := mk(C, 0)
		rc := 8 + r<<2
		n := stS(rc, xp)
		dec(r)
		return string(m.c[rc : rc+n])
	}
	if n == atom {
		switch t {
		case C:
			return c(m.c[8+x<<2])
		case I:
			return int(i(m.k[2+x]))
		case F:
			return m.f[1+x>>1]
		case Z:
			return m.z[1+x>>2]
		case S:
			return str(1 + x>>1)
		case A:
			return [2]interface{}{G(m.k[2+x]), G(m.k[3+x])}
		case N:
			return nil
		case N + 1, N + 2:
			return uint16(m.k[2+x])
		}
	} else {
		switch t {
		case C:
			r := make([]byte, n)
			for i := range r {
				r[i] = c(m.c[8+i+int(x<<2)])
			}
			return r
		case I:
			r := make([]int, n)
			for i := range r {
				r[i] = int(int32(m.k[2+i+int(x)]))
			}
			return r
		case F:
			r := make([]f, n)
			for i := range r {
				r[i] = m.f[1+i+int(x>>1)]
			}
			return r
		case Z:
			r := make([]complex128, n)
			for i := range r {
				r[i] = m.z[1+i+int(x>>2)]
			}
			return r
		case S:
			r := make([]string, n)
			for i := range r {
				r[i] = str(1 + k(i) + x>>1)
			}
			return r
		case L:
			r := make([]interface{}, n)
			for i := range r {
				r[i] = G(m.k[2+i+int(x)])
			}
			return r
		}
	}
	return nil
}
