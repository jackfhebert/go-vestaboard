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

	// settleSteps are taken before any point is plotted, so the walk has
	// converged near the fractal's attractor before it starts drawing.
	settleSteps = 10
)

type point struct{ x, y int }

// corners of the triangle the chaos game jumps toward.
var corners = []point{
	{0, cols / 2},
	{rows - 1, 0},
	{rows - 1, cols - 1},
}

// chaosGame draws a Sierpinski triangle via the "chaos game": each tick it
// jumps halfway from the current point toward a randomly chosen corner and
// lights that cell. No recursion needed - the fractal emerges from pure
// randomness, one tile at a time.
type chaosGame struct {
	grid [rows][cols]bool
	cur  point
	rng  *rand.Rand
}

func newChaosGame(rng *rand.Rand) *chaosGame {
	g := &chaosGame{cur: corners[0], rng: rng}
	for i := 0; i < settleSteps; i++ {
		g.cur = g.midpointTowardRandomCorner()
	}
	return g
}

func (g *chaosGame) midpointTowardRandomCorner() point {
	c := corners[g.rng.Intn(len(corners))]
	return point{(g.cur.x + c.x) / 2, (g.cur.y + c.y) / 2}
}

func (g *chaosGame) step() {
	g.cur = g.midpointTowardRandomCorner()
	g.grid[g.cur.x][g.cur.y] = true
}

func (g *chaosGame) layout() vestaboard.Layout {
	l := vestaboard.NewLayout()
	for x := 0; x < rows; x++ {
		for y := 0; y < cols; y++ {
			if g.grid[x][y] {
				l.SetColor(x, y, vestaboard.Green)
			}
		}
	}
	l.SetColor(g.cur.x, g.cur.y, vestaboard.White)
	return l
}

// Sierpinski triangle via the chaos game.
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
	game := newChaosGame(rng)

	for {
		if _, err := client.SendMessage(ctx, subs.Subscriptions[0].ID, game.layout()); err != nil {
			log.Fatalf("error sending message: %v", err)
		}

		game.step()

		// Vestaboard rate-limits the platform API to about one message
		// every 15 seconds; sending faster than that drops frames.
		time.Sleep(15 * time.Second)
	}
}
