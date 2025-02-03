package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
)

const (
	baseUrl       = "http://localhost:8081"
	createPostfix = "/notes"
	getPostfix    = "/notes/%d"
)

type NoteInfo struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Context  string `json:"context"`
	IsPublic bool   `json:"is_public"`
}

type Note struct {
	ID        int64     `json:"id"`
	Info      NoteInfo  `json:"info"`
	CreatedAt time.Time `json:"createdJ_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func main() {
	note, err := CreateNoteClient()
	if err != nil {
		log.Fatal("Failed to create note:", err)
	}
	log.Println(color.RedString("Note created:\n"), color.BlueString("%+v", note))

	note, err = GetNoteClient(note.ID)
	if err != nil {
		log.Fatal("Failed to get note:", err)
	}
	log.Println(color.GreenString("Note info got:\n"), color.CyanString("%+v", note))

}

func CreateNoteClient() (Note, error) {
	info := &NoteInfo{
		Title:    gofakeit.BeerName(),
		Author:   gofakeit.StreetName(),
		Context:  gofakeit.IPv4Address(),
		IsPublic: gofakeit.Bool(),
	}
	noteData, err := json.Marshal(info)
	if err != nil {
		return Note{}, err
	}
	resp, err := http.Post(baseUrl+createPostfix, "application/json", bytes.NewBuffer(noteData))
	if err != nil {
		return Note{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return Note{}, err
	}

	var createdNote Note
	if err := json.NewDecoder(resp.Body).Decode(&createdNote); err != nil {
		return Note{}, err
	}
	return createdNote, nil
}

func GetNoteClient(noteID int64) (Note, error) {
	resp, err := http.Get(fmt.Sprintf(baseUrl+getPostfix, noteID))
	if err != nil {
		log.Fatal("Failed to get note:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Note{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Note{}, errors.New(fmt.Sprintf("failed to get note: %d", resp.StatusCode))
	}
	var note Note
	if err := json.NewDecoder(resp.Body).Decode(&note); err != nil {
		return Note{}, err
	}
	return note, nil
}
