package analyzer

import (
	"fmt"

	"github.com/k8sgpt-ai/k8sgpt/pkg/common"
	"github.com/k8sgpt-ai/k8sgpt/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PdbAnalyzer struct{}

func (PdbAnalyzer) Analyze(a common.Analyzer) ([]common.Result, error) {

	list, err := a.Client.GetClient().PolicyV1().PodDisruptionBudgets(a.Namespace).List(a.Context, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var preAnalysis = map[string]common.PreAnalysis{}

	for _, pdb := range list.Items {
		var failures []common.Failure

		evt, err := FetchLatestEvent(a.Context, a.Client, pdb.Namespace, pdb.Name)
		if err != nil || evt == nil {
			continue
		}

		if evt.Reason == "NoPods" && evt.Message != "" {
			if pdb.Spec.Selector != nil {
				for k, v := range pdb.Spec.Selector.MatchLabels {
					failures = append(failures, common.Failure{
						Text: fmt.Sprintf("%s, expected label %s=%s", evt.Message, k, v),
						Sensitive: []common.Sensitive{
							{
								Unmasked: k,
								Masked:   util.MaskString(k),
							},
							{
								Unmasked: v,
								Masked:   util.MaskString(v),
							},
						},
					})
				}
				for _, v := range pdb.Spec.Selector.MatchExpressions {
					failures = append(failures, common.Failure{
						Text:      fmt.Sprintf("%s, expected expression %s", evt.Message, v),
						Sensitive: []common.Sensitive{},
					})
				}
			} else {
				failures = append(failures, common.Failure{
					Text:      fmt.Sprintf("%s, selector is nil", evt.Message),
					Sensitive: []common.Sensitive{},
				})
			}
		}

		if len(failures) > 0 {
			preAnalysis[fmt.Sprintf("%s/%s", pdb.Namespace, pdb.Name)] = common.PreAnalysis{
				PodDisruptionBudget: pdb,
				FailureDetails:      failures,
			}
		}
	}

	for key, value := range preAnalysis {
		var currentAnalysis = common.Result{
			Kind:  "PodDisruptionBudget",
			Name:  key,
			Error: value.FailureDetails,
		}

		parent, _ := util.GetParent(a.Client, value.PodDisruptionBudget.ObjectMeta)
		currentAnalysis.ParentObject = parent
		a.Results = append(a.Results, currentAnalysis)
	}

	return a.Results, err
}
