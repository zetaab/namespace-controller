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
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (c *Controller) checkAndUpdate(ns *v1.Namespace) {
	if !Contains(c.config.AdminNamespaces, ns.Name) {
		_, err := c.kclient.CoreV1().LimitRanges(ns.Name).Get("default-limits", metav1.GetOptions{})
		if errors.IsNotFound(err) {
			// default values which can be overriden by config
			limitCPU := "200m"
			limitMemory := "100Mi"
			requestCPU := "25m"
			requestMemory := "100Mi"
			if c.config.LimitCPU != "" {
				limitCPU = c.config.LimitCPU
			}
			if c.config.LimitMemory != "" {
				limitMemory = c.config.LimitMemory
			}
			if c.config.RequestCPU != "" {
				requestCPU = c.config.RequestCPU
			}
			if c.config.RequestMemory != "" {
				requestMemory = c.config.RequestMemory
			}

			limitRange := &v1.LimitRange{
				ObjectMeta: metav1.ObjectMeta{
					Name: "default-limits",
				},
				Spec: v1.LimitRangeSpec{
					Limits: []v1.LimitRangeItem{
						{
							Type: v1.LimitTypeContainer,
							Default: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse(limitCPU),
								v1.ResourceMemory: resource.MustParse(limitMemory),
							},
							DefaultRequest: v1.ResourceList{
								v1.ResourceCPU:    resource.MustParse(requestCPU),
								v1.ResourceMemory: resource.MustParse(requestMemory),
							},
						},
					},
				},
			}
			_, errCreate := c.kclient.CoreV1().LimitRanges(ns.Name).Create(limitRange)
			if errCreate != nil {
				log.Printf("could not create limit-range '%+v'", errCreate)
			}
		} else if err != nil {
			log.Printf("could not fetch limit-ranges '%+v'", err)
		}
	}

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

// Contains is checking does array contain single word
func Contains(array []string, word string) bool {
	for _, item := range array {
		if item == word {
			return true
		}
	}
	return false
}
