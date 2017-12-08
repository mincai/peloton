// Code generated by MockGen. DO NOT EDIT.
// Source: code.uber.internal/infra/peloton/storage (interfaces: JobStore,TaskStore,UpgradeStore,FrameworkInfoStore,ResourcePoolStore,PersistentVolumeStore)

package mocks

import (
	context "context"
	reflect "reflect"

	job "code.uber.internal/infra/peloton/.gen/peloton/api/job"
	peloton "code.uber.internal/infra/peloton/.gen/peloton/api/peloton"
	respool "code.uber.internal/infra/peloton/.gen/peloton/api/respool"
	task "code.uber.internal/infra/peloton/.gen/peloton/api/task"
	upgrade "code.uber.internal/infra/peloton/.gen/peloton/api/upgrade"
	volume "code.uber.internal/infra/peloton/.gen/peloton/api/volume"
	gomock "github.com/golang/mock/gomock"
)

// MockJobStore is a mock of JobStore interface
type MockJobStore struct {
	ctrl     *gomock.Controller
	recorder *MockJobStoreMockRecorder
}

// MockJobStoreMockRecorder is the mock recorder for MockJobStore
type MockJobStoreMockRecorder struct {
	mock *MockJobStore
}

// NewMockJobStore creates a new mock instance
func NewMockJobStore(ctrl *gomock.Controller) *MockJobStore {
	mock := &MockJobStore{ctrl: ctrl}
	mock.recorder = &MockJobStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockJobStore) EXPECT() *MockJobStoreMockRecorder {
	return _m.recorder
}

// CreateJob mocks base method
func (_m *MockJobStore) CreateJob(_param0 context.Context, _param1 *peloton.JobID, _param2 *job.JobConfig, _param3 string) error {
	ret := _m.ctrl.Call(_m, "CreateJob", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateJob indicates an expected call of CreateJob
func (_mr *MockJobStoreMockRecorder) CreateJob(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "CreateJob", reflect.TypeOf((*MockJobStore)(nil).CreateJob), arg0, arg1, arg2, arg3)
}

// DeleteJob mocks base method
func (_m *MockJobStore) DeleteJob(_param0 context.Context, _param1 *peloton.JobID) error {
	ret := _m.ctrl.Call(_m, "DeleteJob", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteJob indicates an expected call of DeleteJob
func (_mr *MockJobStoreMockRecorder) DeleteJob(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "DeleteJob", reflect.TypeOf((*MockJobStore)(nil).DeleteJob), arg0, arg1)
}

// GetAllJobs mocks base method
func (_m *MockJobStore) GetAllJobs(_param0 context.Context) (map[string]*job.RuntimeInfo, error) {
	ret := _m.ctrl.Call(_m, "GetAllJobs", _param0)
	ret0, _ := ret[0].(map[string]*job.RuntimeInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllJobs indicates an expected call of GetAllJobs
func (_mr *MockJobStoreMockRecorder) GetAllJobs(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetAllJobs", reflect.TypeOf((*MockJobStore)(nil).GetAllJobs), arg0)
}

// GetJobConfig mocks base method
func (_m *MockJobStore) GetJobConfig(_param0 context.Context, _param1 *peloton.JobID) (*job.JobConfig, error) {
	ret := _m.ctrl.Call(_m, "GetJobConfig", _param0, _param1)
	ret0, _ := ret[0].(*job.JobConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobConfig indicates an expected call of GetJobConfig
func (_mr *MockJobStoreMockRecorder) GetJobConfig(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetJobConfig", reflect.TypeOf((*MockJobStore)(nil).GetJobConfig), arg0, arg1)
}

// GetJobRuntime mocks base method
func (_m *MockJobStore) GetJobRuntime(_param0 context.Context, _param1 *peloton.JobID) (*job.RuntimeInfo, error) {
	ret := _m.ctrl.Call(_m, "GetJobRuntime", _param0, _param1)
	ret0, _ := ret[0].(*job.RuntimeInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobRuntime indicates an expected call of GetJobRuntime
func (_mr *MockJobStoreMockRecorder) GetJobRuntime(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetJobRuntime", reflect.TypeOf((*MockJobStore)(nil).GetJobRuntime), arg0, arg1)
}

// GetJobsByStates mocks base method
func (_m *MockJobStore) GetJobsByStates(_param0 context.Context, _param1 []job.JobState) ([]peloton.JobID, error) {
	ret := _m.ctrl.Call(_m, "GetJobsByStates", _param0, _param1)
	ret0, _ := ret[0].([]peloton.JobID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetJobsByStates indicates an expected call of GetJobsByStates
func (_mr *MockJobStoreMockRecorder) GetJobsByStates(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetJobsByStates", reflect.TypeOf((*MockJobStore)(nil).GetJobsByStates), arg0, arg1)
}

// QueryJobs mocks base method
func (_m *MockJobStore) QueryJobs(_param0 context.Context, _param1 *peloton.ResourcePoolID, _param2 *job.QuerySpec) ([]*job.JobInfo, uint32, error) {
	ret := _m.ctrl.Call(_m, "QueryJobs", _param0, _param1, _param2)
	ret0, _ := ret[0].([]*job.JobInfo)
	ret1, _ := ret[1].(uint32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// QueryJobs indicates an expected call of QueryJobs
func (_mr *MockJobStoreMockRecorder) QueryJobs(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "QueryJobs", reflect.TypeOf((*MockJobStore)(nil).QueryJobs), arg0, arg1, arg2)
}

// UpdateJobConfig mocks base method
func (_m *MockJobStore) UpdateJobConfig(_param0 context.Context, _param1 *peloton.JobID, _param2 *job.JobConfig) error {
	ret := _m.ctrl.Call(_m, "UpdateJobConfig", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateJobConfig indicates an expected call of UpdateJobConfig
func (_mr *MockJobStoreMockRecorder) UpdateJobConfig(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdateJobConfig", reflect.TypeOf((*MockJobStore)(nil).UpdateJobConfig), arg0, arg1, arg2)
}

// UpdateJobRuntime mocks base method
func (_m *MockJobStore) UpdateJobRuntime(_param0 context.Context, _param1 *peloton.JobID, _param2 *job.RuntimeInfo) error {
	ret := _m.ctrl.Call(_m, "UpdateJobRuntime", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateJobRuntime indicates an expected call of UpdateJobRuntime
func (_mr *MockJobStoreMockRecorder) UpdateJobRuntime(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdateJobRuntime", reflect.TypeOf((*MockJobStore)(nil).UpdateJobRuntime), arg0, arg1, arg2)
}

// MockTaskStore is a mock of TaskStore interface
type MockTaskStore struct {
	ctrl     *gomock.Controller
	recorder *MockTaskStoreMockRecorder
}

// MockTaskStoreMockRecorder is the mock recorder for MockTaskStore
type MockTaskStoreMockRecorder struct {
	mock *MockTaskStore
}

// NewMockTaskStore creates a new mock instance
func NewMockTaskStore(ctrl *gomock.Controller) *MockTaskStore {
	mock := &MockTaskStore{ctrl: ctrl}
	mock.recorder = &MockTaskStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockTaskStore) EXPECT() *MockTaskStoreMockRecorder {
	return _m.recorder
}

// CreateTaskConfigs mocks base method
func (_m *MockTaskStore) CreateTaskConfigs(_param0 context.Context, _param1 *peloton.JobID, _param2 *job.JobConfig) error {
	ret := _m.ctrl.Call(_m, "CreateTaskConfigs", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTaskConfigs indicates an expected call of CreateTaskConfigs
func (_mr *MockTaskStoreMockRecorder) CreateTaskConfigs(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "CreateTaskConfigs", reflect.TypeOf((*MockTaskStore)(nil).CreateTaskConfigs), arg0, arg1, arg2)
}

// CreateTaskRuntime mocks base method
func (_m *MockTaskStore) CreateTaskRuntime(_param0 context.Context, _param1 *peloton.JobID, _param2 uint32, _param3 *task.RuntimeInfo, _param4 string) error {
	ret := _m.ctrl.Call(_m, "CreateTaskRuntime", _param0, _param1, _param2, _param3, _param4)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTaskRuntime indicates an expected call of CreateTaskRuntime
func (_mr *MockTaskStoreMockRecorder) CreateTaskRuntime(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "CreateTaskRuntime", reflect.TypeOf((*MockTaskStore)(nil).CreateTaskRuntime), arg0, arg1, arg2, arg3, arg4)
}

// CreateTaskRuntimes mocks base method
func (_m *MockTaskStore) CreateTaskRuntimes(_param0 context.Context, _param1 *peloton.JobID, _param2 []*task.RuntimeInfo, _param3 string) error {
	ret := _m.ctrl.Call(_m, "CreateTaskRuntimes", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTaskRuntimes indicates an expected call of CreateTaskRuntimes
func (_mr *MockTaskStoreMockRecorder) CreateTaskRuntimes(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "CreateTaskRuntimes", reflect.TypeOf((*MockTaskStore)(nil).CreateTaskRuntimes), arg0, arg1, arg2, arg3)
}

// GetTaskByID mocks base method
func (_m *MockTaskStore) GetTaskByID(_param0 context.Context, _param1 string) (*task.TaskInfo, error) {
	ret := _m.ctrl.Call(_m, "GetTaskByID", _param0, _param1)
	ret0, _ := ret[0].(*task.TaskInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskByID indicates an expected call of GetTaskByID
func (_mr *MockTaskStoreMockRecorder) GetTaskByID(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetTaskByID", reflect.TypeOf((*MockTaskStore)(nil).GetTaskByID), arg0, arg1)
}

// GetTaskConfig mocks base method
func (_m *MockTaskStore) GetTaskConfig(_param0 context.Context, _param1 *peloton.JobID, _param2 uint32, _param3 int64) (*task.TaskConfig, error) {
	ret := _m.ctrl.Call(_m, "GetTaskConfig", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(*task.TaskConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskConfig indicates an expected call of GetTaskConfig
func (_mr *MockTaskStoreMockRecorder) GetTaskConfig(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetTaskConfig", reflect.TypeOf((*MockTaskStore)(nil).GetTaskConfig), arg0, arg1, arg2, arg3)
}

// GetTaskForJob mocks base method
func (_m *MockTaskStore) GetTaskForJob(_param0 context.Context, _param1 *peloton.JobID, _param2 uint32) (map[uint32]*task.TaskInfo, error) {
	ret := _m.ctrl.Call(_m, "GetTaskForJob", _param0, _param1, _param2)
	ret0, _ := ret[0].(map[uint32]*task.TaskInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskForJob indicates an expected call of GetTaskForJob
func (_mr *MockTaskStoreMockRecorder) GetTaskForJob(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetTaskForJob", reflect.TypeOf((*MockTaskStore)(nil).GetTaskForJob), arg0, arg1, arg2)
}

// GetTaskRuntime mocks base method
func (_m *MockTaskStore) GetTaskRuntime(_param0 context.Context, _param1 *peloton.JobID, _param2 uint32) (*task.RuntimeInfo, error) {
	ret := _m.ctrl.Call(_m, "GetTaskRuntime", _param0, _param1, _param2)
	ret0, _ := ret[0].(*task.RuntimeInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskRuntime indicates an expected call of GetTaskRuntime
func (_mr *MockTaskStoreMockRecorder) GetTaskRuntime(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetTaskRuntime", reflect.TypeOf((*MockTaskStore)(nil).GetTaskRuntime), arg0, arg1, arg2)
}

// GetTaskStateSummaryForJob mocks base method
func (_m *MockTaskStore) GetTaskStateSummaryForJob(_param0 context.Context, _param1 *peloton.JobID) (map[string]uint32, error) {
	ret := _m.ctrl.Call(_m, "GetTaskStateSummaryForJob", _param0, _param1)
	ret0, _ := ret[0].(map[string]uint32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTaskStateSummaryForJob indicates an expected call of GetTaskStateSummaryForJob
func (_mr *MockTaskStoreMockRecorder) GetTaskStateSummaryForJob(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetTaskStateSummaryForJob", reflect.TypeOf((*MockTaskStore)(nil).GetTaskStateSummaryForJob), arg0, arg1)
}

// GetTasksForJob mocks base method
func (_m *MockTaskStore) GetTasksForJob(_param0 context.Context, _param1 *peloton.JobID) (map[uint32]*task.TaskInfo, error) {
	ret := _m.ctrl.Call(_m, "GetTasksForJob", _param0, _param1)
	ret0, _ := ret[0].(map[uint32]*task.TaskInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksForJob indicates an expected call of GetTasksForJob
func (_mr *MockTaskStoreMockRecorder) GetTasksForJob(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetTasksForJob", reflect.TypeOf((*MockTaskStore)(nil).GetTasksForJob), arg0, arg1)
}

// GetTasksForJobAndState mocks base method
func (_m *MockTaskStore) GetTasksForJobAndState(_param0 context.Context, _param1 *peloton.JobID, _param2 string) (map[uint32]*task.TaskInfo, error) {
	ret := _m.ctrl.Call(_m, "GetTasksForJobAndState", _param0, _param1, _param2)
	ret0, _ := ret[0].(map[uint32]*task.TaskInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksForJobAndState indicates an expected call of GetTasksForJobAndState
func (_mr *MockTaskStoreMockRecorder) GetTasksForJobAndState(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetTasksForJobAndState", reflect.TypeOf((*MockTaskStore)(nil).GetTasksForJobAndState), arg0, arg1, arg2)
}

// GetTasksForJobByRange mocks base method
func (_m *MockTaskStore) GetTasksForJobByRange(_param0 context.Context, _param1 *peloton.JobID, _param2 *task.InstanceRange) (map[uint32]*task.TaskInfo, error) {
	ret := _m.ctrl.Call(_m, "GetTasksForJobByRange", _param0, _param1, _param2)
	ret0, _ := ret[0].(map[uint32]*task.TaskInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTasksForJobByRange indicates an expected call of GetTasksForJobByRange
func (_mr *MockTaskStoreMockRecorder) GetTasksForJobByRange(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetTasksForJobByRange", reflect.TypeOf((*MockTaskStore)(nil).GetTasksForJobByRange), arg0, arg1, arg2)
}

// QueryTasks mocks base method
func (_m *MockTaskStore) QueryTasks(_param0 context.Context, _param1 *peloton.JobID, _param2 *task.QuerySpec) ([]*task.TaskInfo, uint32, error) {
	ret := _m.ctrl.Call(_m, "QueryTasks", _param0, _param1, _param2)
	ret0, _ := ret[0].([]*task.TaskInfo)
	ret1, _ := ret[1].(uint32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// QueryTasks indicates an expected call of QueryTasks
func (_mr *MockTaskStoreMockRecorder) QueryTasks(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "QueryTasks", reflect.TypeOf((*MockTaskStore)(nil).QueryTasks), arg0, arg1, arg2)
}

// UpdateTaskRuntime mocks base method
func (_m *MockTaskStore) UpdateTaskRuntime(_param0 context.Context, _param1 *peloton.JobID, _param2 uint32, _param3 *task.RuntimeInfo) error {
	ret := _m.ctrl.Call(_m, "UpdateTaskRuntime", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTaskRuntime indicates an expected call of UpdateTaskRuntime
func (_mr *MockTaskStoreMockRecorder) UpdateTaskRuntime(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdateTaskRuntime", reflect.TypeOf((*MockTaskStore)(nil).UpdateTaskRuntime), arg0, arg1, arg2, arg3)
}

// MockUpgradeStore is a mock of UpgradeStore interface
type MockUpgradeStore struct {
	ctrl     *gomock.Controller
	recorder *MockUpgradeStoreMockRecorder
}

// MockUpgradeStoreMockRecorder is the mock recorder for MockUpgradeStore
type MockUpgradeStoreMockRecorder struct {
	mock *MockUpgradeStore
}

// NewMockUpgradeStore creates a new mock instance
func NewMockUpgradeStore(ctrl *gomock.Controller) *MockUpgradeStore {
	mock := &MockUpgradeStore{ctrl: ctrl}
	mock.recorder = &MockUpgradeStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockUpgradeStore) EXPECT() *MockUpgradeStoreMockRecorder {
	return _m.recorder
}

// AddTaskToProcessing mocks base method
func (_m *MockUpgradeStore) AddTaskToProcessing(_param0 context.Context, _param1 *upgrade.WorkflowID, _param2 uint32) error {
	ret := _m.ctrl.Call(_m, "AddTaskToProcessing", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddTaskToProcessing indicates an expected call of AddTaskToProcessing
func (_mr *MockUpgradeStoreMockRecorder) AddTaskToProcessing(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "AddTaskToProcessing", reflect.TypeOf((*MockUpgradeStore)(nil).AddTaskToProcessing), arg0, arg1, arg2)
}

// CreateUpgrade mocks base method
func (_m *MockUpgradeStore) CreateUpgrade(_param0 context.Context, _param1 *upgrade.WorkflowID, _param2 *upgrade.UpgradeSpec) error {
	ret := _m.ctrl.Call(_m, "CreateUpgrade", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUpgrade indicates an expected call of CreateUpgrade
func (_mr *MockUpgradeStoreMockRecorder) CreateUpgrade(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "CreateUpgrade", reflect.TypeOf((*MockUpgradeStore)(nil).CreateUpgrade), arg0, arg1, arg2)
}

// GetWorkflowProgress mocks base method
func (_m *MockUpgradeStore) GetWorkflowProgress(_param0 context.Context, _param1 *upgrade.WorkflowID) ([]uint32, uint32, error) {
	ret := _m.ctrl.Call(_m, "GetWorkflowProgress", _param0, _param1)
	ret0, _ := ret[0].([]uint32)
	ret1, _ := ret[1].(uint32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetWorkflowProgress indicates an expected call of GetWorkflowProgress
func (_mr *MockUpgradeStoreMockRecorder) GetWorkflowProgress(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetWorkflowProgress", reflect.TypeOf((*MockUpgradeStore)(nil).GetWorkflowProgress), arg0, arg1)
}

// RemoveTaskFromProcessing mocks base method
func (_m *MockUpgradeStore) RemoveTaskFromProcessing(_param0 context.Context, _param1 *upgrade.WorkflowID, _param2 uint32) error {
	ret := _m.ctrl.Call(_m, "RemoveTaskFromProcessing", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveTaskFromProcessing indicates an expected call of RemoveTaskFromProcessing
func (_mr *MockUpgradeStoreMockRecorder) RemoveTaskFromProcessing(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "RemoveTaskFromProcessing", reflect.TypeOf((*MockUpgradeStore)(nil).RemoveTaskFromProcessing), arg0, arg1, arg2)
}

// MockFrameworkInfoStore is a mock of FrameworkInfoStore interface
type MockFrameworkInfoStore struct {
	ctrl     *gomock.Controller
	recorder *MockFrameworkInfoStoreMockRecorder
}

// MockFrameworkInfoStoreMockRecorder is the mock recorder for MockFrameworkInfoStore
type MockFrameworkInfoStoreMockRecorder struct {
	mock *MockFrameworkInfoStore
}

// NewMockFrameworkInfoStore creates a new mock instance
func NewMockFrameworkInfoStore(ctrl *gomock.Controller) *MockFrameworkInfoStore {
	mock := &MockFrameworkInfoStore{ctrl: ctrl}
	mock.recorder = &MockFrameworkInfoStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockFrameworkInfoStore) EXPECT() *MockFrameworkInfoStoreMockRecorder {
	return _m.recorder
}

// GetFrameworkID mocks base method
func (_m *MockFrameworkInfoStore) GetFrameworkID(_param0 context.Context, _param1 string) (string, error) {
	ret := _m.ctrl.Call(_m, "GetFrameworkID", _param0, _param1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFrameworkID indicates an expected call of GetFrameworkID
func (_mr *MockFrameworkInfoStoreMockRecorder) GetFrameworkID(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetFrameworkID", reflect.TypeOf((*MockFrameworkInfoStore)(nil).GetFrameworkID), arg0, arg1)
}

// GetMesosStreamID mocks base method
func (_m *MockFrameworkInfoStore) GetMesosStreamID(_param0 context.Context, _param1 string) (string, error) {
	ret := _m.ctrl.Call(_m, "GetMesosStreamID", _param0, _param1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMesosStreamID indicates an expected call of GetMesosStreamID
func (_mr *MockFrameworkInfoStoreMockRecorder) GetMesosStreamID(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetMesosStreamID", reflect.TypeOf((*MockFrameworkInfoStore)(nil).GetMesosStreamID), arg0, arg1)
}

// SetMesosFrameworkID mocks base method
func (_m *MockFrameworkInfoStore) SetMesosFrameworkID(_param0 context.Context, _param1 string, _param2 string) error {
	ret := _m.ctrl.Call(_m, "SetMesosFrameworkID", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetMesosFrameworkID indicates an expected call of SetMesosFrameworkID
func (_mr *MockFrameworkInfoStoreMockRecorder) SetMesosFrameworkID(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "SetMesosFrameworkID", reflect.TypeOf((*MockFrameworkInfoStore)(nil).SetMesosFrameworkID), arg0, arg1, arg2)
}

// SetMesosStreamID mocks base method
func (_m *MockFrameworkInfoStore) SetMesosStreamID(_param0 context.Context, _param1 string, _param2 string) error {
	ret := _m.ctrl.Call(_m, "SetMesosStreamID", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetMesosStreamID indicates an expected call of SetMesosStreamID
func (_mr *MockFrameworkInfoStoreMockRecorder) SetMesosStreamID(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "SetMesosStreamID", reflect.TypeOf((*MockFrameworkInfoStore)(nil).SetMesosStreamID), arg0, arg1, arg2)
}

// MockResourcePoolStore is a mock of ResourcePoolStore interface
type MockResourcePoolStore struct {
	ctrl     *gomock.Controller
	recorder *MockResourcePoolStoreMockRecorder
}

// MockResourcePoolStoreMockRecorder is the mock recorder for MockResourcePoolStore
type MockResourcePoolStoreMockRecorder struct {
	mock *MockResourcePoolStore
}

// NewMockResourcePoolStore creates a new mock instance
func NewMockResourcePoolStore(ctrl *gomock.Controller) *MockResourcePoolStore {
	mock := &MockResourcePoolStore{ctrl: ctrl}
	mock.recorder = &MockResourcePoolStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockResourcePoolStore) EXPECT() *MockResourcePoolStoreMockRecorder {
	return _m.recorder
}

// CreateResourcePool mocks base method
func (_m *MockResourcePoolStore) CreateResourcePool(_param0 context.Context, _param1 *peloton.ResourcePoolID, _param2 *respool.ResourcePoolConfig, _param3 string) error {
	ret := _m.ctrl.Call(_m, "CreateResourcePool", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateResourcePool indicates an expected call of CreateResourcePool
func (_mr *MockResourcePoolStoreMockRecorder) CreateResourcePool(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "CreateResourcePool", reflect.TypeOf((*MockResourcePoolStore)(nil).CreateResourcePool), arg0, arg1, arg2, arg3)
}

// DeleteResourcePool mocks base method
func (_m *MockResourcePoolStore) DeleteResourcePool(_param0 context.Context, _param1 *peloton.ResourcePoolID) error {
	ret := _m.ctrl.Call(_m, "DeleteResourcePool", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteResourcePool indicates an expected call of DeleteResourcePool
func (_mr *MockResourcePoolStoreMockRecorder) DeleteResourcePool(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "DeleteResourcePool", reflect.TypeOf((*MockResourcePoolStore)(nil).DeleteResourcePool), arg0, arg1)
}

// GetAllResourcePools mocks base method
func (_m *MockResourcePoolStore) GetAllResourcePools(_param0 context.Context) (map[string]*respool.ResourcePoolConfig, error) {
	ret := _m.ctrl.Call(_m, "GetAllResourcePools", _param0)
	ret0, _ := ret[0].(map[string]*respool.ResourcePoolConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllResourcePools indicates an expected call of GetAllResourcePools
func (_mr *MockResourcePoolStoreMockRecorder) GetAllResourcePools(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetAllResourcePools", reflect.TypeOf((*MockResourcePoolStore)(nil).GetAllResourcePools), arg0)
}

// GetResourcePool mocks base method
func (_m *MockResourcePoolStore) GetResourcePool(_param0 context.Context, _param1 *peloton.ResourcePoolID) (*respool.ResourcePoolInfo, error) {
	ret := _m.ctrl.Call(_m, "GetResourcePool", _param0, _param1)
	ret0, _ := ret[0].(*respool.ResourcePoolInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResourcePool indicates an expected call of GetResourcePool
func (_mr *MockResourcePoolStoreMockRecorder) GetResourcePool(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetResourcePool", reflect.TypeOf((*MockResourcePoolStore)(nil).GetResourcePool), arg0, arg1)
}

// GetResourcePoolsByOwner mocks base method
func (_m *MockResourcePoolStore) GetResourcePoolsByOwner(_param0 context.Context, _param1 string) (map[string]*respool.ResourcePoolConfig, error) {
	ret := _m.ctrl.Call(_m, "GetResourcePoolsByOwner", _param0, _param1)
	ret0, _ := ret[0].(map[string]*respool.ResourcePoolConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResourcePoolsByOwner indicates an expected call of GetResourcePoolsByOwner
func (_mr *MockResourcePoolStoreMockRecorder) GetResourcePoolsByOwner(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetResourcePoolsByOwner", reflect.TypeOf((*MockResourcePoolStore)(nil).GetResourcePoolsByOwner), arg0, arg1)
}

// UpdateResourcePool mocks base method
func (_m *MockResourcePoolStore) UpdateResourcePool(_param0 context.Context, _param1 *peloton.ResourcePoolID, _param2 *respool.ResourcePoolConfig) error {
	ret := _m.ctrl.Call(_m, "UpdateResourcePool", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateResourcePool indicates an expected call of UpdateResourcePool
func (_mr *MockResourcePoolStoreMockRecorder) UpdateResourcePool(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdateResourcePool", reflect.TypeOf((*MockResourcePoolStore)(nil).UpdateResourcePool), arg0, arg1, arg2)
}

// MockPersistentVolumeStore is a mock of PersistentVolumeStore interface
type MockPersistentVolumeStore struct {
	ctrl     *gomock.Controller
	recorder *MockPersistentVolumeStoreMockRecorder
}

// MockPersistentVolumeStoreMockRecorder is the mock recorder for MockPersistentVolumeStore
type MockPersistentVolumeStoreMockRecorder struct {
	mock *MockPersistentVolumeStore
}

// NewMockPersistentVolumeStore creates a new mock instance
func NewMockPersistentVolumeStore(ctrl *gomock.Controller) *MockPersistentVolumeStore {
	mock := &MockPersistentVolumeStore{ctrl: ctrl}
	mock.recorder = &MockPersistentVolumeStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockPersistentVolumeStore) EXPECT() *MockPersistentVolumeStoreMockRecorder {
	return _m.recorder
}

// CreatePersistentVolume mocks base method
func (_m *MockPersistentVolumeStore) CreatePersistentVolume(_param0 context.Context, _param1 *volume.PersistentVolumeInfo) error {
	ret := _m.ctrl.Call(_m, "CreatePersistentVolume", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePersistentVolume indicates an expected call of CreatePersistentVolume
func (_mr *MockPersistentVolumeStoreMockRecorder) CreatePersistentVolume(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "CreatePersistentVolume", reflect.TypeOf((*MockPersistentVolumeStore)(nil).CreatePersistentVolume), arg0, arg1)
}

// DeletePersistentVolume mocks base method
func (_m *MockPersistentVolumeStore) DeletePersistentVolume(_param0 context.Context, _param1 *peloton.VolumeID) error {
	ret := _m.ctrl.Call(_m, "DeletePersistentVolume", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePersistentVolume indicates an expected call of DeletePersistentVolume
func (_mr *MockPersistentVolumeStoreMockRecorder) DeletePersistentVolume(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "DeletePersistentVolume", reflect.TypeOf((*MockPersistentVolumeStore)(nil).DeletePersistentVolume), arg0, arg1)
}

// GetPersistentVolume mocks base method
func (_m *MockPersistentVolumeStore) GetPersistentVolume(_param0 context.Context, _param1 *peloton.VolumeID) (*volume.PersistentVolumeInfo, error) {
	ret := _m.ctrl.Call(_m, "GetPersistentVolume", _param0, _param1)
	ret0, _ := ret[0].(*volume.PersistentVolumeInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPersistentVolume indicates an expected call of GetPersistentVolume
func (_mr *MockPersistentVolumeStoreMockRecorder) GetPersistentVolume(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetPersistentVolume", reflect.TypeOf((*MockPersistentVolumeStore)(nil).GetPersistentVolume), arg0, arg1)
}

// UpdatePersistentVolume mocks base method
func (_m *MockPersistentVolumeStore) UpdatePersistentVolume(_param0 context.Context, _param1 *volume.PersistentVolumeInfo) error {
	ret := _m.ctrl.Call(_m, "UpdatePersistentVolume", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePersistentVolume indicates an expected call of UpdatePersistentVolume
func (_mr *MockPersistentVolumeStoreMockRecorder) UpdatePersistentVolume(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "UpdatePersistentVolume", reflect.TypeOf((*MockPersistentVolumeStore)(nil).UpdatePersistentVolume), arg0, arg1)
}
