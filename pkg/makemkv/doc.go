// Package makemkv provides a Go interface to makemkvcon, the command-line
// tool bundled with MakeMKV.
//
// The central type is [Con], which wraps makemkvcon and exposes three
// operations: listing drives ([Con.ListDrives]), scanning a disc
// ([Con.ScanDrive]), and backing up a title ([Con.BackupTitle]). All three
// stream progress lines back to the caller via Go's iter.Seq2 so the UI can
// update in real time without waiting for the command to finish.
//
// Disc metadata (titles, streams, attributes) is modeled by [Disc], [Title],
// and [Stream]. Attribute IDs and their meanings are defined in the defs
// subpackage.
package makemkv
