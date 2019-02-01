// Copyright © 2018 NAME HERE <jbonds@jbvm.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jbvmio/kafkactl"
	"github.com/jbvmio/kafkactl/cli/x/out"
)

type followDetails struct {
	topic   string
	partMap map[int32]int64
}

func FollowTopic(flags MSGFlags, outFlags out.OutFlags, topics ...string) {
	exact = true
	var count int
	var details []followDetails
	var timeCheck time.Time
	for _, topic := range topics {
		var parts []int32
		topicSummary := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
		match := true
		switch match {
		case len(topicSummary) != 1:
			closeFatal("Error isolating topic: %v\n", topic)
		case flags.Partition != -1:
			parts = append(parts, flags.Partition)
		case len(flags.Partitions) == 0:
			parts = topicSummary[0].Partitions
		default:
			parts = validateParts(flags.Partitions)
		}
		startMap := make(map[int32]int64, len(parts))
		offset := getTailValue(flags.Tail)
		for _, p := range parts {
			off, err := client.GetOffsetNewest(topic, p)
			if err != nil {
				closeFatal("Error validating Partition: %v for topic: %v\n", p, err)
			}
			startMap[p] = off + offset
			count++
		}
		d := followDetails{
			topic:   topic,
			partMap: startMap,
		}
		details = append(details, d)
	}
	msgChan := make(chan *kafkactl.Message, 100)
	stopChan := make(chan bool, count)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	for _, d := range details {
		for part, offset := range d.partMap {
			go client.ChanPartitionConsume(d.topic, part, offset, msgChan, stopChan)
		}
	}
ConsumeLoop:
	for {
		select {
		case msg := <-msgChan:
			if msg.Timestamp == timeCheck {
				if len(msg.Value) != 0 {
					out.Warnf("%s", msg.Value)
				}
			} else {
				PrintMSG(msg, outFlags)
			}
			continue ConsumeLoop
		case <-sigChan:
			fmt.Printf("signal: interrupt\n  Stopping kafkactl ...\n")
			for i := 0; i < count; i++ {
				stopChan <- true
			}
			break ConsumeLoop
		}
	}
}

/*
func getTopicMsg(topic string, partition int32, offset int64) {
	msg, err := client.ConsumeOffsetMsg("testtopic", 0, 1955)
	if err != nil {
		closeFatal("Error: %v\n", err)
	}
	fmt.Printf("%s", msg.Value)
}

func tailTopic(topic string, relativeOffset int64, partitions ...int32) {
	if relativeOffset > 0 {
		closeFatal("reletive offset must be a negative number")
	}
	exact = true
	tSum := kafkactl.GetTopicSummaries(SearchTopicMeta(topic))
	if len(tSum) != 1 {
		closeFatal("Error finding topic: %v\n", topic)
	}
	if len(partitions) == 0 {
		partitions = tSum[0].Partitions
	}
	pMap := make(map[int32]int64)
	for _, ts := range tSum {
		for _, p := range partitions {
			off, err := client.GetOffsetNewest(ts.Topic, p)
			if err != nil {
				closeFatal("Error validating Partition: %v for topic: %v\n", p, err)
			}
			pMap[p] = off
		}
	}
	msgChan := make(chan *kafkactl.Message, 100)
	stopChan := make(chan bool, len(pMap))
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	for part, offset := range pMap {
		rel := (offset + relativeOffset)
		go client.ChanPartitionConsume(topic, part, rel, msgChan, stopChan)
	}
ConsumeLoop:
	for {
		select {
		case msg := <-msgChan:
			fmt.Printf("%s\n", msg.Value)
		case <-sigChan:
			fmt.Printf("signal: interrupt\n  Stopping kafkactl ...\n")
			for i := 0; i < len(pMap); i++ {
				stopChan <- true
			}
			break ConsumeLoop
		}
	}
}
*/