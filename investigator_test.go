package Investigator

import (
	"context"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	fmt.Println(Investigate(context.Background(), ModelGraphState{
		KimiModel: &Model{},
		Topic:     "gemini3 和gork4的比拼，谁的能力更强？",
	}))

}
