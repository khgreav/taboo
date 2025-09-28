package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
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

func (ws *WordStorage) GetShuffledIds() []uint {
	ids := make([]uint, 0, len(ws.words))
	for id := range ws.words {
		ids = append(ids, id)
	}
	rand.Shuffle(len(ids), func(i, j int) {
		ids[i], ids[j] = ids[j], ids[i]
	})
	return ids
}

func (ws *WordStorage) GetWordsByIds(ids []uint) ([]*TabooWord, error) {
	words := make([]*TabooWord, 0, len(ids))
	for _, id := range ids {
		word, ok := ws.words[id]
		if !ok {
			continue
		}
		words = append(words, word)
	}
	return words, nil
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
