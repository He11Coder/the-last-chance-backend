package nsfwFilter

import (
	"sync"
)

type ParallelResult struct {
	Inf           Inference
	ProcessingErr error
}

func RunInParallel(base64Images ...string) []ParallelResult {
	/*runCommand := func(_wg *sync.WaitGroup, cmd command, inputImage string, job_number int, resSlice []ParallelResult) {
		defer _wg.Done()

		res, err := cmd(inputImage)

		resSlice[job_number] = ParallelResult{
			Inf: res,
			ProcessingErr: err,
		}
	}*/

	results := make([]ParallelResult, len(base64Images))
	wg := &sync.WaitGroup{}
	for job_number, image := range base64Images {
		wg.Add(1)
		go func(i int, img string) {
			defer wg.Done()

			res, err := IsSafeForWork(img)
			results[i] = ParallelResult{Inf: res, ProcessingErr: err}
		}(job_number, image)
	}

	wg.Wait()

	return results
}
