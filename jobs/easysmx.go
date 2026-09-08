package jobs

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	// 2. Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("Error parsing HTML: %v", err)
		return
	}

	// 3. Look for "Sold Out" on the page
	// We'll search the text content of the entire document (or specific buttons if known).
	// A simple check is just searching the text of the body.
	bodyText := doc.Find("body").Text()
	
	// Many Shopify sites use uppercase, capitalize, etc. So we check case-insensitively.
	isSoldOut := strings.Contains(strings.ToLower(bodyText), "sold out")

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
