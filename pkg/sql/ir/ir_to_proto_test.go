// Copyright 2019 The SQLFlow Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ir

import (
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/stretchr/testify/assert"
	pb "sqlflow.org/sqlflow/pkg/proto"
)

func mockTrainIR() *TrainClause {
	return &TrainClause{
		DataSource:       "mysql://root:root@tcp(127.0.0.1:3306)/?maxAllowedPacket=0",
		Select:           "select * from iris.train;",
		ValidationSelect: "select * from iris.test;",
		Estimator:        "DNNClassifier",
		Attributes: map[string]interface{}{
			"train.batch_size":   4,
			"train.epoch":        3,
			"model.hidden_units": []int{10, 20},
			"model.n_classes":    3},
		Features: map[string][]FeatureColumn{
			"feature_columns": {
				&NumericColumn{&FieldMeta{"sepal_length", Float, "", []int{1}, false, nil, 0}},
				&NumericColumn{&FieldMeta{"sepal_width", Float, "", []int{1}, false, nil, 0}},
				&NumericColumn{&FieldMeta{"petal_length", Float, "", []int{1}, false, nil, 0}},
				&NumericColumn{&FieldMeta{"petal_width", Float, "", []int{1}, false, nil, 0}}}},
		Label: &NumericColumn{&FieldMeta{"class", Int, "", []int{1}, false, nil, 0}}}
}

func mockSession() *pb.Session {
	return &pb.Session{
		Token:            "",
		DbConnStr:        "",
		ExitOnSubmit:     false,
		UserId:           "",
		HiveLocation:     "/sqlflowtmp",
		HdfsNamenodeAddr: "192.168.1.1:8020",
		HdfsUser:         "sqlflow_admin",
		HdfsPass:         "sqlflow_pass",
	}
}

func TestTrainProto(t *testing.T) {
	a := assert.New(t)
	sampleTrainIR := mockTrainIR()
	pbIR, err := TrainIRToProto(sampleTrainIR, mockSession())
	a.NoError(err)
	pbtxt := proto.MarshalTextString(pbIR)
	pbIRToTest := &pb.TrainClause{}
	err = proto.UnmarshalText(pbtxt, pbIRToTest)
	a.NoError(err)
	a.Equal(
		sampleTrainIR.Features["feature_columns"][2].GetFieldMeta()[0].Name,
		pbIRToTest.GetFeatures()["feature_columns"].GetFeatureColumns()[2].GetNc().GetFieldMeta().GetName(),
	)
	a.Equal(
		int32(sampleTrainIR.Attributes["train.batch_size"].(int)),
		pbIRToTest.GetAttributes()["train.batch_size"].GetI(),
	)
}

func TestPredictProto(t *testing.T) {
	a := assert.New(t)
	samplePredIR := &PredictClause{
		DataSource:   "mysql://root:root@tcp(127.0.0.1:3306)/?maxAllowedPacket=0",
		Select:       "select * from iris.test;",
		ResultTable:  "predict",
		ResultColumn: "class",
		Attributes:   make(map[string]interface{}), // empty attribute
		TrainIR:      mockTrainIR(),
	}
	pbIR, err := PredictIRToProto(samplePredIR, mockSession())
	a.NoError(err)
	pbtxt := proto.MarshalTextString(pbIR)
	pbIRToTest := &pb.PredictClause{}
	err = proto.UnmarshalText(pbtxt, pbIRToTest)
	a.NoError(err)
	a.Equal(
		samplePredIR.ResultTable,
		pbIRToTest.GetResultTable(),
	)
}

func TestAnalyzeProto(t *testing.T) {
	a := assert.New(t)
	sampleAnalyzeIR := &AnalyzeClause{
		DataSource: "mysql://root:root@tcp(127.0.0.1:3306)/?maxAllowedPacket=0",
		Select:     "select * from iris.train;",
		Attributes: make(map[string]interface{}), // empty attribute
		Explainer:  "TreeExplainer",
		TrainIR:    mockTrainIR(),
	}
	pbIR, err := AnalyzeIRToProto(sampleAnalyzeIR, mockSession())
	a.NoError(err)
	pbtxt := proto.MarshalTextString(pbIR)
	pbIRToTest := &pb.AnalyzeClause{}
	err = proto.UnmarshalText(pbtxt, pbIRToTest)
	a.NoError(err)
	a.Equal(
		sampleAnalyzeIR.Explainer,
		pbIRToTest.GetExplainer(),
	)
}
