package main

import (
	"fmt"
	"os"
	"strconv"

	board "github.com/dannyh79/graphic-raid/quorum/internal"
)

var usageExample = "Usage: ./main <integer number of members>"

func main() {
	if len(os.Args) != 2 {
		fmt.Println(usageExample)
		return
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(usageExample)
		return
	}

	w := os.Stdout
	board.HoldQuorumElection(w, n)
}
