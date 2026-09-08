package main

import (
	"fmt"
	"math"
	"strings"

	deck "github.com/AbePlays/gophercises/09-deck-of-cards"
)

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))

	for i := range h {
		strs[i] = h[i].String()
	}

	return strings.Join(strs, ", ")
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func (h Hand) Score() int {
	minScore := h.MinScore()

	if minScore > 11 {
		return minScore
	}

	for _, c := range h {
		if c.Rank == deck.Ace {
			return minScore + 10
		}
	}

	return minScore
}

func (h Hand) MinScore() int {
	score := 0

	for _, c := range h {
		score += int(math.Min(10, float64(c.Rank)))
	}

	return score
}

func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.New(deck.Deck(3), deck.Shuffle)
	return ret
}

func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make([]deck.Card, 0, 5)
	ret.Dealer = make([]deck.Card, 0, 5)
	var card deck.Card

	for i := 0; i < 2; i++ {
		card, ret.Deck = drawCards(ret.Deck)
		ret.Player = append(ret.Player, card)
		card, ret.Deck = drawCards(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
	}
	ret.State = StatePlayerTurn

	return ret
}

func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++
	return ret
}

func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = drawCards(ret.Deck)
	*hand = append(*hand, card)

	if hand.Score() > 21 {
		return Stand(ret)
	}

	return ret
}

func EndHand(gs GameState) GameState {
	ret := clone(gs)
	playerScore := ret.Player.Score()
	dealerScore := ret.Dealer.Score()

	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:", ret.Player, "\nScore:", playerScore)
	fmt.Println("Dealer:", ret.Dealer, "\nScore:", dealerScore)

	switch {
	case playerScore > 21:
		fmt.Println("Player busts! Dealer wins.")
	case dealerScore > 21:
		fmt.Println("Dealer busts! Player wins.")
	case playerScore > dealerScore:
		fmt.Println("Player wins!")
	case dealerScore > playerScore:
		fmt.Println("Dealer wins!")
	case playerScore == dealerScore:
		fmt.Println("It's a tie!")
	}
	ret.Player = nil
	ret.Dealer = nil

	return ret
}

func main() {
	var gs GameState
	gs = Shuffle(gs)
	gs = Deal(gs)

	var input string
	for gs.State == StatePlayerTurn {
		fmt.Println("Player:", gs.Player)
		fmt.Println("Dealer:", gs.Dealer.DealerString())
		fmt.Println("What will you do? (h)it, (s)tand")
		fmt.Scanf("%s\n", &input)

		switch input {
		case "h":
			gs = Hit(gs)
		case "s":
			gs = Stand(gs)
		default:
			fmt.Println("Invalid input. Please enter 'h' or 's'.")
		}
	}

	for gs.State == StateDealerTurn {
		if gs.Dealer.Score() < 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
			gs = Hit(gs)
		} else {
			gs = Stand(gs)
		}
	}

	gs = EndHand(gs)
}

func drawCards(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

type State int

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("It isnt currently any player's turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make([]deck.Card, len(gs.Player)),
		Dealer: make([]deck.Card, len(gs.Dealer)),
	}

	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)

	return ret
}
