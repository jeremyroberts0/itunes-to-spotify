# iTunes to Spotify Playlist Migration Tool

## Project Overview

This is a REST API service written in Go that enables users to migrate playlists from iTunes to Spotify. The application provides a simple HTTP API that takes an iTunes-exported playlist file and creates an equivalent playlist in a user's Spotify account using the Spotify Web API.

### Key Features

- **iTunes Playlist Parsing**: Parses tab-separated iTunes playlist export files
- **Spotify Integration**: Authenticates with Spotify using OAuth2 and creates playlists
- **Song Matching**: Searches for iTunes songs on Spotify using track name and artist
- **Batch Processing**: Uses a worker pool for concurrent song searches with rate limiting
- **Statistics**: Provides playlist analytics (most common artists, albums, songs)
- **Dockerized**: Fully containerized with Docker and Docker Compose support

### Main Use Case

1. User exports playlist from iTunes as tab-separated file
2. User authorizes the service to access their Spotify account
3. User uploads the iTunes playlist file via REST API
4. Service matches songs and creates new playlist in Spotify

## Architecture

### High-Level Components

- **HTTP API Server**: Built with Gin web framework
- **Authentication Layer**: Spotify OAuth2 integration with cookie-based session management
- **Playlist Parser**: CSV/TSV parser for iTunes playlist exports
- **Spotify Client**: Integration with Spotify Web API for searching and playlist creation
- **Worker Pool**: Concurrent processing for song matching operations

### Request Flow

```
Client Request → Gin Router → Auth Middleware → Playlist Parser → 
Song Matching Pool → Spotify API → Playlist Creation → Response
```

### Key Architectural Decisions

- Uses worker pool pattern for concurrent Spotify API calls
- Implements rate limiting and retry logic for Spotify API
- Cookie-based authentication to maintain session state
- Processes songs in batches of 100 (Spotify API limit)

## Directory Structure

```
├── api/                    # HTTP API handlers and routing
│   ├── auth.go            # Spotify OAuth2 authentication
│   ├── routes.go          # Main router setup
│   ├── spotify.go         # Core playlist migration logic
│   └── stats.go           # Playlist analytics endpoints
├── itunes/                # iTunes playlist parsing
│   └── parsePlaylist.go   # TSV file parser for iTunes exports
├── types/                 # Shared data structures
│   └── types.go           # Song data type definitions
├── main.go               # Application entry point
├── Dockerfile            # Multi-stage Docker build
├── docker-compose.yml    # Local development setup
└── README.md            # Basic usage instructions
```

## Key Files

### Core Application Files

- **`main.go`**: Entry point that starts the Gin server on port 8081
- **`api/routes.go`**: Main router configuration that registers all endpoints
- **`api/spotify.go`**: Core business logic for playlist migration and Spotify integration
- **`api/auth.go`**: OAuth2 flow implementation for Spotify authentication

### Data Processing

- **`itunes/parsePlaylist.go`**: Sophisticated TSV parser that handles iTunes playlist export format
- **`types/types.go`**: Simple data structure for iTunes song representation

### Infrastructure

- **`Dockerfile`**: Multi-stage build using Go 1.9.2 base image (legacy version)
- **`docker-compose.yml`**: Simple service definition with environment file support

## API Endpoints

### Authentication
- `GET /authorize` - Redirect to Spotify OAuth authorization
- `GET /authorized` - OAuth callback handler, sets authentication cookie

### Core Functionality
- `POST /itunes-to-spotify?name=<playlist-name>` - Main playlist migration endpoint
- `POST /stats` - Generate playlist statistics from iTunes file

## Coding Conventions

### Go Style
- Standard Go naming conventions (PascalCase for exported, camelCase for internal)
- Package-level organization with clear separation of concerns
- Error handling follows Go idioms with explicit error returns
- Uses struct embedding for response types

### Code Organization
- Handlers are organized by feature area (auth, spotify, stats)
- Pure functions for data processing (playlist parsing)
- Dependency injection pattern for HTTP handlers
- Configuration through environment variables

### Key Patterns
- **Worker Pool**: Concurrent processing with channels for communication
- **Builder Pattern**: For constructing complex API responses
- **Middleware**: Cookie-based authentication checking
- **Error Response**: Consistent error response format across endpoints

## Dependencies

### Core Libraries
- **Gin** (`github.com/gin-gonic/gin`): HTTP web framework
- **Spotify API** (`github.com/zmb3/spotify`): Official Spotify Web API client
- **OAuth2** (`golang.org/x/oauth2`): OAuth2 client implementation
- **Worker Pool** (`github.com/jeremyroberts0/pool`): Custom concurrency management

### Development Tools
- **Docker**: Containerization (Go 1.9.2 base image)
- **Docker Compose**: Local development orchestration

### Architecture Notes
- This is a **legacy Go project** that predates Go modules (uses GOPATH-style imports)
- Uses older Dockerfile patterns and Go 1.9.2
- No formal dependency management (relies on `go get`)

## Development Workflow

### Prerequisites
- Go 1.9+ (project targets 1.9.2)
- Docker and Docker Compose
- Spotify Developer Account (for client credentials)

### Environment Setup
Create `.env` file with:
```
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
```

### Local Development
```bash
# Using Docker (recommended)
docker-compose up

# Or native Go (requires GOPATH setup)
go get ./...
go build .
./itunes-to-spotify
```

### Building
- **Docker**: `docker build .` (multi-stage build)
- **Native**: `go build .` (outputs binary based on directory name)

### Testing
- No automated tests present in the codebase
- Manual testing through REST API endpoints
- Docker container health can be verified through endpoint access

## Common Tasks

### Adding New Endpoints
1. Add handler function to appropriate file in `api/` package
2. Register route in `api/routes.go` `GetRouter()` function
3. Follow existing error response patterns using `createError()` helper

### Extending Playlist Parsing
1. Modify `itunes/parsePlaylist.go` to handle new iTunes export formats
2. Update `types/types.go` if new song metadata is needed
3. Adjust column parsing logic in `getColPositions()` function

### Spotify API Integration
1. Extend `api/spotify.go` with new Spotify Web API calls
2. Use existing `client` instance with `AutoRetry` enabled
3. Follow batch processing patterns for bulk operations
4. Implement proper rate limiting using `time.Sleep()`

### Authentication Changes
1. Modify OAuth scopes in `api/auth.go` `init()` function
2. Update redirect URI and cookie management as needed
3. Ensure cookie security settings match deployment environment

## Important Context & Gotchas

### Legacy Go Project
- **No Go Modules**: Uses old GOPATH-style dependency management
- **Dockerfile Version**: Targets Go 1.9.2 (very old)
- **Modern Go**: Would need significant updates to work with current Go versions and modules

### Spotify API Limitations
- **Rate Limiting**: Service implements retry logic and delays between requests
- **Batch Size**: Playlists are added in chunks of 100 tracks (Spotify limit)
- **Search Accuracy**: Song matching relies on simple text search and takes first result

### Docker Considerations
- **Multi-stage Build**: Final image is Alpine-based for size optimization
- **Port Exposure**: Hardcoded to port 8081
- **Static Binary**: Uses `CGO_ENABLED=0` for portable binary

### Production Considerations
- **No HTTPS**: Service runs HTTP only (suitable for local/internal use)
- **Cookie Security**: Authentication cookies are not secure-flagged
- **Error Handling**: Limited error recovery for partial playlist failures
- **Logging**: Uses basic `fmt.Printf` rather than structured logging

### Development Notes
- **Worker Pool**: Custom implementation for concurrent processing
- **Memory Usage**: Loads entire playlist into memory during processing
- **No Persistence**: No database - all state is ephemeral
- **Local Development**: Hardcoded localhost redirect URI

This service is designed for local/personal use and would need security and scalability improvements for production deployment.