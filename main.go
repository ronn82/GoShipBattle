package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	emptySpot      int = 0
	computerShip   int = 1
	userShip       int = 2
	computerGotHit int = 3
	computerMissed int = 4
	userGotHit     int = 5
	userMissed     int = 6
)

type coords struct {
	x      int
	y      int
	status int
}

type game struct {
	ocean         []coords
	title         string
	userShips     int
	computerShips int
	isNew         bool
	gridSize      int
}

var newGame game

func main() {
	whosTurn := "user"

	newGame.isNew = true

	for {

		drawGame(&newGame)

		switch whosTurn {
		case "user":
			Attack(false, &newGame)
		case "computer":
			Attack(true, &newGame)
		}

		if whosTurn == "user" {
			whosTurn = "computer"
		} else {
			whosTurn = "user"
		}
		if newGame.computerShips == 0 || newGame.userShips == 0 {

			if newGame.computerShips > newGame.userShips {
				fmt.Println("User wins!")
			} else {
				fmt.Println("Computer wins!")
			}
			os.Exit(0)
		}
	}
}

func drawGame(ngame *game) {

	if ngame.isNew {
		ngame.computerShips = 0
		ngame.userShips = 0
		ngame.title = "~~ Welcome to BattleShip ~~"
		var ocean []coords
		ngame.ocean = ocean
		ngame.gridSize = 10
		ngame.isNew = false
		//Cleaning the grid
		for i := 0; i < ngame.gridSize; i++ {
			for j := 0; j < ngame.gridSize; j++ {
				newCoord := coords{i, j, 0}
				ngame.ocean = append(ngame.ocean, newCoord)
			}
		}
		fmt.Println("The sea is empty, lets deploy the fleet...")
		deployShips(ngame)
	}

	//Drawing game
	fmt.Println()
	fmt.Println(ngame.title)
	fmt.Printf("User ships: %d ---- Computer ships: %d", ngame.userShips, ngame.computerShips)
	fmt.Println()

	for i := 0; i < len(ngame.ocean); i++ {

		switch ngame.ocean[i].status {

		case emptySpot:
			fmt.Print(" [ ]")

		case computerShip:
			fmt.Print(" [ ]")

		case userShip:
			fmt.Print(" [@]")

		case computerGotHit:
			fmt.Print(" [X]")

		case computerMissed:
			fmt.Print(" [ ]")

		case userGotHit:
			fmt.Print(" [-]")

		case userMissed:
			fmt.Print(" [0]")
		}

		if (i+1)%ngame.gridSize == 0 {
			fmt.Printf("-%d-", (i)/10)
			fmt.Println()
		}
	}

	for i := 0; i < 10; i++ {
		fmt.Printf("--%d-", i)
	}
	fmt.Println()
}

func deployShips(ngame *game) {

	for {
		x, y := readCoordsFromStdin("Enter coordinates to attack, Example: 5,3: ")

		if x > ngame.gridSize || y > ngame.gridSize {
			fmt.Println("The ship cannot be placed outside the ocean.")
			continue
		}

		for i := 0; i < len(ngame.ocean); i++ {
			if ngame.ocean[i].x == x && ngame.ocean[i].y == y {
				if ngame.ocean[i].status == emptySpot {
					ngame.ocean[i].status = userShip
					fmt.Println("Deploying ship at: ", x, y)
					ngame.userShips++
				} else {
					fmt.Println("Spot taken, please try again")
				}
			}
		}
		if ngame.userShips == 5 {
			break
		}
	}

	rand.Seed(int64(time.Now().UnixNano()))
	for {
		coord := coords{rand.Intn(ngame.gridSize), rand.Intn(ngame.gridSize), computerShip}

		for i := 0; i < len(ngame.ocean); i++ {
			if ngame.ocean[i].x == coord.x && ngame.ocean[i].y == coord.y && ngame.ocean[i].status == emptySpot {
				ngame.ocean[i].status = computerShip
				ngame.computerShips++
				//For testing only
				fmt.Printf("Deploying computer at: %d, %d", coord.x, coord.y)
			}
		}
		if ngame.computerShips == 5 {
			break
		}
	}
}

func readCoordsFromStdin(title string) (int, int) {
	fmt.Print(title)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	x, errx := strconv.Atoi(strings.Split(text, ",")[0])
	y, erry := strconv.Atoi(strings.Split(text, ",")[1])

	if errx != nil || erry != nil {
		os.Exit(1)
	}
	return x, y
}

func Attack(isComputer bool, ngame *game) {
	var x int
	var y int

	if isComputer {
		rand.Seed(time.Now().UnixNano())
		x = rand.Intn(10)
		y = rand.Intn(10)
	} else {
		x, y = readCoordsFromStdin("Enter coordinates to attack, Example: 5,3: ")
	}

	for i := 0; i < len(ngame.ocean); i++ {
		if ngame.ocean[i].x == x && ngame.ocean[i].y == y {

			switch ngame.ocean[i].status {
			case emptySpot:
				fmt.Println("hit the water")
				if isComputer {
					ngame.ocean[i].status = computerMissed
				} else {
					ngame.ocean[i].status = userMissed
				}

			case userShip:
				if isComputer {
					fmt.Println("You got hit")
					ngame.ocean[i].status = userGotHit
					ngame.userShips--
				} else {
					fmt.Println("Friendly fire, you sank your own")
					ngame.ocean[i].status = userGotHit
					ngame.userShips--
				}

			case computerShip:
				if isComputer {
					fmt.Println("Your enemy is doing your job")
					ngame.computerShips--
					ngame.ocean[i].status = computerGotHit
				} else {
					fmt.Println("Good job, hit the target")
					ngame.computerShips--
					ngame.ocean[i].status = computerGotHit
				}
			default:
				fmt.Println("Hitting nada")

			}
		}
	}
}
