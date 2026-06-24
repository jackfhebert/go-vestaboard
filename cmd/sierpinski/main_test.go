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

func TestChaosGameStepIsMidpointOfACorner(t *testing.T) {
	t.Parallel()

	g := &chaosGame{cur: point{3, 11}, rng: rand.New(rand.NewSource(2))}
	prev := g.cur
	g.step()

	found := false
	for _, c := range corners {
		if want := (point{(prev.x + c.x) / 2, (prev.y + c.y) / 2}); g.cur == want {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("cur=%+v is not the midpoint of %+v and any corner", g.cur, prev)
	}
	if !g.grid[g.cur.x][g.cur.y] {
		t.Errorf("step() did not mark %+v as lit", g.cur)
	}
}

func TestChaosGameStaysInBounds(t *testing.T) {
	t.Parallel()

	g := newChaosGame(rand.New(rand.NewSource(1)))
	for i := 0; i < 1000; i++ {
		g.step()
		if g.cur.x < 0 || g.cur.x >= rows || g.cur.y < 0 || g.cur.y >= cols {
			t.Fatalf("point out of bounds: %+v", g.cur)
		}
	}
}

func TestNewChaosGameDoesNotPlotSettlingSteps(t *testing.T) {
	t.Parallel()

	g := newChaosGame(rand.New(rand.NewSource(3)))
	for x := 0; x < rows; x++ {
		for y := 0; y < cols; y++ {
			if g.grid[x][y] {
				t.Fatalf("expected no cells lit before the first step(), found (%d,%d)", x, y)
			}
		}
	}
}
