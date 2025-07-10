# Airline Voucher Seat Assignment

A full-stack web application for generating airline crew vouchers with seat assignments.

## Project Structure

```
airline-voucher-seat-assignment/
├── backend/          # Go backend API
├── frontend/         # React frontend
└── README.md
```

## Prerequisites

Before running the application, ensure you have the following installed:

- **Go** (version 1.19 or higher) - [Download here](https://golang.org/dl/)
- **Node.js** (version 16 or higher) - [Download here](https://nodejs.org/)
- **npm** (comes with Node.js)

## Backend Setup & Running

### 1. Navigate to Backend Directory
```bash
cd backend
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Environment Configuration
The backend uses an `app.env` file for configuration. Make sure it contains:
```env
ENVIRONMENT=development
SERVER_PORT=8080
DB_PATH=vouchers.db
LOG_LEVEL=info
```

### 4. Initialize Database
The application will automatically create the SQLite database and tables on first run.

### 5. Build and Run Backend
```bash
# Build the application
go build -o main .

# Run the application
./main
```

Or run directly without building:
```bash
go run main.go
```

The backend API will be available at `http://localhost:8080`

### Backend API Endpoints

- `POST /api/check` - Check if vouchers exist for a flight
- `POST /api/generate` - Generate new vouchers with seat assignments

## Frontend Setup & Running

### 1. Navigate to Frontend Directory
```bash
cd frontend
```

### 2. Install Dependencies
```bash
npm install
```

### 3. Start Development Server
```bash
npm start
```

The frontend will be available at `http://localhost:3000`

## Running the Complete Application

### Step-by-Step Instructions

1. **Start the Backend** (in one terminal):
   ```bash
   cd backend
   go run main.go
   ```
   
2. **Start the Frontend** (in another terminal):
   ```bash
   cd frontend
   npm start
   ```

3. **Access the Application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## Features

- **Voucher Generation**: Generate crew vouchers with random seat assignments
- **Duplicate Prevention**: Check if vouchers already exist for a flight/date
- **Aircraft Support**: Supports ATR, Airbus 320, and Boeing 737 Max
- **Responsive UI**: Modern, mobile-friendly interface built with React and Tailwind CSS

## Usage

1. Fill in the crew details:
   - Crew Name
   - Crew ID
   - Flight Number
   - Flight Date
   - Aircraft Type

2. Click "Generate Vouchers" to create seat assignments

3. The system will:
   - Check if vouchers already exist for the flight/date
   - Generate 3 random seat assignments if none exist
   - Display the assigned seats

## Aircraft Seat Configurations

- **ATR**: 18 rows, 4 seats per row (A, C, D, F)
- **Airbus 320**: 32 rows, 6 seats per row (A, B, C, D, E, F)
- **Boeing 737 Max**: 32 rows, 6 seats per row (A, B, C, D, E, F)

## Development

### Backend Development
- The backend uses Go with a clean architecture pattern
- SQLite database for simplicity
- Environment-based configuration

### Frontend Development
- React with functional components and hooks
- Tailwind CSS for styling
- Axios for API communication

## API Examples

### Check Voucher Exists
```bash
curl -X POST http://localhost:8080/api/check \
  -H "Content-Type: application/json" \
  -d '{"flightNumber":"AI101","date":"2024-01-15"}'
```

### Generate Voucher
```bash
curl -X POST http://localhost:8080/api/generate \
  -H "Content-Type: application/json" \
  -d '{
    "name":"John Doe",
    "id":"CR001",
    "flightNumber":"AI101",
    "date":"2024-01-15",
    "aircraft":"ATR"
  }'
```