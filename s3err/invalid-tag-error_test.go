// Copyright 2026 Versity Software
// This file is licensed under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package s3err

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidTagErrorXMLBody(t *testing.T) {
	tests := []struct {
		name        string
		err         InvalidTagError
		contains    []string
		notContains []string
	}{
		{
			name: "empty tag key",
			err:  GetInvalidTagErr(ErrInvalidTagKey, "", ""),
			contains: []string{
				"<Code>InvalidTag</Code>",
				"<Message>The TagKey you have provided is invalid</Message>",
				"<TagKey></TagKey>",
			},
			notContains: []string{"<TagValue>"},
		},
		{
			name: "invalid tag key",
			err:  GetInvalidTagErr(ErrInvalidTagKey, "te&st", ""),
			contains: []string{
				"<Code>InvalidTag</Code>",
				"<TagKey>te&amp;st</TagKey>",
			},
			notContains: []string{"<TagValue>"},
		},
		{
			name: "invalid tag value",
			err:  GetInvalidTagErr(ErrInvalidTagValue, "key", "val\tue"),
			contains: []string{
				"<Code>InvalidTag</Code>",
				"<Message>The TagValue you have provided is invalid</Message>",
				"<TagKey>key</TagKey>",
				"<TagValue>val&#x9;ue</TagValue>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := string(tt.err.XMLBody("reqid", "hostid"))
			for _, s := range tt.contains {
				assert.Contains(t, body, s)
			}
			for _, s := range tt.notContains {
				assert.NotContains(t, body, s)
			}
			assert.Contains(t, body, "<RequestId>reqid</RequestId>")
			assert.Contains(t, body, "<HostId>hostid</HostId>")
		})
	}
}
