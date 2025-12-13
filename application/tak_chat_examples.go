package main

import (
	"bytes"
	"fmt"

	"darbelis.eu/stabas/tak"
	"github.com/google/uuid"
)

// Example usage
func main() {
	// Configuration
	const (
		takServerHost  = "localhost"
		httpsPort      = 8443
		senderUID      = "GOLANG-12345" // Use uuid.New().String() for unique ID
		senderCallsign = "GoClient"
		messageText    = "Hello from Go TAK client!"
		chatroom       = "All Chat Rooms"
	)

	// Generate unique sender UID
	uniqueUID := fmt.Sprintf("GOLANG-%s", uuid.New().String()[:8])

	// Create chat CoT message
	fmt.Println("Creating chat CoT message...")
	cotXML, err := tak.CreateChatCoTXML(uniqueUID, senderCallsign, messageText, chatroom)
	if err != nil {
		fmt.Printf("Error creating CoT XML: %v\n", err)
		return
	}

	fmt.Println("\nGenerated CoT XML:")
	fmt.Println(cotXML)
	fmt.Println("\n" + string(bytes.Repeat([]byte("="), 80)))

	// Example 1: Send via REST API with certificates
	fmt.Println("\nExample 1: Send via REST API with client certificates")
	fmt.Println("Usage:")
	fmt.Printf(`
err := tak.SendChatViaRESTAPI(
	"%s",
	%d,
	cotXML,
	"/path/to/giedrius.pem",
	"/path/to/giedrius-key.pem",
	"/path/to/ca.pem",
)
if err != nil {
	fmt.Printf("Error: %%v\n", err)
}
`, takServerHost, httpsPort)

	// Example 2: Send via REST API (insecure, for testing)
	fmt.Println("\nExample 2: Send via REST API (skip TLS verification - testing only)")
	fmt.Println("Usage:")
	fmt.Printf(`
err := tak.SendChatViaRESTAPIInsecure(
	"%s",
	%d,
	cotXML,
	"/path/to/giedrius.pem",
	"/path/to/giedrius-key.pem",
)
if err != nil {
	fmt.Printf("Error: %%v\n", err)
}
`, takServerHost, httpsPort)

	// Uncomment to actually send a message (requires valid certificate paths)
	/*
		err = tak.SendChatViaRESTAPIInsecure(
			takServerHost,
			httpsPort,
			cotXML,
			"/opt/tak/certs/files/giedrius.pem",
			"/opt/tak/certs/files/giedrius-key.pem",
		)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
		} else {
			fmt.Println("✓ Message sent successfully!")
		}
	*/

	fmt.Println("\n" + string(bytes.Repeat([]byte("="), 80)))
	fmt.Println("\nTo use this code:")
	fmt.Println("1. Install dependencies:")
	fmt.Println("   go get github.com/google/uuid")
	fmt.Println("\n2. Update certificate paths in the code")
	fmt.Println("\n3. Run:")
	fmt.Println("   go run tak_chat_examples.go")
	fmt.Println("\n4. Or build:")
	fmt.Println("   go build -o tak_chat tak_chat_examples.go")
	fmt.Println("   ./tak_chat")
}
