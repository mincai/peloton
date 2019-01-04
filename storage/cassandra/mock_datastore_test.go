// +build !unit

package cassandra

import (
	"context"
	"testing"
	"time"

	"github.com/uber/peloton/.gen/peloton/api/v0/job"
	"github.com/uber/peloton/.gen/peloton/api/v0/peloton"
	"github.com/uber/peloton/.gen/peloton/api/v0/task"
	"github.com/uber/peloton/.gen/peloton/api/v0/update"
	"github.com/uber/peloton/.gen/peloton/private/models"
	"github.com/uber/peloton/common"
	"github.com/uber/peloton/storage"
	datastore "github.com/uber/peloton/storage/cassandra/api"
	datastoremocks "github.com/uber/peloton/storage/cassandra/api/mocks"
	datastoreimpl "github.com/uber/peloton/storage/cassandra/impl"

	"github.com/golang/mock/gomock"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

const (
	testJobName  = "uber"
	testJob      = "941ff353-ba82-49fe-8f80-fb5bc649b04d"
	testUpdateID = "141ff353-ba82-49fe-8f80-fb5bc649b042"
)

type MockDatastoreTestSuite struct {
	suite.Suite
	testJobID *peloton.JobID

	ctrl            *gomock.Controller
	mockedDataStore *datastoremocks.MockDataStore
	store           *Store
}

func (suite *MockDatastoreTestSuite) SetupTest() {
	var result datastore.ResultSet

	suite.testJobID = &peloton.JobID{Value: testJob}

	suite.ctrl = gomock.NewController(suite.T())
	suite.mockedDataStore = datastoremocks.NewMockDataStore(suite.ctrl)

	suite.store = &Store{
		DataStore:   suite.mockedDataStore,
		metrics:     storage.NewMetrics(testScope.SubScope("storage")),
		Conf:        &Config{},
		retryPolicy: nil,
	}

	queryBuilder := &datastoreimpl.QueryBuilder{}
	// Mock datastore execute to fail
	suite.mockedDataStore.EXPECT().Execute(gomock.Any(), gomock.Any()).
		Return(result, errors.New("my-error")).AnyTimes()
	suite.mockedDataStore.EXPECT().NewQuery().Return(queryBuilder).AnyTimes()
}

func TestMockDatastoreTestSuite(t *testing.T) {
	suite.Run(t, new(MockDatastoreTestSuite))
}

// TestDataStoreDeleteJob test delete job
func (suite *MockDatastoreTestSuite) TestDataStoreDeleteJob() {
	ctx, cancelFunc := context.WithTimeout(
		context.Background(),
		time.Second)
	defer cancelFunc()
	jobID := &peloton.JobID{
		Value: uuid.New(),
	}

	// Failure test for GetJobConfig
	suite.mockedDataStore.EXPECT().Execute(ctx, gomock.Any()).
		Return(nil, errors.New("my-error"))
	suite.Error(suite.store.DeleteJob(ctx, jobID))
}

// TestDataStoreFailureGetJobConfig tests datastore failures in getting job cfg
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetJobConfig() {
	_, _, err := suite.store.GetJobConfigWithVersion(
		context.Background(), suite.testJobID, 0)
	suite.Error(err)

	_, _, err = suite.store.GetJobConfig(
		context.Background(), suite.testJobID)
	suite.Error(err)

}

// TestDataStoreFailureGetJobRuntime tests datastore failures in getting
// job runtime
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetJobRuntime() {
	suite.Error(suite.store.CreateJobRuntimeWithConfig(
		context.Background(),
		suite.testJobID,
		&job.RuntimeInfo{},
		&job.JobConfig{}))

	_, err := suite.store.GetJobRuntime(
		context.Background(), suite.testJobID)
	suite.Error(err)
}

// TestGetFrameworkID tests the fetch for framework ID
func (suite *MockDatastoreTestSuite) TestGetFrameworkID() {
	_, err := suite.store.GetFrameworkID(context.Background(), common.PelotonRole)
	suite.Error(err)
}

// TestGetStreamID test the fetch for stream ID
func (suite *MockDatastoreTestSuite) TestGetStreamID() {
	_, err := suite.store.GetMesosStreamID(context.Background(), common.PelotonRole)
	suite.Error(err)
}

// TestDataStoreFailureGetJobSummary tests datastore failures in getting
// job summary
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetJobSummary() {
	_, err := suite.store.GetJobSummaryFromIndex(
		context.Background(), suite.testJobID)
	suite.Error(err)

	_, err = suite.store.getJobSummaryFromConfig(
		context.Background(), suite.testJobID)
	suite.Error(err)
}

// TestDataStoreFailureGetJob tests datastore failures in getting job
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetJob() {
	_, err := suite.store.GetJobsByStates(
		context.Background(), []job.JobState{job.JobState_RUNNING})
	suite.Error(err)

	_, err = suite.store.GetMaxJobConfigVersion(
		context.Background(), suite.testJobID)
	suite.Error(err)
}

// TestDataStoreFailureGetTasks tests datastore failures in getting tasks
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetTasks() {
	_, err := suite.store.GetTasksForJobAndStates(
		context.Background(), suite.testJobID, []task.TaskState{
			task.TaskState(task.TaskState_PENDING)})
	suite.Error(err)

	_, err = suite.store.GetTasksForJobResultSet(
		context.Background(), suite.testJobID)
	suite.Error(err)

	_, err = suite.store.GetTasksForJob(
		context.Background(), suite.testJobID)
	suite.Error(err)

	_, err = suite.store.GetTaskForJob(
		context.Background(), suite.testJobID.GetValue(), 0)
	suite.Error(err)

	_, err = suite.store.GetTaskIDsForJobAndState(
		context.Background(), suite.testJobID, task.TaskState_PENDING.String())
	suite.Error(err)

	_, err = suite.store.getTaskStateCount(
		context.Background(), suite.testJobID, task.TaskState_PENDING.String())
	suite.Error(err)

	_, err = suite.store.getTask(context.Background(), testJob, 0)
	suite.Error(err)
}

// TestDataStoreFailureGetTaskConfig tests datastore failures in getting task cfg
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetTaskConfig() {
	_, _, err := suite.store.GetTaskConfig(
		context.Background(), suite.testJobID, 0, 0)
	suite.Error(err)

	_, _, err = suite.store.GetTaskConfigs(
		context.Background(), suite.testJobID, []uint32{0}, 0)
	suite.Error(err)

	_, err = suite.store.GetTaskStateSummaryForJob(
		context.Background(), suite.testJobID)
	suite.Error(err)
}

// TestDataStoreFailureGetTaskRuntime tests datastore failures in getting
// task runtime
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetTaskRuntime() {
	_, err := suite.store.GetTaskRuntimesForJobByRange(
		context.Background(), suite.testJobID, &task.InstanceRange{
			From: uint32(0),
			To:   uint32(3),
		})
	suite.Error(err)

	_, err = suite.store.GetTaskRuntime(
		context.Background(), suite.testJobID, 0)
	suite.Error(err)

	_, err = suite.store.getTaskRuntimeRecord(context.Background(), testJob, 0)
	suite.Error(err)
}

// TestDataStoreFailureJobQuery tests datastore failures in job query
func (suite *MockDatastoreTestSuite) TestDataStoreFailureJobQuery() {
	_, _, _, err := suite.store.QueryJobs(
		context.Background(), nil, &job.QuerySpec{}, false)
	suite.Error(err)
}

// TestDataStoreFailureTaskQuery tests datastore failures in task query
func (suite *MockDatastoreTestSuite) TestDataStoreFailureTaskQuery() {
	_, _, err := suite.store.QueryTasks(
		context.Background(), suite.testJobID, &task.QuerySpec{})
	suite.Error(err)
}

// TestDataStoreFailureGetRespools tests datastore failures in get respools
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetRespools() {
	_, err := suite.store.GetAllResourcePools(context.Background())
	suite.Error(err)

	_, err = suite.store.GetResourcePoolsByOwner(context.Background(), "dummy")
	suite.Error(err)
}

// TestDataStoreFailureFramework tests datastore failures in get frameworks
func (suite *MockDatastoreTestSuite) TestDataStoreFailureFramework() {
	_, err := suite.store.getFrameworkInfo(context.Background(), "framwork-id")
	suite.Error(err)
}

// TestDataStoreFailureGetPersistentVolume tests datastore failures in get
// persistent volume
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetPersistentVolume() {
	_, err := suite.store.GetPersistentVolume(
		context.Background(), &peloton.VolumeID{Value: "test"})
	suite.Error(err)
}

// TestDataStoreFailureGetSecret tests datastore failures in get secret
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetSecret() {
	_, err := suite.store.GetSecret(
		context.Background(), &peloton.SecretID{Value: "test"})
	suite.Error(err)
}

// TestDataStoreFailureGetUpdate tests datastore failures in get update
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetUpdate() {
	_, err := suite.store.GetUpdate(
		context.Background(), &peloton.UpdateID{Value: "test"})
	suite.Error(err)

	_, err = suite.store.GetUpdateProgress(
		context.Background(), &peloton.UpdateID{Value: "test"})
	suite.Error(err)

	_, err = suite.store.GetUpdatesForJob(
		context.Background(),
		suite.testJobID.GetValue())
	suite.Error(err)
}

// TestDataStoreFailureDeleteJobCfgVersion tests datastore failures in delete
// job config version
func (suite *MockDatastoreTestSuite) TestDataStoreFailureDeleteJobCfgVersion() {
	ctx := context.Background()
	var result datastore.ResultSet

	// Setup mocks for this context

	// Simulate failure to delete task config
	suite.mockedDataStore.EXPECT().Execute(ctx, gomock.Any()).
		Return(result, errors.New("my-error"))

	err := suite.store.deleteJobConfigVersion(ctx, suite.testJobID, 0)
	suite.Error(err)

	// Simulate success to to delete task cfg and failure to delete job cfg
	suite.mockedDataStore.EXPECT().Execute(ctx, gomock.Any()).
		Return(result, nil)
	suite.mockedDataStore.EXPECT().Execute(ctx, gomock.Any()).
		Return(result, errors.New("my-error"))

	err = suite.store.deleteJobConfigVersion(ctx, suite.testJobID, 0)
	suite.Error(err)
}

// TestDataStoreFailureActiveJobs tests datastore failures add/get/delete jobID
// from active jobs
func (suite *MockDatastoreTestSuite) TestDataStoreFailureActiveJobs() {
	err := suite.store.AddActiveJob(context.Background(), suite.testJobID)
	suite.Error(err)

	_, err = suite.store.GetActiveJobs(context.Background())
	suite.Error(err)

	err = suite.store.DeleteActiveJob(context.Background(), suite.testJobID)
	suite.Error(err)
}

// TestJobNameToIDMapFailures tests failure scenarios for job name to job uuid
func (suite *MockDatastoreTestSuite) TestJobNameToIDMapFailures() {
	jobID := &peloton.JobID{
		Value: testJob,
	}

	jobConfig := &job.JobConfig{
		Name: testJobName,
		Type: job.JobType_SERVICE,
	}
	err := suite.store.addJobNameToJobIDMapping(context.Background(), jobID, jobConfig)
	suite.Error(err)

	_, err = suite.store.GetJobIDFromJobName(context.Background(), testJobName)
	suite.Error(err)
}

// TestCreateTaskConfigFailures tests failure scenarios for create task configs
func (suite *MockDatastoreTestSuite) TestCreateTaskConfigFailures() {

	jobID := &peloton.JobID{
		Value: testJob,
	}

	err := suite.store.CreateTaskConfig(context.Background(), jobID,
		0, &task.TaskConfig{Name: "dummy-task"}, nil, 0)
	suite.Error(err)
}

// TestWorkflowEventsFailures tests failure scenarios for workflow events
func (suite *MockDatastoreTestSuite) TestWorkflowEventsFailures() {
	updateID := &peloton.UpdateID{
		Value: testUpdateID,
	}

	err := suite.store.AddWorkflowEvent(
		context.Background(),
		updateID,
		0,
		models.WorkflowType_UPDATE,
		update.State_ROLLING_FORWARD)
	suite.Error(err)

	err = suite.store.deleteWorkflowEvents(context.Background(), updateID, 0)
	suite.Error(err)

	_, err = suite.store.GetWorkflowEvents(context.Background(), updateID, 0)
	suite.Error(err)
}

// TestDataStoreFailureGetTaskConfigs tests datastore failures in get task
// config from legacy/v2 tables
func (suite *MockDatastoreTestSuite) TestDataStoreFailureGetTaskConfigs() {
	ctx := context.Background()
	var result datastore.ResultSet

	// Setup mocks for this context
	// Simulate success for the first query and failure for the second query
	suite.mockedDataStore.EXPECT().Execute(ctx, gomock.Any()).
		Return(result, nil)
	suite.mockedDataStore.EXPECT().Execute(ctx, gomock.Any()).
		Return(result, errors.New("my-error"))
	_, _, err := suite.store.GetTaskConfig(ctx, suite.testJobID, uint32(0), 0)
	suite.Error(err)

	suite.mockedDataStore.EXPECT().Execute(ctx, gomock.Any()).
		Return(result, nil)
	suite.mockedDataStore.EXPECT().Execute(ctx, gomock.Any()).
		Return(result, errors.New("my-error"))
	_, _, err = suite.store.GetTaskConfigs(ctx, suite.testJobID, []uint32{0}, 0)
	suite.Error(err)
}
