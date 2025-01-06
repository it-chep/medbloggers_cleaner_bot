package main

import (
	"context"
	"log"
	"medbloggers_cleaner_bot/internal"
)

func main() {
	ctx := context.Background()
	log.Fatal(internal.NewApp(ctx).Run(ctx))
}
