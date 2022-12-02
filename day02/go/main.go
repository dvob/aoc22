package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Sign int

const (
	Rock Sign = iota
	Paper
	Scissors
)

func (s Sign) String() string {
	switch s {
	case Rock:
		return "Rock"
	case Paper:
		return "Paper"
	case Scissors:
		return "Scissors"
	default:
		panic(fmt.Sprintf("invalid sign: %d", int(s)))
	}
}

var (
	Signs = map[string]Sign{
		"A": Rock,
		"B": Paper,
		"C": Scissors,
		"X": Rock,
		"Y": Paper,
		"Z": Scissors,
	}
	// mapping for second part
	Results = map[Sign]Result{
		Rock:     Lose,
		Paper:    Draw,
		Scissors: Win,
	}
	SignPoints = map[Sign]int{
		Rock:     1,
		Paper:    2,
		Scissors: 3,
	}
	ResultPoints = map[Result]int{
		Win:  6,
		Draw: 3,
		Lose: 0,
	}
)

type Game struct {
	My       Sign
	Opponent Sign
}

type Result int

const (
	Win Result = iota
	Draw
	Lose
)

func (r Result) String() string {
	switch r {
	case Win:
		return "Win"
	case Draw:
		return "Draw"
	case Lose:
		return "Lose"
	default:
		panic("invalid result")
	}
}

func getResult(game Game) Result {
	// Draw
	if game.My == game.Opponent {
		return Draw
	}

	// Win
	if game.My == Rock && game.Opponent == Scissors {
		return Win
	}
	if game.My == Paper && game.Opponent == Rock {
		return Win
	}
	if game.My == Scissors && game.Opponent == Paper {
		return Win
	}

	// Lose
	return Lose
}

func getScore(game Game) int {
	result := getResult(game)
	return SignPoints[game.My] + ResultPoints[result]
}

func getExpectedSign(sign Sign, result Result) Sign {
	var expectedSign Sign
	switch result {
	case Win:
		expectedSign = (sign + 1) % 3
	case Draw:
		expectedSign = sign
	case Lose:
		expectedSign = (sign - 1) % 3
		if expectedSign == -1 {
			expectedSign = 2
		}
	}
	return expectedSign
}

func getScore2(game Game) int {
	result := Results[game.My]
	expectedSign := getExpectedSign(game.Opponent, result)
	score := SignPoints[expectedSign] + ResultPoints[result]
	//fmt.Printf("op=%s, my=%s, result=%s, expected sign=%s, sign points=%d, result points=%d, score=%d\n", game.Opponent, game.My, result, expectedSign, SignPoints[expectedSign], ResultPoints[result], score)
	return score
}

func getScores(games []Game, scoreFn func(Game) int) int {
	score := 0
	for _, game := range games {
		score += scoreFn(game)
	}
	return score
}

func readGame(line string) (Game, error) {
	var game Game

	line = strings.TrimSpace(line)
	op, my, found := strings.Cut(line, " ")
	if !found {
		return game, fmt.Errorf("separator not found")
	}

	game.My, found = Signs[my]
	if !found {
		return game, fmt.Errorf("unknown sign %s", my)
	}

	game.Opponent, found = Signs[op]
	if !found {
		return game, fmt.Errorf("unknown sign %s", my)
	}

	return game, nil
}

func readGames(in io.Reader) ([]Game, error) {
	games := []Game{}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		if scanner.Err() != nil {
			return nil, scanner.Err()
		}

		game, err := readGame(scanner.Text())
		if err != nil {
			return nil, err
		}

		games = append(games, game)
	}

	return games, nil
}

func readGamesFile(file string) ([]Game, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return readGames(f)
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("missing argument: filename")
	}

	file := os.Args[1]

	games, err := readGamesFile(file)
	if err != nil {
		return err
	}

	score := getScores(games, getScore)

	fmt.Println(score)

	score2 := getScores(games, getScore2)

	fmt.Println(score2)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
