// Code generated by "stringer -type=Attr"; DO NOT EDIT.

package defs

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[Type-1]
	_ = x[Name-2]
	_ = x[LangCode-3]
	_ = x[LangName-4]
	_ = x[CodecID-5]
	_ = x[CodecShort-6]
	_ = x[CodecLong-7]
	_ = x[ChapterCount-8]
	_ = x[Duration-9]
	_ = x[DiskSize-10]
	_ = x[DiscSizeBytes-11]
	_ = x[StreamTypeExtension-12]
	_ = x[Bitrate-13]
	_ = x[AudioChannelsCount-14]
	_ = x[AngleInfo-15]
	_ = x[SourceFileName-16]
	_ = x[AudioSampleRate-17]
	_ = x[AudioSampleSize-18]
	_ = x[VideoSize-19]
	_ = x[VideoAspectRatio-20]
	_ = x[VideoFrameRate-21]
	_ = x[StreamFlags-22]
	_ = x[DateTime-23]
	_ = x[OriginalTitleID-24]
	_ = x[SegmentsCount-25]
	_ = x[SegmentsMap-26]
	_ = x[OutputFileName-27]
	_ = x[MetadataLanguageCode-28]
	_ = x[MetadataLanguageName-29]
	_ = x[TreeInfo-30]
	_ = x[PanelTitle-31]
	_ = x[VolumeName-32]
	_ = x[OrderWeight-33]
	_ = x[OutputFormat-34]
	_ = x[OutputFormatDescription-35]
	_ = x[SeamlessInfo-36]
	_ = x[PanelText-37]
	_ = x[MkvFlags-38]
	_ = x[MkvFlagsText-39]
	_ = x[AudioChannelLayoutName-40]
	_ = x[OutputCodecShort-41]
	_ = x[OutputConversionType-42]
	_ = x[OutputAudioSampleRate-43]
	_ = x[OutputAudioSampleSize-44]
	_ = x[OutputAudioChannelsCount-45]
	_ = x[OutputAudioChannelLayoutName-46]
	_ = x[OutputAudioChannelLayout-47]
	_ = x[OutputAudioMixDescription-48]
	_ = x[Comment-49]
	_ = x[OffsetSequenceID-50]
}

const _Attr_name = "UnknownTypeNameLangCodeLangNameCodecIDCodecShortCodecLongChapterCountDurationDiskSizeDiscSizeBytesStreamTypeExtensionBitrateAudioChannelsCountAngleInfoSourceFileNameAudioSampleRateAudioSampleSizeVideoSizeVideoAspectRatioVideoFrameRateStreamFlagsDateTimeOriginalTitleIDSegmentsCountSegmentsMapOutputFileNameMetadataLanguageCodeMetadataLanguageNameTreeInfoPanelTitleVolumeNameOrderWeightOutputFormatOutputFormatDescriptionSeamlessInfoPanelTextMkvFlagsMkvFlagsTextAudioChannelLayoutNameOutputCodecShortOutputConversionTypeOutputAudioSampleRateOutputAudioSampleSizeOutputAudioChannelsCountOutputAudioChannelLayoutNameOutputAudioChannelLayoutOutputAudioMixDescriptionCommentOffsetSequenceID"

var _Attr_index = [...]uint16{0, 7, 11, 15, 23, 31, 38, 48, 57, 69, 77, 85, 98, 117, 124, 142, 151, 165, 180, 195, 204, 220, 234, 245, 253, 268, 281, 292, 306, 326, 346, 354, 364, 374, 385, 397, 420, 432, 441, 449, 461, 483, 499, 519, 540, 561, 585, 613, 637, 662, 669, 685}

func (i Attr) String() string {
	if i < 0 || i >= Attr(len(_Attr_index)-1) {
		return "Attr(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Attr_name[_Attr_index[i]:_Attr_index[i+1]]
}
