package processor

import "time"

// Entity представляет входные данные для обработки
type Entity struct {
	Link  string
	Title string
}

// ProcessResult содержит результат обработки одной Entity
type ProcessResult struct {
	Entity         Entity
	StatusCode     int
	HTMLContent    string
	Links          []string
	M3U8Links      []string
	HasKeywords    bool
	ShouldContinue bool
	ProcessTime    time.Duration
	Error          error
}

// Config содержит настройки для процессора
type Config struct {
	Keywords       []string
	MaxWorkers     int
	RequestTimeout time.Duration
	MaxRedirects   int
	UserAgent      string
}

// ProcessorStats содержит статистику работы
type ProcessorStats struct {
	TotalProcessed int
	SuccessCount   int
	ErrorCount     int
	AvgProcessTime time.Duration
}
