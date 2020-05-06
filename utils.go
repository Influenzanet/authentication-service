package main

import (
	"context"
	"time"
)

func getContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(conf.DB.Timeout)*time.Second)
}
