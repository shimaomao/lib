package api

import (
	"github.com/speedland/lib/models"
)

func (c *ApiClient) GetTvRecords() ([]*models.TvRecord, error) {
	endpoint := buildUrl("/api/pt/records/")
	if resp, err := c.Get(endpoint); err != nil {
		return nil, err
	} else {
		var records []*models.TvRecord
		if err := handleAsJson(resp, &records); err != nil {
			return nil, err
		} else {
			return records, nil
		}
	}
}

func (c *ApiClient) GetTvChannels() ([]*models.TvChannel, error) {
	endpoint := buildUrl("/api/pt/channels/")
	if resp, err := c.Get(endpoint); err != nil {
		return nil, err
	} else {
		var channels []*models.TvChannel
		if err := handleAsJson(resp, &channels); err != nil {
			return nil, err
		} else {
			return channels, nil
		}
	}
}

func GetTvRecords() ([]*models.TvRecord, error) {
	return DefaultApiClient.GetTvRecords()
}
