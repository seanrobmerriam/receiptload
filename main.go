// main.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Receipt represents a user's receipt
type Receipt struct {
	ID           string    `json:"id"`
	Store        string    `json:"store"`
	PurchaseDate time.Time `json:"purchaseDate"`
	TotalAmount  float64   `json:"totalAmount"`
	Notes        string    `json:"notes,omitempty"`
	FileName     string    `json:"fileName"`
	UploadDate   time.Time `json:"uploadDate"`
	Status       string    `json:"status"`
}

// Server represents the API server
type Server struct {
	Receipts     map[string]*Receipt
	UploadFolder string
}

// NewServer creates a new server instance
func NewServer() *Server {
	// Create uploads directory if it doesn't exist
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	return &Server{
		Receipts:     make(map[string]*Receipt),
		UploadFolder: uploadDir,
	}
}

// RegisterRoutes sets up the routes for the server
func (s *Server) RegisterRoutes(r *mux.Router) {
	// Static file server for the frontend
	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/receipts", s.handleGetReceipts).Methods("GET")
	api.HandleFunc("/receipts", s.handleUploadReceipt).Methods("POST")
	api.HandleFunc("/receipts/{id}", s.handleGetReceipt).Methods("GET")
	api.HandleFunc("/receipts/{id}", s.handleDeleteReceipt).Methods("DELETE")

	// Serve the main HTML page
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
}

// handleGetReceipts returns all receipts
func (s *Server) handleGetReceipts(w http.ResponseWriter, r *http.Request) {
	receipts := make([]*Receipt, 0, len(s.Receipts))
	for _, receipt := range s.Receipts {
		receipts = append(receipts, receipt)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipts)
}

// handleGetReceipt returns a specific receipt
func (s *Server) handleGetReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	receipt, exists := s.Receipts[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(receipt)
}

// handleDeleteReceipt deletes a receipt
func (s *Server) handleDeleteReceipt(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	receipt, exists := s.Receipts[id]
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	// Delete the file
	filePath := filepath.Join(s.UploadFolder, receipt.FileName)
	if err := os.Remove(filePath); err != nil {
		log.Printf("Failed to delete file: %v", err)
	}

	// Delete the receipt from memory
	delete(s.Receipts, id)

	w.WriteHeader(http.StatusNoContent)
}

// handleUploadReceipt handles the receipt upload
func (s *Server) handleUploadReceipt(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 10MB max memory
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get form values
	store := r.FormValue("store")
	purchaseDateStr := r.FormValue("purchaseDate")
	totalAmountStr := r.FormValue("totalAmount")
	notes := r.FormValue("notes")

	// Validate required fields
	if store == "" || purchaseDateStr == "" || totalAmountStr == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Parse purchase date
	purchaseDate, err := time.Parse("2006-01-02", purchaseDateStr)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	// Parse total amount
	totalAmount, err := strconv.ParseFloat(totalAmountStr, 64)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, fileHeader, err := r.FormFile("receiptFile")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	fileType := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if fileType != ".jpg" && fileType != ".jpeg" && fileType != ".png" && fileType != ".pdf" {
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	uniqueID := uuid.New().String()
	filename := uniqueID + fileType
	filePath := filepath.Join(s.UploadFolder, filename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Create receipt object
	receipt := &Receipt{
		ID:           uniqueID,
		Store:        store,
		PurchaseDate: purchaseDate,
		TotalAmount:  totalAmount,
		Notes:        notes,
		FileName:     filename,
		UploadDate:   time.Now(),
		Status:       "processed", // In a real app, this might be "pending" until OCR is done
	}

	// Save receipt to memory (in production, save to database)
	s.Receipts[uniqueID] = receipt

	// Process receipt (in a real app, this would be done asynchronously)
	go s.processReceipt(receipt)

	// Return receipt data
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(receipt)
}

// processReceipt simulates processing a receipt
// In a real app, this would perform OCR, data extraction, etc.
func (s *Server) processReceipt(receipt *Receipt) {
	// Simulate processing time
	time.Sleep(2 * time.Second)

	// Update receipt status
	receipt.Status = "processed"

	// In a real app, you would extract data from the receipt here
	// This could include line items, vendor information, etc.

	log.Printf("Processed receipt %s from %s", receipt.ID, receipt.Store)
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Create a new server
	server := NewServer()

	// Create router
	router := mux.NewRouter()
	
	// Register routes
	server.RegisterRoutes(router)

	// Add middleware
	handler := corsMiddleware(router)

	// Start the server
	addr := ":8080"
	fmt.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
