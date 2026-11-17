package checks

import "golang.org/x/sys/unix"

var mountFlagNames = []mountFlagName{
	{unix.MNT_SYNCHRONOUS, "sync"},
	{unix.MNT_NOEXEC, "noexec"},
	{unix.MNT_NOSUID, "nosuid"},
	{unix.MNT_UNION, "union"},
	{unix.MNT_ASYNC, "async"},
	{unix.MNT_SUIDDIR, "suiddir"},
	{unix.MNT_SOFTDEP, "softdep"},
	{unix.MNT_NOSYMFOLLOW, "nosymfollow"},
	{unix.MNT_GJOURNAL, "gjournal"},
	{unix.MNT_MULTILABEL, "multilabel"},
	{unix.MNT_ACLS, "acls"},
	{unix.MNT_NOATIME, "noatime"},
	{unix.MNT_NOCLUSTERR, "noclusterr"},
	{unix.MNT_NOCLUSTERW, "noclusterw"},
	{unix.MNT_NFS4ACLS, "nfsv4acls"},
}

// notRespondingMounts is not required on FreeBSD, statfs of a not responding file system blocks
func notRespondingMounts() map[string]bool {
	return nil
}
