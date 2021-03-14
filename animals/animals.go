package animals

type (
	animals struct {
		animals       []*Animal
		count         int
		lastIndex     int
		internalIndex int
		startCount    int
		newAnimal     func() *Animal
		animalType    string
		size          int
	}
	Animals interface {
		Len() *int
		Append(*Animal)
		Remove(int)
		StartCount() int
		Count() int
		InternalIndex() int
		Animal(int) *Animal
		NewAnimal() *Animal
		Type() string
	}
)

func (a *animals) Len() *int {
	return &a.count
}

func (a *animals) Append(newAnimal *Animal) {
	if a.size <= a.count {
		a.animals = append(a.animals, newAnimal)
		a.size++
	} else {
		a.animals[a.lastIndex] = newAnimal
	}
	newAnimal.index = a.lastIndex
	a.count++
	a.lastIndex++
	a.internalIndex++
}

func (a *animals) Remove(index int) {
	lastIndex := a.lastIndex - 1
	if lastIndex <= 0 || lastIndex == index {
		a.animals[index] = nil
	} else {
		a.animals[index] = a.animals[lastIndex]
		a.animals[lastIndex] = nil
		a.animals[index].index = index
	}

	a.count--
	a.lastIndex--
}

func (a *animals) StartCount() int {
	return a.startCount
}

func (a *animals) NewAnimal() *Animal {
	return a.newAnimal()
}

func (a *animals) InternalIndex() int {
	return a.internalIndex
}

func (a *animals) Count() int {
	return a.count
}

func (a *animals) Animal(index int) *Animal {
	return a.animals[index]
}

func (a *animals) Type() string {
	return a.animalType
}

func newAnimals(startCount int, newAnimal func() *Animal, animalType string) *animals {
	return &animals{startCount: startCount, newAnimal: newAnimal, animalType: animalType}
}
