/*
 * JuiceFS, Copyright (C) 2020 Juicedata, Inc.
 *
 * This program is free software: you can use, redistribute, and/or modify
 * it under the terms of the GNU Affero General Public License, version 3
 * or later ("AGPL"), as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package fuse

import (
	"sync"
	"time"

	"github.com/juicedata/juicefs/pkg/meta"
	"github.com/juicedata/juicefs/pkg/vfs"

	"github.com/hanwen/go-fuse/v2/fuse"
)

type Ino = meta.Ino
type Attr = meta.Attr
type Context = vfs.LogContext

type fuseContext struct {
	start    time.Time
	header   *fuse.InHeader
	canceled bool
	cancel   <-chan struct{}
}

var contextPool = sync.Pool{
	New: func() interface{} {
		return &fuseContext{}
	},
}

func newContext(cancel <-chan struct{}, header *fuse.InHeader) *fuseContext {
	ctx := contextPool.Get().(*fuseContext)
	ctx.start = time.Now()
	ctx.canceled = false
	ctx.cancel = cancel
	ctx.header = header
	return ctx
}

func releaseContext(ctx *fuseContext) {
	contextPool.Put(ctx)
}

func (c *fuseContext) Uid() uint32 {
	return uint32(c.header.Uid)
}

func (c *fuseContext) Gid() uint32 {
	return uint32(c.header.Gid)
}

func (c *fuseContext) Pid() uint32 {
	return uint32(c.header.Pid)
}

func (c *fuseContext) Duration() time.Duration {
	return time.Since(c.start)
}

func (c *fuseContext) Cancel() {
	c.canceled = true
}

func (c *fuseContext) Canceled() bool {
	if c.canceled {
		return true
	}
	select {
	case <-c.cancel:
		return true
	default:
		return false
	}
}
