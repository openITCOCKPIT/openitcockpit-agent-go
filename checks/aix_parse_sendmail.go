//go:build aix || linux
// +build aix linux

// Parser for the mail queue of sendmail (mailq), the mail system of AIX. The file is also built on Linux,
// so the unit tests run in the CI.

package checks

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// resultSendmail is the mail queue of sendmail
type resultSendmail struct {
	Messages  int64 `json:"messages"`
	Deferred  int64 `json:"deferred"` // Messages with a failed delivery attempt, e.g. (Deferred: Connection refused by ...)
	SizeBytes int64 `json:"size_bytes"`
	OldestAge int64 `json:"oldest_age"` // Age of the oldest message in seconds, 0 if the queue is empty
}

// Line of a message: queue id (with the status *, X or - at the end), size, queue time and sender
//
//	5AG9xKq1012345*     512 Mon Nov 16 10:42 <root@node1>
var sendmailMessageLine = regexp.MustCompile(`^(\w+)[*X-]?\s+(\d+)\s+(\w{3} \w{3} [ \d]\d \d\d:\d\d)\s`)

// Line with the reason of a failed delivery after the line of the message
//
//	(Deferred: Connection refused by mail.example.com.)
var sendmailReasonLine = regexp.MustCompile(`^\s+\(.+\)\s*$`)

// mailq prints the queue time without year
const sendmailTimeLayout = "Mon Jan _2 15:04"

// parseMailq parses the output of mailq, the times must be printed in the time zone of now
//
//			/var/spool/mqueue (2 requests)
//	-----Q-ID----- --Size-- -----Q-Time----- ------------Sender/Recipient-----------
//	5AG9xKq1012345      512 Mon Nov 16 10:42 <root@node1>
//	                 (Deferred: Connection refused by mail.example.com.)
//	                                         <user@example.com>
//			Total requests: 2
func parseMailq(output string, now time.Time) *resultSendmail {
	result := &resultSendmail{}
	var oldest time.Time
	inMessage := false
	for _, line := range strings.Split(output, "\n") {
		match := sendmailMessageLine.FindStringSubmatch(line)
		if match != nil {
			result.Messages++
			size, _ := strconv.ParseInt(match[2], 10, 64)
			result.SizeBytes += size
			oldest = olderTime(oldest, parseSendmailTime(match[3], now))
			inMessage = true
			continue
		}
		if inMessage && sendmailReasonLine.MatchString(line) {
			result.Deferred++
			// Only the first reason of a message is counted
			inMessage = false
		}
	}
	if !oldest.IsZero() && oldest.Before(now) {
		result.OldestAge = int64(now.Sub(oldest).Seconds())
	}
	return result
}

// parseSendmailTime returns the queue time in the year of now. A time after now is from the previous year
// (e.g. a message of December 31 at January 1).
func parseSendmailTime(value string, now time.Time) time.Time {
	parsed, err := time.ParseInLocation(sendmailTimeLayout, value, now.Location())
	if err != nil {
		return time.Time{}
	}
	queued := time.Date(now.Year(), parsed.Month(), parsed.Day(), parsed.Hour(), parsed.Minute(), 0, 0, now.Location())
	// mailq prints the time without seconds, so a message of the current minute can be up to 59 seconds "in the future"
	if queued.After(now.Add(time.Minute)) {
		queued = queued.AddDate(-1, 0, 0)
	}
	return queued
}

// olderTime returns the older time, a zero time is ignored
func olderTime(oldest time.Time, value time.Time) time.Time {
	if value.IsZero() {
		return oldest
	}
	if oldest.IsZero() || value.Before(oldest) {
		return value
	}
	return oldest
}
