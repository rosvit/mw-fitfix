# :bicyclist: MyWhoosh FIT Fix

**This is a personal project and is in no way affiliated with or endorsed by MyWhoosh.**

[MyWhoosh](https://www.mywhoosh.com/) is an excellent (and free!) virtual cycling platform.
Currently, the FIT files produced by MyWhoosh contain Garmin Edge 1030 Plus as the recording device,
and they contain two laps (50 min and 55 min long) for activities in the free ride mode, regardless of the actual length
of the activity.

This tool allows you to change one of these properties (or both) so that your activity can be uploaded and viewed
correctly on Strava and other platforms.

## Features

- CLI version
- WebAssembly version
- Replace `Garmin Edge 1030 Plus` with `MyWhoosh` as FIT file creator
- Fix laps to contain only a single lap for the entire activity
- Fix timestamps in the FIT file

## Try it!

Live WebAssembly version is available [HERE](https://rosvit.com/mwfitfix).

## CLI version

To run CLI version, clone this repository and run:

```
go run cmd/cli/main.go <FIT_FILE>
```

To see all available options, use `--help` flag:

```
go run cmd/cli/main.go --help
```

## WebAssembly version

Build WASM binary:

```
GOOS=js GOARCH=wasm go build -o web/main.wasm ./cmd/wasm/...
```

Then start development HTTP server:

```
go run cmd/devserver/main.go
```

To access the tool, visit `http://localhost:8080` in your browser.
