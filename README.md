# SteelSeries Sonar Go API

**This is the Go implementation of the Python package [steelseries-sonar-py](https://github.com/Mark7888/steelseries-sonar-py)**

## Overview

This Go package provides a convenient interface for interacting with the SteelSeries Sonar application API.    
The Sonar application allows users to control and display volumes for various audio channels.

## Installation

To use this package, follow these steps:

1. Install the package using go get:

   ```bash
   go get github.com/mark7888/steelseries-sonar-go
   ```

2. Import the package in your Go application:

   ```go
   import "github.com/mark7888/steelseries-sonar-go"
   ```

## Usage

### Initializing the Sonar Object

The `New` function accepts two optional parameters:   
`appDataPath`: Specify a custom path for the SteelSeries Engine 3 coreProps.json file   
(default is the default installation path: `C:\ProgramData\SteelSeries\SteelSeries Engine 3\coreProps.json`).   
`streamerMode`: Set to true to use streamer mode (default is auto-detected).

```go
// Default initialization
sonar, err := sonar.New(nil, nil)
if err != nil {
    log.Fatal(err)
}
```

```go
// Custom app data path
appDataPath := "C:\\path\\to\\coreProps.json"
sonar, err := sonar.New(&appDataPath, nil)
if err != nil {
    log.Fatal(err)
}
```

```go
// Force streamer mode
streamerMode := true
sonar, err := sonar.New(nil, &streamerMode)
if err != nil {
    log.Fatal(err)
}
```

### Streamer Mode

The SteelSeries Sonar Go API supports streamer mode, which allows users to manage two separate sliders: `streaming` and `monitoring`. These sliders enable fine-tuned control over different audio channels.

To check if the streamer mode is enabled, use:

```go
isStreaming, err := sonar.IsStreamerMode()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Is Streamer Mode: %t\n", isStreaming)
```

To enable or disable streamer mode, use:

```go
// Enable streamer mode
enabled, err := sonar.SetStreamerMode(true)
if err != nil {
    log.Fatal(err)
}

// Disable streamer mode
disabled, err := sonar.SetStreamerMode(false)
if err != nil {
    log.Fatal(err)
}
```

### Retrieving Volume Information

Retrieve information about the current volume settings for all channels:

```go
volumeData, err := sonar.GetVolumeData()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Volume Data: %+v\n", volumeData)
```

### Setting Volume for a Channel

Set the volume for a specific channel. The `channel` parameter should be one of the following:   
`master`, `game`, `chatRender`, `media`, `aux`, `chatCapture`. The `volume` parameter should be a float64 between 0 and 1.   
Additionally, an optional `streamerSlider` parameter can be provided, with values "streaming" (default) or "monitoring":

```go
channel := "master"
volume := 0.75
streamerSlider := "streaming" // or "monitoring"

result, err := sonar.SetVolume(channel, volume, streamerSlider)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Set Volume Result: %+v\n", result)
```

### Muting/Unmuting a Channel

Toggle mute status for a specific channel. The `channel` parameter should be one of the following:   
`master`, `game`, `chatRender`, `media`, `aux`, `chatCapture`. The `muted` parameter should be a boolean indicating whether to mute (`true`) or unmute (`false`) the channel.   
Additionally, an optional `streamerSlider` parameter can be provided, with values "streaming" (default) or "monitoring":

```go
channel := "game"
muted := true
streamerSlider := "monitoring"

result, err := sonar.MuteChannel(channel, muted, streamerSlider)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Mute Result: %+v\n", result)
```

### Chatmix

Retrieve chat-mix data:

```go
chatmixData, err := sonar.GetChatMixData()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Chatmix Data: %+v\n", chatmixData)
```

Set chat-mix value between `-1 and 1` to focus sound from the `game` or `chatRender` channel:

```go
result, err := sonar.SetChatMix(0.5)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Set Chatmix Result: %+v\n", result)
```

## Error Handling

The package provides custom error types that correspond to the Python package exceptions.    
It is advisable to handle these errors accordingly in your code. Here is the list of potential errors:

- `EnginePathNotFoundError`: Raised when SteelSeries Engine 3 is not installed or not in the default location.
- `ServerNotAccessibleError`: Raised when the SteelSeries server is not accessible. Provides the HTTP status code.
- `SonarNotEnabledError`: Raised when SteelSeries Sonar is not enabled.
- `ServerNotReadyError`: Raised when SteelSeries Sonar is not ready.
- `ServerNotRunningError`: Raised when SteelSeries Sonar is not running.
- `WebServerAddressNotFoundError`: Raised when the web server address is not found.
- `ChannelNotFoundError`: Raised when the specified channel is not found.
- `InvalidVolumeError`: Raised when an invalid volume value is provided.
- `InvalidMixVolumeError`: Raised when an invalid mix volume value is provided.
- `SliderNotFoundError`: Raised when an unknown slider name is provided as `streamerSlider` value.

Example error handling:

```go
sonar, err := sonar.New(nil, nil)
if err != nil {
    switch e := err.(type) {
    case *sonar.EnginePathNotFoundError:
        fmt.Println("Engine not found!")
        return
    case *sonar.ServerNotAccessibleError:
        fmt.Printf("Server not accessible, status code: %d\n", e.StatusCode)
        return
    default:
        log.Fatal(err)
    }
}
```

## Example

Here is a complete example demonstrating the usage of the SteelSeries Sonar Go API:

```go
package main

import (
    "fmt"
    "log"
    "github.com/mark7888/steelseries-sonar-go"
)

func main() {
    // Initialize Sonar object
    customPath := "C:\\path\\to\\coreProps.json"
    sonar, err := sonar.New(&customPath, nil)
    if err != nil {
        if _, ok := err.(*sonar.EnginePathNotFoundError); ok {
            fmt.Println("Engine not found!")
            return
        }
        log.Fatal(err)
    }

    // Retrieve volume data
    volumeData, err := sonar.GetVolumeData()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Volume Data: %+v\n", volumeData)

    // Set volume for the 'master' channel
    channel := "master"
    volume := 0.8
    streamerSlider := "streaming"
    result, err := sonar.SetVolume(channel, volume, streamerSlider)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Set volume for %s: %+v\n", channel, result)

    // Mute the 'game' channel
    channel = "game"
    muted := true
    streamerSlider = "monitoring"
    result, err = sonar.MuteChannel(channel, muted, streamerSlider)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Mute %s: %+v\n", channel, result)

    // Retrieve chat-mix data
    chatmixData, err := sonar.GetChatMixData()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Chatmix Data: %+v\n", chatmixData)

    // Set chat-mix value
    result, err = sonar.SetChatMix(0.5)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Set Chatmix: %+v\n", result)
}
```

## Testing

Run the tests using:

```bash
go test -v
```

## Special Thanks

This Go package is based on the Python implementation [steelseries-sonar-py](https://github.com/Mark7888/steelseries-sonar-py).

Thanks to all contributors who made the original package possible - [wex](https://github.com/wex/sonar-rev) for figuring out the API, [TotalPanther317](https://github.com/TotalPanther317/steelseries-sonar-py) for understanding streamer mode and [cookie](https://github.com/cookie0o) for features like chat mix and streamer mode detection. Grateful for their efforts!

This documentation reflects the Go implementation of the SteelSeries Sonar API with the same functionality as the Python version.
