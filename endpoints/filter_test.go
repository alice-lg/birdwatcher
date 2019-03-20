package endpoints

import (
	"testing"
)

func TestValidateProtocol(t *testing.T) {

	validProtocols := []string{
		"ID421_AS11171_123.8.127.19",
		"ID429_AS12240_2222:7af8:8:05:01:30bb:0:1",
		"AI421_AS11171_123..8..127..19",
	}

	invalidProtocols := []string{
		"ID421_AS11171_123.8.127.lö19",
		"ThisValueIsTooLong12345678901234567890123456789012345678901234567890123456789012345678901234567890",
	}

	// Valid protocol values
	for _, param := range validProtocols {
		t.Log("Testing valid protocol:", param)
		_, err := ValidateProtocolParam(param)
		if err != nil {
			t.Error(param, "should be a valid protocol param")
		}
	}

	// Invalid protocol values
	for _, param := range invalidProtocols {
		t.Log("Testing invalid protocol:", param)
		_, err := ValidateProtocolParam(param)
		if err == nil {
			t.Error(param, "should be an invalid protocol param")
		}
	}

}
