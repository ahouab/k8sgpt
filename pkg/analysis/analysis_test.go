package analysis

import (
	"encoding/json"
	"fmt"
	"github.com/k8sgpt-ai/k8sgpt/pkg/analyzer"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAnalysis_NoProblemJsonOutput(t *testing.T) {

	analysis := Analysis{
		Results:   []analyzer.Result{},
		Namespace: "default",
	}

	expected := JsonOutput{
		Status:   StateOK,
		Problems: 0,
		Results:  []analyzer.Result{},
	}

	gotJson, err := analysis.JsonOutput()
	if err != nil {
		t.Error(err)
	}

	got := JsonOutput{}
	err = json.Unmarshal(gotJson, &got)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(got)
	fmt.Println(expected)

	require.Equal(t, got, expected)
}

func TestAnalysis_ProblemJsonOutput(t *testing.T) {
	analysis := Analysis{
		Results: []analyzer.Result{
			{
				"Deployment",
				"test-deployment",
				[]string{"test-problem"},
				"test-solution",
				"parent-resource"},
		},
		Namespace: "default",
	}

	expected := JsonOutput{
		Status:   StateProblemDetected,
		Problems: 1,
		Results: []analyzer.Result{
			{"Deployment",
				"test-deployment",
				[]string{"test-problem"},
				"test-solution",
				"parent-resource"},
		},
	}

	gotJson, err := analysis.JsonOutput()
	if err != nil {
		t.Error(err)
	}

	got := JsonOutput{}
	err = json.Unmarshal(gotJson, &got)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(got)
	fmt.Println(expected)

	require.Equal(t, got, expected)
}

func TestAnalysis_MultipleProblemJsonOutput(t *testing.T) {
	analysis := Analysis{
		Results: []analyzer.Result{
			{
				"Deployment",
				"test-deployment",
				[]string{"test-problem", "another-test-problem"},
				"test-solution",
				"parent-resource"},
		},
		Namespace: "default",
	}

	expected := JsonOutput{
		Status:   StateProblemDetected,
		Problems: 2,
		Results: []analyzer.Result{
			{"Deployment",
				"test-deployment",
				[]string{"test-problem", "another-test-problem"},
				"test-solution",
				"parent-resource"},
		},
	}

	gotJson, err := analysis.JsonOutput()
	if err != nil {
		t.Error(err)
	}

	got := JsonOutput{}
	err = json.Unmarshal(gotJson, &got)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(got)
	fmt.Println(expected)

	require.Equal(t, got, expected)
}
