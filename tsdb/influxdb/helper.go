package influxdb

import (
	"fmt"
	"github.com/3th1nk/easygo/util/arrUtil"
)

func (this *Client) ShowRetentionPolicies(db string) ([]*RetentionPolicy, error) {
	series, err := this.RawQuery(db, "show retention policies")
	if err != nil {
		return nil, err
	}
	if len(series) == 0 {
		return nil, nil
	}
	return series[0].toRetentionPolicies(), nil
}

func (this *Client) CreateRetentionPolicy(db string, rp *RetentionPolicy) (bool, error) {
	rps, err := this.ShowRetentionPolicies(db)
	if err != nil {
		return false, err
	}
	for _, r := range rps {
		// 如果已经存在，则不创建
		if r.Name == rp.Name {
			return false, nil
		}
	}

	// CREATE RETENTION POLICY <retention_policy_name> ON <database_name> DURATION <duration> REPLICATION <n> [SHARD DURATION <duration>] [DEFAULT]
	sql := fmt.Sprintf("create retention policy \"%s\" on \"%s\" duration %s replication %d", rp.Name, db, rp.Duration, rp.Replication)
	if rp.ShardGroupDuration != "" {
		sql += " shard duration " + rp.ShardGroupDuration
	}
	if rp.Default {
		sql += " default"
	}

	if _, err = this.RawQuery(db, sql); err != nil {
		return false, err
	}
	return true, nil
}

func (this *Client) AlterRetentionPolicy(db string, rp *RetentionPolicy) error {
	// ALTER RETENTION POLICY <retention_policy_name> ON <database_name> DURATION <duration> REPLICATION <n> [SHARD DURATION <duration>] [DEFAULT]
	sql := fmt.Sprintf("alter retention policy \"%s\" on \"%s\" duration %s replication %d", rp.Name, db, rp.Duration, rp.Replication)
	if rp.ShardGroupDuration != "" {
		sql += " shard duration " + rp.ShardGroupDuration
	}
	if rp.Default {
		sql += " default"
	}
	_, err := this.RawQuery(db, sql)
	return err
}

func (this *Client) DropRetentionPolicy(db, rp string) error {
	// 当前启用的 retention policy 不允许删除
	rps, err := this.ShowRetentionPolicies(db)
	if err != nil {
		return err
	}
	for _, r := range rps {
		if r.Name == rp && r.Default {
			return fmt.Errorf("retention policy \"%s\" is enabled, can't drop", rp)
		}
	}

	// DROP RETENTION POLICY <retention_policy_name> ON <database_name>
	_, err = this.RawQuery("", fmt.Sprintf("drop retention policy \"%s\" on \"%s\"", rp, db))
	return err
}

func (this *Client) ShowDatabases() ([]string, error) {
	series, err := this.RawQuery("", "show databases")
	if err != nil {
		return nil, err
	}
	if len(series) == 0 {
		return nil, nil
	}

	return series[0].toStringSlice(), nil
}

func (this *Client) CreateDatabase(db string) (bool, error) {
	if db == "" {
		return false, nil
	}
	databases, err := this.ShowDatabases()
	if err != nil {
		return false, err
	}

	if arrUtil.ContainsString(databases, db) {
		return false, nil
	}

	if _, err = this.RawQuery("", "create database "+db); err != nil {
		return false, err
	}
	return true, nil
}

func (this *Client) DropDatabase(db string) error {
	if db == "" {
		return nil
	}
	_, err := this.RawQuery("", "drop database "+db)
	return err
}

func (this *Client) ShowMeasurements(db string) ([]string, error) {
	series, err := this.RawQuery(db, "show measurements")
	if err != nil {
		return nil, err
	}
	if len(series) == 0 {
		return nil, nil
	}

	return series[0].toStringSlice(), nil
}

func (this *Client) DropMeasurement(db, measurement string) error {
	_, err := this.RawQuery(db, "drop measurement "+measurement)
	return err
}

func (this *Client) ShowTagKeys(db, measurement string) ([]string, error) {
	series, err := this.RawQuery(db, "show tag keys from "+measurement)
	if err != nil {
		return nil, err
	}
	if len(series) == 0 {
		return nil, nil
	}

	return series[0].toStringSlice(), nil
}

func (this *Client) ShowFieldKeys(db, measurement string) ([]string, error) {
	series, err := this.RawQuery(db, "show field keys from "+measurement)
	if err != nil {
		return nil, err
	}
	if len(series) == 0 {
		return nil, nil
	}

	return series[0].toStringSlice(), nil
}

func (this *Client) ShowSeries(db, measurement string) ([]string, error) {
	series, err := this.RawQuery(db, "show series from "+measurement)
	if err != nil {
		return nil, err
	}
	if len(series) == 0 {
		return nil, nil
	}

	return series[0].toStringSlice(), nil
}
