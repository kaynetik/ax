package ax

// TestCase - base case struct used in testing Suite.
type TestCase struct {
	Name          string
	PreRequisites func()
	Assert        func()

	SubTests []TestCase
}

// SuiteRunner - interface for running testing Suite.
type SuiteRunner interface {
	Run(name string, subtest func()) bool
}

// RunTestCases - func executing the passed testing Suite.
func RunTestCases(s SuiteRunner, testCases []TestCase) {
	if len(testCases) == 0 {
		return
	}

	for _, tc := range testCases {
		tcInit := tc

		s.Run(tcInit.Name, func() {
			if tcInit.PreRequisites != nil {
				tcInit.PreRequisites()
			}

			if tcInit.Assert != nil {
				tcInit.Assert()
			}

			RunTestCases(s, tcInit.SubTests)
		})
	}
}
