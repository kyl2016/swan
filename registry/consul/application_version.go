package consul

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Dataman-Cloud/swan/types"
	"github.com/Sirupsen/logrus"
	consul "github.com/hashicorp/consul/api"
	//"gopkg.in/mgo.v2/bson"
)

// RegisterApplicationVersion is used to register a application version in consul. Use applicationId as
// key, and application version information as value.
func (c *Consul) RegisterApplicationVersion(applicationVersion *types.ApplicationVersion) error {
	versionId := time.Now().UnixNano()
	data, err := json.Marshal(applicationVersion)
	if err != nil {
		logrus.Infof("Marshal application failed: %s", err.Error())
		return err
	}

	// application version key is as format: applications/applicaitonId/versions/versionId
	appVer := consul.KVPair{
		Key:   fmt.Sprintf("applications/%s/versions/%d", applicationVersion.ID, versionId),
		Value: data,
	}

	_, err = c.client.KV().Put(&appVer, nil)
	if err != nil {
		logrus.Info("Register application %s in consul failed: %s", applicationVersion.ID, err.Error())
		return err
	}

	return nil
}

// ListApplicationVersions is used to retrieve all version ids for a application from consul by application id.
func (c *Consul) ListApplicationVersions(applicationId string) ([]string, error) {
	vers, _, err := c.client.KV().Keys(fmt.Sprintf("applications/%s/versions", applicationId), "", nil)
	if err != nil {
		logrus.Errorf("Fetch appliction failed: %s", err.Error())
		return nil, err
	}

	versions := make([]string, 0)
	for _, ver := range vers {
		v := strings.TrimPrefix(ver, fmt.Sprintf("applications/%s/versions/", applicationId))
		versions = append(versions, v)
	}

	return versions, nil
}

// FetchApplicationVersion is used to fetch specified version by version id from consul.
func (c *Consul) FetchApplicationVersion(applicationId, versionId string) (*types.ApplicationVersion, error) {
	ver, _, err := c.client.KV().Get(fmt.Sprintf("applications/%s/versions/%s", applicationId, versionId), nil)
	if err != nil {
		logrus.Errorf("Fetch appliction version failed: %s", err.Error())
		return nil, err
	}

	if ver == nil {
		logrus.Errorf("Application version %s not found in consul", versionId)
		return nil, err
	}

	var version types.ApplicationVersion
	logrus.Info("=====", string(ver.Value))
	if err := json.Unmarshal(ver.Value, &version); err != nil {
		logrus.Errorf("Unmarshal application version failed: %s", err.Error())
		return nil, err
	}

	return &version, nil
}
