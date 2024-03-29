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

func (ts *Token_stream) putback(t Token) {
	if ts.full {
		panic("putback() into a full buffer")
	}
	ts.buffer = t
	ts.full = true
}

var reader = bufio.NewReader(os.Stdin)

func (ts *Token_stream) get() Token {
	if ts.full {
		ts.full = false
		return ts.buffer
	}

	fmt.Println("TAKING INPUT ")

	var ch rune
	//reader := bufio.NewReader(os.Stdin)
	ch, _, _ = reader.ReadRune()

	fmt.Printf("%q\n", ch)

	if ch == ' ' {
		ch, _, _ = reader.ReadRune()
	}

	if ch == '\n' {
		fmt.Println("END OF LINE")
	}

	switch ch {
	case quit, print, '(', ')', '+', '-', '*', '/', '=':
		{
			fmt.Println("CASE 1 ")
			return Token{kind: ch}
		}
	case '.', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		{
			fmt.Println("CASE 2 ")
			if reader.UnreadRune() != nil {
				panic("unreadRune() failed")
			}
			var input float64
			fmt.Println("NUMBER", input)
			//input, err := reader.ReadString(' ') //вписываем вернувшийся на консоль символ в value
			_, err := fmt.Fscan(reader, &input)
			//input = strings.TrimSpace(input)
			if err != nil {
				panic(err)
			}
			fmt.Println("NUMBER", input)
			//value, _ := strconv.ParseFloat(input, 64)

			//fmt.Println("NUMBER", value)
			return Token{kind: number, value: input}

		}
	default:
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
	//reader := bufio.NewReader(os.Stdin)
	for {
		ch, _, _ = reader.ReadRune()
		if ch == c {
			return
		}
	}

}

var str Token_stream //глобальная переменная

func primary() float64 {

	t := str.get()
	fmt.Println("PRIMARY: ", t)

	switch t.kind {
	case '(':
		{

			fmt.Println("PRIMARY PARENTHESES ", t)

			f := expression()
			t = str.get()
			if t.kind != ')' {
				panic("')' expected")
			}
			return f
		}
	case number:

		fmt.Println("PRIMARY NUMBER ", t)

		return t.value
	case name:
		{
			next := str.get()
			if next.kind == '=' {
				d := expression()
				return d

			}
		}
	case '-':
		return -primary()
	case '+':
		return +primary()
	default:
		fmt.Println("FALLTHRU IN PRIMARY ", t)
		panic("primary expected")
	}
	panic("primary expected")
}

func term() float64 {

	left := primary()
	t := str.get()

	fmt.Println("TERM: ", t)

	for {
		switch t.kind {
		case '*':
			{
				fmt.Println("*: ", t)

				left *= primary()
				t = str.get()
			}
		case '/':
			{
				fmt.Println("///: ", t)

				f := primary()
				if f == 0.0 {
					panic("divide by zero")
				}
				left /= f
				t = str.get()
			}
		default:
			fmt.Println("FALLTHRU IN TERM ", t)
			str.putback(t)
			return left
		}
	}

}

func expression() float64 {
	left := term()
	t := str.get()

	fmt.Println("EXPRESSION: ", t)

	for {
		switch t.kind {
		case '+':

			fmt.Println("+++ ", t)

			left += term()
			fmt.Println("Value ", left)
			t = str.get()
		case '-':

			fmt.Println("---- ", t)

			left -= term()
			t = str.get()
		default:

			fmt.Println("FALLTHRU IN EXPRESSION ", t)

			str.putback(t)
			return left
		}
	}
}

func statement() float64 {

	t := str.get()
	switch t.kind {
	case let:
	case con:
	default:
		str.putback(t)
		return expression()
	}
	panic("statement expected")
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
