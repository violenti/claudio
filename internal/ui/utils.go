package ui

import (
	"math/rand/v2"
)

func Quote() string {

	pepe := []string{
		"I say, pay attention, son!",
		"That's a joke, son! I say, that's a joke!",
		"Look at me when I'm talkin' to ya, boy.",
		"Fast... fast! Fast as a... I say, fast as a chicken hawk!",
		"You're in a heap of trouble, boy.",
		"I'm not talking to you, I'm talking to myself!",
	}
	legth := len(pepe)
	randomIndex := rand.IntN(legth)
	return pepe[randomIndex]
}
