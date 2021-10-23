package api_service

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/stretchr/testify/require"
)

func TestLoadAppMetaData(t *testing.T) {
	t.Run("invalid 1", func(t *testing.T) {
		const fileName = "../data/invalid_appmetadata1.yaml"
		_, err := LoadAppMeta(fileName)
		require.Error(t, err)
	})
	t.Run("invalid 2", func(t *testing.T) {
		const fileName = "../data/invalid_appmetadata2.yaml"
		_, err := LoadAppMeta(fileName)
		require.Error(t, err)
	})
	validResponseArr := make([]*AppMetaData, 0)
	t.Run("valid 1", func(t *testing.T) {
		const fileName = "../data/valid_metadata1.yaml"
		res, err := LoadAppMeta(fileName)
		require.NoError(t, err)
		validResponseArr = append(validResponseArr, res)
		res.Print()
		require.Equal(t, res.Company, "Random Inc.")
		require.Len(t, res.Maintainers, 2)
		for i, mt := range res.Maintainers {
			if i == 0 {
				require.Equal(t, mt.Name, "firstmaintainer app1")
				require.Equal(t, mt.Email, "firstmaintainer@hotmail.com")
			} else {
				require.Equal(t, mt.Name, "secondmaintainer app1")
				require.Equal(t, mt.Email, "secondmaintainer@gmail.com")
			}
		}
		encodedStr, err := yaml.Marshal(res)
		require.NoError(t, err)
		t.Logf("encode back to yaml string:\n%v\n", string(encodedStr))
	})
	t.Run("valid 2", func(t *testing.T) {
		const fileName = "../data/valid_metadata2.yaml"
		res, err := LoadAppMeta(fileName)
		require.NoError(t, err)
		validResponseArr = append(validResponseArr, res)
		res.Print()
		encodedStr, err := yaml.Marshal(res)
		require.NoError(t, err)
		t.Logf("encode back to yaml string:\n%v\n", string(encodedStr))
	})
	t.Run("encode Array to yaml", func(t *testing.T) {
		encodedStr, err := yaml.Marshal(validResponseArr)
		require.NoError(t, err)
		t.Logf("encoded array to yaml:\n%v\n", string(encodedStr))
	})
}

func TestIsValidEmail(t *testing.T) {
	res := IsValidEmail("abc@hotmail.com")
	require.True(t, res)
	res = IsValidEmail("apptwohotmail.com")
	require.False(t, res)
	res = IsValidEmail("apptwo@hotmailcom")
	require.False(t, res)
	res = IsValidEmail("app#two@hotmail.com")
	require.False(t, res)
}
