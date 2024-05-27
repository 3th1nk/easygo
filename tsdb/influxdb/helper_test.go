package influxdb

import (
	"fmt"
	"github.com/3th1nk/easygo/util/jsonUtil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_CreateDatabase(t *testing.T) {
	_, err := influx.CreateDatabase(testDB)
	assert.NoError(t, err)

	databases, err := influx.ShowDatabases()
	assert.NoError(t, err)
	t.Log(databases)

	err = influx.DropDatabase(testDB)
	assert.NoError(t, err)
}

func initTestDbRp() error {
	_, err := influx.CreateDatabase(testDB)
	if err != nil {
		return err
	}

	_, err = influx.CreateRetentionPolicy(testDB, testRP)
	return err
}

func TestClient_CreateRetentionPolicies(t *testing.T) {
	err := initTestDbRp()
	assert.NoError(t, err)

	data, err := influx.ShowRetentionPolicies(testDB)
	fmt.Println(jsonUtil.MustMarshalToStringIndent(data))

	//err = influx.DropRetentionPolicy(testDB, testRP.Name)
	//assert.NoError(t, err)
}
