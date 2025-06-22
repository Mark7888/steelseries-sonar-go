package sonar

import (
	"fmt"
	"testing"
)

func TestClassicMode(t *testing.T) {
	sonar, err := New(nil, nil)
	if err != nil {
		if _, ok := err.(*EnginePathNotFoundError); ok {
			fmt.Println("Engine not found!")
			return
		}
		if serverErr, ok := err.(*ServerNotAccessibleError); ok {
			fmt.Printf("Server not accessible, status code: %d\n", serverErr.StatusCode)
			return
		}
		t.Fatalf("Failed to initialize Sonar: %v", err)
	}

	// Disable streamer mode
	streamerDisabled, err := sonar.SetStreamerMode(false)
	if err != nil {
		if serverErr, ok := err.(*ServerNotAccessibleError); ok {
			fmt.Printf("Server not accessible, status code: %d\n", serverErr.StatusCode)
			return
		}
		t.Fatalf("Failed to set streamer mode: %v", err)
	}
	fmt.Printf("Disabling streamer mode: %t\n", !streamerDisabled)

	// Get volume data
	volumeData, err := sonar.GetVolumeData()
	if err != nil {
		t.Fatalf("Failed to get volume data: %v", err)
	}
	fmt.Printf("Classic Mode - Volume Data: %+v\n", volumeData)

	// Set volume
	channel := "master"
	volume := 0.5
	result, err := sonar.SetVolume(channel, volume)
	if err != nil {
		t.Fatalf("Failed to set volume: %v", err)
	}
	fmt.Printf("Classic Mode - Set volume for %s: %+v\n", channel, result)

	// Mute channel
	channel = "game"
	muted := true
	result, err = sonar.MuteChannel(channel, muted)
	if err != nil {
		t.Fatalf("Failed to mute channel: %v", err)
	}
	fmt.Printf("Classic Mode - Mute %s: %+v\n", channel, result)
}

func TestStreamerMode(t *testing.T) {
	sonar, err := New(nil, nil)
	if err != nil {
		if _, ok := err.(*EnginePathNotFoundError); ok {
			fmt.Println("Engine not found!")
			return
		}
		if serverErr, ok := err.(*ServerNotAccessibleError); ok {
			fmt.Printf("Server not accessible, status code: %d\n", serverErr.StatusCode)
			return
		}
		t.Fatalf("Failed to initialize Sonar: %v", err)
	}

	// Enable streamer mode
	streamerEnabled, err := sonar.SetStreamerMode(true)
	if err != nil {
		if serverErr, ok := err.(*ServerNotAccessibleError); ok {
			fmt.Printf("Server not accessible, status code: %d\n", serverErr.StatusCode)
			return
		}
		t.Fatalf("Failed to set streamer mode: %v", err)
	}
	fmt.Printf("Enabling streamer mode: %t\n", streamerEnabled)

	sliders := []string{"streaming", "monitoring"}
	for _, slider := range sliders {
		// Get volume data
		volumeData, err := sonar.GetVolumeData()
		if err != nil {
			t.Fatalf("Failed to get volume data: %v", err)
		}
		fmt.Printf("Streamer Mode (%s) - Volume Data: %+v\n", slider, volumeData)

		// Set volume
		channel := "master"
		volume := 0.5
		result, err := sonar.SetVolume(channel, volume, slider)
		if err != nil {
			t.Fatalf("Failed to set volume: %v", err)
		}
		fmt.Printf("Streamer Mode (%s) - Set volume for %s: %+v\n", slider, channel, result)

		// Mute channel
		channel = "game"
		muted := true
		result, err = sonar.MuteChannel(channel, muted, slider)
		if err != nil {
			t.Fatalf("Failed to mute channel: %v", err)
		}
		fmt.Printf("Streamer Mode (%s) - Mute %s: %+v\n", slider, channel, result)
	}
}

func TestChatMix(t *testing.T) {
	sonar, err := New(nil, nil)
	if err != nil {
		if _, ok := err.(*EnginePathNotFoundError); ok {
			fmt.Println("Engine not found!")
			return
		}
		if serverErr, ok := err.(*ServerNotAccessibleError); ok {
			fmt.Printf("Server not accessible, status code: %d\n", serverErr.StatusCode)
			return
		}
		t.Fatalf("Failed to initialize Sonar: %v", err)
	}

	// Get chat mix data
	chatMixData, err := sonar.GetChatMixData()
	if err != nil {
		t.Fatalf("Failed to get chat mix data: %v", err)
	}
	fmt.Printf("Chat Mix Data: %+v\n", chatMixData)

	// Set chat mix
	mixVolume := 0.5
	result, err := sonar.SetChatMix(mixVolume)
	if err != nil {
		t.Fatalf("Failed to set chat mix: %v", err)
	}
	fmt.Printf("Set Chat Mix: %+v\n", result)
}
