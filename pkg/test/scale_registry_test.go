package test

import (
	"testing"

	mc "github.com/eth-easl/loader/pkg/metric"
	tc "github.com/eth-easl/loader/pkg/trace"
	"github.com/stretchr/testify/assert"
)

func TestGetColdStartCount(t *testing.T) {
	functions := []tc.Function{}
	functions = append(functions, tc.Function{Name: "func-1"})
	functions = append(functions, tc.Function{Name: "func-2"})
	functions = append(functions, tc.Function{Name: "func-3"})
	registry := mc.ScaleRegistry{}

	/** Initialisation */
	records := []mc.DeploymentScale{
		//* Scale up NOT from 0.
		{Deployment: "func-1", Scale: 0},
		//* Scale up from 0.
		{Deployment: "func-2", Scale: 0},
		//* Haven't scaled.
		{Deployment: "func-2", Scale: 0},
	}
	registry.Init(records)

	assert.Equal(t, 0, registry.UpdateAndGetColdStartCount(records))
	/** Cold start. */
	records = []mc.DeploymentScale{
		{Deployment: "func-1", Scale: 10},
	}
	assert.Equal(t, 1, registry.UpdateAndGetColdStartCount(records))

	/** Mixing cold start and normal scaling up. */
	records = []mc.DeploymentScale{
		//* Scale up NOT from 0.
		{Deployment: "func-1", Scale: 100},
		//* Scale up from 0.
		{Deployment: "func-2", Scale: 100},
		//* Haven't scaled.
		{Deployment: "func-2", Scale: 0},
	}
	assert.Equal(t, 1, registry.UpdateAndGetColdStartCount(records))

	//* Scale down to 0.
	records = []mc.DeploymentScale{
		{Deployment: "func-1", Scale: 0},
		{Deployment: "func-2", Scale: 0},
	}
	assert.Equal(t, 0, registry.UpdateAndGetColdStartCount(records))

	/** All cold starts */
	records = []mc.DeploymentScale{
		{Deployment: "func-1", Scale: 200},
		{Deployment: "func-2", Scale: 200},
		{Deployment: "func-3", Scale: 200},
	}
	assert.Equal(t, 3, registry.UpdateAndGetColdStartCount(records))

}