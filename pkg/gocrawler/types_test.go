package gocrawler

import (
    "testing"
)

func TestScope(t *testing.T) {
    
    includes := []string{
        `^https?://([\w]+\.)*example\.com`,
    }
    scope := NewScope(includes, nil)

    _testSample(t, scope, "http://example.com", true)

    _testSample(t, scope, "https://www.google.com/?search=https://www.example.com", false)

    _testSample(t, scope, "ftp://example.com", false)
    _testSample(t, scope, "https://www.example.com", true)
}

func _testSample(t *testing.T, scope *Scope, test string, inScope bool) {
    if scope.InScope(test) != inScope {
        if inScope {
            t.Errorf("Expected '%s' to be in scope but got out of scope", test)
        } else {
            t.Errorf("Expected '%s' to be out of scope but got in scope", test)
        }
    }
}
