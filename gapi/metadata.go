package gapi

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayAgentHeader = "grpcgateway-user-agent"
	userAgentHeader        = "user-agent"
	xForwardedForHeader    = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (server *Server) ectractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		fmt.Printf("md: %v", md)
		if userAgent := md.Get(grpcGatewayAgentHeader); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}

		if userAgent := md.Get(userAgentHeader); len(userAgent) > 0 {
			mtdt.UserAgent = userAgent[0]
		}

		if userHeader := md.Get(xForwardedForHeader); len(userHeader) > 0 {
			mtdt.ClientIP = userHeader[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
