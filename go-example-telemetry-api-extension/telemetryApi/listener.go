// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT-0

package telemetryApi

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-collections/go-datastructures/queue"
	"github.com/valyala/fasthttp"
)

const defaultListenerPort = "4323"
const initialQueueSize = 5

// Used to listen to the Telemetry API
type TelemetryApiListener struct {
	fasthttpServer *fasthttp.Server
	// LogEventsQueue is a synchronous queue and is used to put the received log events to be dispatched later
	LogEventsQueue *queue.Queue
}

func NewTelemetryApiListener() *TelemetryApiListener {
	return &TelemetryApiListener{
		fasthttpServer: nil,
		LogEventsQueue: queue.New(initialQueueSize),
	}
}

func listenOnAddress() string {
	env_aws_local, ok := os.LookupEnv("AWS_SAM_LOCAL")
	var addr string
	if ok && env_aws_local == "true" {
		addr = ":" + defaultListenerPort
	} else {
		addr = "sandbox:" + defaultListenerPort
	}

	return addr
}

// Starts the server in a goroutine where the log events will be sent
func (s *TelemetryApiListener) Start() (string, error) {
	address := listenOnAddress()
	l.Info("[listener:Start] Starting on address", address)
	s.fasthttpServer = &fasthttp.Server{
		Handler: s.http_handler,
	}
	go func() {
		if err := s.fasthttpServer.ListenAndServe(address); err != nil {
			l.Error("[listener:goroutine] Unexpected stop on Http Server:", err)
			s.Shutdown()
		} else {
			l.Info("[listener:goroutine] Http Server closed:", err)
		}
	}()
	return fmt.Sprintf("http://%s/", address), nil
}

// http_handler handles the requests coming from the Telemetry API.
// Everytime Telemetry API sends log events, this function will read them from the response body
// and put into a synchronous queue to be dispatched later.
// Logging or printing besides the error cases below is not recommended if you have subscribed to
// receive extension logs. Otherwise, logging here will cause Telemetry API to send new logs for
// the printed lines which may create an infinite loop.
func (s *TelemetryApiListener) http_handler(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()

	if strings.Contains(string(body), "logsDropped") {
		l.Info("Dropped: %v", body)
	}
}

// Terminates the HTTP server listening for logs
func (s *TelemetryApiListener) Shutdown() {
	if s.fasthttpServer != nil {
		_, _ = context.WithTimeout(context.Background(), 1*time.Second)
		err := s.fasthttpServer.Shutdown()
		if err != nil {
			l.Error("[listener:Shutdown] Failed to shutdown http server gracefully:", err)
		} else {
			s.fasthttpServer = nil
		}
	}
}
