/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqltypes

import (
	"google.golang.org/protobuf/proto"

	"github.com/vedadiyan/sqlparser/pkg/vterrors"

	querypb "github.com/vedadiyan/sqlparser/pkg/query"
)

// This file contains the proto3 conversion functions for the structures
// defined here.

// RowToProto3 converts []Value to proto3.
func RowToProto3(row []Value) *querypb.Row {
	result := &querypb.Row{}
	_ = RowToProto3Inplace(row, result)
	return result
}

// RowToProto3Inplace converts []Value to proto3 and stores the conversion in the provided Row
func RowToProto3Inplace(row []Value, result *querypb.Row) int {
	if result.Lengths == nil {
		result.Lengths = make([]int64, 0, len(row))
	} else {
		result.Lengths = result.Lengths[:0]
	}
	total := 0
	for _, c := range row {
		if c.IsNull() {
			result.Lengths = append(result.Lengths, -1)
			continue
		}
		length := c.Len()
		result.Lengths = append(result.Lengths, int64(length))
		total += length
	}
	if result.Values == nil {
		result.Values = make([]byte, 0, total)
	} else {
		result.Values = result.Values[:0]
	}
	for _, c := range row {
		if c.IsNull() {
			continue
		}
		result.Values = append(result.Values, c.Raw()...)
	}
	return total
}

// RowsToProto3 converts [][]Value to proto3.
func RowsToProto3(rows [][]Value) []*querypb.Row {
	if len(rows) == 0 {
		return nil
	}

	result := make([]*querypb.Row, len(rows))
	for i, r := range rows {
		result[i] = RowToProto3(r)
	}
	return result
}

// proto3ToRows converts a proto3 rows to [][]Value. The function is private
// because it uses the trusted API.
func proto3ToRows(fields []*querypb.Field, rows []*querypb.Row) [][]Value {
	if len(rows) == 0 {
		// TODO(sougou): This is needed for backward compatibility.
		// Remove when it's not needed any more.
		return [][]Value{}
	}

	result := make([][]Value, len(rows))
	for i, r := range rows {
		result[i] = MakeRowTrusted(fields, r)
	}
	return result
}

// ResultToProto3 converts Result to proto3.
func ResultToProto3(qr *Result) *querypb.QueryResult {
	if qr == nil {
		return nil
	}
	return &querypb.QueryResult{
		Fields:              qr.Fields,
		RowsAffected:        qr.RowsAffected,
		InsertId:            qr.InsertID,
		InsertIdChanged:     qr.InsertIDChanged,
		Rows:                RowsToProto3(qr.Rows),
		Info:                qr.Info,
		SessionStateChanges: qr.SessionStateChanges,
	}
}

// Proto3ToResult converts a proto3 Result to an internal data structure. This function
// should be used only if the field info is populated in qr.
func Proto3ToResult(qr *querypb.QueryResult) *Result {
	if qr == nil {
		return nil
	}
	return &Result{
		Fields:              qr.Fields,
		RowsAffected:        qr.RowsAffected,
		InsertID:            qr.InsertId,
		InsertIDChanged:     qr.InsertIdChanged,
		Rows:                proto3ToRows(qr.Fields, qr.Rows),
		Info:                qr.Info,
		SessionStateChanges: qr.SessionStateChanges,
	}
}

// CustomProto3ToResult converts a proto3 Result to an internal data structure. This function
// takes a separate fields input because not all QueryResults contain the field info.
// In particular, only the first packet of streaming queries contain the field info.
func CustomProto3ToResult(fields []*querypb.Field, qr *querypb.QueryResult) *Result {
	if qr == nil {
		return nil
	}
	return &Result{
		Fields:              qr.Fields,
		RowsAffected:        qr.RowsAffected,
		InsertID:            qr.InsertId,
		InsertIDChanged:     qr.InsertIdChanged,
		Rows:                proto3ToRows(fields, qr.Rows),
		Info:                qr.Info,
		SessionStateChanges: qr.SessionStateChanges,
	}
}

// ResultsToProto3 converts []Result to proto3.
func ResultsToProto3(qr []*Result) []*querypb.QueryResult {
	if len(qr) == 0 {
		return nil
	}
	result := make([]*querypb.QueryResult, len(qr))
	for i, q := range qr {
		result[i] = ResultToProto3(q)
	}
	return result
}

// Proto3ToResults converts proto3 results to []Result.
func Proto3ToResults(qr []*querypb.QueryResult) []*Result {
	if len(qr) == 0 {
		return nil
	}
	result := make([]*Result, len(qr))
	for i, q := range qr {
		result[i] = Proto3ToResult(q)
	}
	return result
}

// QueryResponseToProto3 converts QueryResponse to proto3.
func QueryResponseToProto3(qr QueryResponse) *querypb.ResultWithError {
	return &querypb.ResultWithError{
		Result: ResultToProto3(qr.QueryResult),
		Error:  vterrors.ToVTRPC(qr.QueryError),
	}
}

// QueryResponsesToProto3 converts []QueryResponse to proto3.
func QueryResponsesToProto3(qr []QueryResponse) []*querypb.ResultWithError {
	if len(qr) == 0 {
		return nil
	}
	result := make([]*querypb.ResultWithError, len(qr))
	for i, q := range qr {
		result[i] = QueryResponseToProto3(q)
	}
	return result
}

// Proto3ToQueryReponses converts proto3 queryResponse to []QueryResponse.
func Proto3ToQueryReponses(qr []*querypb.ResultWithError) []QueryResponse {
	if len(qr) == 0 {
		return nil
	}
	result := make([]QueryResponse, len(qr))
	for i, q := range qr {
		result[i] = QueryResponse{
			QueryResult: Proto3ToResult(q.Result),
			QueryError:  vterrors.FromVTRPC(q.Error),
		}
	}
	return result
}

// Proto3ResultsEqual compares two arrays of proto3 Result.
// reflect.DeepEqual shouldn't be used because of the protos.
func Proto3ResultsEqual(r1, r2 []*querypb.QueryResult) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i, r := range r1 {
		if !proto.Equal(r, r2[i]) {
			return false
		}
	}
	return true
}

// Proto3QueryResponsesEqual compares two arrays of proto3 QueryResponse.
// reflect.DeepEqual shouldn't be used because of the protos.
func Proto3QueryResponsesEqual(r1, r2 []*querypb.ResultWithError) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i, r := range r1 {
		if !proto.Equal(r, r2[i]) {
			return false
		}
	}
	return true
}

// Proto3ValuesEqual compares two arrays of proto3 Value.
func Proto3ValuesEqual(v1, v2 []*querypb.Value) bool {
	if len(v1) != len(v2) {
		return false
	}
	for i, v := range v1 {
		if !proto.Equal(v, v2[i]) {
			return false
		}
	}
	return true
}
