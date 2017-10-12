package main

/***************************************
   TODO:
    Unit tests!
   Randomizer deck
***************************************/

import (
	"fmt"
	"io/ioutil"
	"os"
)

import pl "github.com/nickaubert/dominionator/players"

// import bs "github.com/nickaubert/dominionator/basic"
import yaml "gopkg.in/yaml.v2"

func main() {

	cnf := checkConfig("dom.yaml")

	fmt.Println("Dominion!")
	fmt.Println()

	pg := pl.InitializePlaygroup(cnf)
	for n := range pg.Players {
		p := &pg.Players[n]
		pl.ShuffleDeck(p)
		nc := pl.Draw(p, 5)
		p.Hand.Cards = append(p.Hand.Cards, nc...)
	}

	fmt.Println("starting supply:")
	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Println()

	turnCount := 0
	for {

		turnCount++
		fmt.Printf("Turn %d: ", turnCount)
		endGame := pl.PlayTurn(&pg)

		if endGame == true {
			break
		}

		if turnCount > 200 {
			fmt.Println("Interrupted game at turn 201")
			break
		}

	}

	fmt.Println(turnCount, "turns")
	fmt.Println("ending supply:")
	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Print("ending trash: ")
	for _, c := range pg.Trash.Cards {
		fmt.Print(c.Name, ", ")
	}
	fmt.Print("\n")
	fmt.Println()

	pl.CheckScores(pg)

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

	for _, kc := range cnf.Kingdom {
		fmt.Println("cards", kc)
	}

	/*
	   for _, c := range pl.initializeRandomizer(10) {
	       s.Piles = append(s.Piles, cd.SupplyPile{Card: c, Count: 10})
	   }
	*/

	return cnf

}
