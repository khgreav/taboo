package main

import "github.com/google/uuid"

func generateUUID() string {
	return uuid.NewString()
}
