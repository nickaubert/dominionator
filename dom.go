package main

/***************************************
   TODO:
   Randomizer deck
   Throne Room (make sure card doesnt duplicate)
***************************************/

import (
	"fmt"
)

import pl "github.com/nickaubert/dominionator/players"

func main() {

	fmt.Println("Dominion!")
	fmt.Println()

	pg := pl.InitializePlaygroup(3)
	for n := range pg.Players {
		pl.ShuffleDeck(&pg.Players[n])
		nc := pl.Draw(&pg.Players[n], 5)
		pl.AddHand(&pg.Players[n], nc)
	}

	fmt.Println("starting supply:")
	for _, p := range pg.Supply.Piles {
		fmt.Println("pile", p.Count, p.Card.Name)
	}
	fmt.Println()

	turnCount := 0
	for {

		turnCount++
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
	for _, c := range pg.Trash {
		fmt.Print(c.Name, ", ")
	}
	fmt.Print("\n")
	fmt.Println()

	pl.CheckScores(pg)

}
