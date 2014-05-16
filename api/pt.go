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

func GetTvRecords() ([]*models.TvRecord, error) {
	return DefaultApiClient.GetTvRecords()
}
