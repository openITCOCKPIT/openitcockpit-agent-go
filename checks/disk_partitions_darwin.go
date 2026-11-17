package checks

import "golang.org/x/sys/unix"

var mountFlagNames = []mountFlagName{
	{unix.MNT_SYNCHRONOUS, "sync"},
	{unix.MNT_NOEXEC, "noexec"},
	{unix.MNT_NOSUID, "nosuid"},
	{unix.MNT_UNION, "union"},
	{unix.MNT_ASYNC, "async"},
	{unix.MNT_DONTBROWSE, "nobrowse"},
	{unix.MNT_AUTOMOUNTED, "automounted"},
	{unix.MNT_JOURNALED, "journaled"},
	{unix.MNT_MULTILABEL, "multilabel"},
	{unix.MNT_NOATIME, "noatime"},
	{unix.MNT_NODEV, "nodev"},
	{unix.MNT_LOCAL, "local"},
	{unix.MNT_CPROTECT, "protect"},
}
