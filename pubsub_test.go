// Copyright 2025 Google LLC
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

package bugbash

import (
	"fmt"
	"testing"

	pstest "cloud.google.com/go/pubsub/v2/pstest"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
)

func newFake(t *testing.T) (*pstest.Server, []option.ClientOption) {
	srv := pstest.NewServer()
	opts := []option.ClientOption{
		option.WithEndpoint(srv.Addr),
		option.WithoutAuthentication(),
		option.WithGRPCDialOption(grpc.WithInsecure()),
		option.WithTelemetryDisabled(),
	}
	return srv, opts
}

func TestProduce(t *testing.T) {
	_, opts := newFake(t)
	if err := produceMessage(opts...); err != nil {
		t.Fatal(err)
	}
}

func TestConsume(t *testing.T) {
	fake, opts := newFake(t)
	fullTopicName := fmt.Sprintf("projects/%s/topics/%s", projectID, topicID)
	fake.Publish(fullTopicName, []byte("a"), nil)
	if err := consumeMessage(opts...); err != nil {
		t.Fatal(err)
	}
}

func TestAdmin(t *testing.T) {
	_, opts := newFake(t)
	if err := setupAdmin(opts...); err != nil {
		t.Fatal(err)
	}
}
