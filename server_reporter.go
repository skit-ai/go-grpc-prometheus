// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_prometheus

import (
	"time"

	"google.golang.org/grpc/codes"
)

type serverReporter struct {
	metrics     *ServerMetrics
	rpcType     grpcType
	serviceName string
	methodName  string
	startTime   time.Time
	clientUuid  string
	flowUuid    string
	first       bool
}

func newServerReporter(m *ServerMetrics, rpcType grpcType, fullMethod string) *serverReporter {
	r := &serverReporter{
		metrics: m,
		rpcType: rpcType,
		first:   true,
	}
	if r.metrics.serverHandledHistogramEnabled {
		r.startTime = time.Now()
	}
	r.serviceName, r.methodName = splitMethodName(fullMethod)
	return r
}

func (r *serverReporter) ReceivedMessage() {
	if r.first {
		r.metrics.serverStartedCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.clientUuid, r.flowUuid).Inc()
		r.first = false
	}
	r.metrics.serverStreamMsgReceived.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.clientUuid, r.flowUuid).Inc()
	if r.metrics.serverHandledHistogramEnabled {
		r.startTime = time.Now()
	}
}

func (r *serverReporter) SentMessage() {
	r.metrics.serverStreamMsgSent.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.clientUuid, r.flowUuid).Inc()
}

func (r *serverReporter) Handled(code codes.Code) {
	r.metrics.serverHandledCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, code.String(), r.clientUuid, r.flowUuid).Inc()
	if r.metrics.serverHandledHistogramEnabled {
		r.metrics.serverHandledHistogram.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, r.clientUuid, r.flowUuid).Observe(time.Since(r.startTime).Seconds())
	}
}
