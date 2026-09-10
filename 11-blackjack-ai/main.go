package main

import (
	"fmt"

	deck "github.com/AbePlays/gophercises/09-deck-of-cards"
	"github.com/AbePlays/gophercises/11-blackjack-ai/blackjack"
)

type basicAi struct {
	score int
	seen  int
	decks int
}

func (ai *basicAi) Bet(shuffled bool) int {
	if shuffled {
		ai.score = 0
		ai.seen = 0
	}

	trueScore := ai.score / ((ai.decks*52 - ai.seen) / 52)
	switch {
	case trueScore >= 14:
		return 10_000
	case trueScore >= 8:
		return 500
	default:
		return 100
	}
}

func (ai *basicAi) Play(hand []deck.Card, dealer deck.Card) blackjack.Move {
	score := blackjack.Score(hand...)
	if len(hand) == 2 {
		if hand[0] == hand[1] {
			cardScore := blackjack.Score(hand[0])
			if cardScore >= 8 && cardScore != 10 {
				return blackjack.MoveSplit
			}
		}

		if (score == 10 || score == 11) && !blackjack.Soft(hand...) {
			return blackjack.MoveDouble
		}
	}

	dealerScore := blackjack.Score(dealer)
	if dealerScore >= 5 && dealerScore <= 6 {
		return blackjack.MoveStand
	}

	if score < 13 {
		return blackjack.MoveHit
	}

	return blackjack.MoveStand
}

func (ai *basicAi) Results(hands [][]deck.Card, dealer []deck.Card) {
	for _, card := range dealer {
		ai.count(card)
	}

	for _, hand := range hands {
		for _, card := range hand {
			ai.count(card)
		}
	}
}

func (ai *basicAi) count(card deck.Card) {
	score := blackjack.Score(card)

	switch {
	case score >= 10:
		ai.score--
	case score <= 6:
		ai.score++
	}
	ai.seen++
}

func main() {
	options := blackjack.Options{
		Decks:           4,
		Hands:           50_000,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(options)
	winnings := game.Play(&basicAi{
		decks: 4,
	})

	fmt.Println("You won:", winnings)
}
