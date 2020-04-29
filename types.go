package rebuilderd

import (
	"net/url"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var err error
	t.Time, err = time.Parse("2006-01-02T15:04:05.999999999", string(b[1:len(b)-1]))
	return err
}

// // https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L265
// type BuildReport struct {
// 	Queue QueueItem `json:"queue"`
// 	// good, bad, fail
// 	// https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L258
// 	Status string `json:"status"`
// }
//
// // https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L249
// type DropQueueItem struct {
// 	Name         string  `json:"name"`
// 	Version      *string `json:"version"`
// 	Distro       string  `json:"distro"`
// 	Suite        string  `json:"suite"`
// 	Architecture *string `json:"architecture"`
// }
//
// // https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L195
// type JobAssignment struct {
// 	Rebuild *QueueItem `json:"rebuild"`
// }

// https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L209
type ListPkgs struct {
	Name         *string `json:"name"`
	Status       *string `json:"status"`
	Distro       *string `json:"distro"`
	Suite        *string `json:"suite"`
	Architecture *string `json:"architecture"`
}

func (l *ListPkgs) Values() url.Values {
	v := url.Values{}
	if l != nil {
		if l.Name != nil {
			v.Set("name", *l.Name)
		}
		if l.Status != nil {
			v.Set("status", *l.Status)
		}
		if l.Distro != nil {
			v.Set("distro", *l.Distro)
		}
		if l.Suite != nil {
			v.Set("suite", *l.Suite)
		}
		if l.Architecture != nil {
			v.Set("architecture", *l.Architecture)
		}
	}
	return v
}

// https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L235
type ListQueue struct {
	Limit *int64 `json:"limit"`
}

// https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/lib.rs#L23
type PkgRelease struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	// GOOD, BAD, UNKWN
	Status       string `json:"status"`
	Distro       string `json:"distro"`
	Suite        string `json:"suite"`
	Architecture string `json:"architecture"`
	URL          string `json:"url"`
}

// // https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L240
// type PushQueue struct {
// 	Name         string  `json:"name"`
// 	Version      *string `json:"version"`
// 	Distro       string  `json:"distro"`
// 	Suite        string  `json:"suite"`
// 	Architecture *string `json:"architecture"`
// }
//
// // https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L218
type QueueList struct {
	Now   string      `json:"now"`
	Queue []QueueItem `json:"queue"`
}

//
// // https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L224
type QueueItem struct {
	ID        int32      `json"id"`
	Package   PkgRelease `json:"package"`
	Version   string     `json"version"`
	QueuedAt  string     `json:"queued_at"`
	WorkerID  *int32     `json:"worker_id"`
	StartedAt *string    `json:"started_at"`
	LastPing  *string    `json:"last_ping"`
}

//
// // https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L201
// type SuiteImport struct {
// 	// "archlinux", "debian"
// 	// https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/lib.rs#L17:10
// 	Distro       string `json:"distro"`
// 	Suite        string `json:"suite"`
// 	Architecture string `json:"architecture"`
// 	// https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/lib.rs#L23
// 	Pkgs []PkgRelease `json:"pkgs"`
// }

// https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L182
type Worker struct {
	Key      string  `json:"key"`
	Addr     string  `json:"addr"`
	Status   *string `json:"status"`
	LastPing Time    `json:"last_ping"`
	Online   bool    `json:"online"`
}

// // https://github.com/kpcyrd/rebuilderd/blob/6bf8e2219c87fe053af9daf6223666342108bc44/common/src/api.rs#L191
// type WorkQuery struct{}
