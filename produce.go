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
	"context"
	"log"

	// TODO: change this import
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	// TODO(developer): replace with your own project
	projectID = "alxh-pubsub"

	topicID = "bugbash-topic"
)

func produceMessage(opts ...option.ClientOption) error {
	ctx := context.Background()
	c, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}
	defer c.Close()

	// TODO: change this line
	topic := c.Topic(topicID)

	err = publishSingleMessage(ctx, topic)
	if err != nil {
		e, ok := status.FromError(err)
		if !ok {
			return err
		}
		// If the publish failed because the topic doesn't exist
		// create the topic and try publishing again.
		if e.Code() == codes.NotFound {
			// TODO: change this line
			topic, err := c.CreateTopic(ctx, topicID)
			if err != nil {
				return err
			}
			return publishSingleMessage(ctx, topic)
		}
	}
	return err
}

// TODO: change argument
func publishSingleMessage(ctx context.Context, topic *pubsub.Topic) error {
	res := topic.Publish(ctx, &pubsub.Message{
		Data: []byte("a single message"),
	})
	// The publish happens asynchronously.
	_, err := res.Get(ctx)
	return err
}
