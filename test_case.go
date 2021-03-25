package ax

type TestCase struct {
	Name          string
	PreRequisites func()
	Assert        func()

	SubTests []TestCase
}

type SuiteRunner interface {
	Run(name string, subtest func()) bool
}

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
