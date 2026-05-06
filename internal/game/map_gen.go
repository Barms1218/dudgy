package game

import (
	t "github.com/Barms1218/dudgy/internal/types"
	"math/rand"
)

func GenerateMap(width, height, seed int16) (t.Map, t.Position) {
	tiles := make([]int8, width*height)
	for i := range tiles {
		tiles[i] = 1
	}

	r := rand.New(rand.NewSource(int64(seed)))
	x, y := int(width/2), int(height/2)
	steps := int(width * height / 3)
	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

	for range steps {
		tiles[y*int(width)+x] = 0
		d := dirs[r.Intn(4)]
		nx, ny := x+d[0], y+d[1]
		if nx > 0 && nx < int(width)-1 && ny > 0 && ny < int(height)-1 {
			x, y = nx, ny
		}
	}

	spawn := t.Position{X: float32(width / 2), Y: float32(height / 2)}
	return t.Map{Width: width, Height: height, Tiles: tiles}, spawn
}
