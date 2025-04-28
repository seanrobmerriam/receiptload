# Receipt Uploader Application

A web application that allows users to upload their cash receipts from various department stores including Target, Home Depot, Best Buy, CVS, Rite Aid, Safeway, and others.

## Project Structure

```
receipt-uploader/
├── main.go                  # Main Go application
├── static/                  # Static frontend files
│   ├── index.html           # Main HTML file
│   ├── css/                 # CSS files (optional)
│   └── js/                  # JavaScript files (optional)
├── uploads/                 # Directory for uploaded receipt files
├── templates/               # HTML templates (optional)
├── go.mod                   # Go module file
├── go.sum                   # Go dependencies lock file
└── README.md                # Project documentation
```

## Features

- Upload receipt images or PDFs
- Select store from popular retailers
- Add purchase date and amount
- Add optional notes
- View uploaded receipts
- Simple, responsive user interface

## Installation & Setup

### Prerequisites

- Go 1.17 or higher
- Web browser

### Steps

1. Clone the repository
   ```
   git clone https://github.com/seanrobmerriam/receipt-uploader.git
   cd receipt-uploader
   ```

2. Install dependencies
   ```
   go mod tidy
   ```

3. Create needed directories
   ```
   mkdir -p static uploads
   ```

4. Copy the HTML code into `static/index.html`

5. Run the application
   ```
   go run main.go
   ```

6. Open your browser and navigate to `http://localhost:8080`

## Development

### Frontend

The frontend is built with:
- HTML5
- CSS
- JavaScript (vanilla)

### Backend

The backend is written in Go and uses:
- Gorilla Mux for routing
- Standard library for file handling
- UUID generation for unique identifiers

### Future Enhancements

- User authentication
- Receipt OCR processing
- Data extraction from receipts
- Categorization of purchases
- Export functionality
- Mobile application integration

## License

GNU General Public License v3.0

## Support

For support or questions, please contact seanrobmerriam@gmail.com(mailto:
