package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"daily_jobs/jobs"
	"daily_jobs/scheduler"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error reading it. Relying on system environment variables.")
	}

	// The email address to notify
	notifyEmail := os.Getenv("NOTIFY_EMAIL")
	if notifyEmail == "" {
		log.Fatal("NOTIFY_EMAIL environment variable is required")
	}

	s := scheduler.NewScheduler()

	// Register the EasySMX stock check job.
	// "@daily" runs once a day at midnight.
	// Alternatively, use "0 9 * * *" to run at 9 AM every day.
	schedule := os.Getenv("EASYSMX_SCHEDULE")
	if schedule == "" {
		schedule = "@daily" 
	}

	jobFunc := func() {
		jobs.CheckEasySMXStock(notifyEmail)
	}

	err = s.AddJob(schedule, "Check EasySMX Stock", jobFunc)
	if err != nil {
		log.Fatalf("Error adding job: %v", err)
	}

	// Start the scheduler
	s.Start()

	// Run the job immediately once on startup
	go func() {
		log.Println("Running job immediately on startup...")
		jobFunc()
	}()

	// Block main thread until an interrupt signal is received
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	<-sigChan
	log.Println("Shutting down gracefully...")
	s.Stop()
}
