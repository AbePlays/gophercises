package blackjack

import (
	"errors"
	"math"

	deck "github.com/AbePlays/gophercises/09-deck-of-cards"
)

const (
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type state int

type Options struct {
	Decks           int
	Hands           int
	BlackjackPayout float64
}

func New(options Options) Game {
	game := Game{
		state:    statePlayerTurn,
		dealerAi: dealerAi{},
		balance:  0,
	}

	if options.Decks == 0 {
		options.Decks = 3
	}

	if options.Hands == 0 {
		options.Hands = 100
	}

	if options.BlackjackPayout == 0.0 {
		options.BlackjackPayout = 1.5
	}

	game.nDecks = options.Decks
	game.nHands = options.Hands
	game.blackjackPayout = options.BlackjackPayout

	return game
}

type Game struct {
	nDecks          int
	nHands          int
	blackjackPayout float64

	deck  []deck.Card
	state state

	player    []hand
	handIdx   int
	playerBet int
	balance   int

	dealer   []deck.Card
	dealerAi Ai
}

func (g *Game) currentHand() *[]deck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player[g.handIdx].cards
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("It isnt currently any player's turn")
	}
}

type hand struct {
	cards []deck.Card
	bet   int
}

func bet(g *Game, ai Ai, shuffled bool) {
	bet := ai.Bet(shuffled)

	if bet < 100 {
		panic("Bet must be at least 100")
	}

	g.playerBet = bet
}

func deal(g *Game) {
	playerHand := make([]deck.Card, 0, 5)
	g.handIdx = 0
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card

	for range 2 {
		card, g.deck = drawCards(g.deck)
		playerHand = append(playerHand, card)
		card, g.deck = drawCards(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.player = []hand{{cards: playerHand, bet: g.playerBet}}
	g.state = statePlayerTurn
}

func (g *Game) Play(ai Ai) int {
	g.deck = nil
	min := 52 * g.nDecks / 3

	for i := 0; i < g.nHands; i++ {
		shuffled := false
		if len(g.deck) < min {
			g.deck = deck.New(deck.Deck(g.nDecks), deck.Shuffle)
			shuffled = true
		}

		bet(g, ai, shuffled)
		deal(g)

		if Blackjack(g.dealer...) {
			endRound(g, ai)
			continue
		}

		for g.state == statePlayerTurn {
			hand := make([]deck.Card, len(*g.currentHand()))
			copy(hand, *g.currentHand())
			move := ai.Play(hand, g.dealer[0])
			error := move(g)

			switch error {
			case errorBust:
				MoveStand(g)
			case nil:
				// noop
			default:
				panic(error)
			}
		}

		for g.state == stateDealerTurn {
			hand := make([]deck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAi.Play(hand, g.dealer[0])
			move(g)
		}

		endRound(g, ai)
	}

	return g.balance
}

var (
	errorBust = errors.New("Hand Score exceeded 21")
)

type Move func(*Game) error

func MoveSplit(g *Game) error {
	cards := g.currentHand()
	if len(*cards) != 2 {
		return errors.New("You can only split with two cards in your hand")
	}

	if (*cards)[0].Rank != (*cards)[1].Rank {
		return errors.New("Both cards must have the same rank to split")
	}

	g.player = append(g.player, hand{
		cards: []deck.Card{(*cards)[1]},
		bet:   g.player[g.handIdx].bet,
	})

	g.player[g.handIdx].cards = (*cards)[:1]
	return nil
}

func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = drawCards(g.deck)
	*hand = append(*hand, card)

	if Score(*hand...) > 21 {
		return errorBust
	}

	return nil
}

func MoveDouble(g *Game) error {
	if len(*g.currentHand()) != 2 {
		return errors.New("You can only double on a hand with 2 cards")
	}
	g.playerBet *= 2
	MoveHit(g)
	return MoveStand(g)
}

func MoveStand(g *Game) error {
	if g.state == stateDealerTurn {
		g.state++
		return nil
	}

	if g.state == statePlayerTurn {
		g.handIdx++
		if g.handIdx >= len(g.player) {
			g.state++
		}

		return nil
	}

	return errors.New("Invalid State")
}

func drawCards(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

func Score(hand ...deck.Card) int {
	minScore := minScore(hand...)

	if minScore > 11 {
		return minScore
	}

	for _, c := range hand {
		if c.Rank == deck.Ace {
			return minScore + 10
		}
	}

	return minScore
}

func Soft(hand ...deck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)

	return minScore != score
}

func minScore(hand ...deck.Card) int {
	score := 0

	for _, c := range hand {
		score += int(math.Min(10, float64(c.Rank)))
	}

	return score
}

func Blackjack(hand ...deck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

func endRound(g *Game, ai Ai) {
	dealerScore := Score(g.dealer...)
	dealerBlackjack := Blackjack(g.dealer...)
	allHands := make([][]deck.Card, 0, len(g.player))

	for i, hand := range g.player {
		cards := hand.cards
		allHands[i] = cards
		winnings := hand.bet

		playerScore := Score(cards...)
		playerBlackjack := Blackjack(cards...)

		switch {
		case playerBlackjack && dealerBlackjack:
			winnings = 0
		case dealerBlackjack:
			winnings *= -1
		case playerBlackjack:
			winnings = int(float64(winnings) * g.blackjackPayout)
		case playerScore > 21:
			winnings *= -1
		case dealerScore > 21:
			// win
		case playerScore > dealerScore:
			// win
		case dealerScore > playerScore:
			winnings *= -1
		case playerScore == dealerScore:
			winnings = 0
		}
		g.balance += winnings

	}

	ai.Results(allHands, g.dealer)
	g.player = nil
	g.dealer = nil
}
