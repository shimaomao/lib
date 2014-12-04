package api

import (
	"fmt"
	"github.com/speedland/lib/models/tv"
	"io"
)

func (c *ApiClient) GetTvRecords() ([]*tv.TvRecord, error) {
	endpoint := buildUrl("/api/pt/records/")
	if resp, err := c.Get(endpoint); err != nil {
		return nil, err
	} else {
		var records []*tv.TvRecord
		if err := handleAsJson(resp, &records); err != nil {
			return nil, err
		} else {
			return records, nil
		}
	}
}

func GetTvRecords() ([]*tv.TvRecord, error) {
	return DefaultApiClient.GetTvRecords()
}

func (c *ApiClient) GetTvChannels() ([]*tv.TvChannel, error) {
	endpoint := buildUrl("/api/pt/channels/")
	if resp, err := c.Get(endpoint); err != nil {
		return nil, err
	} else {
		var channels []*tv.TvChannel
		if err := handleAsJson(resp, &channels); err != nil {
			return nil, err
		} else {
			return channels, nil
		}
	}
}

func GetTvChannels() ([]*tv.TvChannel, error) {
	return DefaultApiClient.GetTvChannels()
}

func (c *ApiClient) UploadPrograms(cid string, jsondata io.Reader) (map[string][]interface{}, error) {
	endpoint := buildUrl(fmt.Sprintf("/api/pt/epgs/%s", cid))
	if resp, err := c.Post(endpoint, "application/json", jsondata); err != nil {
		return nil, err
	} else {
		var result map[string][]interface{}
		if err = handleAsJson(resp, &result); err != nil {
			return nil, err
		} else {
			return result, nil
		}
	}
}

func UploadPrograms(cid string, jsondata io.Reader) (map[string][]interface{}, error) {
	return DefaultApiClient.UploadPrograms(cid, jsondata)
}
