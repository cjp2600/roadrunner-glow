package glow

import (
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog"
	"gopkg.in/workanator/go-floc.v2"
	"gopkg.in/workanator/go-floc.v2/run"
	"os"
)

// Glow
type Glow struct {
	svc *Service
}

// Execute
func (s *Glow) Execute(input string, output *string) error {
	var j []floc.Job

	logger := zerolog.New(os.Stdout).With().Timestamp().Str("service", "glow").Logger()
	request, err := s.getInputRequest(input)
	if err != nil {
		return err
	}
	if request == nil {
		return errors.New("empty request")
	}

	ctx := floc.NewContext()
	jobs := make(map[string]string)
	for _, sequence := range request.Sequence {

		switch sequence.Type {
		case Parallel:
			var pj []floc.Job
			for _, job := range sequence.Jobs {
				pj = append(pj, NewRequestTranslator(job, logger, request.Options).FlocExecute())
			}
			j = append(j, run.Parallel(pj...))
		case Sync:
			var sj []floc.Job
			for _, job := range sequence.Jobs {
				sj = append(sj, NewRequestTranslator(job, logger, request.Options).FlocExecute())
			}
			j = append(j, sj...)
		}
	}

	flow := run.Sequence(j...)
	_, _, err = floc.RunWith(ctx, floc.NewControl(ctx), flow)
	if err != nil {
		return err
	}
	for _, sequence := range request.Sequence {
		for _, job := range sequence.Jobs {
			if v, ok := ctx.Value(job.Id).(string); ok {
				jobs[job.Id] = v
			}
		}
	}

	err = s.getResponse(jobs, output)
	if err != nil {
		return err
	}

	return nil
}

// getInputRequest
func (s *Glow) getInputRequest(input string) (*ExecuteRequest, error) {
	var request ExecuteRequest
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	err := json.Unmarshal([]byte(input), &request)
	if err != nil {
		return nil, err
	}
	return &request, nil
}

// getResponse
func (s *Glow) getResponse(jobs map[string]string, output *string) error {
	response := new(ExecuteResponse)
	response.Jobs = jobs
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.Marshal(response)
	if err != nil {
		return err
	}
	*output = string(b)

	return nil
}
