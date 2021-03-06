// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bigquery

import "golang.org/x/net/context"

// RecordsPerRequest returns a ReadOption that sets the number of records to fetch per request when streaming data from BigQuery.
func RecordsPerRequest(n int64) ReadOption { return recordsPerRequest(n) }

type recordsPerRequest int64

func (opt recordsPerRequest) customizeRead(conf *pagingConf) {
	conf.recordsPerRequest = int64(opt)
	conf.setRecordsPerRequest = true
}

// StartIndex returns a ReadOption that sets the zero-based index of the row to start reading from.
func StartIndex(i uint64) ReadOption { return startIndex(i) }

type startIndex uint64

func (opt startIndex) customizeRead(conf *pagingConf) {
	conf.startIndex = uint64(opt)
}

type tableFetcher struct {
	c      *Client
	cursor *readTableCursor
}

func (tf *tableFetcher) fetch(ctx context.Context, token string) (*readDataResult, error) {
	tf.cursor.paging.pageToken = token
	return tf.c.service.readTabledata(ctx, tf.cursor)
}

func (c *Client) readTable(t *Table, options []ReadOption) (*Iterator, error) {
	cursor := &readTableCursor{}
	t.customizeReadSrc(cursor)

	for _, o := range options {
		o.customizeRead(&cursor.paging)
	}

	tf := &tableFetcher{
		c:      c,
		cursor: cursor,
	}
	return &Iterator{pf: tf}, nil
}

type queryFetcher struct {
	c      *Client
	cursor *readQueryCursor
}

func (qf *queryFetcher) fetch(ctx context.Context, token string) (*readDataResult, error) {
	qf.cursor.paging.pageToken = token
	return qf.c.service.readQuery(ctx, qf.cursor)
}

func (c *Client) readQueryResults(job *Job, options []ReadOption) (*Iterator, error) {
	cursor := &readQueryCursor{}
	if err := job.customizeReadQuery(cursor); err != nil {
		return nil, err
	}

	for _, o := range options {
		o.customizeRead(&cursor.paging)
	}

	qf := &queryFetcher{
		c:      c,
		cursor: cursor,
	}
	return &Iterator{pf: qf}, nil
}
