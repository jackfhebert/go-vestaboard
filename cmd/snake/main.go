// Copyright 2021 Mike Helmick
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/mikehelmick/go-vestaboard"
	"github.com/mikehelmick/go-vestaboard/internal/config"
)

const (
	rows = 6
	cols = 22
)

type point struct {
	x, y int
}

var directions = []point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

// snakeGame is a self-playing snake: the head always greedily steps toward
// the food, growing by one segment each time it's eaten.
type snakeGame struct {
	body []point // body[0] is the head, last element is the tail
	food point
	rng  *rand.Rand
}

func newSnakeGame(rng *rand.Rand) *snakeGame {
	g := &snakeGame{
		body: []point{{2, 12}, {2, 11}, {2, 10}},
		rng:  rng,
	}
	g.placeFood()
	return g
}

func (g *snakeGame) occupied(p point) bool {
	for _, b := range g.body {
		if b == p {
			return true
		}
	}
	return false
}

func (g *snakeGame) placeFood() {
	for {
		p := point{g.rng.Intn(rows), g.rng.Intn(cols)}
		if !g.occupied(p) {
			g.food = p
			return
		}
	}
}

func inBounds(p point) bool {
	return p.x >= 0 && p.x < rows && p.y >= 0 && p.y < cols
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func distance(a, b point) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

// step advances the game by one tick, moving the head into whichever
// in-bounds, non-colliding neighbor cell is closest to the food. It returns
// false if the snake has nowhere left to move.
func (g *snakeGame) step() bool {
	head := g.body[0]
	tail := g.body[len(g.body)-1]

	found := false
	var best point
	bestDist := 0
	for _, d := range directions {
		next := point{head.x + d.x, head.y + d.y}
		if !inBounds(next) {
			continue
		}
		// The tail vacates this tick, so it's always safe to move into,
		// even though it's still in g.body.
		if next != tail && g.occupied(next) {
			continue
		}
		if dist := distance(next, g.food); !found || dist < bestDist {
			best, bestDist, found = next, dist, true
		}
	}
	if !found {
		return false
	}

	ate := best == g.food
	g.body = append([]point{best}, g.body...)
	if !ate {
		g.body = g.body[:len(g.body)-1]
	} else if len(g.body) < rows*cols {
		g.placeFood()
	}
	return true
}

func (g *snakeGame) layout() vestaboard.Layout {
	l := vestaboard.NewLayout()
	l.SetColor(g.food.x, g.food.y, vestaboard.PoppyRed)
	for i, b := range g.body {
		c := vestaboard.Green
		if i == 0 {
			c = vestaboard.White
		}
		l.SetColor(b.x, b.y, c)
	}
	return l
}

// Self-playing snake.
func main() {
	flag.Parse()

	ctx := context.Background()
	c, err := config.New(ctx)
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	client := vestaboard.New(c.APIKey, c.Secret)

	subs, err := client.Subscriptions(ctx)
	if err != nil {
		log.Fatalf("error calling Subscriptions: %v", err)
	}
	log.Printf("result: %+v", subs)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	game := newSnakeGame(rng)

	for {
		if _, err := client.SendMessage(ctx, subs.Subscriptions[0].ID, game.layout()); err != nil {
			log.Fatalf("error sending message: %v", err)
		}

		if !game.step() {
			game = newSnakeGame(rng)
		}

		// Vestaboard rate-limits the platform API to about one message
		// every 15 seconds; sending faster than that causes frames to be
		// dropped.
		time.Sleep(15 * time.Second)
	}
}
