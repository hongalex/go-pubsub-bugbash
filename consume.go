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
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	subID       = "bugbash-sub"
	fullSubName = fmt.Sprintf("projects/%s/subscriptions/%s", projectID, subID)
)

// Starts pulling messages from a subscription.
// Optimistically creates this subscription if it doesn't exist.
func consumeMessage(opts ...option.ClientOption) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID, opts...)
	if err != nil {
		log.Fatalf("Could not create pubsub Client: %v", err)
	}
	defer client.Close()

	// TODO: change this line
	sub := client.Subscription(subID)

	// Receive messages for 10 seconds, which simplifies testing.
	// Comment this out in production, since `Receive` should
	// be used as a long running operation.
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// TODO: change these settings
	sub.ReceiveSettings = pubsub.ReceiveSettings{
		MaxExtension:         30 * time.Minute,
		MinExtensionPeriod:   1 * time.Minute,
		MaxExtensionPeriod:   5 * time.Minute,
		Synchronous:          true,
		UseLegacyFlowControl: true,
	}

	err = sub.Receive(ctx, func(_ context.Context, msg *pubsub.Message) {
		log.Printf("Got from existing subscription: %q\n", string(msg.Data))
		msg.Ack()
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return err
		}
		if st.Code() == codes.NotFound {
			// TODO: change this function invocation.
			s, err := client.CreateSubscription(ctx, subID, pubsub.SubscriptionConfig{
				Topic: client.Topic(topicID),
			})
			if err != nil {
				return err
			}

			// Pull from the new subscription.
			ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			s.Receive(ctx2, func(_ context.Context, msg *pubsub.Message) {
				log.Printf("Got from new subscriber: %q\n", string(msg.Data))
				msg.Ack()
			})
		}
	}
	return nil
}
