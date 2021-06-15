package commands

// // TestGitCommandGetCommitDifferences is a function.
// func TestGitCommandGetCommitDifferences(t *testing.T) {
// 	type scenario struct {
// 		testName string
// 		command  func(string, ...string) *exec.Cmd
// 		test     func(string, string)
// 	}

// 	scenarios := []scenario{
// 		{
// 			"Can't retrieve pushable count",
// 			func(string, ...string) *exec.Cmd {
// 				return secureexec.Command("test")
// 			},
// 			func(pushableCount string, pullableCount string) {
// 				assert.EqualValues(t, "?", pushableCount)
// 				assert.EqualValues(t, "?", pullableCount)
// 			},
// 		},
// 		{
// 			"Can't retrieve pullable count",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				if args[1] == "HEAD..@{u}" {
// 					return secureexec.Command("test")
// 				}

// 				return secureexec.Command("echo")
// 			},
// 			func(pushableCount string, pullableCount string) {
// 				assert.EqualValues(t, "?", pushableCount)
// 				assert.EqualValues(t, "?", pullableCount)
// 			},
// 		},
// 		{
// 			"Retrieve pullable and pushable count",
// 			func(cmd string, args ...string) *exec.Cmd {
// 				if args[1] == "HEAD..@{u}" {
// 					return secureexec.Command("echo", "10")
// 				}

// 				return secureexec.Command("echo", "11")
// 			},
// 			func(pushableCount string, pullableCount string) {
// 				assert.EqualValues(t, "11", pushableCount)
// 				assert.EqualValues(t, "10", pullableCount)
// 			},
// 		},
// 	}

// 	for _, s := range scenarios {
// 		t.Run(s.testName, func(t *testing.T) {
// 			gitCmd := NewDummyGit()
// 			gitCmd.GetOSCommand().Command = s.command
// 			s.test(gitCmd.GetCommitDifferences("HEAD", "@{u}"))
// 		})
// 	}
// }
