package processcreds_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/processcreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/awstesting"
)

func TestProcessProviderFromSessionCfg(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	if runtime.GOOS == "windows" {
		os.Setenv("AWS_CONFIG_FILE", "testdata\\shconfig_win.ini")
	} else {
		os.Setenv("AWS_CONFIG_FILE", "testdata/shconfig.ini")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("region")},
	)

	if err != nil {
		t.Errorf("error getting session: %v", err)
	}

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "accessKey", creds.AccessKeyID; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "secret", creds.SecretAccessKey; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "tokenDefault", creds.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessProviderFromSessionWithProfileCfg(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_PROFILE", "non_expire")
	if runtime.GOOS == "windows" {
		os.Setenv("AWS_CONFIG_FILE", "testdata\\shconfig_win.ini")
	} else {
		os.Setenv("AWS_CONFIG_FILE", "testdata/shconfig.ini")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("region")},
	)

	if err != nil {
		t.Errorf("error getting session: %v", err)
	}

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "nonDefaultToken", creds.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessProviderNotFromCredProcCfg(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
	os.Setenv("AWS_PROFILE", "not_alone")
	if runtime.GOOS == "windows" {
		os.Setenv("AWS_CONFIG_FILE", "testdata\\shconfig_win.ini")
	} else {
		os.Setenv("AWS_CONFIG_FILE", "testdata/shconfig.ini")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("region")},
	)

	if err != nil {
		t.Errorf("error getting session: %v", err)
	}

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "notFromCredProcAccess", creds.AccessKeyID; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "notFromCredProcSecret", creds.SecretAccessKey; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessProviderFromSessionCrd(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	if runtime.GOOS == "windows" {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata\\shcred_win.ini")
	} else {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata/shcred.ini")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("region")},
	)

	if err != nil {
		t.Errorf("error getting session: %v", err)
	}

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "accessKey", creds.AccessKeyID; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "secret", creds.SecretAccessKey; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "tokenDefault", creds.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessProviderFromSessionWithProfileCrd(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	os.Setenv("AWS_PROFILE", "non_expire")
	if runtime.GOOS == "windows" {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata\\shcred_win.ini")
	} else {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata/shcred.ini")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("region")},
	)

	if err != nil {
		t.Errorf("error getting session: %v", err)
	}

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "nonDefaultToken", creds.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessProviderNotFromCredProcCrd(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	os.Setenv("AWS_PROFILE", "not_alone")
	if runtime.GOOS == "windows" {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata\\shcred_win.ini")
	} else {
		os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "testdata/shcred.ini")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("region")},
	)

	if err != nil {
		t.Errorf("error getting session: %v", err)
	}

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		t.Errorf("error getting credentials: %v", err)
	}

	if e, a := "notFromCredProcAccess", creds.AccessKeyID; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "notFromCredProcSecret", creds.SecretAccessKey; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}

}

func TestProcessProviderBadCommand(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	creds := processcreds.NewCredentials("/bad/process")
	_, err := creds.Get()
	if err.(awserr.Error).Code() != processcreds.ErrCodeProcessProviderExecution {
		t.Errorf("expected %v, got %v", processcreds.ErrCodeProcessProviderExecution, err)
	}
}

func TestProcessProviderMoreEmptyCommands(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	creds := processcreds.NewCredentials("")
	_, err := creds.Get()
	if err.(awserr.Error).Code() != processcreds.ErrCodeProcessProviderExecution {
		t.Errorf("expected %v, got %v", processcreds.ErrCodeProcessProviderExecution, err)
	}

}

func TestProcessProviderExpectErrors(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	creds := processcreds.NewCredentials(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "malformed.json"},
				string(os.PathSeparator))))
	_, err := creds.Get()
	if err.(awserr.Error).Code() != processcreds.ErrCodeProcessProviderParse {
		t.Errorf("expected %v, got %v", processcreds.ErrCodeProcessProviderParse, err)
	}

	creds = processcreds.NewCredentials(
		fmt.Sprintf("%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "wrongversion.json"},
				string(os.PathSeparator))))
	_, err = creds.Get()
	if err.(awserr.Error).Code() != processcreds.ErrCodeProcessProviderVersion {
		t.Errorf("expected %v, got %v", processcreds.ErrCodeProcessProviderVersion, err)
	}

	creds = processcreds.NewCredentials(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "missingkey.json"},
				string(os.PathSeparator))))
	_, err = creds.Get()
	if err.(awserr.Error).Code() != processcreds.ErrCodeProcessProviderRequired {
		t.Errorf("expected %v, got %v", processcreds.ErrCodeProcessProviderRequired, err)
	}

	creds = processcreds.NewCredentials(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "missingsecret.json"},
				string(os.PathSeparator))))
	_, err = creds.Get()
	if err.(awserr.Error).Code() != processcreds.ErrCodeProcessProviderRequired {
		t.Errorf("expected %v, got %v", processcreds.ErrCodeProcessProviderRequired, err)
	}

}

func TestProcessProviderTimeout(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	command := "/bin/sleep 2"
	if runtime.GOOS == "windows" {
		// "timeout" command does not work due to pipe redirection
		command = "ping -n 2 127.0.0.1>nul"
	}

	creds := processcreds.NewCredentialsTimeout(
		command,
		time.Duration(1)*time.Second)
	if _, err := creds.Get(); err == nil || err.(awserr.Error).Code() != processcreds.ErrCodeProcessProviderExecution || err.(awserr.Error).Message() != "credential process timed out" {
		t.Errorf("expected %v, got %v", processcreds.ErrCodeProcessProviderExecution, err)
	}

}

func TestProcessProviderWithLongSessionToken(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	creds := processcreds.NewCredentials(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "longsessiontoken.json"},
				string(os.PathSeparator))))
	v, err := creds.Get()
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}

	// Text string same length as session token returned by AWS for AssumeRoleWithWebIdentity
	e := "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
	if a := v.SessionToken; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}
}

type credentialTest struct {
	Version         int
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string
	Expiration      string
}

func TestProcessProviderStatic(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	// static
	creds := processcreds.NewCredentials(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "static.json"},
				string(os.PathSeparator))))
	_, err := creds.Get()
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if creds.IsExpired() {
		t.Errorf("expected %v, got %v", "static credentials/not expired", "expired")
	}

}

func TestProcessProviderNotExpired(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	// non-static, not expired
	exp := &credentialTest{}
	exp.Version = 1
	exp.AccessKeyID = "accesskey"
	exp.SecretAccessKey = "secretkey"
	exp.Expiration = time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)
	b, err := json.Marshal(exp)
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}

	tmpFile := strings.Join(
		[]string{"testdata", "tmp_expiring.json"},
		string(os.PathSeparator))
	if err = ioutil.WriteFile(tmpFile, b, 0644); err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	defer func() {
		if err = os.Remove(tmpFile); err != nil {
			t.Errorf("expected %v, got %v", "no error", err)
		}
	}()
	creds := processcreds.NewCredentials(
		fmt.Sprintf("%s %s", getOSCat(), tmpFile))
	_, err = creds.Get()
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if creds.IsExpired() {
		t.Errorf("expected %v, got %v", "not expired", "expired")
	}
}

func TestProcessProviderExpired(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	// non-static, expired
	exp := &credentialTest{}
	exp.Version = 1
	exp.AccessKeyID = "accesskey"
	exp.SecretAccessKey = "secretkey"
	exp.Expiration = time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	b, err := json.Marshal(exp)
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}

	tmpFile := strings.Join(
		[]string{"testdata", "tmp_expired.json"},
		string(os.PathSeparator))
	if err = ioutil.WriteFile(tmpFile, b, 0644); err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	defer func() {
		if err = os.Remove(tmpFile); err != nil {
			t.Errorf("expected %v, got %v", "no error", err)
		}
	}()
	creds := processcreds.NewCredentials(
		fmt.Sprintf("%s %s", getOSCat(), tmpFile))
	_, err = creds.Get()
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if !creds.IsExpired() {
		t.Errorf("expected %v, got %v", "expired", "not expired")
	}
}

func TestProcessProviderForceExpire(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	// non-static, not expired

	// setup test credentials file
	exp := &credentialTest{}
	exp.Version = 1
	exp.AccessKeyID = "accesskey"
	exp.SecretAccessKey = "secretkey"
	exp.Expiration = time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)
	b, err := json.Marshal(exp)
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	tmpFile := strings.Join(
		[]string{"testdata", "tmp_force_expire.json"},
		string(os.PathSeparator))
	if err = ioutil.WriteFile(tmpFile, b, 0644); err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	defer func() {
		if err = os.Remove(tmpFile); err != nil {
			t.Errorf("expected %v, got %v", "no error", err)
		}
	}()

	// get credentials from file
	creds := processcreds.NewCredentials(
		fmt.Sprintf("%s %s", getOSCat(), tmpFile))
	if _, err = creds.Get(); err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if creds.IsExpired() {
		t.Errorf("expected %v, got %v", "not expired", "expired")
	}

	// force expire creds
	creds.Expire()
	if !creds.IsExpired() {
		t.Errorf("expected %v, got %v", "expired", "not expired")
	}

	// renew creds
	if _, err = creds.Get(); err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if creds.IsExpired() {
		t.Errorf("expected %v, got %v", "not expired", "expired")
	}

}

func TestProcessProviderAltConstruct(t *testing.T) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	// constructing with exec.Cmd instead of string
	myCommand := exec.Command(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "static.json"},
				string(os.PathSeparator))))
	creds := processcreds.NewCredentialsCommand(myCommand, func(opt *processcreds.ProcessProvider) {
		opt.Timeout = time.Duration(1) * time.Second
	})
	_, err := creds.Get()
	if err != nil {
		t.Errorf("expected %v, got %v", "no error", err)
	}
	if creds.IsExpired() {
		t.Errorf("expected %v, got %v", "static credentials/not expired", "expired")
	}
}

func BenchmarkProcessProvider(b *testing.B) {
	oldEnv := preserveImportantStashEnv()
	defer awstesting.PopEnv(oldEnv)

	creds := processcreds.NewCredentials(
		fmt.Sprintf(
			"%s %s",
			getOSCat(),
			strings.Join(
				[]string{"testdata", "static.json"},
				string(os.PathSeparator))))
	_, err := creds.Get()
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := creds.Get()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func preserveImportantStashEnv() []string {
	envsToKeep := []string{"PATH"}

	if runtime.GOOS == "windows" {
		envsToKeep = append(envsToKeep, "ComSpec")
		envsToKeep = append(envsToKeep, "SYSTEM32")
	}

	extraEnv := getEnvs(envsToKeep)

	oldEnv := awstesting.StashEnv() //clear env

	for key, val := range extraEnv {
		os.Setenv(key, val)
	}

	return oldEnv
}

func getEnvs(envs []string) map[string]string {
	extraEnvs := make(map[string]string)
	for _, env := range envs {
		if val, ok := os.LookupEnv(env); ok && len(val) > 0 {
			extraEnvs[env] = val
		}
	}
	return extraEnvs
}

func getOSCat() string {
	if runtime.GOOS == "windows" {
		return "type"
	}
	return "cat"
}
