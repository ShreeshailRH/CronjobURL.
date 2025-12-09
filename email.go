package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"time"
)

// URLs to monitor
var urls = []string{
	"https://github.com/ShreeshailRH/CronjobURL./blob/main/email.go",
	"https://www.tutorialspoint.com/http/http_status_codes.htm",
	"https://www.guvi.in/sqlkata/sql/2/2232"
}

// Threshold (seconds)
const maxResponseSeconds = 2

// Gmail SMTP
const (
	smtpHost      = "smtp.gmail.com"
	smtpPort      = "587"
	senderEmail   = "shreeshailrh92@gmail.com"
	senderPass    = "xhme wdxw qucg bkmd" // app password
	receiverEmail = "shreeshail@krazybee.com"
)

func main() {
	for _, link := range urls {
		checkURL(link)
	}
}

// Check a single URL
func checkURL(url string) {
	start := time.Now()

	resp, err := http.Get(url)
	elapsed := time.Since(start).Seconds()

	// If request completely failed
	if err != nil {
		saveToCSV(url, "DOWN", 0)
		sendAlert(fmt.Sprintf("❌ URL DOWN: %s\nError: %v", url, err))
		return
	}

	defer resp.Body.Close()

	// Non-200 response
	if resp.StatusCode != 200 {
		saveToCSV(url, "BAD_STATUS", elapsed)
		sendAlert(fmt.Sprintf("❌ BAD STATUS for %s\nStatus Code: %d", url, resp.StatusCode))
		return
	}

	// Slow response
	if elapsed > maxResponseSeconds {
		saveToCSV(url, "SLOW", elapsed)
		sendAlert(fmt.Sprintf("⚠️ SLOW RESPONSE for %s\nTime: %.2f seconds", url, elapsed))
		return
	}

	// Everything OK
	saveToCSV(url, "OK", elapsed)
	sendAlert(fmt.Sprintf("✅ URL Working: %s\nResponse Time: %.2f sec", url, elapsed))
}

// Send Email Alert
func sendAlert(message string) {
	auth := smtp.PlainAuth("", senderEmail, senderPass, smtpHost)

	subject := "Subject: Monitoring Alert ⚠️\n"
	msg := []byte(subject + "\n" + message)

	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		senderEmail,
		[]string{receiverEmail},
		msg,
	)

	if err != nil {
		fmt.Println("SMTP ERROR:", err)
		return
	}

	fmt.Println("EMAIL SENT:", message)
}

// Save results to CSV file
func saveToCSV(url, status string, timeTaken float64) {
	file, err := os.OpenFile("monitoring.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("CSV ERROR:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	row := []string{
		time.Now().Format("2006-01-02 15:04:05"),
		url,
		status,
		fmt.Sprintf("%.2f", timeTaken),
	}

	writer.Write(row)
}
