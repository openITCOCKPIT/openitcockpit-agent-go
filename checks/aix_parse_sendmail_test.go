//go:build aix || linux
// +build aix linux

package checks

import (
	"reflect"
	"testing"
	"time"
)

// Output of "mailq" (requires root) on AIX 7.3
func TestParseMailqEmpty(t *testing.T) {
	output := `/var/spool/mqueue is empty
                Total requests: 0
`
	now := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)
	if result := parseMailq(output, now); !reflect.DeepEqual(result, &resultSendmail{}) {
		t.Errorf("unexpected result %+v", result)
	}
}

// Output of "mailq" (requires root, runs with TZ=UTC0 in the check). The queue of the AIX system was empty, the
// messages have the format of sendmail 8 (queue id, size, queue time, sender, reason, recipients).
func TestParseMailq(t *testing.T) {
	// The second message is being delivered (*), the third one has two recipients
	output := `		/var/spool/mqueue (3 requests)
-----Q-ID----- --Size-- -----Q-Time----- ------------Sender/Recipient-----------
5AG9xKq1012345      512 Mon Nov 16 10:42 <root@node1>
                 (Deferred: Connection refused by mail.example.com.)
					 <user@example.com>
5AGA1Bc2012346*    2048 Mon Nov 16 11:58 <root@node1>
					 <user@example.com>
5AF81Cd3012347     1000 Sun Nov  8 09:00 <daemon@node1>
                 (host map: lookup (example.org): deferred)
					 <admin@example.org>
					 <other@example.org>
		Total requests: 3
`
	now := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)
	want := &resultSendmail{
		Messages:  3,
		Deferred:  2,
		SizeBytes: 3560,
		OldestAge: int64(now.Sub(time.Date(2026, 11, 8, 9, 0, 0, 0, time.UTC)).Seconds()),
	}
	if result := parseMailq(output, now); !reflect.DeepEqual(result, want) {
		t.Errorf("got %+v, want %+v", result, want)
	}
}

func TestParseSendmailTime(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 30, 0, 0, time.UTC)
	if got := parseSendmailTime("Wed Dec 31 23:00", now); !got.Equal(time.Date(2025, 12, 31, 23, 0, 0, 0, time.UTC)) {
		t.Errorf("a time after now must be from the previous year: %s", got)
	}
	if got := parseSendmailTime("Thu Jan  1 00:30", now); !got.Equal(time.Date(2026, 1, 1, 0, 30, 0, 0, time.UTC)) {
		t.Errorf("unexpected time %s", got)
	}
	if got := parseSendmailTime("invalid", now); !got.IsZero() {
		t.Errorf("expected a zero time for an invalid value: %s", got)
	}
}
