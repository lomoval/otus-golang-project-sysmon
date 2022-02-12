package client

import (
	"context"
	"fmt"
	"github.com/lomoval/otus-golang-project-sysmon/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net"
	"strings"
	"time"
)

type Client struct {
	host  string
	port  string
	table *table
}

func New(host string, port string) *Client {
	return &Client{host: host, port: port}
}

func (c *Client) Start(ctx context.Context, groupName string, avgInterval int, notifyInterval int) error {
	conn, err := grpc.Dial(
		net.JoinHostPort(c.host, c.port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	client := api.NewMetricsClient(conn)
	m, err := client.GetMetrics(ctx, &api.GetMetricsRequest{
		NotifyInterval:      int32(notifyInterval),
		AverageCalcInterval: int32(avgInterval),
	})
	if err != nil {
		s := status.Convert(err)
		if s.Code() == codes.Canceled {
			return nil
		}
		return err
	}

	println("waiting first notification from server...")
	for {
		response, err := m.Recv()
		if err != nil {
			s := status.Convert(err)
			if s.Code() == codes.Canceled {
				return nil
			}
			return err
		}

		found := false
		for _, group := range response.GetGroups() {
			if strings.EqualFold(group.GetName(), strings.ToLower(groupName)) {
				if err := c.printTable(group); err != nil {
					return err
				}
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("no metrics group with name '%s'\n", groupName)
		}
	}
}

func (c *Client) printTable(group *api.MetricGroup) error {
	switch c.table {
	case nil:
		var err error
		c.table, err = newTable(group)
		if err != nil {
			return err
		}
	default:
		// Move cursor to start of table to redraw
		for i := 0; i < c.table.height(); i++ {
			print("\033[F")
		}
	}

	values := make([]interface{}, len(group.GetMetrics())+1)
	values[0] = group.Timestamp.AsTime().Format(time.RFC3339)
	for _, metric := range group.Metrics {
		values[c.table.columnsIndexes[metric.Name]+1] = metric.Value
	}
	c.table.addLine(c.table.buildLine(values))
	c.table.print()
	return nil
}
