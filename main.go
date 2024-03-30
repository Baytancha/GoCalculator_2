package main

import (
	"bufio"
	"fmt"
	"os"

	//"strconv"
	//"strings"
	"unicode"
)

const number = '8'       // t.kind==number means that t is a number Token
const quit = 'q'         // t.kind==quit means that t is a quit Token
const print = ';'        // t.kind==print means that t is a print Token
const name = 'a'         // name token
const let = 'L'          // declaration token
const con = 'C'          // const declaration token
const declkey = "let"    // declaration keyword
const constkey = "const" // const keyword
const prompt = '>'
const result = '='

type Token struct {
	kind  rune
	value float64
	name  string
}

type Token_stream struct {
	buffer Token
	full   bool
}

type Variable struct {
	name  string
	value float64
	v     bool
}
type Symbol_table struct {
	var_table []Variable
}

func (s *Symbol_table) get(str string) float64 {
	for i := 0; i < len(s.var_table); i++ {
		if s.var_table[i].name == str {
			return s.var_table[i].value
		}
	}
	panic("undefined variable")
}

func (s *Symbol_table) set(str string, value float64) {
	for i := 0; i < len(s.var_table); i++ {
		if s.var_table[i].name == str {
			if !s.var_table[i].v {
				panic("it's a constant")
			}
			s.var_table[i].value = value
			return
		}
	}
	panic("undefined variable")
}

func (s *Symbol_table) is_declared(str string) bool {
	for i := 0; i < len(s.var_table); i++ {
		if s.var_table[i].name == str {
			return true
		}
	}
	return false
}

func (s *Symbol_table) define(name string, value float64, v bool) float64 {
	if s.is_declared(name) {
		panic("variable already declared")
	}
	s.var_table = append(s.var_table, Variable{name, value, v})
	return value
}

func (ts *Token_stream) putback(t Token) {
	if ts.full {
		panic("putback() into a full buffer")
	}
	ts.buffer = t
	ts.full = true
}

// глобальный буфер
var reader = bufio.NewReader(os.Stdin)

func (ts *Token_stream) get() Token {
	if ts.full {
		ts.full = false
		return ts.buffer
	}

	var ch rune
	//считывает в т.ч. пробелы
	ch, _, _ = reader.ReadRune()

	//проскочить все пробелы вводимые пользователем
	for ch == ' ' {
		ch, _, _ = reader.ReadRune()
	}

	switch ch {
	case quit, print, '(', ')', '+', '-', '*', '/', '=':
		{
			return Token{kind: ch}
		}
	case '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		{
			if reader.UnreadRune() != nil {
				panic("unreadRune() failed")
			}
			var input float64
			//с помощью Fscan считываем возвращенный символ и другие символы числа
			_, err := fmt.Fscan(reader, &input)
			if err != nil {
				panic(err)
			}
			return Token{kind: number, value: input}

		}
	default:
		//поскольку bufio.ReadRune() считывает пробелы, то их нужно обрабатывать
		if ch == '\r' || ch == '\n' {
			return Token{kind: print}
		} else if unicode.IsLetter(ch) {
			var s string
			s += string(ch)
			for { //обработка переменных
				ch, _, _ = reader.ReadRune()
				if !unicode.IsLetter(ch) {
					break
				}
				s += string(ch)
			}
			if reader.UnreadRune() != nil { //сохраняем следующий символ
				panic("unreadRune() failed")
			}
			if s == declkey {
				return Token{kind: let}
			}
			if s == constkey {
				return Token{kind: con}
			}

			return Token{kind: name, name: s}

		}
	}
	panic("bad token")
}

func (ts *Token_stream) ignore(c rune) {

	var ch rune
	for {
		ch, _, _ = reader.ReadRune()
		if ch == c {
			return
		}
	}

}

var str Token_stream //глобальная переменная
var tbl Symbol_table

func primary() float64 {
	t := str.get()
	switch t.kind {
	case '(':
		{
			f := expression()
			t = str.get()
			if t.kind != ')' {
				panic("')' expected")
			}
			return f
		}
	case number:
		return t.value
	case name:
		{
			next := str.get()
			if next.kind == '=' {
				d := expression()
				tbl.set(t.name, d)
				return d
			} else {
				str.putback(next)
				return tbl.get(t.name)
			}
		}
	case '-':
		return -primary()
	case '+':
		return +primary()
	default:
		panic("primary expected")
	}
}

func term() float64 {

	left := primary()
	t := str.get()

	for {
		switch t.kind {
		case '*':
			{
				left *= primary()
				t = str.get()
			}
		case '/':
			{
				f := primary()
				if f == 0.0 {
					panic("divide by zero")
				}
				left /= f
				t = str.get()
			}
		default:
			str.putback(t)
			return left
		}
	}

}

func expression() float64 {
	left := term()
	t := str.get()

	for {
		switch t.kind {
		case '+':
			left += term()
			t = str.get()
		case '-':
			left -= term()
			t = str.get()
		default:
			str.putback(t)
			return left
		}
	}
}

func declaration(t Token) float64 {
	t = str.get()

	if t.kind != name {
		panic("name expected")
	}
	var_name := t.name

	t2 := str.get()
	if t2.kind != '=' {
		panic("'=' expected")
	}
	//переменная может быть результатом мат выражения
	d := expression()
	tbl.define(var_name, d, t.kind == let)
	return d

}

func statement() float64 {

	t := str.get()
	switch t.kind {
	case let, con:
		return declaration(t)
	default:
		str.putback(t)
		return expression()
	}
}

func cleanup() {
	str.ignore(print)
}

func main() {

	// for {
	// 	var ch rune
	// 	reader := bufio.NewReader(os.Stdin)
	// 	ch, _, _ = reader.ReadRune()
	// 	fmt.Println(string(ch))

	// 	reader.UnreadRune()

	// 	var ch2 float64
	// 	fmt.Fscan(reader, &ch2)
	// 	ch2 = ch2 + 150
	// 	//ch2, _ = reader.ReadString(' ')
	// 	fmt.Println(ch2)
	// 	//ch2 = strings.TrimSpace(ch2)
	// 	//value, _ := strconv.ParseFloat(ch2, 64)

	// 	//fmt.Println(value)
	// }

	for {
		fmt.Printf("%c", prompt)

		t := str.get()

		for t.kind == print {
			t = str.get()
		}
		if t.kind == quit {
			break
		}
		str.putback(t)
		fmt.Println("RESULT: ", statement())

	}

}
