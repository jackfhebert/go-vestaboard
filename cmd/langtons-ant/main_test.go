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

import "testing"

func TestWrap(t *testing.T) {
	t.Parallel()

	cases := []struct{ v, m, want int }{
		{-1, 6, 5},
		{6, 6, 0},
		{0, 6, 0},
		{7, 22, 7},
		{-22, 22, 0},
	}
	for _, c := range cases {
		if got := wrap(c.v, c.m); got != c.want {
			t.Errorf("wrap(%d, %d) = %d, want %d", c.v, c.m, got, c.want)
		}
	}
}

func TestDirectionTurns(t *testing.T) {
	t.Parallel()

	d := north
	for _, want := range []direction{east, south, west, north} {
		d = d.turnRight()
		if d != want {
			t.Errorf("turnRight chain: got %d, want %d", d, want)
		}
	}
	for _, want := range []direction{west, south, east, north} {
		d = d.turnLeft()
		if d != want {
			t.Errorf("turnLeft chain: got %d, want %d", d, want)
		}
	}
}

// The first few steps of Langton's ant carve a small, deterministic square
// loop before behavior turns chaotic; this pins that known-good sequence.
func TestLangtonsAntInitialLoop(t *testing.T) {
	t.Parallel()

	a := &langtonsAnt{x: 2, y: 2, dir: north}

	a.step()
	if !a.grid[2][2] || a.x != 2 || a.y != 3 || a.dir != east {
		t.Fatalf("step 1: grid[2][2]=%v pos=(%d,%d) dir=%d", a.grid[2][2], a.x, a.y, a.dir)
	}

	a.step()
	if !a.grid[2][3] || a.x != 3 || a.y != 3 || a.dir != south {
		t.Fatalf("step 2: grid[2][3]=%v pos=(%d,%d) dir=%d", a.grid[2][3], a.x, a.y, a.dir)
	}

	a.step()
	if !a.grid[3][3] || a.x != 3 || a.y != 2 || a.dir != west {
		t.Fatalf("step 3: grid[3][3]=%v pos=(%d,%d) dir=%d", a.grid[3][3], a.x, a.y, a.dir)
	}

	a.step()
	if !a.grid[3][2] || a.x != 2 || a.y != 2 || a.dir != north {
		t.Fatalf("step 4: grid[3][2]=%v pos=(%d,%d) dir=%d", a.grid[3][2], a.x, a.y, a.dir)
	}

	// The ant returns to its already-lit starting cell, so this step turns
	// left and switches it back off instead of repeating the loop.
	a.step()
	if a.grid[2][2] || a.x != 2 || a.y != 1 || a.dir != west {
		t.Fatalf("step 5: grid[2][2]=%v pos=(%d,%d) dir=%d", a.grid[2][2], a.x, a.y, a.dir)
	}
}

func TestLangtonsAntWrapsAtEdges(t *testing.T) {
	t.Parallel()

	// Each ant starts on an off cell, so step() always turns right; the
	// initial direction is chosen so that turn lands the ant facing
	// straight into the edge it's sitting on.
	northAnt := &langtonsAnt{x: 0, y: 5, dir: west} // west -> north
	northAnt.step()
	if northAnt.x != rows-1 {
		t.Errorf("north wrap: x=%d, want %d", northAnt.x, rows-1)
	}

	southAnt := &langtonsAnt{x: rows - 1, y: 5, dir: east} // east -> south
	southAnt.step()
	if southAnt.x != 0 {
		t.Errorf("south wrap: x=%d, want 0", southAnt.x)
	}

	eastAnt := &langtonsAnt{x: 3, y: cols - 1, dir: north} // north -> east
	eastAnt.step()
	if eastAnt.y != 0 {
		t.Errorf("east wrap: y=%d, want 0", eastAnt.y)
	}

	westAnt := &langtonsAnt{x: 3, y: 0, dir: south} // south -> west
	westAnt.step()
	if westAnt.y != cols-1 {
		t.Errorf("west wrap: y=%d, want %d", westAnt.y, cols-1)
	}
}
