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

package xgboost

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	pb "sqlflow.org/sqlflow/pkg/proto"
	"sqlflow.org/sqlflow/pkg/sql/ir"
)

func TestTrainAndPredict(t *testing.T) {
	a := assert.New(t)
	tir := mockTrainIR()
	_, err := Train(tir)
	a.NoError(err)

	pir := mockPrdcIR(tir)
	sess := &pb.Session{
		Token:            "",
		DbConnStr:        "",
		ExitOnSubmit:     false,
		UserId:           "",
		HiveLocation:     "/sqlflowtmp",
		HdfsNamenodeAddr: "192.168.1.1:8020",
		HdfsUser:         "sqlflow_admin",
		HdfsPass:         "sqlflow_pass",
	}
	code, err := Pred(pir, sess)

	r, _ := regexp.Compile(`hdfs_user='''(.*)'''`)
	a.Equal(r.FindStringSubmatch(code)[1], "sqlflow_admin")
	r, _ = regexp.Compile(`hdfs_pass='''(.*)'''`)
	a.Equal(r.FindStringSubmatch(code)[1], "sqlflow_pass")

	a.NoError(err)
}

func mockPrdcIR(trainIR *ir.TrainClause) *ir.PredictClause {
	return &ir.PredictClause{
		DataSource:  trainIR.DataSource,
		Select:      "select * from iris.test;",
		ResultTable: "iris.predict",
		TrainIR:     trainIR,
	}
}
func mockTrainIR() *ir.TrainClause {
	cfg := &mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		AllowNativePasswords: true,
	}
	_ = `SELECT *
		FROM iris.train
	TO TRAIN xgboost.gbtree
	WITH
		objective = "multi:softprob"
		eta = 3.1,
		num_class = 3,
		train.num_boost_round = 30
	COLUMN sepal_length, sepal_width, petal_length, petal_width
	LABEL class
	INTO sqlflow_models.my_xgboost_model;`
	return &ir.TrainClause{
		DataSource:       fmt.Sprintf("mysql://%s", cfg.FormatDSN()),
		Select:           "select * from iris.train;",
		ValidationSelect: "select * from iris.test;",
		Estimator:        "xgboost.gbtree",
		Attributes: map[string]interface{}{
			"train.num_boost_round": 10,
			"objective":             "multi:softprob",
			"eta":                   float32(0.1),
			"num_class":             3},
		Features: map[string][]ir.FeatureColumn{
			"feature_columns": {
				&ir.NumericColumn{&ir.FieldMeta{"sepal_length", ir.Float, "", []int{1}, false, nil, 0}},
				&ir.NumericColumn{&ir.FieldMeta{"sepal_width", ir.Float, "", []int{1}, false, nil, 0}},
				&ir.NumericColumn{&ir.FieldMeta{"petal_length", ir.Float, "", []int{1}, false, nil, 0}},
				&ir.NumericColumn{&ir.FieldMeta{"petal_width", ir.Float, "", []int{1}, false, nil, 0}}}},
		Label: &ir.NumericColumn{&ir.FieldMeta{"class", ir.Int, "", []int{1}, false, nil, 0}}}
}
