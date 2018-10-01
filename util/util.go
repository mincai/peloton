package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"

	mesos "code.uber.internal/infra/peloton/.gen/mesos/v1"
	"code.uber.internal/infra/peloton/.gen/peloton/api/v0/job"
	"code.uber.internal/infra/peloton/.gen/peloton/api/v0/peloton"
	"code.uber.internal/infra/peloton/.gen/peloton/api/v0/task"
	"code.uber.internal/infra/peloton/.gen/peloton/private/hostmgr/hostsvc"
)

const (
	// ResourceEpsilon is the minimum epsilon mesos resource;
	// This is because Mesos internally uses a fixed point precision. See MESOS-4687 for details.
	ResourceEpsilon = 0.0009
)

// UUIDLength represents the length of a 16 byte v4 UUID as a string
var UUIDLength = len(uuid.New())

// Min returns the minimum value of x, y
func Min(x, y uint32) uint32 {
	if x < y {
		return x
	}
	return y
}

// Max returns the maximum value of x, y
func Max(x, y uint32) uint32 {
	if x > y {
		return x
	}
	return y
}

// PtrPrintf returns a pointer to a string format
func PtrPrintf(format string, a ...interface{}) *string {
	str := fmt.Sprintf(format, a...)
	return &str
}

// GetOfferScalarResourceSummary generates a summary for all the scalar values: role -> offerName-> Value
// first level : role -> map(resource type-> resouce value)
func GetOfferScalarResourceSummary(offer *mesos.Offer) map[string]map[string]float64 {
	var result = make(map[string]map[string]float64)
	for _, resource := range offer.Resources {
		if resource.Scalar != nil {
			var role = "*"
			if resource.Role != nil {
				role = *resource.Role
			}
			if _, ok := result[role]; !ok {
				result[role] = make(map[string]float64)
			}
			result[role][*resource.Name] = result[role][*resource.Name] + *resource.Scalar.Value
		}
	}
	return result
}

// CreateMesosScalarResources is a helper function to convert resource values into Mesos resources.
func CreateMesosScalarResources(values map[string]float64, role string) []*mesos.Resource {
	var rs []*mesos.Resource
	for name, value := range values {
		// Skip any value smaller than Espilon.
		if math.Abs(value) < ResourceEpsilon {
			continue
		}

		rs = append(rs, NewMesosResourceBuilder().WithName(name).WithValue(value).WithRole(role).Build())
	}

	return rs
}

// MesosResourceBuilder is the helper to build a mesos resource
type MesosResourceBuilder struct {
	Resource mesos.Resource
}

// NewMesosResourceBuilder creates a MesosResourceBuilder
func NewMesosResourceBuilder() *MesosResourceBuilder {
	defaultRole := "*"
	defaultType := mesos.Value_SCALAR
	return &MesosResourceBuilder{
		Resource: mesos.Resource{
			Role: &defaultRole,
			Type: &defaultType,
		},
	}
}

// WithName sets name
func (o *MesosResourceBuilder) WithName(name string) *MesosResourceBuilder {
	o.Resource.Name = &name
	return o
}

// WithType sets type
func (o *MesosResourceBuilder) WithType(t mesos.Value_Type) *MesosResourceBuilder {
	o.Resource.Type = &t
	return o
}

// WithRole sets role
func (o *MesosResourceBuilder) WithRole(role string) *MesosResourceBuilder {
	o.Resource.Role = &role
	return o
}

// WithValue sets value
func (o *MesosResourceBuilder) WithValue(value float64) *MesosResourceBuilder {
	scalarVal := mesos.Value_Scalar{
		Value: &value,
	}
	o.Resource.Scalar = &scalarVal
	return o
}

// WithRanges sets ranges
func (o *MesosResourceBuilder) WithRanges(ranges *mesos.Value_Ranges) *MesosResourceBuilder {
	o.Resource.Ranges = ranges
	return o
}

// WithReservation sets reservation info.
func (o *MesosResourceBuilder) WithReservation(
	reservation *mesos.Resource_ReservationInfo) *MesosResourceBuilder {

	o.Resource.Reservation = reservation
	return o
}

// WithDisk sets disk info.
func (o *MesosResourceBuilder) WithDisk(
	diskInfo *mesos.Resource_DiskInfo) *MesosResourceBuilder {

	o.Resource.Disk = diskInfo
	return o
}

// WithRevocable sets resource as revocable resource type
func (o *MesosResourceBuilder) WithRevocable(
	revocable *mesos.Resource_RevocableInfo) *MesosResourceBuilder {
	o.Resource.Revocable = revocable
	return o
}

// TODO: add other building functions when needed

// Build returns the mesos resource
func (o *MesosResourceBuilder) Build() *mesos.Resource {
	res := o.Resource
	return &res
}

// MesosStateToPelotonState translates mesos task state to peloton task state
// TODO: adjust in case there are additional peloton states
func MesosStateToPelotonState(mstate mesos.TaskState) task.TaskState {
	switch mstate {
	case mesos.TaskState_TASK_STAGING:
		return task.TaskState_LAUNCHED
	case mesos.TaskState_TASK_STARTING:
		return task.TaskState_STARTING
	case mesos.TaskState_TASK_RUNNING:
		return task.TaskState_RUNNING
		// NOTE: This should only be sent when the framework has
		// the TASK_KILLING_STATE capability.
	case mesos.TaskState_TASK_KILLING:
		return task.TaskState_RUNNING
	case mesos.TaskState_TASK_FINISHED:
		return task.TaskState_SUCCEEDED
	case mesos.TaskState_TASK_FAILED:
		return task.TaskState_FAILED
	case mesos.TaskState_TASK_KILLED:
		return task.TaskState_KILLED
	case mesos.TaskState_TASK_LOST:
		return task.TaskState_LOST
	case mesos.TaskState_TASK_ERROR:
		return task.TaskState_FAILED
	default:
		log.Errorf("Unknown mesos taskState %v", mstate)
		return task.TaskState_INITIALIZED
	}
}

// IsPelotonStateTerminal returns true if state is terminal
// otherwise false
func IsPelotonStateTerminal(state task.TaskState) bool {
	switch state {
	case task.TaskState_SUCCEEDED, task.TaskState_FAILED,
		task.TaskState_KILLED, task.TaskState_LOST:
		return true
	default:
		return false
	}
}

// IsPelotonJobStateTerminal returns true if job state is terminal
// otherwise false
func IsPelotonJobStateTerminal(state job.JobState) bool {
	switch state {
	case job.JobState_SUCCEEDED, job.JobState_FAILED, job.JobState_KILLED:
		return true
	default:
		return false
	}
}

// IsTaskHasValidVolume returns true if a task is stateful and has a valid volume
func IsTaskHasValidVolume(taskInfo *task.TaskInfo) bool {
	if taskInfo.GetConfig().GetVolume() != nil &&
		len(taskInfo.GetRuntime().GetVolumeID().GetValue()) != 0 {
		return true
	}
	return false
}

// CreateMesosTaskID creates mesos task id given jobID, instanceID and runID
func CreateMesosTaskID(jobID *peloton.JobID,
	instanceID uint32,
	runID uint64) *mesos.TaskID {
	mesosID := fmt.Sprintf(
		"%s-%d-%d",
		jobID.GetValue(),
		instanceID,
		runID)
	return &mesos.TaskID{Value: &mesosID}
}

// ParseRunID parse the runID from mesosTaskID
func ParseRunID(mesosTaskID string) (uint64, error) {
	splitMesosTaskID := strings.Split(mesosTaskID, "-")
	if len(mesosTaskID) == 0 { // prev mesos task id is nil
		return 0, errors.New("mesosTaskID provided is empty")
	} else if len(splitMesosTaskID) == 7 {
		if runID, err := strconv.ParseUint(
			splitMesosTaskID[len(splitMesosTaskID)-1], 10, 64); err == nil {
			return runID, nil
		}
	}

	return 0, errors.New("unable to parse mesos task id: " + mesosTaskID)
}

// ParseTaskID parses the jobID and instanceID from peloton taskID
func ParseTaskID(taskID string) (string, uint32, error) {
	pos := strings.LastIndex(taskID, "-")
	if len(taskID) < UUIDLength || pos == -1 {
		return "", 0, fmt.Errorf("invalid pelotonTaskID %v", taskID)
	}
	jobID := taskID[0:pos]
	ins := taskID[pos+1:]

	instanceID, err := strconv.ParseUint(ins, 10, 32)
	if err != nil {
		log.WithFields(log.Fields{
			"task_id": taskID,
			"job_id":  jobID,
		}).WithError(err).Error("failed to parse taskID")
		log.Info(err)
		return "",
			0,
			fmt.Errorf("unable to parse instanceID %v", taskID)
	}
	return jobID, uint32(instanceID), err
}

// ParseTaskIDFromMesosTaskID parses the taskID from mesosTaskID
func ParseTaskIDFromMesosTaskID(mesosTaskID string) (string, error) {
	// mesos task id would be "(jobID)-(instanceID)-(runID)" form
	if len(mesosTaskID) < UUIDLength+1 {
		return "", fmt.Errorf("invalid mesostaskID %v", mesosTaskID)
	}

	// TODO: deprecate the check once mesos task id migration is complete from
	// uuid-int-uuid -> uuid(job ID)-int(instance ID)-int(monotonically incremental)
	// If uuid has all digits from uuid-int-uuid then it will increment from
	// that value and not default to 1.
	var pelotonTaskID string
	if len(mesosTaskID) > 2*UUIDLength {
		pelotonTaskID = mesosTaskID[:len(mesosTaskID)-(UUIDLength+1)]
	} else {
		pelotonTaskID = mesosTaskID[:strings.LastIndex(mesosTaskID, "-")]
	}

	_, _, err := ParseTaskID(pelotonTaskID)
	if err != nil {
		log.WithError(err).
			WithField("mesos_task_id", mesosTaskID).
			Error("Invalid mesos task ID, cannot parse jobID / instance")
		return "", err
	}
	return pelotonTaskID, nil
}

// ParseJobAndInstanceID return jobID and instanceID from given mesos task id.
func ParseJobAndInstanceID(mesosTaskID string) (string, uint32, error) {
	pelotonTaskID, err := ParseTaskIDFromMesosTaskID(mesosTaskID)
	if err != nil {
		return "", 0, err
	}
	return ParseTaskID(pelotonTaskID)
}

// UnmarshalToType unmarshal a string to a typed interface{}
func UnmarshalToType(jsonString string, resultType reflect.Type) (interface{},
	error) {
	result := reflect.New(resultType)
	err := json.Unmarshal([]byte(jsonString), result.Interface())
	if err != nil {
		log.Errorf("Unmarshal failed with error %v, type %v, jsonString %v",
			err, resultType, jsonString)
		return nil, err
	}
	return result.Interface(), nil
}

// ConvertLabels will convert Peloton labels to Mesos labels.
func ConvertLabels(pelotonLabels []*peloton.Label) *mesos.Labels {
	mesosLabels := &mesos.Labels{
		Labels: make([]*mesos.Label, 0, len(pelotonLabels)),
	}
	for _, label := range pelotonLabels {
		mesosLabel := &mesos.Label{
			Key:   &label.Key,
			Value: &label.Value,
		}
		mesosLabels.Labels = append(mesosLabels.Labels, mesosLabel)
	}
	return mesosLabels
}

// Contains checks whether an item contains in a list
func Contains(list []string, item string) bool {
	for _, name := range list {
		if name == item {
			return true
		}
	}
	return false
}

// CreateHostInfo takes the agent Info and create the hostsvc.HostInfo
func CreateHostInfo(hostname string,
	agentInfo *mesos.AgentInfo) *hostsvc.HostInfo {
	// if agentInfo is nil , return nil HostInfo
	if agentInfo == nil {
		log.WithField("host", hostname).
			Warn("Agent Info is nil")
		return nil
	}
	// return the HostInfo object
	return &hostsvc.HostInfo{
		Hostname:   hostname,
		AgentId:    agentInfo.Id,
		Attributes: agentInfo.Attributes,
		Resources:  agentInfo.Resources,
	}
}

// CreateSecretVolume builds a mesos volume of type secret
// from the given secret path and secret value string
// This volume will be added to the job's default config
func CreateSecretVolume(secretPath string, secretStr string) *mesos.Volume {
	volumeMode := mesos.Volume_RO
	volumeSourceType := mesos.Volume_Source_SECRET
	secretType := mesos.Secret_VALUE
	return &mesos.Volume{
		Mode:          &volumeMode,
		ContainerPath: &secretPath,
		Source: &mesos.Volume_Source{
			Type: &volumeSourceType,
			Secret: &mesos.Secret{
				Type: &secretType,
				Value: &mesos.Secret_Value{
					Data: []byte(secretStr),
				},
			},
		},
	}
}

// IsSecretVolume returns true if the given volume is of type secret
func IsSecretVolume(volume *mesos.Volume) bool {
	return volume.GetSource().GetType() == mesos.Volume_Source_SECRET
}

// ConfigHasSecretVolumes returns true if config contains secret volumes
func ConfigHasSecretVolumes(config *task.TaskConfig) bool {
	for _, v := range config.GetContainer().GetVolumes() {
		if ok := IsSecretVolume(v); ok {
			return true
		}
	}
	return false
}

// RemoveSecretVolumesFromConfig removes secret volumes from the task config
// in place and returns the secret volumes
// Secret volumes are added internally at the time of creating a job with
// secrets by handleSecrets method. They are not supplied in the config in
// job create/update requests. Consequently, they should not be displayed
// as part of Job Get API response. This is necessary to achieve the broader
// goal of using the secrets proto message in Job Create/Update/Get API to
// describe secrets and not allow users to checkin secrets as part of config
func RemoveSecretVolumesFromConfig(config *task.TaskConfig) []*mesos.Volume {
	if config.GetContainer().GetVolumes() == nil {
		return nil
	}
	secretVolumes := []*mesos.Volume{}
	volumes := []*mesos.Volume{}
	for _, volume := range config.GetContainer().GetVolumes() {
		if ok := IsSecretVolume(volume); ok {
			secretVolumes = append(secretVolumes, volume)
		} else {
			volumes = append(volumes, volume)
		}
	}
	config.GetContainer().Volumes = nil
	if len(volumes) > 0 {
		config.GetContainer().Volumes = volumes
	}
	return secretVolumes
}

// RemoveSecretVolumesFromJobConfig removes secret volumes from the default
// config as well as instance config in place and returns the secret volumes
func RemoveSecretVolumesFromJobConfig(cfg *job.JobConfig) []*mesos.Volume {
	// remove secret volumes if present from default config
	secretVolumes := RemoveSecretVolumesFromConfig(cfg.GetDefaultConfig())

	// remove secret volumes if present from instance config
	for _, config := range cfg.GetInstanceConfig() {
		// instance config contains the same secret volumes as default config,
		// so no need to operate on them
		_ = RemoveSecretVolumesFromConfig(config)
	}
	return secretVolumes
}
