package extractor

import (
	"regexp"

	"golang.org/x/xerrors"
)

func ExtractAndReplaceText(text string, regexString string, outputFormat string) ([]string, error) {
	filter, err := regexp.Compile(regexString)
	if err != nil {
		return nil, xerrors.Errorf("failed to compile regex: %w", err)
	}

	outputs := filter.FindAllString(text, -1)
	if outputFormat == "" {
		return outputs, nil
	}

	for i, output := range outputs {
		outputs[i] = filter.ReplaceAllString(output, outputFormat)
	}

	return outputs, nil
}
