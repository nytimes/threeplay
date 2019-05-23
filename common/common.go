package common

const ThreePlayHost = "api.3playmedia.com"
const ThreePlayStaticHost = "static.3playmedia.com"

// CaptionsFormat is supported output format for captions
type CaptionsFormat string

const (
	// SRT format for captions file
	SRT CaptionsFormat = "srt"
	// WebVTT format for captions file
	WebVTT CaptionsFormat = "vtt"
	// DFX format for captions file
	DFX CaptionsFormat = "pdfxp"
	// SMI format for captions file
	SMI CaptionsFormat = "smi"
	// STL format for captions file
	STL CaptionsFormat = "stl"
	// QT format for captions file
	QT CaptionsFormat = "qt"
	// QTXML format for captions file
	QTXML CaptionsFormat = "qtxml"
	// CPTXML format for captions file
	CPTXML CaptionsFormat = "cptxml"
	// ADBE format for captions file
	ADBE CaptionsFormat = "adbe"
)
