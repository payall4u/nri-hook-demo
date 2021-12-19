package config

import (
	"encoding/json"

	v1 "github.com/containerd/nri/types/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

type NRIConfig struct {
	MatchAll       bool                   `json:"MatchAll,omitempty"`
	LabelSelectors []metav1.LabelSelector `json:"LabelSelectors,omitempty"`
	States         []string               `json:"States"`
	Address        string                 `json:"Address"`
	URI            string                 `json:"URI"`

	selectors []labels.Selector
}

func FromRequest(request *v1.Request) (nc NRIConfig, err error) {
	if err := json.Unmarshal(request.Conf, &nc); err != nil {
		return nc, err
	}
	for _, ls := range nc.LabelSelectors {
		selector, err := metav1.LabelSelectorAsSelector(&ls)
		if err != nil {
			return nc, err
		}
		nc.selectors = append(nc.selectors, selector)
	}
	return nc, nil
}

func (nc *NRIConfig) InterestedState(state string) bool {
	for _, s := range nc.States {
		if s == state {
			return true
		}
	}
	return false
}

func (nc *NRIConfig) InterestedLabels(labelMap map[string]string) bool {
	if nc.MatchAll {
		return true
	}
	for _, selector := range nc.selectors {
		if selector.Matches(labels.Set(labelMap)) {
			return true
		}
	}
	return false
}
