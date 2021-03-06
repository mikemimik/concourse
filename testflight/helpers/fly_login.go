package helpers

import (
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func FlyLogin(atcURL, concourseAlias, flyBinaryPath string, teamName, username, password string) error {
	err := flyLogin(flyBinaryPath, []string{
		"-c", atcURL,
		"-t", concourseAlias,
		"-u", username,
		"-p", password,
	})
	if err != nil {
		return err
	}

	err = flyCreateTeam(flyBinaryPath, []string{
		"-t", concourseAlias,
		"-n", teamName,
		"--local-user", username,
		"--non-interactive",
	})
	if err != nil {
		return err
	}

	return flyLogin(flyBinaryPath, []string{
		"-c", atcURL,
		"-t", concourseAlias,
		"-u", username,
		"-p", password,
		"-n", teamName,
	})

}

func flyLogin(flyBinaryPath string, loginArgs []string) error {
	args := []string{"login"}

	caCertContents, certProvided := os.LookupEnv("FLY_CA_CERT")
	if certProvided {
		pathToCaCert, err := ioutil.TempFile("", "testflight-ca-cert")
		if err != nil {
			return err
		}

		defer os.Remove(pathToCaCert.Name())

		_, err = pathToCaCert.WriteString(caCertContents)
		if err != nil {
			return err
		}

		args = append(args, "--ca-cert", pathToCaCert.Name())
	}

	loginCmd := exec.Command(flyBinaryPath, append(args, loginArgs...)...)
	loginProcess, err := gexec.Start(loginCmd, GinkgoWriter, GinkgoWriter)
	if err != nil {
		return err
	}

	Eventually(loginProcess, time.Minute).Should(gexec.Exit(0))

	return nil
}

func flyCreateTeam(flyBinaryPath string, createTeamArgs []string) error {
	args := []string{"set-team"}

	createTeamCmd := exec.Command(flyBinaryPath, append(args, createTeamArgs...)...)
	createTeamProcess, err := gexec.Start(createTeamCmd, GinkgoWriter, GinkgoWriter)
	if err != nil {
		return err
	}

	Eventually(createTeamProcess, time.Minute).Should(gexec.Exit(0))

	return nil
}
