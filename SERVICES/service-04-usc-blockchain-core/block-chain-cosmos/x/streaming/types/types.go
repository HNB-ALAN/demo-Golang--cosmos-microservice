package types

import (
	"fmt"
	"time"
)

const (
	// ModuleName defines the module name
	ModuleName = "streaming"

	// RouterKey defines the message route for the streaming module
	RouterKey = ModuleName

	// QuerierRoute defines the querier route for the streaming module
	QuerierRoute = ModuleName
)

// Event types
const (
	EventTypeStreamCreated    = "stream_created"
	EventTypeStreamUpdated    = "stream_updated"
	EventTypeStreamStarted    = "stream_started"
	EventTypeStreamStopped    = "stream_stopped"
	EventTypeStreamPaused     = "stream_paused"
	EventTypeStreamResumed    = "stream_resumed"
	EventTypeStreamQuality    = "stream_quality"
	EventTypeStreamViewer     = "stream_viewer"
	EventTypeStreamChat       = "stream_chat"
	EventTypeStreamDonation   = "stream_donation"
	EventTypeStreamModeration = "stream_moderation"

	AttributeKeyStreamID   = "stream_id"
	AttributeKeyStreamerID = "streamer_id"
	AttributeKeyViewerID   = "viewer_id"
	AttributeKeyMessageID  = "message_id"
	AttributeKeyQuality    = "quality"
	AttributeKeyModule     = ModuleName
)

// StreamStatus represents the status of a stream
type StreamStatus string

const (
	StreamStatusActive   StreamStatus = "active"
	StreamStatusPaused   StreamStatus = "paused"
	StreamStatusStopped  StreamStatus = "stopped"
	StreamStatusInactive StreamStatus = "inactive"
	StreamStatusError    StreamStatus = "error"
)

// Stream represents a streaming session
type Stream struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	CreatorID   string            `json:"creator_id"`
	StreamerID  string            `json:"streamer_id"`
	Status      StreamStatus      `json:"status"`
	Category    string            `json:"category"`
	Quality     string            `json:"quality"`
	StartTime   time.Time         `json:"start_time"`
	StartedAt   time.Time         `json:"started_at"`
	EndTime     *time.Time        `json:"end_time,omitempty"`
	EndedAt     *time.Time        `json:"ended_at,omitempty"`
	UpdatedAt   time.Time         `json:"updated_at"`
	ViewerCount int               `json:"viewer_count"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata"`
}

// StreamViewer represents a viewer in a stream
type StreamViewer struct {
	ID        string            `json:"id"`
	StreamID  string            `json:"stream_id"`
	UserID    string            `json:"user_id"`
	ViewerID  string            `json:"viewer_id"`
	Username  string            `json:"username"`
	JoinedAt  time.Time         `json:"joined_at"`
	LeftAt    *time.Time        `json:"left_at,omitempty"`
	WatchTime int64             `json:"watch_time"` // seconds
	Metadata  map[string]string `json:"metadata"`
}

// StreamChat represents a chat message in a stream
type StreamChat struct {
	ID        string            `json:"id"`
	StreamID  string            `json:"stream_id"`
	UserID    string            `json:"user_id"`
	Username  string            `json:"username"`
	Message   string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
	Type      string            `json:"type"` // message, donation, moderation
	Metadata  map[string]string `json:"metadata"`
}

// StreamDonation represents a donation in a stream
type StreamDonation struct {
	ID        string            `json:"id"`
	StreamID  string            `json:"stream_id"`
	DonorID   string            `json:"donor_id"`
	Amount    int64             `json:"amount"` // USC tokens
	Message   string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
}

// StreamQualityMetrics represents quality metrics for a stream
type StreamQualityMetrics struct {
	ID           string            `json:"id"`
	StreamID     string            `json:"stream_id"`
	Timestamp    time.Time         `json:"timestamp"`
	Bitrate      int64             `json:"bitrate"`
	FPS          int64             `json:"fps"`
	Resolution   string            `json:"resolution"`
	BufferHealth int64             `json:"buffer_health"` // 0 to 100 (percentage)
	Latency      int64             `json:"latency"`       // milliseconds
	PacketLoss   int64             `json:"packet_loss"`   // 0 to 100 (percentage)
	Jitter       int64             `json:"jitter"`        // milliseconds
	QualityScore int64             `json:"quality_score"` // 0 to 100 (percentage)
	Metadata     map[string]string `json:"metadata"`
}

// StreamAnalytics represents analytics for a stream
type StreamAnalytics struct {
	ID             string            `json:"id"`
	StreamID       string            `json:"stream_id"`
	Period         string            `json:"period"` // daily, weekly, monthly
	StartTime      time.Time         `json:"start_time"`
	EndTime        time.Time         `json:"end_time"`
	TotalViewers   int               `json:"total_viewers"`
	PeakViewers    int               `json:"peak_viewers"`
	AvgViewers     int               `json:"avg_viewers"`
	TotalWatchTime int64             `json:"total_watch_time"` // seconds
	Engagement     int64             `json:"engagement"`       // 0 to 100 (percentage)
	Retention      int64             `json:"retention"`        // 0 to 100 (percentage)
	ChatMessages   int               `json:"chat_messages"`
	Donations      int               `json:"donations"`
	TotalDonations int64             `json:"total_donations"` // USC tokens
	Metadata       map[string]string `json:"metadata"`
}

// StreamModeration represents moderation actions in a stream
type StreamModeration struct {
	ID          string            `json:"id"`
	StreamID    string            `json:"stream_id"`
	ModeratorID string            `json:"moderator_id"`
	Action      string            `json:"action"` // timeout, ban, mute, unmute
	TargetID    string            `json:"target_id"`
	Reason      string            `json:"reason"`
	Duration    int64             `json:"duration"` // seconds
	Timestamp   time.Time         `json:"timestamp"`
	Metadata    map[string]string `json:"metadata"`
}

// StreamEvent represents a streaming event
type StreamEvent struct {
	ID        string            `json:"id"`
	StreamID  string            `json:"stream_id"`
	Type      string            `json:"type"` // created, updated, started, stopped, viewer_joined, viewer_left, quality_changed
	Data      map[string]string `json:"data"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
}

// StreamQuality represents stream quality settings
type StreamQuality struct {
	ID           string            `json:"id"`
	StreamID     string            `json:"stream_id"`
	Resolution   string            `json:"resolution"`    // 720p, 1080p, 4K
	Bitrate      int64             `json:"bitrate"`       // bits per second
	FPS          int64             `json:"fps"`           // frames per second
	Codec        string            `json:"codec"`         // h264, h265, vp9
	AudioCodec   string            `json:"audio_codec"`   // aac, mp3, opus
	AudioBitrate int64             `json:"audio_bitrate"` // bits per second
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Metadata     map[string]string `json:"metadata"`
}

// GenesisState defines the streaming module's genesis state
type GenesisState struct {
	Streams        []Stream               `json:"streams"`
	Viewers        []StreamViewer         `json:"viewers"`
	Chats          []StreamChat           `json:"chats"`
	Donations      []StreamDonation       `json:"donations"`
	QualityMetrics []StreamQualityMetrics `json:"quality_metrics"`
	Analytics      []StreamAnalytics      `json:"analytics"`
	Moderations    []StreamModeration     `json:"moderations"`
	Events         []StreamEvent          `json:"events"`
	Qualities      []StreamQuality        `json:"qualities"`
	Params         Params                 `json:"params"`
}

// Params defines the streaming module parameters
type Params struct {
	MaxStreamsPerUser    int64 `json:"max_streams_per_user"`
	MaxViewersPerStream  int64 `json:"max_viewers_per_stream"`
	MaxChatMessageLength int64 `json:"max_chat_message_length"`
	StreamTimeout        int64 `json:"stream_timeout"`         // seconds
	QualityCheckInterval int64 `json:"quality_check_interval"` // seconds
}

// DefaultParams returns default streaming module parameters
func DefaultParams() Params {
	return Params{
		MaxStreamsPerUser:    5,
		MaxViewersPerStream:  10000,
		MaxChatMessageLength: 500,
		StreamTimeout:        3600, // 1 hour
		QualityCheckInterval: 30,   // 30 seconds
	}
}

// Validate validates the streaming module parameters
func (p Params) Validate() error {
	if p.MaxStreamsPerUser <= 0 {
		return fmt.Errorf("max_streams_per_user must be positive")
	}
	if p.MaxViewersPerStream <= 0 {
		return fmt.Errorf("max_viewers_per_stream must be positive")
	}
	if p.MaxChatMessageLength <= 0 {
		return fmt.Errorf("max_chat_message_length must be positive")
	}
	if p.StreamTimeout <= 0 {
		return fmt.Errorf("stream_timeout must be positive")
	}
	if p.QualityCheckInterval <= 0 {
		return fmt.Errorf("quality_check_interval must be positive")
	}
	return nil
}

// Validate validates a stream
func (s Stream) Validate() error {
	if s.ID == "" {
		return fmt.Errorf("stream ID cannot be empty")
	}
	if s.Title == "" {
		return fmt.Errorf("stream title cannot be empty")
	}
	if s.CreatorID == "" {
		return fmt.Errorf("creator ID cannot be empty")
	}
	return nil
}

// Validate validates a stream viewer
func (sv StreamViewer) Validate() error {
	if sv.ID == "" {
		return fmt.Errorf("viewer ID cannot be empty")
	}
	if sv.StreamID == "" {
		return fmt.Errorf("stream ID cannot be empty")
	}
	if sv.UserID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	return nil
}

// Validate validates a stream chat message
func (sc StreamChat) Validate() error {
	if sc.ID == "" {
		return fmt.Errorf("chat message ID cannot be empty")
	}
	if sc.StreamID == "" {
		return fmt.Errorf("stream ID cannot be empty")
	}
	if sc.UserID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if sc.Message == "" {
		return fmt.Errorf("message cannot be empty")
	}
	return nil
}

// Validate validates a stream donation
func (sd StreamDonation) Validate() error {
	if sd.ID == "" {
		return fmt.Errorf("donation ID cannot be empty")
	}
	if sd.StreamID == "" {
		return fmt.Errorf("stream ID cannot be empty")
	}
	if sd.DonorID == "" {
		return fmt.Errorf("donor ID cannot be empty")
	}
	if sd.Amount <= 0 {
		return fmt.Errorf("donation amount must be positive")
	}
	return nil
}

// Validate validates stream quality metrics
func (sqm StreamQualityMetrics) Validate() error {
	if sqm.ID == "" {
		return fmt.Errorf("quality metrics ID cannot be empty")
	}
	if sqm.StreamID == "" {
		return fmt.Errorf("stream ID cannot be empty")
	}
	if sqm.BufferHealth < 0 || sqm.BufferHealth > 100 {
		return fmt.Errorf("buffer health must be between 0 and 100")
	}
	if sqm.PacketLoss < 0 || sqm.PacketLoss > 100 {
		return fmt.Errorf("packet loss must be between 0 and 100")
	}
	if sqm.QualityScore < 0 || sqm.QualityScore > 100 {
		return fmt.Errorf("quality score must be between 0 and 100")
	}
	return nil
}

// Validate validates stream analytics
func (sa StreamAnalytics) Validate() error {
	if sa.ID == "" {
		return fmt.Errorf("analytics ID cannot be empty")
	}
	if sa.StreamID == "" {
		return fmt.Errorf("stream ID cannot be empty")
	}
	if sa.Engagement < 0 || sa.Engagement > 100 {
		return fmt.Errorf("engagement must be between 0 and 100")
	}
	if sa.Retention < 0 || sa.Retention > 100 {
		return fmt.Errorf("retention must be between 0 and 100")
	}
	return nil
}

// Validate validates stream moderation
func (sm StreamModeration) Validate() error {
	if sm.ID == "" {
		return fmt.Errorf("moderation ID cannot be empty")
	}
	if sm.StreamID == "" {
		return fmt.Errorf("stream ID cannot be empty")
	}
	if sm.ModeratorID == "" {
		return fmt.Errorf("moderator ID cannot be empty")
	}
	if sm.TargetID == "" {
		return fmt.Errorf("target ID cannot be empty")
	}
	return nil
}

// Validate validates stream event
func (se StreamEvent) Validate() error {
	if se.ID == "" {
		return fmt.Errorf("event ID cannot be empty")
	}
	if se.StreamID == "" {
		return fmt.Errorf("stream ID cannot be empty")
	}
	if se.Type == "" {
		return fmt.Errorf("event type cannot be empty")
	}
	return nil
}

// Validate validates stream quality
func (sq StreamQuality) Validate() error {
	if sq.ID == "" {
		return fmt.Errorf("quality ID cannot be empty")
	}
	if sq.StreamID == "" {
		return fmt.Errorf("stream ID cannot be empty")
	}
	if sq.Resolution == "" {
		return fmt.Errorf("resolution cannot be empty")
	}
	if sq.Bitrate <= 0 {
		return fmt.Errorf("bitrate must be positive")
	}
	if sq.FPS <= 0 {
		return fmt.Errorf("fps must be positive")
	}
	return nil
}
