package main

import (
	"context"
	"docstar_cleaner_bot/internal"
	"log"
)

func main() {
	ctx := context.Background()
	log.Fatal(internal.NewApp(ctx).Run(ctx))
}
