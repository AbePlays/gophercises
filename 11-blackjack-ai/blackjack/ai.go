package blackjack

import (
	"fmt"

	deck "github.com/AbePlays/gophercises/09-deck-of-cards"
)

type Ai interface {
	Bet(shuffled bool) int
	Play(hand []deck.Card, dealer deck.Card) Move
	Results(hands [][]deck.Card, dealer []deck.Card)
}

type dealerAi struct{}

func (ai dealerAi) Bet(shuffled bool) int {
	return 1
}

func (ai dealerAi) Play(hand []deck.Card, dealer deck.Card) Move {
	dealerScore := Score(hand...)
	if dealerScore <= 16 || (dealerScore == 17 && Soft(hand...)) {
		return MoveHit
	}

	return MoveStand
}

func (ai dealerAi) Results(hands [][]deck.Card, dealer []deck.Card) {
}

func HumanAi() Ai {
	return humanAi{}
}

type humanAi struct{}

func (ai humanAi) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("Deck was just shuffled")
	}
	fmt.Println("What would you like to bet?")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

func (ai humanAi) Play(hand []deck.Card, dealer deck.Card) Move {
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)

		fmt.Println("What will you do? (h)it, (s)tand, (d)ouble, (sp)lit")

		var input string
		fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		case "sp":
			return MoveSplit
		default:
			fmt.Println("Invalid input. Please enter 'h', 's', 'd', or 'sp'.")
		}
	}
}

func (ai humanAi) Results(hands [][]deck.Card, dealer []deck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:")
	for _, hand := range hands {
		fmt.Println(" ", hand)
	}
	fmt.Println("Dealer:", dealer)
}
