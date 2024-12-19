package env_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/datngo2sgtech/go-packages/env"
)

func TestEnvironment(t *testing.T) {
	testCases := []struct {
		name           string
		envValue       string
		expectedResult string
	}{
		{
			name:           "Not set",
			envValue:       "",
			expectedResult: "dev",
		},
		{
			name:           "Test",
			envValue:       "test",
			expectedResult: "test",
		},
		{
			name:           "Production",
			envValue:       "production",
			expectedResult: "production",
		},
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				t.Setenv(env.EnvironmentVariable, tc.envValue)
			}
			actualResult := env.Environment()
			assert.Equal(t, tc.expectedResult, actualResult)
		})
	}
}

func TestIsTestEnv(t *testing.T) {
	testCases := []struct {
		name           string
		envValue       string
		expectedResult bool
	}{
		{
			name:           "Not set",
			envValue:       "",
			expectedResult: false,
		},
		{
			name:           "Test",
			envValue:       "test",
			expectedResult: true,
		},
		{
			name:           "Production",
			envValue:       "production",
			expectedResult: false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.envValue != "" {
				t.Setenv(env.EnvironmentVariable, testCase.envValue)
			}
			actualResult := env.IsTestEnv()
			assert.Equal(t, testCase.expectedResult, actualResult)
		})
	}
}
