package domain

import (
	"net"
	"strings"
	"testing"
	"time"
)

func testIPParser(t *testing.T, testValue string, expectedValue net.IPNet) {
	result := parseIP(testValue)
	println(result.String())
	println(expectedValue.String())
	if strings.Compare(result.String(), expectedValue.String()) != 0 {
		t.Errorf("parseIP gave unexpected output, got %v wanted %v", result, expectedValue)
	}
}

func TestIpParsingNormalIpv4(t *testing.T) {
	expectedNet := net.IPNet{IP: net.IP{192, 168, 1, 12}, Mask: net.IPMask{255, 255, 255, 255}}
	testIPParser(t, "192.168.1.12", expectedNet)
}

func TestIpParsingCidrIpv4(t *testing.T) {
	expectedNet := net.IPNet{IP: net.IP{192, 168, 1, 12}, Mask: net.IPMask{255, 255, 255, 0}}
	testIPParser(t, "192.168.1.12/24", expectedNet)
}

func TestIpParsingShortenedIpv6(t *testing.T) {
	expectedNet := net.IPNet{
		IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
	}
	testIPParser(t, "::1", expectedNet)
}

func TestIpParsingNormalIpv6(t *testing.T) {
	expectedNet := net.IPNet{
		IP:   net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
		Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
	}
	testIPParser(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", expectedNet)
}

func TestIpParsingCidrIpv6(t *testing.T) {
	expectedNet := net.IPNet{
		IP:   net.ParseIP("2001:0db8:85a3:0000:0000:8a2e:0370:7334"),
		Mask: net.IPMask{255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	testIPParser(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334/16", expectedNet)
}

// Handling invalid data is not a responsibility of the struct.
func TestIpParsingInvalidValue(t *testing.T) {
	expectedNet := net.IPNet{
		IP:   net.IP{},
		Mask: net.IPMask{},
	}
	testIPParser(t, "invalid", expectedNet)
}

func TestTimeoutConvertsIntToSecondDuration(t *testing.T) {
	if parseTimeout(123) != (time.Duration(123) * time.Second) {
		t.Errorf("parseTimeout gave unexpected output")
	}
}

func TestNewOpenCommandWillFuzeIpAndRules(t *testing.T) {
	inputIP := parseIP("192.168.2.24/16")
	rules := make([]Rule, 1)
	direction, _ := NewDirectionFromString("outbound")
	rules = append(rules, Rule{Direction: direction})

	result := parseRules(inputIP, rules)
	if result[0].IPNet.String() != inputIP.String() {
		t.Errorf("parseRules failed to apply the input")
	}
}
