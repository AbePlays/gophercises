package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Heart})
	fmt.Println(Card{Rank: Two, Suit: Heart})
	fmt.Println(Card{Rank: Three, Suit: Spade})
	fmt.Println(Card{Rank: Four, Suit: Spade})
	fmt.Println(Card{Rank: Nine, Suit: Club})
	fmt.Println(Card{Rank: Queen, Suit: Diamond})

	// Output:
	// Ace of Hearts
	// Two of Hearts
	// Three of Spades
	// Four of Spades
	// Nine of Clubs
	// Queen of Diamonds
}

func TestNew(t *testing.T) {
	cards := New()

	if len(cards) != 52 {
		t.Errorf("Expected 52 cards, got %d", len(cards))
	}
}

func TestDefaultSort(t *testing.T) {
	cards := New(DefaultSort)
	expected := Card{Rank: Ace, Suit: Spade}
	if cards[0] != expected {
		t.Errorf("Expected %v, got %v", expected, cards[0])
	}
}

func TestSort(t *testing.T) {
	cards := New(Sort(Less))
	expected := Card{Rank: Ace, Suit: Spade}
	if cards[0] != expected {
		t.Errorf("Expected %v, got %v", expected, cards[0])
	}
}

func TestJokers(t *testing.T) {
	cards := New(Jokers(3))
	count := 0
	for _, card := range cards {
		if card.Suit == Joker {
			count++
		}
	}
	if count != 3 {
		t.Errorf("Expected 3 jokers, got %d", count)
	}
}

func TestFilter(t *testing.T) {
	filter := func(card Card) bool {
		return card.Rank == Two || card.Rank == Three
	}

	cards := New(Filter(filter))
	for _, card := range cards {
		if card.Rank == Two || card.Rank == Three {
			t.Errorf("Expected rank not to be Two or Three, got %v", card.Rank)
		}
	}
}

func TestDeck(t *testing.T) {
	cards := New(Deck(3))
	if len(cards) != 156 {
		t.Errorf("Expected 156 cards, got %d", len(cards))
	}
}
