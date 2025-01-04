package statusdetails

import (
	"errors"
	"fmt"
	"strings"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

// ErrDetails holds the google/rpc/error_details.proto messages.
type ErrDetails struct {
	ErrorInfo           *errdetails.ErrorInfo
	BadRequest          *errdetails.BadRequest
	PreconditionFailure *errdetails.PreconditionFailure
	QuotaFailure        *errdetails.QuotaFailure
	RetryInfo           *errdetails.RetryInfo
	ResourceInfo        *errdetails.ResourceInfo
	RequestInfo         *errdetails.RequestInfo
	DebugInfo           *errdetails.DebugInfo
	Help                *errdetails.Help
	LocalizedMessage    *errdetails.LocalizedMessage

	// Unknown stores unidentifiable error details.
	Unknown []interface{}
}

// ErrMessageNotFound is used to signal ExtractProtoMessage found no matching messages.
var ErrMessageNotFound = errors.New("message not found")

// ExtractProtoMessage provides a mechanism for extracting protobuf messages from the
// Unknown error details. If ExtractProtoMessage finds an unknown message of the same type,
// the content of the message is copied to the provided message.
//
// ExtractProtoMessage will return ErrMessageNotFound if there are no message matching the
// protocol buffer type of the provided message.
func (e ErrDetails) ExtractProtoMessage(v proto.Message) error {
	if v == nil {
		return ErrMessageNotFound
	}
	for _, elem := range e.Unknown {
		if elemProto, ok := elem.(proto.Message); ok {
			if v.ProtoReflect().Type() == elemProto.ProtoReflect().Type() {
				proto.Merge(v, elemProto)
				return nil
			}
		}
	}
	return ErrMessageNotFound
}

func (e ErrDetails) String() string {
	var d strings.Builder
	if e.ErrorInfo != nil {
		d.WriteString(
			fmt.Sprintf("error details: name = ErrorInfo reason = %s domain = %s metadata = %s\n",
				e.ErrorInfo.GetReason(), e.ErrorInfo.GetDomain(), e.ErrorInfo.GetMetadata()),
		)
	}

	if e.BadRequest != nil {
		v := e.BadRequest.GetFieldViolations()
		var f []string
		var desc []string
		for _, x := range v {
			f = append(f, x.GetField())
			desc = append(desc, x.GetDescription())
		}
		d.WriteString(fmt.Sprintf("error details: name = BadRequest field = %s desc = %s\n",
			strings.Join(f, " "), strings.Join(desc, " ")))
	}

	if e.PreconditionFailure != nil {
		v := e.PreconditionFailure.GetViolations()
		var t []string
		var s []string
		var desc []string
		for _, x := range v {
			t = append(t, x.GetType())
			s = append(s, x.GetSubject())
			desc = append(desc, x.GetDescription())
		}
		d.WriteString(
			fmt.Sprintf(
				"error details: name = PreconditionFailure type = %s subj = %s desc = %s\n",
				strings.Join(t, " "),
				strings.Join(s, " "),
				strings.Join(desc, " "),
			),
		)
	}

	if e.QuotaFailure != nil {
		v := e.QuotaFailure.GetViolations()
		var s []string
		var desc []string
		for _, x := range v {
			s = append(s, x.GetSubject())
			desc = append(desc, x.GetDescription())
		}
		d.WriteString(fmt.Sprintf("error details: name = QuotaFailure subj = %s desc = %s\n",
			strings.Join(s, " "), strings.Join(desc, " ")))
	}

	if e.RequestInfo != nil {
		d.WriteString(fmt.Sprintf("error details: name = RequestInfo id = %s data = %s\n",
			e.RequestInfo.GetRequestId(), e.RequestInfo.GetServingData()))
	}

	if e.ResourceInfo != nil {
		d.WriteString(
			fmt.Sprintf(
				"error details: name = ResourceInfo type = %s resourcename = %s owner = %s desc = %s\n",
				e.ResourceInfo.GetResourceType(),
				e.ResourceInfo.GetResourceName(),
				e.ResourceInfo.GetOwner(),
				e.ResourceInfo.GetDescription(),
			),
		)

	}
	if e.RetryInfo != nil {
		d.WriteString(
			fmt.Sprintf("error details: retry in %s\n", e.RetryInfo.GetRetryDelay().AsDuration()),
		)

	}
	if e.Unknown != nil {
		var s []string
		for _, x := range e.Unknown {
			s = append(s, fmt.Sprintf("%v", x))
		}
		d.WriteString(
			fmt.Sprintf("error details: name = Unknown  desc = %s\n", strings.Join(s, " ")),
		)
	}

	if e.DebugInfo != nil {
		d.WriteString(
			fmt.Sprintf(
				"error details: name = DebugInfo detail = %s stack = %s\n",
				e.DebugInfo.GetDetail(),
				strings.Join(e.DebugInfo.GetStackEntries(), " "),
			),
		)
	}
	if e.Help != nil {
		var desc []string
		var url []string
		for _, x := range e.Help.Links {
			desc = append(desc, x.GetDescription())
			url = append(url, x.GetUrl())
		}
		d.WriteString(fmt.Sprintf("error details: name = Help desc = %s url = %s\n",
			strings.Join(desc, " "), strings.Join(url, " ")))
	}
	if e.LocalizedMessage != nil {
		d.WriteString(fmt.Sprintf("error details: name = LocalizedMessage locale = %s msg = %s\n",
			e.LocalizedMessage.GetLocale(), e.LocalizedMessage.GetMessage()))
	}

	return d.String()
}

// StatusError wraps either a gRPC Status error or a HTTP googleapi.Error. It
// implements error and Status interfaces.
type StatusError struct {
	status  *status.Status
	details ErrDetails
}

// New returns a StatusError with the given code and message.
func New(code codes.Code, message string) *StatusError {
	st := status.New(code, message)
	return &StatusError{
		status:  st,
		details: ErrDetails{},
	}
}

// Newf returns a StatusError with the given code and formatted message.
func Newf(code codes.Code, format string, args ...interface{}) *StatusError {
	return New(code, fmt.Sprintf(format, args...))
}

// Error returns a StatusError representing an error with the given code and message.
func Error(code codes.Code, message string) error {
	return New(code, message)
}

// Errorf returns a StatusError representing an error with the given code and formatted message.
func Errorf(code codes.Code, format string, args ...interface{}) error {
	return New(code, fmt.Sprintf(format, args...))
}

// Details presents the error details of the StatusError.
func (a *StatusError) Details() ErrDetails {
	return a.details
}

// Error returns a readable representation of the StatusError.
func (a *StatusError) Error() string {
	return a.status.Message()
}

// GRPCStatus extracts the underlying gRPC Status error.
// This method is necessary to fulfill the interface
// described in https://pkg.go.dev/google.golang.org/grpc/status#FromError.
func (a *StatusError) GRPCStatus() *status.Status {
	return a.status
}

// Reason returns the reason in an ErrorInfo.
// If ErrorInfo is nil, it returns an empty string.
func (a *StatusError) Reason() string {
	return a.details.ErrorInfo.GetReason()
}

// Domain returns the domain in an ErrorInfo.
// If ErrorInfo is nil, it returns an empty string.
func (a *StatusError) Domain() string {
	return a.details.ErrorInfo.GetDomain()
}

// Metadata returns the metadata in an ErrorInfo.
// If ErrorInfo is nil, it returns nil.
func (a *StatusError) Metadata() map[string]string {
	return a.details.ErrorInfo.GetMetadata()
}

// WithErrorInfo sets both status detail and errdetail for ErrorInfo
func (a *StatusError) WithErrorInfo(info *errdetails.ErrorInfo) *StatusError {
	newDetails := a.details
	newDetails.ErrorInfo = info
	newStatus, err := a.status.WithDetails(info)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

// WithBadRequest sets both status detail and errdetail for BadRequest
func (a *StatusError) WithBadRequest(badRequest *errdetails.BadRequest) *StatusError {
	newDetails := a.details
	newDetails.BadRequest = badRequest
	newStatus, err := a.status.WithDetails(badRequest)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

// WithPreconditionFailure sets both status detail and errdetail for PreconditionFailure
func (a *StatusError) WithPreconditionFailure(
	precondition *errdetails.PreconditionFailure,
) *StatusError {
	newDetails := a.details
	newDetails.PreconditionFailure = precondition
	newStatus, err := a.status.WithDetails(precondition)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

// WithQuotaFailure sets both status detail and errdetail for QuotaFailure
func (a *StatusError) WithQuotaFailure(quota *errdetails.QuotaFailure) *StatusError {
	newDetails := a.details
	newDetails.QuotaFailure = quota
	newStatus, err := a.status.WithDetails(quota)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

// WithRetryInfo sets both status detail and errdetail for RetryInfo
func (a *StatusError) WithRetryInfo(retry *errdetails.RetryInfo) *StatusError {
	newDetails := a.details
	newDetails.RetryInfo = retry
	newStatus, err := a.status.WithDetails(retry)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

// WithResourceInfo sets both status detail and errdetail for ResourceInfo
func (a *StatusError) WithResourceInfo(resource *errdetails.ResourceInfo) *StatusError {
	newDetails := a.details
	newDetails.ResourceInfo = resource
	newStatus, err := a.status.WithDetails(resource)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

// WithRequestInfo sets both status detail and errdetail for RequestInfo
func (a *StatusError) WithRequestInfo(request *errdetails.RequestInfo) *StatusError {
	newDetails := a.details
	newDetails.RequestInfo = request
	newStatus, err := a.status.WithDetails(request)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

// WithDebugInfo sets both status detail and errdetail for DebugInfo
func (a *StatusError) WithDebugInfo(debug *errdetails.DebugInfo) *StatusError {
	newDetails := a.details
	newDetails.DebugInfo = debug
	newStatus, err := a.status.WithDetails(debug)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

// WithHelp sets both status detail and errdetail for Help
func (a *StatusError) WithHelp(help *errdetails.Help) *StatusError {
	newDetails := a.details
	newDetails.Help = help
	newStatus, err := a.status.WithDetails(help)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

// WithLocalizedMessage sets both status detail and errdetail for LocalizedMessage
func (a *StatusError) WithLocalizedMessage(message *errdetails.LocalizedMessage) *StatusError {
	newDetails := a.details
	newDetails.LocalizedMessage = message
	newStatus, err := a.status.WithDetails(message)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

func (a *StatusError) WithCustomDetail(details ...protoiface.MessageV1) *StatusError {
	newDetails := a.details
	newDetails.Unknown = append(newDetails.Unknown, sliceToInterface(details)...)
	newStatus, err := a.status.WithDetails(details...)
	if err == nil {
		return &StatusError{details: newDetails, status: newStatus}
	}
	return a
}

// FromError parses a Status error or a googleapi.Error and builds an
// StatusError, wrapping the provided error in the new APIError. It
// returns false if neither Status nor googleapi.Error can be parsed.
func FromError(err error) (*StatusError, bool) {
	return ParseError(err, true)
}

// ParseError parses a Status error  and builds an
// StatusError. If wrap is true, it wraps the error in the new APIError.
// It returns false if neither Status nor googleapi.Error can be parsed.
func ParseError(err error, wrap bool) (*StatusError, bool) {
	if err == nil {
		return nil, false
	}
	ae := StatusError{}
	if !ae.setDetailsFromError(err) {
		return nil, false
	}
	return &ae, true
}

// setDetailsFromError parses a Status error
// and sets status and details
// It returns false if its not Status error.
func (a *StatusError) setDetailsFromError(err error) bool {
	if st, isStatus := status.FromError(err); isStatus {
		a.status = st
		a.details = parseDetails(st.Details())
		return true
	}
	return false
}

// parseDetails accepts a slice of interface{} that should be backed by some
// sort of proto.Message that can be cast to the google/rpc/error_details.proto
// types.
//
// This is for internal use only.
func parseDetails(details []interface{}) ErrDetails {
	var ed ErrDetails
	for _, d := range details {
		switch d := d.(type) {
		case *errdetails.ErrorInfo:
			ed.ErrorInfo = d
		case *errdetails.BadRequest:
			ed.BadRequest = d
		case *errdetails.PreconditionFailure:
			ed.PreconditionFailure = d
		case *errdetails.QuotaFailure:
			ed.QuotaFailure = d
		case *errdetails.RetryInfo:
			ed.RetryInfo = d
		case *errdetails.ResourceInfo:
			ed.ResourceInfo = d
		case *errdetails.RequestInfo:
			ed.RequestInfo = d
		case *errdetails.DebugInfo:
			ed.DebugInfo = d
		case *errdetails.Help:
			ed.Help = d
		case *errdetails.LocalizedMessage:
			ed.LocalizedMessage = d
		default:
			ed.Unknown = append(ed.Unknown, d)
		}
	}

	return ed
}

func sliceToInterface[T any](ar []T) []interface{} {
	anyAr := make([]interface{}, len(ar))
	for i := range ar {
		anyAr[i] = ar[i]
	}

	return anyAr
}
