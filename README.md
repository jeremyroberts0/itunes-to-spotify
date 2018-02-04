# Itunes To Spotify
This is a small REST service that takes a playlist exported from iTunes and creates it on your spotify account.

## Running
- You can easily build and run the app with docker.  A docker-compose file is present to make this as easy as `docker-compose up`
- The app requires a few environment variables to run, these can be set via a `.env` file at the root of the repo or by passing them in via standard docker means.  The required variables are:
-- `SPOTIFY_CLIENT_ID`: Assigned by Spotify
-- `SPOTIFY_CLIENT_SECRET`: Assigned by Spotify
- When creating the app on Spotify, be sure that the hostname you're running the app on is added as a redirect URI in your app settings, if this isn't configured correctly then Authentication will fail.

## Usage
1. Export the playlist as a tab separated file by selecting the playlist in itunes and selecting File -> Library -> Export Playlist
2. Authorize the service to access your spotify account by visiting the /authorize endpoint (GET).  When complete, a cookie is returned with your auth info.  If you're using Postman you may have to manually copy the cookie for the next step.
3. POST the playlist as plain text to `/itunes-to-spotify?name=my-playlist-name` where `my-playlist-name` is the name you'd the playlist will be in your spotify account.  This step may take a while depending on the size of your playlist (generally 1 second for every 10 songs)