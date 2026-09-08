package jobs

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"daily_jobs/notifier"
)

const EasySMXURL = "https://www.easysmx.com/products/easysmx-x05pro-charging-dock?srsltid=AfmBOoqhvzEO5c9V02dgbaJJp5ktonGv2w351pRqDIS3gCdgSmQUx8Zh"

// CheckEasySMXStock fetches the EasySMX product page and checks if it's sold out.
// If it is NOT sold out, it sends an email notification.
func CheckEasySMXStock(toAddress string) {
	log.Println("Running job: CheckEasySMXStock")

	// 1. Fetch the HTML
	res, err := http.Get(EasySMXURL)
	if err != nil {
		log.Printf("Error fetching EasySMX page: %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Error: status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	// 2. Read the raw HTML response
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading HTML body: %v", err)
		return
	}
	htmlContent := strings.ToLower(string(bodyBytes))

	// 3. Look for "sold out" or the standard Shopify/Schema.org "outofstock"
	// This is much more reliable since Shopify often renders the visual "Sold Out" text using JavaScript,
	// but the raw HTML always contains the metadata for availability.
	isSoldOut := strings.Contains(htmlContent, "sold out") || strings.Contains(htmlContent, "outofstock")

	if isSoldOut {
		log.Println("EasySMX X05 Pro is still Sold Out.")
	} else {
		log.Println("EasySMX X05 Pro MIGHT BE IN STOCK! Sending email...")
	
		
		subject := "EasySMX X05 Pro In Stock Alert!"
		body := fmt.Sprintf("The 'Sold Out' text was not found on the page.\n\nCheck the link: %s", EasySMXURL)

		err := notifier.SendEmail(subject, body, toAddress)
		if err != nil {
			log.Printf("Failed to send email alert: %v", err)
		} else {
			log.Println("Alert email sent successfully.")
		}
	}
}
