# Spotify Sync

## Overview

The goal of this repository is to keep in sync two spotify playlists

## Usage

### Prequisities

- Create app on [spotify developer website](https://developer.spotify.com/)
- Add `http://localhost:8000/callback` to `Redirect URIs`
- Set `CLIENT_ID` and `CLIENT_SECRET` variables
- Obtain `REFRESH_TOKEN` variable by running following command: `go run cmd/get-token/get-token.go`. Refesh token value should be printed to console:
```
Refresh token: <VALUE>
```

### Running locally

- Set [environment variables](#environment)
- Run following command: `go run cmd/sync/sync.go`

### Environment
| Name | Description |
| ------------- | ------------- |
| CLIENT_ID  | Client ID from your spotify app, you can find it in app settings |
| CLIENT_SECRET | Client secret from your spotify app, you can find it in app settings |
| DST_URI | Destination playlist URI in following format `https://api.spotify.com/v1/playlists/<playlist ID>/tracks`. If you want to sync liked songs use `https://api.spotify.com/v1/me/tracks` |
| REFRESH_TOKEN | Refresh token, you can obtain it by running `go run cmd/get-token/get-token.go` |
| SRC_URI | Source playlist URI in following format `https://api.spotify.com/v1/playlists/<playlist ID>/tracks`. If you want to sync liked songs use `https://api.spotify.com/v1/me/tracks` |
