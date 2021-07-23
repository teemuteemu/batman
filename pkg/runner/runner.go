package runner

import (
	"fmt"

	"github.com/teemuteemu/batman/pkg/client"
	"github.com/teemuteemu/batman/pkg/env"
	"github.com/teemuteemu/batman/pkg/files"
	"github.com/teemuteemu/batman/pkg/renderer"
	"github.com/teemuteemu/batman/pkg/vm"
)

type Run struct {
	Collection *files.Collection
	Env        *env.Env
	VM         *vm.VM
}

func New(c *files.Collection, e *env.Env) (*Run, error) {
	vm := vm.New(e)

	return &Run{
		Collection: c,
		Env:        e,
		VM:         vm,
	}, nil
}

func (r *Run) RunRequests(requestNames []string) error {
	for _, requestName := range requestNames {
		resp, err := r.runRequest(requestName)
		if err != nil {
			return err
		}

		jsonOutput, err := renderer.FormatJSON(resp.Body)
		if err != nil {
			return err
		}

		fmt.Println(jsonOutput)
	}

	return nil
}

func (r *Run) RunScript(script *files.Script) error {
	fmt.Printf("%s:\n", script.Name)

	i := 0
	for {
		step := script.Steps[i]

		fmt.Printf("%s\t", step.Name)

		if step.Before != nil {
			fmt.Printf("\n")

			gotoStep, err := r.runJS(step.Before, nil)
			if err != nil {
				return err
			}

			if len(gotoStep) > 0 {
				jumpIndx := 0

				for ji, step := range script.Steps {
					if step.Name == gotoStep {
						jumpIndx = ji
						break
					}
				}

				i = jumpIndx
				continue
			}
		}

		resp, err := r.runRequest(step.Request)
		if err != nil {
			fmt.Printf("%s\n", err)
			i++
			continue
		}

		if step.After != nil {
			gotoStep, err := r.runJS(step.After, resp)
			if err != nil {
				return err
			}

			if len(gotoStep) > 0 {
				jumpIndx := 0

				for ji, step := range script.Steps {
					if step.Name == gotoStep {
						jumpIndx = ji
						break
					}
				}

				i = jumpIndx
				continue
			}
		}

		if step.Output {
			jsonOutput, err := renderer.FormatJSON(resp.Body)
			if err != nil {
				return err
			}

			fmt.Println(jsonOutput)
		}

		i++
		if i >= len(script.Steps) {
			break
		}
	}

	return nil
}

func (r *Run) runRequest(requestName string) (*client.Response, error) {
	request, err := r.Collection.FindRequest(requestName)
	if err != nil {
		return nil, err
	}
	if request == nil {
		return nil, fmt.Errorf(`Request "%s" not found`, requestName)
	}

	url, err := r.Env.PrepareURL(request.URL)
	if err != nil {
		return nil, err
	}

	headers, err := r.Env.PrepareHeader(request.Header)
	if err != nil {
		return nil, err
	}

	body, err := r.Env.PrepareBody(request.Body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s %s ", request.Method, url)

	resp, err := client.ExecuteRequest(request.Method, url, headers, body)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\t%d\n", resp.StatusCode)

	return resp, nil
}

func (r *Run) runJS(script *string, resp *client.Response) (string, error) {
	gotoStep, err := r.VM.ExecuteScript(script, resp)
	if err != nil {
		return "", err
	}

	return gotoStep, nil
}
