package monitor

/*import (
	"git.jd.com/jcloud-api-gateway/jcloud-sdk-go/core"
	"git.jd.com/jcloud-api-gateway/jcloud-sdk-go/services/monitor/apis"
	"git.jd.com/jcloud-api-gateway/jcloud-sdk-go/services/monitor/client"
	"git.jd.com/jcloud-api-gateway/jcloud-sdk-go/services/monitor/models"
	uuid "github.com/satori/go.uuid"
	"jdcloud.com/dbaas/common/utils/log/golog"
	"jdcloud.com/dbaas/task_manager/context"
	"time"
)

var (
	HandleTagCode = []int64{1, 2}
	OperationFlag = []int64{1, 2, 3}
	TagsFlag      = []string{"instance_id", "pin", "role", "cluster_id"}
	VMMetrics     = []string{"cpu.util", "memory.pused", "network.bytes.incoming", "network.bytes.outgoing", "disk.iops.write", "disk.iops.read"}
)

type MonitorConf struct {
	MonitorUrl     string `toml:"monitorUrl"`
	AppCode        string `toml:"appCode"`
	SrcServiceCode string `toml:"srcServiceCode"`
	ServiceCode    string `toml:"serviceCode"`
}

type TagsInfo struct {
	RegionID   string
	ResourceID string
	Role       string
	InstanceID string
	ClusterGID string
	Pin        string
}

func AddAuxiliaryTags(requestID string, tagsInfo *TagsInfo) (*apis.MaintainAddTagResponse, error) {
	var dataTagList []models.DataTag
	conf := context.ContextInstance.Config
	accessKey := conf.AuxiliaryTags.AccessKey
	secretKey := conf.AuxiliaryTags.SecretKey
	config := core.Config{
		Endpoint: conf.AuxiliaryTags.MonitorURL,
		Scheme:   core.SchemeHttp,
	}
	config.SetTimeout(time.Duration(conf.AuxiliaryTags.Timeout) * time.Second)
	//new client
	credential := core.NewCredential(accessKey, secretKey)
	monitorClient := client.NewMonitorClient(credential)
	monitorClient.SetConfig(&config)
	monitorClient.SetLogger(core.DefaultLogger{core.LogFatal})
	//build request
	regionID := tagsInfo.RegionID
	appCode := conf.AuxiliaryTags.AppCode
	groupCode := conf.AuxiliaryTags.SrcServiceCode + "_" + uuid.NewV4().String()
	if tagsInfo.Role == "R" {
		tagsInfo.Role = "M"
	}
	resourceIds := []string{tagsInfo.ResourceID}
	serviceCode := conf.AuxiliaryTags.ServiceCode
	srcServiceCode := conf.AuxiliaryTags.SrcServiceCode
	dataTag := new(models.DataTag)
	tagsValue := make(map[string]string)
	tagsValue["instance_id"] = tagsInfo.InstanceID
	tagsValue["role"] = tagsInfo.Role
	tagsValue["pin"] = tagsInfo.Pin
	tagsValue["cluster_id"] = tagsInfo.ClusterGID
	for _, tag := range TagsFlag {
		for key, value := range tagsValue {
			if tag == key {
				operation := OperationFlag[0]
				tagKey := tag
				tagValue := value
				dataTag.Operation = &operation
				dataTag.TagKey = &tagKey
				dataTag.TagValue = &tagValue
				dataTagList = append(dataTagList, *dataTag)
			}
		}
	}
	tags := new(models.HandleTags)
	tags.HandleTagCode = &HandleTagCode[0]
	tags.HandleTags = dataTagList
	tags.PrefixMetric = &conf.AuxiliaryTags.MetricPrefix

	request := apis.NewMaintainAddTagRequest(regionID, appCode, groupCode, resourceIds, serviceCode, srcServiceCode, tags)
	userPin := tagsInfo.Pin
	request.AddHeader("x-jdcloud-pin", userPin)
	golog.Infof(requestID, "MaintainAddTag request: %+v", request)
	golog.Debugf(requestID, "MaintainAddTag request tags: %+v, tagCode: %d, prefixMetric: %s", request.Tags.HandleTags, *request.Tags.HandleTagCode, *request.Tags.PrefixMetric)
	for _, handleTags := range request.Tags.HandleTags {
		golog.Debugf(requestID, "operation: %d, tagKey: %s, tagValue: %s", *handleTags.Operation, *handleTags.TagKey, *handleTags.TagValue)
	}
	//send request
	response, err := monitorClient.MaintainAddTag(request)
	golog.Infof(requestID, "MaintainAddTag response: %+v", response)
	if err != nil {
		return nil, err
	}

	return response, nil
}*/
