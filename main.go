package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var fp = flag.String("problems", "./problems.csv", "Path to CSV containing problems for the quiz")

func main() {
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	done := make(chan bool, 1)

	go func() {
		<-sigs
		done <- true
	}()

	go func() {
		quiz()
		done <- true
	}()

	<-done

	fmt.Println("\n\033[1;34mThanks for playing!\033[0m")
}

func quiz() {
	f, err := os.Open(*fp)
	if err != nil {
		fail("reading csv", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	qs, err := r.ReadAll()
	if err != nil {
		fail("parsing csv", err)
	}

	var errs [][]string

	fmt.Printf("\n\033[1;33mIt’s quiz time!\033[0m\n\n")

	for _, q := range qs {
		fmt.Printf("%v: ", q[0])

		var a string
		fmt.Scanf("%s", &a)

		if a == q[1] {
			fmt.Println("\033[1;32mCorrect!\033[0m")
		} else {
			errs = append(errs, append(q, a))
			fmt.Println("\033[1;31mIncorrect!\033[0m")
		}
	}

	fmt.Printf("\n\033[1mYou answered %d/%d questions correctly!\033[0m\n\n", len(qs)-len(errs), len(qs))

	if len(errs) == 0 {
		return
	}

	fmt.Println("Errors:")
	for _, e := range errs {
		fmt.Printf("%v: %v (you answered %v)\n", e[0], e[1], e[2])
	}
}

func fail(msg string, err error) {
	fmt.Fprintf(os.Stderr, "quizgame: %v: %v\n", msg, err)
	os.Exit(1)
}
