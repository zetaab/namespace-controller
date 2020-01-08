package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/mattbaird/jsonpatch"
	"gopkg.in/yaml.v2"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (c *Controller) checkAndUpdate(ns *v1.Namespace) {
	team, err := FindTeam(ns.Name, c.config)
	if err != nil {
		log.Printf("could not find team '%+v'", err)
		return
	}
	for labelKey, labelValue := range team.Labels {
		val, ok := ns.Labels[labelKey]
		if !ok || val != labelValue {
			log.Printf("Updating namespace %s label %s to value %s", ns.Name, labelKey, labelValue)
			err := c.patchNameSpace(ns, labelKey, labelValue)
			if err != nil {
				log.Printf("Got error while patching namespace %v", err)
			}
		}
	}
}

func (c *Controller) patchNameSpace(ns *v1.Namespace, label string, value string) error {
	oldJSON, err := json.Marshal(ns)
	if err != nil {
		return err
	}

	if len(ns.Labels) == 0 {
		ns.Labels = make(map[string]string)
	}
	ns.Labels[label] = value

	newJSON, err := json.Marshal(ns)
	if err != nil {
		return err
	}

	patch, err := jsonpatch.CreatePatch(oldJSON, newJSON)
	if err != nil {
		return err
	}

	pb, err := json.MarshalIndent(patch, "", "  ")
	if err != nil {
		return err
	}

	_, err = c.kclient.CoreV1().Namespaces().Patch(ns.Name, types.JSONPatchType, pb)
	return err
}

// FindTeam finds correct team for namespace
func FindTeam(name string, config *Config) (*Team, error) {
	for _, team := range config.Maintainers {
		for _, ns := range team.NameSpaces {
			match, err := regexp.MatchString(fmt.Sprintf("^%s$", ns), name)
			if err == nil && match {
				return &team, nil
			}
		}
	}
	return nil, fmt.Errorf("Team '%s' not found", name)
}

func makeConfig(path string) (*Config, error) {
	config := &Config{}
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
