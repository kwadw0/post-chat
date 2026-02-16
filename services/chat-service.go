package services

import "fmt"

func ChatService(msg string) string {
	fmt.Printf("Chat Service received: %s\n", msg)
	return "Processed: " + msg
}
