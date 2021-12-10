package function

import (
	"os/exec"

	// fc "github.com/eth-easl/easyloader/internal/function"
	tc "github.com/eth-easl/easyloader/internal/trace"
	log "github.com/sirupsen/logrus"
)

func Deploy(functions []tc.Function, serviceConfigPath string, deploymentConcurrency int) []string {
	var urls []string
	/**
	 * Limit the number of parallel deployments
	 * using a channel (like semaphore).
	 */
	sem := make(chan bool, deploymentConcurrency)

	// log.Info("funcSlice: ", funcSlice)
	for idx, function := range functions {
		sem <- true

		go func(function tc.Function, idx int) {
			defer func() { <-sem }()

			has_deployed := deployFunction(&function, serviceConfigPath)
			function.SetDeployed(has_deployed)
			if has_deployed {
				urls = append(urls, function.GetUrl())
			}

			functions[idx] = function
		}(function, idx)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return urls
}

func deployFunction(function *tc.Function, workloadPath string) bool {
	cmd := exec.Command(
		"kn",
		"service",
		"apply",
		function.GetName(),
		"-f",
		workloadPath,
		"--concurrency-target",
		"1",
	)
	stdoutStderr, err := cmd.CombinedOutput()
	log.Debug("CMD response: ", string(stdoutStderr))

	if err != nil {
		log.Warnf("Failed to deploy function %s: %v\n%s\n", function.GetName(), err, stdoutStderr)
		return false
	}

	// assemble function url from response from kubectl and the standard port
	log.Info("Deployed function ", function.GetUrl())
	return true
}
