package main

import "log/slog"

type config struct {
	addr   string
	logger *slog.Logger
}
