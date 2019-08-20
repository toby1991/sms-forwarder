package interfaces

type Smser interface {
	Index() uint
	Retrieve(chip Chipper)
}

