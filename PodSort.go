package main

import v1 "k8s.io/api/core/v1"

type byAge []v1.Pod

func (s byAge) Len() int {
	return len(s)
}

func (s byAge) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byAge) Less(i, j int) bool {
	if s[i].Status.StartTime == nil || s[j].Status.StartTime == nil {
		return false
	}
	return s[i].Status.StartTime.Time.UnixNano() < s[j].Status.StartTime.Time.UnixNano()
}
