package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type TabooWord struct {
	ID     uint     `json:"id"`
	Word   string   `json:"word"`
	Taboos []string `json:"taboo"`
}

type WordStorage struct {
	words map[uint]*TabooWord
}

var (
	wordStorage *WordStorage
	wOnce       sync.Once
)

func GetWordStorage() (*WordStorage, error) {
	var err error
	wOnce.Do(func() {
		wordStorage = &WordStorage{
			words: make(map[uint]*TabooWord),
		}
		err = wordStorage.loadWords("words.json")
	})
	if err != nil {
		return nil, err
	}
	return wordStorage, nil
}

func (ws *WordStorage) loadWords(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read word file %s; %w", file, err)
	}

	var list []*TabooWord
	err = json.Unmarshal(data, &list)
	if err != nil {
		return fmt.Errorf("failed to unmarshal word file contents: %w", err)
	}

	for _, word := range list {
		wordStorage.words[word.ID] = word
	}

	return nil
}
