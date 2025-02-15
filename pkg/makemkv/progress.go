package makemkv

// SubtaskProgress returns the progress of the current sub-task as a
// percentage.
func (l *Progress) SubtaskProgress() float64 {
	return float64(l.SubtaskValue) / float64(l.Max)
}

// TaskProgress returns the progress of the overall task as a percentage.
func (l *Progress) TaskProgress() float64 {
	return float64(l.TaskValue) / float64(l.Max)
}
