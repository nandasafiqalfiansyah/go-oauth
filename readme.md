# OAuth2 Authentication API

This package provides a simple API for authenticating users using OAuth2 with Google as the provider. It includes endpoints to initiate the authentication process and handle callback responses.

## Features

- **OAuth2 Integration**: Supports Google OAuth2 for user authentication.
- **User Information Retrieval**: Fetches and returns authenticated user information such as email, name, and profile picture.
- **Scalable Design**: Implements the Gin web framework for efficient routing and request handling.

## Requirements

- Go 1.16+
- Google OAuth2 credentials
- `.env` file for storing sensitive credentials

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Create a `.env` file in the root directory with the following content:
   ```plaintext
   CLIENT_ID=<your-google-client-id>
   CLIENT_SECRET=<your-google-client-secret>
   ```

## Usage

### Local Development

1. Run the application locally:
   ```bash
   go run main.go
   ```

2. Open your browser and navigate to `http://localhost:8000`.

### Deployment

For deploying on a cloud provider (e.g., AWS, Google Cloud, etc.), use the provided `Handler` function as the entry point for handling requests.

### Endpoints

#### 1. Root Endpoint (`/`)
Redirects users to Google's OAuth2 authorization page.

#### 2. Callback Endpoint (`/auth/callback`)
Handles OAuth2 callback, retrieves the access token, and fetches user information.

## Example Workflow

1. Visit the root endpoint (`http://localhost:8000`) to start the OAuth2 process.
2. Log in with your Google account.
3. Upon successful authentication, the callback endpoint fetches and displays user information as JSON.

### Example Response
```json
{
  "id": "1234567890",
  "email": "user@example.com",
  "verified_email": true,
  "name": "John Doe",
  "given_name": "John",
  "family_name": "Doe",
  "picture": "https://example.com/photo.jpg",
  "locale": "en"
}
```

## File Structure

```plaintext
├── main.go         # Main application entry point
├── go.mod          # Module dependencies
├── go.sum          # Dependency checksums
├── .env            # Environment variables
└── README.md       # Documentation
```

## Libraries Used

- [Gin](https://github.com/gin-gonic/gin): HTTP web framework.
- [oauth2](https://pkg.go.dev/golang.org/x/oauth2): OAuth2 implementation for Go.
- [godotenv](https://github.com/joho/godotenv): Loads environment variables from a `.env` file.

