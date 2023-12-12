package code

import (
    "testing"
    "regex"
)


func TestPromQLRegex(t *testing.T) {
    regStr := `^(?:[a-zA-Z_][a-zA-Z0-9_]*(?:\{(?:[a-zA-Z_][a-zA-Z0-9_]*="[^"]*"(?:, )?)*\})?(?:\[(?:[0-9]+|"[^"]*")(?:, )?\])*(?:\.(?:[a-zA-Z_][a-zA-Z0-9_]*|\*))*(?:\{(?:[a-zA-Z_][a-zA-Z0-9_]*="[^"]*"(?:, )?)*\})?(?:\[(?:[0-9]+|"[^"]*")(?:, )?\])*(?:\.(?:[a-zA-Z_][a-zA-Z0-9_]*|\*))*)$`
    testCases := []struct {
        input    string
        expected bool
    }{
        {"http_requests_total", true},
        {"http_requests_total[5m]", true},
        {"http_requests_total[5m][10]", true},
        {"http_requests_total{status=\"200\"}", true},
        {"sum(rate(http_requests_total[5m])) by (status)", true},
        {"sum(rate(http_requests_total[5m])) by status", true},
        {"http_requests_total{status=\"200\", method=\"GET\"}", true},
        {"http_requests_total{status=\"200\", method=\"GET\"}[5m]", true},
        {"http_requests_total{status=\"200\", method=\"GET\"}[5m][10]", true},
        //
        {"http_requests_total[5m", false},
        {"http_requests_total]", false},
        {"http_requests_total[5m][10", false},
        {"http_requests_total[5m][10][20]", false},
        {"http_requests_total{status=\"200\", method=\"GET\"", false},
        {"http_requests_total{status=\"200\", method=\"GET\"}}", false},
        {"http_requests_total{status=\"200\", method=\"GET\"}[5m", false},
        {"http_requests_total{status=\"200\", method=\"GET\"}[5m][10", false},
        {"http_requests_total{status=\"200\", method=\"GET\"}[5m][10][20]", false},
        {"sum(rate(http_requests_total[5m])) by (status", false},
        {"sum(rate(http_requests_total[5m])) by (status}", false},
        {"sum(rate(http_requests_total[5m)) by (status)", false},
        {"sum(rate(http_requests_total[5m])) by (status})", false},
        {"sum(rate(http_requests_total[5m)) by status", false},
    }

    for _, tc := range testCases {
        t.Run(tc.input, func(t *testing.T) {
            re := regexp.MustCompile(regStr)
            if re.MatchString(tc.input) != tc.expected {
                t.Errorf("expected %v but got %v", tc.expected, !tc.expected)
            }
        })
    }
}