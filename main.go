package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var fp = flag.String("problems", "./problems.csv", "Path to CSV containing problems for the quiz")
var fd = flag.Int("duration", 30, "Duration of the quiz")

type game struct {
	questions [][]string
	answers   []string
}

func main() {
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	done := make(chan bool, 1)

	game, err := setup()
	if err != nil {
		fail(err)
	}

	// wait for signals
	go func() {
		<-sigs
		done <- true
	}()

	// wait for game to complete
	go func() {
		game.play()
		done <- true
	}()

	// wait for clock to run down
	go func() {
		c := 0
		for range time.Tick(time.Second) {
			c += 1
			if c >= *fd {
				break
			}
		}
		fmt.Printf("\n\n\033[1mTime’s up!\033[0m\n")
		done <- true
	}()

	<-done

	game.printResults()

	printThanks()
}

func printThanks() {
	fmt.Println("\033[1mThanks for playing!\033[0m")
}

func setup() (*game, error) {
	f, err := os.Open(*fp)
	if err != nil {
		return nil, fmt.Errorf("parsing csv: %v", err)
	}
	defer f.Close()

	r := csv.NewReader(f)
	qs, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parsing csv: %v", err)
	}

	return &game{questions: qs}, nil
}

func (g *game) play() {
	fmt.Printf("\n\033[1mIt’s quiz time!\033[0m\n\n")

	for i, q := range g.questions {
		fmt.Printf("%v: ", q[0])

		var a string
		fmt.Scanf("%s", &a)

		g.answers = append(g.answers, a)

		if a == g.questions[i][1] {
			fmt.Println("\033[1;32mCorrect!\033[0m")
		} else {
			fmt.Println("\033[1;31mIncorrect!\033[0m")
		}
	}
}

func (g *game) printResults() {
	total := len(g.questions)
	correct := 0

	var errs [][]string

	for i, q := range g.questions {
		if len(g.answers) <= i {
			continue
		}
		a := g.answers[i]

		if a != q[1] {
			errs = append(errs, append(q, a))
		} else {
			correct += 1
		}
	}

	fmt.Printf("\n\033[1mYou answered %d/%d questions correctly!\033[0m\n\n", correct, total)

	if len(errs) == 0 {
		return
	}

	fmt.Println("\033[1mErrors:\033[0m")
	for _, e := range errs {
		fmt.Printf("%v: %v (you answered %v)\n", e[0], e[1], e[2])
	}
	fmt.Println()
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "quizgame: %v\n", err)
	os.Exit(1)
}
