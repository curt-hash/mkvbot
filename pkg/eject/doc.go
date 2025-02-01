/*
Package eject provides a cross-platform Eject function that attempts to eject a
disc given the volume name.

The implementation currently relies on platform-specific command execution:
`eject` on Linux, `drutil` on Darwin, and `powershell` on Windows.

For example:

	err := Eject(context.Background(), "D:")
*/
package eject
