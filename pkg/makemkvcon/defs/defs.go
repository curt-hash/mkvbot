package defs

// Originally from apdefs.h in the makemkv for Linux source tarball.

//go:generate stringer -type=Attr
type Attr int

const (
	Unknown Attr = iota
	Type
	Name
	LangCode
	LangName
	CodecID // 5
	CodecShort
	CodecLong
	ChapterCount
	Duration
	DiskSize // 10
	DiscSizeBytes
	StreamTypeExtension
	Bitrate
	AudioChannelsCount
	AngleInfo // 15
	SourceFileName
	AudioSampleRate
	AudioSampleSize
	VideoSize
	VideoAspectRatio // 20
	VideoFrameRate
	StreamFlags
	DateTime
	OriginalTitleID
	SegmentsCount // 25
	SegmentsMap
	OutputFileName
	MetadataLanguageCode
	MetadataLanguageName
	TreeInfo // 30
	PanelTitle
	VolumeName
	OrderWeight
	OutputFormat
	OutputFormatDescription // 35
	SeamlessInfo
	PanelText
	MkvFlags
	MkvFlagsText
	AudioChannelLayoutName // 40
	OutputCodecShort
	OutputConversionType
	OutputAudioSampleRate
	OutputAudioSampleSize
	OutputAudioChannelsCount // 45
	OutputAudioChannelLayoutName
	OutputAudioChannelLayout
	OutputAudioMixDescription
	Comment
	OffsetSequenceID // 50
)

//go:generate stringer -type=TypeCode -trimprefix=TypeCode
type TypeCode int

const (
	TypeCodeTitle TypeCode = iota + 6200
	TypeCodeVideo
	TypeCodeAudio
	TypeCodeSubtitles
)
