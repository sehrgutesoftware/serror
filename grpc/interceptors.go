package grpc

import (
	"context"
	"encoding/json"

	"github.com/sehrgutesoftware/serror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const trailerKey = "X-Serror"

// EncodeErrors is a UnaryServerInterceptor that serializes the error returned
// from the handler and attaches it to the gRPC response as a trailer.
func EncodeError(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	if err == nil {
		return resp, err
	}

	// Encode the error
	errTree := serror.Encode(err)
	if errTree == nil {
		return resp, nil
	}
	encodedErr, err := json.Marshal(errTree)
	if err != nil {
		return nil, err
	}

	// Attach metadata as trailer
	trailer := metadata.Pairs(trailerKey, string(encodedErr))
	grpc.SetTrailer(ctx, trailer)

	return resp, err
}

// HydrateError returns a UnaryClientInterceptor that deserializes errors
// from the gRPC trailer using the provided code resolver.
func HydrateError(resolve serror.CodeResolver) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		var trailer metadata.MD
		opts = append(opts, grpc.Trailer(&trailer))

		// Make the RPC call
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			return err
		}

		// Retrieve the trailers
		trailers := trailer.Get(trailerKey)
		if len(trailers) == 0 {
			return nil
		}

		// Decode the error
		var errTree serror.Tree
		if err := json.Unmarshal([]byte(trailers[0]), &errTree); err != nil {
			return err
		}
		serror.Hydrate(&errTree, resolve)

		return &errTree
	}
}
