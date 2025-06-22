package sonar

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Sonar struct {
	ChannelNames        []string
	StreamerSliderNames []string
	VolumePath          string
	AppDataPath         string
	BaseURL             string
	WebServerAddress    string
	StreamerMode        bool
	httpClient          *http.Client
}

type CoreProps struct {
	GGEncryptedAddress string `json:"ggEncryptedAddress"`
}

type SubApps struct {
	SubApps map[string]SubApp `json:"subApps"`
}

type SubApp struct {
	IsEnabled bool                   `json:"isEnabled"`
	IsReady   bool                   `json:"isReady"`
	IsRunning bool                   `json:"isRunning"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func New(appDataPath *string, streamerMode *bool) (*Sonar, error) {
	s := &Sonar{
		ChannelNames:        []string{"master", "game", "chatRender", "media", "aux", "chatCapture"},
		StreamerSliderNames: []string{"streaming", "monitoring"},
		VolumePath:          "/volumeSettings/classic",
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}

	if appDataPath != nil {
		s.AppDataPath = *appDataPath
	} else {
		programData := os.Getenv("ProgramData")
		if programData == "" {
			programData = "C:\\ProgramData"
		}
		s.AppDataPath = filepath.Join(programData, "SteelSeries", "SteelSeries Engine 3", "coreProps.json")
	}

	if err := s.loadBaseURL(); err != nil {
		return nil, err
	}

	if err := s.loadServerAddress(); err != nil {
		return nil, err
	}

	if streamerMode == nil {
		isStreamer, err := s.IsStreamerMode()
		if err != nil {
			return nil, err
		}
		s.StreamerMode = isStreamer
	} else {
		s.StreamerMode = *streamerMode
	}

	if s.StreamerMode {
		s.VolumePath = "/volumeSettings/streamer"
	}

	return s, nil
}

func (s *Sonar) IsStreamerMode() (bool, error) {
	resp, err := s.httpClient.Get(s.WebServerAddress + "/mode/")
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, NewServerNotAccessibleError(resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var mode string
	if err := json.Unmarshal(body, &mode); err != nil {
		return false, err
	}

	return mode == "stream", nil
}

func (s *Sonar) SetStreamerMode(streamerMode bool) (bool, error) {
	var mode string
	if streamerMode {
		mode = "stream"
	} else {
		mode = "classic"
	}

	url := fmt.Sprintf("%s/mode/%s", s.WebServerAddress, mode)
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return false, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return false, NewServerNotAccessibleError(resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var resultMode string
	if err := json.Unmarshal(body, &resultMode); err != nil {
		return false, err
	}

	s.StreamerMode = resultMode == "stream"
	if s.StreamerMode {
		s.VolumePath = "/volumeSettings/streamer"
	} else {
		s.VolumePath = "/volumeSettings/classic"
	}

	return s.StreamerMode, nil
}

func (s *Sonar) loadBaseURL() error {
	if _, err := os.Stat(s.AppDataPath); os.IsNotExist(err) {
		return NewEnginePathNotFoundError()
	}

	file, err := os.Open(s.AppDataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	var coreProps CoreProps
	if err := json.NewDecoder(file).Decode(&coreProps); err != nil {
		return err
	}

	s.BaseURL = fmt.Sprintf("https://%s", coreProps.GGEncryptedAddress)
	return nil
}

func (s *Sonar) loadServerAddress() error {
	resp, err := s.httpClient.Get(s.BaseURL + "/subApps")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return NewServerNotAccessibleError(resp.StatusCode)
	}

	var subApps SubApps
	if err := json.NewDecoder(resp.Body).Decode(&subApps); err != nil {
		return err
	}

	sonarApp, exists := subApps.SubApps["sonar"]
	if !exists {
		return NewSonarNotEnabledError()
	}

	if !sonarApp.IsEnabled {
		return NewSonarNotEnabledError()
	}

	if !sonarApp.IsReady {
		return NewServerNotReadyError()
	}

	if !sonarApp.IsRunning {
		return NewServerNotRunningError()
	}

	webServerAddress, exists := sonarApp.Metadata["webServerAddress"]
	if !exists {
		return NewWebServerAddressNotFoundError()
	}

	s.WebServerAddress = webServerAddress.(string)
	if s.WebServerAddress == "" || s.WebServerAddress == "null" {
		return NewWebServerAddressNotFoundError()
	}

	return nil
}

func (s *Sonar) GetVolumeData() (map[string]interface{}, error) {
	volumeInfoURL := s.WebServerAddress + s.VolumePath

	resp, err := s.httpClient.Get(volumeInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewServerNotAccessibleError(resp.StatusCode)
	}

	var volumeData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&volumeData); err != nil {
		return nil, err
	}

	return volumeData, nil
}

func (s *Sonar) SetVolume(channel string, volume float64, streamerSlider ...string) (map[string]interface{}, error) {
	if !contains(s.ChannelNames, channel) {
		return nil, NewChannelNotFoundError(channel)
	}

	slider := "streaming"
	if len(streamerSlider) > 0 {
		slider = streamerSlider[0]
	}

	if s.StreamerMode && !contains(s.StreamerSliderNames, slider) {
		return nil, NewSliderNotFoundError(slider)
	}

	if volume < 0 || volume > 1 {
		return nil, NewInvalidVolumeError(volume)
	}

	fullVolumePath := s.VolumePath
	if s.StreamerMode {
		fullVolumePath += "/" + slider
	}

	volumeJSON, _ := json.Marshal(volume)
	url := fmt.Sprintf("%s%s/%s/Volume/%s", s.WebServerAddress, fullVolumePath, channel, string(volumeJSON))

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewServerNotAccessibleError(resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Sonar) MuteChannel(channel string, muted bool, streamerSlider ...string) (map[string]interface{}, error) {
	if !contains(s.ChannelNames, channel) {
		return nil, NewChannelNotFoundError(channel)
	}

	slider := "streaming"
	if len(streamerSlider) > 0 {
		slider = streamerSlider[0]
	}

	if s.StreamerMode && !contains(s.StreamerSliderNames, slider) {
		return nil, NewSliderNotFoundError(slider)
	}

	fullVolumePath := s.VolumePath
	if s.StreamerMode {
		fullVolumePath += "/" + slider
	}

	muteKeyword := "Mute"
	if s.StreamerMode {
		muteKeyword = "isMuted"
	}

	mutedJSON, _ := json.Marshal(muted)
	url := fmt.Sprintf("%s%s/%s/%s/%s", s.WebServerAddress, fullVolumePath, channel, muteKeyword, string(mutedJSON))

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewServerNotAccessibleError(resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *Sonar) GetChatMixData() (map[string]interface{}, error) {
	chatMixURL := s.WebServerAddress + "/chatMix"

	resp, err := s.httpClient.Get(chatMixURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewServerNotAccessibleError(resp.StatusCode)
	}

	var chatMixData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&chatMixData); err != nil {
		return nil, err
	}

	return chatMixData, nil
}

func (s *Sonar) SetChatMix(mixVolume float64) (map[string]interface{}, error) {
	if mixVolume < -1 || mixVolume > 1 {
		return nil, NewInvalidMixVolumeError(mixVolume)
	}

	mixVolumeJSON, _ := json.Marshal(mixVolume)
	url := fmt.Sprintf("%s/chatMix?balance=%s", s.WebServerAddress, string(mixVolumeJSON))

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, NewServerNotAccessibleError(resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
