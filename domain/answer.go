package domain

// Answer represent answers's data model
type Answer struct {
	Type           string
	SequenceNumber int
	IsCorrect      bool
	IsLastSequence bool
	FileLocation   string
	Text           string
}
