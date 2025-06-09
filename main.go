package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ObjReponse/web-content-processor-interview/processor"
)

func main() {
	// Конфигурация процессора
	config := processor.Config{
		Keywords:       []string{"streaming", "video", "download", "watch", "movie"},
		MaxWorkers:     3,
		RequestTimeout: 10 * time.Second,
		UserAgent:      "WebContentProcessor/1.0",
	}

	// Создаем процессор
	webProcessor := processor.NewWebProcessor(config)

	// Создаем каналы
	input := make(chan processor.Entity, 10)
	output := make(chan processor.ProcessResult, 10)

	// Контекст для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запускаем процессор в отдельной горутине
	go func() {
		if err := webProcessor.Process(ctx, input, output); err != nil {
			log.Printf("Processor error: %v", err)
		}
	}()

	// Горутина для обработки результатов
	go func() {
		for result := range output {
			printResult(result)
		}
	}()

	// Отправляем тестовые данные
	testEntities := []processor.Entity{
		{Link: "https://httpbin.org/html", Title: "Test HTML Page"},
		{Link: "https://httpbin.org/json", Title: "Test JSON Response"},
		{Link: "https://example.com", Title: "Example Domain"},
		{Link: "https://httpstat.us/404", Title: "Not Found Page"},
	}

	fmt.Println("🚀 Starting web content processor...")
	fmt.Println("📝 Processing test entities...")

	for _, entity := range testEntities {
		select {
		case input <- entity:
			fmt.Printf("📤 Queued: %s\n", entity.Link)
		case <-ctx.Done():
			break
		}
	}

	// Ждем сигнал завершения или таймаут
	select {
	case <-sigChan:
		fmt.Println("\n🛑 Received shutdown signal...")
	case <-time.After(30 * time.Second):
		fmt.Println("\n⏰ Timeout reached...")
	}

	// Graceful shutdown
	fmt.Println("🔄 Shutting down gracefully...")
	close(input)
	cancel()

	// Даем время на завершение
	time.Sleep(2 * time.Second)

	// Выводим статистику
	stats := webProcessor.GetStats()
	fmt.Printf("\n📊 Final Statistics:\n")
	fmt.Printf("   Total Processed: %d\n", stats.TotalProcessed)
	fmt.Printf("   Success Count: %d\n", stats.SuccessCount)
	fmt.Printf("   Error Count: %d\n", stats.ErrorCount)
	fmt.Printf("   Average Process Time: %v\n", stats.AvgProcessTime)

	fmt.Println("✅ Shutdown complete")
}

func printResult(result processor.ProcessResult) {
	fmt.Printf("\n🔍 Processing Result for: %s\n", result.Entity.Link)
	fmt.Printf("   Status Code: %d\n", result.StatusCode)
	fmt.Printf("   Content Length: %d bytes\n", len(result.HTMLContent))
	fmt.Printf("   Links Found: %d\n", len(result.Links))
	fmt.Printf("   M3U8 Links: %d\n", len(result.M3U8Links))
	fmt.Printf("   Has Keywords: %v\n", result.HasKeywords)
	fmt.Printf("   Should Continue: %v\n", result.ShouldContinue)
	fmt.Printf("   Process Time: %v\n", result.ProcessTime)

	if result.Error != nil {
		fmt.Printf("   ❌ Error: %v\n", result.Error)
	} else {
		fmt.Printf("   ✅ Success\n")
	}

	// Показываем найденные M3U8 ссылки
	if len(result.M3U8Links) > 0 {
		fmt.Printf("   🎥 M3U8 Links:\n")
		for _, link := range result.M3U8Links {
			fmt.Printf("      - %s\n", link)
		}
	}
}
