package unit_file

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ReadFirstLine() string {
	open, err := os.Open("log")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	defer open.Close()
	if err != nil {
		return ""
	}
	scanner := bufio.NewScanner(open)
	for scanner.Scan() {
		return scanner.Text()
	}
	return ""

}

func ProcessFirstLine() string {

	line := ReadFirstLine()
	destLine := strings.ReplaceAll(line, "11", "00")
	return destLine
}
