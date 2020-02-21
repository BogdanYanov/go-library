package library

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

const yearOfFirstBook = 868

// Wastepaper describes the methods necessary for working with the library
type Wastepaper interface {
	Take()
	Put()
	Lock()
	Unlock()
	GetText() string
	IsTaken() bool
	Info()
}

// Book example implementation of wastepaper
type Book struct {
	author         string
	publishingYear int
	text           string
	mu             *sync.Mutex
	isTaken        bool
}

// NewBook creates new book
func NewBook(author string, publishingYear int, text string) Wastepaper {
	if len(strings.TrimSpace(author)) == 0 {
		return nil
	}
	if len(strings.TrimSpace(text)) == 0 {
		return nil
	}
	yearNow := time.Now().Year()
	if publishingYear < yearOfFirstBook || publishingYear > yearNow {
		return nil
	}
	return &Book{
		author:         author,
		publishingYear: publishingYear,
		text:           text,
		mu:             &sync.Mutex{},
		isTaken:        false,
	}
}

// Take forbids other readers to read this book
func (b *Book) Take() {
	b.isTaken = true
}

// Put allows other readers to read this book.
func (b *Book) Put() {
	b.isTaken = false
}

// IsTaken checks if the book is taken from the library
func (b *Book) IsTaken() bool {
	return b.isTaken
}

// Lock locks mu for book
func (b *Book) Lock() {
	b.mu.Lock()
}

// Unlock unlocks mu for book
func (b *Book) Unlock() {
	b.mu.Unlock()
}

// GetText return text of book
func (b *Book) GetText() string {
	return b.text
}

// Info displays information about book
func (b *Book) Info() {
	fmt.Printf("Book by %s (%d publishing year)\n", b.author, b.publishingYear)
}
