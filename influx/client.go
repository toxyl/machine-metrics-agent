package influx

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/query"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Client struct {
	url      string
	org      string
	bucket   string
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
	queryAPI api.QueryAPI
}

func NewClient(url, org, bucket, token string, verifyTLS bool) *Client {
	client := influxdb2.NewClientWithOptions(url, token, influxdb2.DefaultOptions().SetTLSConfig(&tls.Config{
		InsecureSkipVerify: !verifyTLS,
	}))
	c := &Client{
		url:      url,
		org:      org,
		bucket:   bucket,
		client:   client,
		writeAPI: client.WriteAPIBlocking(org, bucket),
		queryAPI: client.QueryAPI(org),
	}
	return c
}

func (c *Client) write(measurement string, moment time.Time, tags map[string]string, fields map[string]interface{}) error {
	return c.writeAPI.WritePoint(context.Background(), write.NewPoint(measurement, tags, fields, moment))
}

func (c *Client) WriteSingle(measurement string, moment time.Time, tagName, tagValue string, fieldName string, fieldValue interface{}) error {
	return c.write(measurement,
		moment,
		map[string]string{
			tagName: tagValue,
		}, map[string]interface{}{
			fieldName: fieldValue,
		})
}

func (c *Client) WriteMultiple(measurement string, moment time.Time, tagName, tagValue string, fields map[string]interface{}) error {
	return c.write(measurement, moment, map[string]string{tagName: tagValue}, fields)
}

func (c *Client) Query(q string) ([]*query.FluxRecord, error) {
	q = strings.ReplaceAll(q, "%B%", c.bucket)
	results, err := c.queryAPI.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	records := []*query.FluxRecord{}
	for results.Next() {
		records = append(records, results.Record())
	}
	return records, results.Err()
}
