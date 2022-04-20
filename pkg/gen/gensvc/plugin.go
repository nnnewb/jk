package gensvc

import "github.com/nnnewb/jk/pkg/gen/gencore"

func GenerateEndpoint(data *gencore.PluginData) error {
	if err := genParamsResultsStruct(data); err != nil {
		return err
	}

	if err := genEndpointMaker(data); err != nil {
		return err
	}

	return nil
}
