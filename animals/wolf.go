package animals

import "math/rand"

const WolfDefCount = 40
const WolfType = "wolf"

func NewWolfs(startCount int) *animals {
	return NewAnimals(startCount, newWolf, WolfType)
}

func newWolf() *animal {
	return &animal{
		speed: speed{
			run:     3,
			step:    2,
			stealth: 1,
		},
		typeOf: WolfType,
		sex:    rand.Intn(99) > 49,
		age:    0,
		breedAge: breedAge{
			min: 3,
			max: 13,
		},
		deathAge:    15,
		maxBirths:   3,
		minBirths:   1,
		birthRadius: 3,
	}
}
