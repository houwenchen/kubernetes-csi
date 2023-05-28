package driver

import (
	"context"
	"fmt"
	"os"
	"strings"

	"google.golang.org/grpc"
	"k8s.io/klog/v2"
)

func parseEndpoint(ep string) (string, string, error) {
	if strings.HasPrefix(strings.ToLower(ep), "unix://") {
		s := strings.SplitN(ep, "://", 2)
		if s[1] != "" {
			return s[0], s[1], nil
		}
	}
	return "", "", fmt.Errorf("invalid endpoint: %v", ep)
}

func logGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	klog.V(2).Infof("GRPC call: %s", info.FullMethod)

	resp, err := handler(ctx, req)
	if err != nil {
		klog.Errorf("GRPC error: %v", err)
	}

	return resp, err
}

func makeVolumeDir(volDir string) error {
	_, err := os.Stat(volDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err = os.MkdirAll(volDir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
