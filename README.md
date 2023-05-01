# WatchParty
Join a theater and watch your favourite movies or videos together, and clap along.

## Dependencies
This project requires `go` and `npm` to build, and `ffmpeg` during runtime.

## Build
First, build the Vue based client with npm.
```
    $ cd client
    $ npm run install
    $ npm run build
    $ cd ..
```

Finally, you can start the Go server.
```
    $ cd server
    $ go run .
```

Or build a binary for distribution.

Use the `-help` command line option on the server for configuring.

## Screenshot
![WatchParty](https://raw.githubusercontent.com/BenJilks/WatchParty/master/docs/screenshot.png)

