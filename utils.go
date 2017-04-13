package rrd

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

// MD5 return md5
func MD5(raw string) string {
	h := md5.New()
	io.WriteString(h, raw)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func formatDBfile(rrdDir, ds, dsType, pk string) string {
	uid := fmt.Sprintf("%s%s%s", ds, dsType, pk)
	md5 := MD5(uid)
	dbfile := fmt.Sprintf("%s/%s/%s.rrd", rrdDir, md5[0:2], md5)

	return dbfile
}

// Unix returns the local Time corresponding to the given Unix time
func Unix(sec int64) time.Time {
	return time.Unix(sec, 0)
}
