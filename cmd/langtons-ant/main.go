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
	"time"

	"github.com/mikehelmick/go-vestaboard"
	"github.com/mikehelmick/go-vestaboard/internal/config"
)

const (
	rows = 6
	cols = 22
)

type direction int

const (
	north direction = iota
	east
	south
	west
)

func (d direction) turnRight() direction { return (d + 1) % 4 }
func (d direction) turnLeft() direction  { return (d + 3) % 4 }

func (d direction) delta() (int, int) {
	switch d {
	case north:
		return -1, 0
	case east:
		return 0, 1
	case south:
		return 1, 0
	default: // west
		return 0, -1
	}
}

func wrap(v, m int) int {
	return ((v % m) + m) % m
}

// langtonsAnt is the classic single-ant cellular automaton: on an "off"
// cell it turns right and switches the cell on, on an "on" cell it turns
// left and switches the cell off, then steps forward. Exactly one cell
// flips per step. The board wraps at the edges, since the physical board
// is far smaller than the space this pattern usually needs to roam in
// before any structure emerges.
type langtonsAnt struct {
	grid [rows][cols]bool
	x, y int
	dir  direction
}

func newLangtonsAnt() *langtonsAnt {
	return &langtonsAnt{x: rows / 2, y: cols / 2, dir: north}
}

func (a *langtonsAnt) step() {
	if a.grid[a.x][a.y] {
		a.dir = a.dir.turnLeft()
		a.grid[a.x][a.y] = false
	} else {
		a.dir = a.dir.turnRight()
		a.grid[a.x][a.y] = true
	}

	dx, dy := a.dir.delta()
	a.x = wrap(a.x+dx, rows)
	a.y = wrap(a.y+dy, cols)
}

func (a *langtonsAnt) layout() vestaboard.Layout {
	l := vestaboard.NewLayout()
	for x := 0; x < rows; x++ {
		for y := 0; y < cols; y++ {
			if a.grid[x][y] {
				l.SetColor(x, y, vestaboard.Green)
			}
		}
	}
	l.SetColor(a.x, a.y, vestaboard.White)
	return l
}

// Langton's Ant.
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

	ant := newLangtonsAnt()

	for {
		if _, err := client.SendMessage(ctx, subs.Subscriptions[0].ID, ant.layout()); err != nil {
			log.Fatalf("error sending message: %v", err)
		}

		ant.step()

		// Vestaboard rate-limits the platform API to about one message
		// every 15 seconds; sending faster than that drops frames.
		time.Sleep(15 * time.Second)
	}
}
