package statusdetails

import (
	"log/slog"

	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

func StatusWithDetails(statusPb *status.Status, details ...protoadapt.MessageV1) *status.Status {
	statusWithDetailsPb, err := statusPb.WithDetails(details...)
	if err != nil {
		slog.Error("could not add details to error",
			"error", err,
			"status_proto", statusPb,
			"details", details)

		return statusPb
	}

	return statusWithDetailsPb
}

func MustStatusWithDetails(
	statusPb *status.Status,
	details ...protoadapt.MessageV1,
) *status.Status {
	statusWithDetailsPb, err := statusPb.WithDetails(details...)
	if err != nil {
		slog.Error("could not add details to error",
			"error", err,
			"status_proto", statusPb,
			"details", details)
		panic(err)
	}

	return statusWithDetailsPb
}
