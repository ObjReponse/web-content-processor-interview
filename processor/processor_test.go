package processor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebProcessor_extractLinks(t *testing.T) {
	processor := NewWebProcessor(Config{})

	html := `
		<html>
			<body>
				<a href="https://example.com">Link 1</a>
				<a href="/relative">Relative Link</a>
				<a href="https://test.com/video.m3u8">M3U8 Link</a>
			</body>
		</html>
	`

	links := processor.extractLinks(html)

	if len(links) < 2 {
		t.Errorf("Expected at least 2 links, got %d", len(links))
	}
}

func TestWebProcessor_findM3U8Links(t *testing.T) {
	processor := NewWebProcessor(Config{})

	links := []string{
		"https://example.com",
		"https://test.com/playlist.m3u8",
		"https://stream.com/video.m3u8",
		"https://site.com/page.html",
	}

	m3u8Links := processor.findM3U8Links(links)

	if len(m3u8Links) != 2 {
		t.Errorf("Expected 2 m3u8 links, got %d", len(m3u8Links))
	}
}

func TestWebProcessor_hasKeywords(t *testing.T) {
	config := Config{
		Keywords: []string{"streaming", "video", "download"},
	}
	processor := NewWebProcessor(config)

	tests := []struct {
		title    string
		body     string
		expected bool
	}{
		{"Video streaming site", "Watch movies online", true},
		{"Download center", "Get your files here", true},
		{"Regular website", "Welcome to our site", false},
	}

	for _, test := range tests {
		result := processor.hasKeywords(test.title, test.body)
		if result != test.expected {
			t.Errorf("hasKeywords(%q, %q) = %v, expected %v",
				test.title, test.body, result, test.expected)
		}
	}
}

func TestWebProcessor_shouldContinueProcessing(t *testing.T) {
	processor := NewWebProcessor(Config{})

	tests := []struct {
		hasKeywords bool
		m3u8Count   int
		expected    bool
	}{
		{true, 0, true},   // есть ключевые слова
		{false, 1, true},  // есть m3u8 ссылки
		{true, 2, true},   // есть и то, и другое
		{false, 0, false}, // ничего нет
	}

	for _, test := range tests {
		result := processor.shouldContinueProcessing(test.hasKeywords, test.m3u8Count)
		if result != test.expected {
			t.Errorf("shouldContinueProcessing(%v, %d) = %v, expected %v",
				test.hasKeywords, test.m3u8Count, result, test.expected)
		}
	}
}

func TestWebProcessor_fetchContent(t *testing.T) {
	// Создаем тестовый HTTP сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Test content</body></html>"))
	}))
	defer server.Close()

	config := Config{
		RequestTimeout: 5 * time.Second,
		UserAgent:      "Test Agent",
	}
	processor := NewWebProcessor(config)

	ctx := context.Background()
	statusCode, content, err := processor.fetchContent(ctx, server.URL)

	if err != nil {
		t.Errorf("fetchContent() error = %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("fetchContent() statusCode = %d, expected %d", statusCode, http.StatusOK)
	}

	if len(content) == 0 {
		t.Error("fetchContent() returned empty content")
	}
}

func TestWebProcessor_Process_Integration(t *testing.T) {
	// Создаем тестовый HTTP сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
			<html>
				<head><title>Streaming Video Site</title></head>
				<body>
					<h1>Watch streaming videos</h1>
					<a href="https://cdn.example.com/playlist.m3u8">Stream 1</a>
					<a href="https://example.com/page">Other link</a>
				</body>
			</html>
		`
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(html))
	}))
	defer server.Close()

	config := Config{
		Keywords:       []string{"streaming", "video"},
		MaxWorkers:     2,
		RequestTimeout: 5 * time.Second,
		UserAgent:      "Test Agent",
	}

	processor := NewWebProcessor(config)

	// Создаем каналы
	input := make(chan Entity, 1)
	output := make(chan ProcessResult, 1)

	// Запускаем процессор
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		err := processor.Process(ctx, input, output)
		if err != nil {
			t.Errorf("Process() error = %v", err)
		}
	}()

	// Отправляем тестовую Entity
	input <- Entity{
		Link:  server.URL,
		Title: "Test Page",
	}
	close(input)

	// Получаем результат
	select {
	case result := <-output:
		if result.Error != nil {
			t.Errorf("ProcessResult error = %v", result.Error)
		}

		if result.StatusCode != http.StatusOK {
			t.Errorf("ProcessResult StatusCode = %d, expected %d", result.StatusCode, http.StatusOK)
		}

		if !result.HasKeywords {
			t.Error("ProcessResult should have keywords = true")
		}

		if len(result.M3U8Links) == 0 {
			t.Error("ProcessResult should have found m3u8 links")
		}

		if !result.ShouldContinue {
			t.Error("ProcessResult should continue processing = true")
		}

	case <-time.After(5 * time.Second):
		t.Error("Test timed out waiting for result")
	}
}
