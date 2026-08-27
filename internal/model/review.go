package model

import "time"

type ReviewDecision string

const (
	DecisionApprove ReviewDecision = "approve"
	DecisionReturn  ReviewDecision = "return"
)

type Review struct {
	Reviewer   string         `json:"reviewer"`
	Decision   ReviewDecision `json:"decision"`
	Note       string         `json:"note"`
	ReviewedAt time.Time      `json:"reviewed_at"`
}

func (r Review) Approved() bool { return r.Decision == DecisionApprove }
func (r Review) Valid() bool {
	return r.Reviewer != "" && (r.Decision == DecisionApprove || r.Decision == DecisionReturn)
}
