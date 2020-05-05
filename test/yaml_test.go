package test

import (
	"fmt"
	"github.com/jinzhu/configor"
	"testing"
)

func TestConfigYaml(t *testing.T) {
	var Test = struct {
		ServerZones []struct {
			Name string
			Idx  int
			Num  int
		}

		InsteadTypes []struct {
			Name string
			Idx  int
		}

		Levels []struct {
			Name  string
			Idx   int
			Stars []struct {
				Name  string
				Idx   int
				Price int
			}
		}
	}{}
	configor.Load(&Test, "yaml.yml")
	fmt.Printf("config: %#v", Test)
}
