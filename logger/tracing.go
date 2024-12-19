package logger

import (
	"context"

	"go.uber.org/zap"
)

const (
	OpentracingContextKeyTraceID  = "trace-id"
	OpentracingContextKeySpanID   = "span-id"
	OpentracingContextKeyParentID = "parent-id"

	OpentracingLogKeyTraceID  = "trace_id"
	OpentracingLogKeySpanID   = "span_id"
	OpentracingLogKeyParentID = "parent_id"
)

type OpentracingContextKey = string

type tracingData struct {
	traceID  string
	spanID   string
	parentID string
}

func (t *tracingData) ToInterfaceSlice() []interface{} {
	return []interface{}{
		OpentracingLogKeyTraceID,
		t.traceID,
		OpentracingLogKeySpanID,
		t.spanID,
		OpentracingLogKeyParentID,
		t.parentID,
	}
}

func (t *tracingData) ToFieldSlice() []Field {
	return []Field{
		zap.String(OpentracingLogKeyTraceID, t.traceID),
		zap.String(OpentracingLogKeySpanID, t.spanID),
		zap.String(OpentracingLogKeyParentID, t.parentID),
	}
}

func extractTracingDataFromContext(ctx context.Context) *tracingData {
	traceID, ok := ctx.Value(OpentracingContextKey(OpentracingContextKeyTraceID)).(string)
	if !ok || traceID == "" {
		return nil
	}
	spanID, ok := ctx.Value(OpentracingContextKey(OpentracingContextKeySpanID)).(string)
	if !ok || spanID == "" {
		return nil
	}
	parentID, _ := ctx.Value(OpentracingContextKey(OpentracingContextKeyParentID)).(string)
	if !ok {
		parentID = ""
	}
	return &tracingData{
		traceID:  traceID,
		spanID:   spanID,
		parentID: parentID,
	}
}
