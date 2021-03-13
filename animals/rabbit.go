package animals

import "math/rand"

const RabbitDefCount = 50
const RabbitType = "rabbit"

func NewRabbits(startCount int) Animals {
	return NewAnimals(startCount, newRabbit, RabbitType)
}

func newRabbit() *animal {
	return &animal{
		speed: speed{
			run:     4,
			step:    2,
			stealth: 2,
		},
		typeOf: RabbitType,
		sex:    rand.Intn(1) > 0,
		age:    0,
		breedAge: breedAge{
			min: 1,
			max: 8,
		},
		deathAge:    9,
		maxBirths:   5,
		minBirths:   2,
		birthRadius: 5,
	}
}
