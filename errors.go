package sonar

import "fmt"

type EnginePathNotFoundError struct{}

func NewEnginePathNotFoundError() *EnginePathNotFoundError {
	return &EnginePathNotFoundError{}
}

func (e *EnginePathNotFoundError) Error() string {
	return "SteelSeries Engine 3 not installed or not in the default location!"
}

type ServerNotAccessibleError struct {
	StatusCode int
}

func NewServerNotAccessibleError(statusCode int) *ServerNotAccessibleError {
	return &ServerNotAccessibleError{StatusCode: statusCode}
}

func (e *ServerNotAccessibleError) Error() string {
	return fmt.Sprintf("SteelSeries server not accessible! Status code: %d", e.StatusCode)
}

type SonarNotEnabledError struct{}

func NewSonarNotEnabledError() *SonarNotEnabledError {
	return &SonarNotEnabledError{}
}

func (e *SonarNotEnabledError) Error() string {
	return "SteelSeries Sonar is not enabled!"
}

type ServerNotReadyError struct{}

func NewServerNotReadyError() *ServerNotReadyError {
	return &ServerNotReadyError{}
}

func (e *ServerNotReadyError) Error() string {
	return "SteelSeries Sonar is not ready yet!"
}

type ServerNotRunningError struct{}

func NewServerNotRunningError() *ServerNotRunningError {
	return &ServerNotRunningError{}
}

func (e *ServerNotRunningError) Error() string {
	return "SteelSeries Sonar is not running!"
}

type WebServerAddressNotFoundError struct{}

func NewWebServerAddressNotFoundError() *WebServerAddressNotFoundError {
	return &WebServerAddressNotFoundError{}
}

func (e *WebServerAddressNotFoundError) Error() string {
	return "Web server address not found"
}

type ChannelNotFoundError struct {
	Channel string
}

func NewChannelNotFoundError(channel string) *ChannelNotFoundError {
	return &ChannelNotFoundError{Channel: channel}
}

func (e *ChannelNotFoundError) Error() string {
	return fmt.Sprintf("Channel '%s' not found", e.Channel)
}

type SliderNotFoundError struct {
	Slider string
}

func NewSliderNotFoundError(slider string) *SliderNotFoundError {
	return &SliderNotFoundError{Slider: slider}
}

func (e *SliderNotFoundError) Error() string {
	return fmt.Sprintf("Slider '%s' not found", e.Slider)
}

type InvalidVolumeError struct {
	Volume float64
}

func NewInvalidVolumeError(volume float64) *InvalidVolumeError {
	return &InvalidVolumeError{Volume: volume}
}

func (e *InvalidVolumeError) Error() string {
	return fmt.Sprintf("Invalid volume '%f'! Value must be between 0 and 1!", e.Volume)
}

type InvalidMixVolumeError struct {
	MixVolume float64
}

func NewInvalidMixVolumeError(mixVolume float64) *InvalidMixVolumeError {
	return &InvalidMixVolumeError{MixVolume: mixVolume}
}

func (e *InvalidMixVolumeError) Error() string {
	return fmt.Sprintf("Invalid mix volume '%f'! Value must be between -1 and 1!", e.MixVolume)
}
