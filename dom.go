package main

/***************************************
   TODO:
    Unit tests!
***************************************/

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

import pl "github.com/nickaubert/dominionator/players"
import yaml "gopkg.in/yaml.v2"

func main() {

	confFile := flag.String("c", "dom.yaml", "config.yaml")
	rounds := flag.Int("r", 1, "rounds to play")
	flag.Parse()

	cnf := checkConfig(*confFile)
	defer cnf.Buffer.Flush()

	fmt.Println("Dominion!")
	for r := 0; r < *rounds; r++ {
		pg := playDominion(cnf)
		pl.CheckScores(pg)
	}

}

func checkConfig(file string) pl.Config {

	var cnf pl.Config

	yf, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading config: ", err)
		os.Exit(1)
	}

	err = yaml.Unmarshal(yf, &cnf)
	if err != nil {
		fmt.Println("Error parsing config: ", err)
		os.Exit(1)
	}

	if pl.ValidateCards(cnf.Kingdom) != true {
		fmt.Println("Error: Invalid Dominion card!")
		os.Exit(2)
	}

	fmt.Println("logfile:", cnf.Logfile)
	w := bufio.NewWriter(ioutil.Discard)
	if cnf.Logfile != "" {
		f, err := os.Create(cnf.Logfile)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(2)
		}
		w = bufio.NewWriter(f)
	}
	cnf.Buffer = w

	return cnf

}

func playDominion(cnf pl.Config) pl.Playgroup {

	pg := pl.InitializePlaygroup(cnf)
	for n := range pg.Players {
		p := &pg.Players[n]
		pl.ShuffleDeck(p)
		nc := pl.Draw(p, 5)
		p.Hand.Cards = append(p.Hand.Cards, nc...)
	}

	fmt.Fprintln(cnf.Buffer, "starting supply:")
	cnf.Buffer.Flush()
	for _, p := range pg.Supply.Piles {
		fmt.Fprintln(cnf.Buffer, "pile", p.Count, p.Card.Name)
	}

	turnCount := 0
	for {

		turnCount++
		fmt.Fprintf(cnf.Buffer, "Turn %d: ", turnCount)
		endGame := pl.PlayTurn(&pg, cnf)

		if endGame == true {
			break
		}

		if turnCount > 200 {
			fmt.Fprintln(cnf.Buffer, "Interrupted game at turn 201")
			break
		}

	}

	fmt.Fprintln(cnf.Buffer, "ending supply:")
	for _, p := range pg.Supply.Piles {
		fmt.Fprintln(cnf.Buffer, "pile", p.Count, p.Card.Name)
	}
	fmt.Fprintln(cnf.Buffer, "ending trash: ")
	for _, c := range pg.Trash.Cards {
		fmt.Fprint(cnf.Buffer, c.Name, ", ")
	}
	fmt.Fprint(cnf.Buffer, "\n")
	fmt.Fprintln(cnf.Buffer)

	return pg
}
