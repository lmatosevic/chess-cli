package utils

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
	"syscall"
)

func ReadStringFromStdin(prompt string) (string, error) {
	if prompt != "" {
		fmt.Printf(prompt)
	}

	reader := bufio.NewReader(os.Stdin)
	text, e := reader.ReadString('\n')
	if e != nil {
		return "", e
	}

	text = strings.Replace(text, "\n", "", -1)
	text = strings.Replace(text, "\r", "", -1)

	return strings.TrimSpace(text), nil
}

func ReadPasswordFromStdin(prompt string) (string, error) {
	if prompt != "" {
		fmt.Printf(prompt)
	}

	text, e := terminal.ReadPassword(int(syscall.Stdin))
	if e != nil {
		return "", e
	}

	return string(text), nil
}

func WaitForAnyKey() error {
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadByte()
	return err
}

func PrintStruct(obj any) {
	bytes, _ := json.MarshalIndent(obj, "", "    ")
	fmt.Printf("%s\n", string(bytes))
}

func PrintTable(title string, headers table.Row, rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle(title)
	t.SetIndexColumn(0)
	t.AppendHeader(headers)
	t.AppendRows(rows)
	t.Render()
}

func PrintChessBoard(tiles string) {
	fmt.Print("     a    b    c    d    e    f    g    h   \n")
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if j == 0 {
				fmt.Printf("%d  ", 8-i)
			}

			figure := string(tiles[j+(i*8)])

			bracketOpen := "["
			bracketClose := "]"
			if (i+j)%2 == 0 {
				// white tile
				bracketOpen = "("
				bracketClose = ")"
			}

			if figure == "0" {
				fmt.Printf("%s   %s", bracketOpen, bracketClose)
			} else {
				fmt.Printf("%s %s %s", bracketOpen, figure, bracketClose)
			}
		}
		fmt.Println()
	}
}
