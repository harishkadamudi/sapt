package snapshot

import (
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
)

func makeRoute(routeName, clusterName string) *route.RouteConfiguration {
	return &route.RouteConfiguration{
		Name:         routeName,
		VirtualHosts: []*route.VirtualHost{},
	}
}
