// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package driver

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/signadot/hotrod/pkg/kafka"
	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/tracing"
	"go.uber.org/zap"
)

type Processor struct {
	logger log.Factory
}

func NewProcessor(logger log.Factory) *Processor {
	return &Processor{
		logger: logger,
	}
}

// Run starts the Driver server
func (p *Processor) Run() error {
	ctx, cancel := context.WithCancel(context.Background())

	p.logger.For(ctx).Info("Starting a new consumer")

	// get a tracer provider for the driver
	tracerProvider := tracing.InitOTEL("driver", p.logger)

	// create a consumer handler
	consumer := newConsumer(ctx, tracerProvider, p.logger)

	// create a new tracer provider for kafka
	kafkaTracerProvider := tracing.InitOTEL("kafka", p.logger)

	var (
		consumerGroup sarama.ConsumerGroup
		handler       sarama.ConsumerGroupHandler
		err           error
	)
	ticker := time.NewTicker(time.Second / 2)
	defer ticker.Stop()
	for {
		// create a consumer group
		consumerGroup, handler, err = kafka.GetConsumerGroup(
			"hotrod-driver", "driver", kafkaTracerProvider, consumer)
		if err == nil {
			break
		}
		select {
		case <-ctx.Done():
			cancel()
			return fmt.Errorf("error creating consumer group client: %v", err)
		case <-ticker.C:
			p.logger.For(ctx).Error("error creating consumer group client", zap.Error(err))
			p.logger.For(ctx).Info("retrying")
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			err := consumerGroup.Consume(ctx, []string{kafka.DispatchDriverTopic}, handler)
			if err != nil {
				p.logger.For(ctx).Error("Error from consumer", zap.Error(err))
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	p.logger.For(ctx).Info("Consumer up and running!")
	http.Handle("/healthz", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		p.logger.For(ctx).Info("handling /healthz")
		resp.Write([]byte("ok"))
	}))
	go http.ListenAndServe(":8082", http.DefaultServeMux)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	keepRunning := true
	for keepRunning {
		select {
		case <-ctx.Done():
			p.logger.For(ctx).Info("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			p.logger.For(ctx).Info("terminating: via signal")
			keepRunning = false
		}
	}
	cancel()
	wg.Wait()
	if err := consumerGroup.Close(); err != nil {
		p.logger.For(ctx).Error("could not close kafka consumer group", zap.Error(err))
	}
	return nil
}
