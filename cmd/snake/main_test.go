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
	"math/rand"
	"testing"
)

func TestSnakeMovesTowardFood(t *testing.T) {
	t.Parallel()

	g := &snakeGame{
		body: []point{{0, 5}, {0, 4}, {0, 3}},
		food: point{0, 8},
		rng:  rand.New(rand.NewSource(1)),
	}
	if !g.step() {
		t.Fatalf("expected step to succeed")
	}
	if got, want := g.body[0], (point{0, 6}); got != want {
		t.Errorf("head = %+v, want %+v", got, want)
	}
	if len(g.body) != 3 {
		t.Errorf("length = %d, want 3 (no growth)", len(g.body))
	}
}

func TestSnakeGrowsWhenEatingFood(t *testing.T) {
	t.Parallel()

	g := &snakeGame{
		body: []point{{0, 5}, {0, 4}, {0, 3}},
		food: point{0, 6},
		rng:  rand.New(rand.NewSource(1)),
	}
	if !g.step() {
		t.Fatalf("expected step to succeed")
	}
	if got, want := g.body[0], (point{0, 6}); got != want {
		t.Errorf("head = %+v, want %+v", got, want)
	}
	if len(g.body) != 4 {
		t.Errorf("length = %d, want 4 (grew by one)", len(g.body))
	}
	if g.food == (point{0, 6}) {
		t.Errorf("expected new food to be placed off the snake's head")
	}
}

func TestSnakeGameOverWhenTrapped(t *testing.T) {
	t.Parallel()

	// Head is boxed into the corner by its own body, with the tail too far
	// away to vacate a cell the head could move into.
	g := &snakeGame{
		body: []point{{0, 0}, {0, 1}, {1, 1}, {1, 0}, {2, 0}},
		food: point{5, 21},
		rng:  rand.New(rand.NewSource(1)),
	}
	if g.step() {
		t.Fatalf("expected step to fail, snake is trapped")
	}
}
