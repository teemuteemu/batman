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
	vm, err := vm.New(e)
	if err != nil {
		return nil, err
	}

	return &Run{
		Collection: c,
		Env:        e,
		VM:         vm,
	}, nil
}

func (r *Run) ProcessRequests(requestNames []string, formatter client.Formatter, execute bool) error {
	for _, requestName := range requestNames {
		request, err := r.prepareRequest(requestName)
		if err != nil {
			return err
		}

		if !execute {
			fmt.Print(formatter.RenderRequest(request))
		} else {
			response, err := client.ExecuteRequest(request)
			if err != nil {
				return err
			}

			call := client.Call{
				Request:  request,
				Response: response,
			}

			fmt.Print(formatter.Render(&call))
		}
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

		request, err := r.prepareRequest(step.Request)
		if err != nil {
			fmt.Printf("%s\n", err)
			i++
			continue
		}

		response, err := client.ExecuteRequest(request)
		if err != nil {
			fmt.Printf("%s\n", err)
			i++
			continue
		}

		if step.After != nil {
			gotoStep, err := r.runJS(step.After, response)
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
			jsonOutput, err := renderer.FormatJSON(response.Body)
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

func (r *Run) prepareRequest(requestName string) (*client.Request, error) {
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

	renderedRequest := *request
	renderedRequest.URL = url
	renderedRequest.Header = headers
	renderedRequest.Body = body.String()

	return &renderedRequest, nil
}

func (r *Run) runJS(script *string, response *client.Response) (string, error) {
	gotoStep, err := r.VM.ExecuteScript(script, response)
	if err != nil {
		return "", err
	}

	return gotoStep, nil
}
