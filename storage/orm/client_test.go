package orm

import (
	"context"
	"testing"

	"github.com/uber/peloton/storage/objects/base"
	ormmocks "github.com/uber/peloton/storage/orm/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ORMTestSuite struct {
	suite.Suite

	ctrl *gomock.Controller
	ctx  context.Context
}

// testRow in DB representation looks like this:
//
// 	*	id	| 	 name	|  	 data
//  1.   1  	"test"	   "testdata"
var testRow = []base.Column{
	{
		Name:  "id",
		Value: uint64(1),
	},
	{
		Name:  "name",
		Value: "test",
	},
	{
		Name:  "data",
		Value: "testdata",
	},
}

// keyRow in DB representation looks like this:
//
// 	*	id	| 	 name
//  1.   1  	"test"
var keyRow = []base.Column{
	{
		Name:  "id",
		Value: uint64(1),
	},
	{
		Name:  "name",
		Value: "test",
	},
}

// testValidObject is the storage object representation of the above row
var testValidObject = &ValidObject{
	ID:   uint64(1),
	Name: "test",
	Data: "testdata",
}

// test if both rows are equal even if out of order
func (suite *ORMTestSuite) ensureRowsEqual(row1, row2 []base.Column) {

	suite.Equal(len(row1), len(row2))
	rowMap := make(map[string]interface{})
	for _, col1 := range row1 {
		rowMap[col1.Name] = col1.Value
	}

	for _, col2 := range row2 {
		val1, ok := rowMap[col2.Name]
		suite.True(ok)
		suite.Equal(val1, col2.Value)
	}
}

func (suite *ORMTestSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.ctx = context.Background()
}

func TestORMTestSuite(t *testing.T) {
	suite.Run(t, new(ORMTestSuite))
}

// TestNewClientc tests creating new base client with base objects
func (suite *ORMTestSuite) TestNewClient() {
	defer suite.ctrl.Finish()
	conn := ormmocks.NewMockConnector(suite.ctrl)
	_, err := NewClient(conn, &ValidObject{})
	suite.NoError(err)

	_, err = NewClient(conn, &InvalidObject1{})
	suite.Error(err)
}

// TestClientCreate tests client create operation on valid and invalid entities
func (suite *ORMTestSuite) TestClientCreate() {
	defer suite.ctrl.Finish()
	conn := ormmocks.NewMockConnector(suite.ctrl)

	conn.EXPECT().Create(suite.ctx, gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, _ *base.Definition, row []base.Column) {
			suite.ensureRowsEqual(row, testRow)
		}).Return(nil)

	client, err := NewClient(conn, &ValidObject{})
	suite.NoError(err)

	err = client.Create(suite.ctx, testValidObject)
	suite.NoError(err)

	err = client.Create(suite.ctx, &InvalidObject1{})
	suite.Error(err)
}

// TestClientGet tests client get operation on valid and invalid entities
func (suite *ORMTestSuite) TestClientGet() {
	defer suite.ctrl.Finish()
	conn := ormmocks.NewMockConnector(suite.ctrl)

	// ValidObject instance with only primary key set
	e := &ValidObject{
		ID: uint64(1),
	}

	conn.EXPECT().Get(suite.ctx, gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, _ *base.Definition,
			row []base.Column) {
			suite.Equal("id", row[0].Name)
			suite.Equal(e.ID, row[0].Value)
		}).Return(testRow, nil)

	client, err := NewClient(conn, &ValidObject{})
	suite.NoError(err)

	// Do a get on the ValidObject instance and verify that the expected
	// fields in the object are set as per testRow
	err = client.Get(suite.ctx, e)
	suite.NoError(err)

	// compare the values from testRow to that of the entity fields
	suite.Equal(testRow[1].Value, e.Name)
	suite.Equal(testRow[2].Value, e.Data)

	err = client.Get(suite.ctx, &InvalidObject1{})
	suite.Error(err)
}

// TestClientUpdate tests client update operation on valid and invalid entities
func (suite *ORMTestSuite) TestClientUpdate() {
	defer suite.ctrl.Finish()
	conn := ormmocks.NewMockConnector(suite.ctrl)

	conn.EXPECT().Update(suite.ctx, gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, _ *base.Definition,
			row []base.Column, keyRow []base.Column) {
			if "name" == row[0].Name {
				suite.Equal("test", row[0].Value)
				suite.Equal("testdata", row[1].Value)
			} else {
				suite.Equal("testdata", row[0].Value)
				suite.Equal("test", row[1].Value)
			}
			suite.Equal("id", keyRow[0].Name)
			suite.Equal(uint64(1), keyRow[0].Value)
		}).Return(nil)

	client, err := NewClient(conn, &ValidObject{})
	suite.NoError(err)

	// Update Data and Name fields in testValidObject
	err = client.Update(suite.ctx, testValidObject, "Name", "Data")
	suite.NoError(err)

	conn.EXPECT().Update(suite.ctx, gomock.Any(), gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, _ *base.Definition,
			row []base.Column, keyRow []base.Column) {
			suite.Equal("data", row[0].Name)
			suite.Equal("testdata", row[0].Value)
			suite.Equal("id", keyRow[0].Name)
			suite.Equal(uint64(1), keyRow[0].Value)
		}).Return(nil)
	// Update Data field in testValidObject
	err = client.Update(suite.ctx, testValidObject, "Data")
	suite.NoError(err)

	err = client.Update(suite.ctx, &InvalidObject1{})
	suite.Error(err)
}

// TestClientDelete tests client delete operation on valid and invalid entities
func (suite *ORMTestSuite) TestClientDelete() {
	defer suite.ctrl.Finish()
	conn := ormmocks.NewMockConnector(suite.ctrl)

	conn.EXPECT().Delete(suite.ctx, gomock.Any(), gomock.Any()).
		Do(func(_ context.Context, _ *base.Definition, row []base.Column) {
			suite.ensureRowsEqual(row, keyRow)
		}).Return(nil)

	client, err := NewClient(conn, &ValidObject{})
	suite.NoError(err)

	err = client.Delete(suite.ctx, testValidObject)
	suite.NoError(err)

	err = client.Delete(suite.ctx, &InvalidObject1{})
	suite.Error(err)
}
