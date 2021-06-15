package commands

// // TestGitCommandResetHard is a function.
// func TestGitCommandResetHard(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		ref      string
// 		command  func(string, ...string) *exec.Cmd
// 		test     func(error)
// 	}

// 	scenarios := []scenario{
// 		{
// 			"valid case",
// 			"HEAD",
// 			test.CreateMockCommand(t, []*test.CommandSwapper{
// 				{
// 					Expect:  `git reset --hard HEAD`,
// 					Replace: "echo",
// 				},
// 			}),
// 			func(err error) {
// 				assert.NoError(t, err)
// 			},
// 		},
// 	}

// 	gitCmd := NewDummyGit()

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd.GetOSCommand().Command = s.command
// 			s.test(gitCmd.ResetHard(s.ref))
// 		})
// 	}
// }
