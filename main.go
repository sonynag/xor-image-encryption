package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// XOR encryption
func xorEncrypt(data []byte, key byte) []byte {
	for i := 0; i < len(data); i++ {
		data[i] ^= key
	}
	return data
}

// Serve HTML page
func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

// Handle file upload
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Increase max upload size to 50 MB
	r.ParseMultipartForm(50 << 20)

	// Get file
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "File upload error: "+err.Error(), http.StatusBadRequest)
		fmt.Println("File upload error:", err)
		return
	}
	defer file.Close()

	// Read file
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Error reading file:", err)
		return
	}

	// Get key
	keyInt, err := strconv.Atoi(r.FormValue("key"))
	if err != nil || keyInt < 0 || keyInt > 255 {
		http.Error(w, "Invalid key: "+err.Error(), http.StatusBadRequest)
		fmt.Println("Invalid key:", err)
		return
	}
	key := byte(keyInt)

	// Apply XOR
	data = xorEncrypt(data, key)

	// Ensure folder exists
	os.MkdirAll("encrypted", os.ModePerm)

	// Safe file name
	safeName := strings.ReplaceAll(handler.Filename, " ", "_")
	outputFile := "encrypted/encrypted_" + safeName

	// Save file
	err = os.WriteFile(outputFile, data, 0644)
	if err != nil {
		http.Error(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Fprintf(w, "✅ Image encrypted/decrypted successfully!<br>Saved as: %s", outputFile)
	fmt.Println("File saved:", outputFile)
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Server running at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server failed:", err)
	}
}
