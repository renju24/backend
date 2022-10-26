package main

import (
	"fmt"
	"log"

	"github.com/renju24/backend/pkg/game"
)

func main() {
	g := game.NewGame()
	var user int
	for {
		var x, y int
		fmt.Print("Enter move: ")
		if _, err := fmt.Scanf("%d %d", &x, &y); err != nil {
			log.Fatalln(err)
		}
		switch user {
		case 1:
			user = 2
		case 2:
			user = 1
		default:
			user = 1
		}
		winner, err := g.ApplyMove(game.NewMove(x, y, user))
		if err != nil {
			log.Fatalln(err)
		}
		if winner != 0 {
			fmt.Println(winner, "won!")
			return
		}
	}
}
