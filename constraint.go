package bento

import "fmt"

type Constraint interface {
	fmt.Stringer

	isConstraint()
}

type (
	ConstraintMin        int
	ConstraintMax        int
	ConstraintLen        int
	ConstraintPercentage int
	ConstraintRatio      struct{ Num, Den int }
	ConstraintFill       int
)

func (m ConstraintMin) String() string { return fmt.Sprintf("Min(%d)", m) }
func (ConstraintMin) isConstraint()    {}

func (m ConstraintMax) String() string { return fmt.Sprintf("Max(%d)", m) }
func (ConstraintMax) isConstraint()    {}

func (l ConstraintLen) String() string { return fmt.Sprintf("Len(%d)", l) }
func (ConstraintLen) isConstraint()    {}

func (p ConstraintPercentage) String() string { return fmt.Sprintf("Percentage(%d)", p) }
func (ConstraintPercentage) isConstraint()    {}

func (r ConstraintRatio) String() string { return fmt.Sprintf("Ratio(%d / %d)", r.Num, r.Den) }
func (ConstraintRatio) isConstraint()    {}

func (f ConstraintFill) String() string { return fmt.Sprintf("Fill(%d)", f) }
func (ConstraintFill) isConstraint()    {}
