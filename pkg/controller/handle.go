package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/mattbaird/jsonpatch"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	// PodSecurityModeEnforcing ...
	PodSecurityModeEnforcing = "pod-security.kubernetes.io/enforce"
	// PodSecurityModeWarn ...
	PodSecurityModeWarn = "pod-security.kubernetes.io/warn"
)

func (c *Controller) checkAndUpdate(ns *v1.Namespace) {
	ctx := context.Background()
	if !Contains(c.config.AdminNamespaces, ns.Name) {
		_, err := c.kclient.CoreV1().LimitRanges(ns.Name).Get(ctx, "default-limits", metav1.GetOptions{})
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
			_, errCreate := c.kclient.CoreV1().LimitRanges(ns.Name).Create(ctx, limitRange, metav1.CreateOptions{})
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
		newValue := c.handleLabels(ctx, ns, labelKey, labelValue)
		if len(newValue) > 0 {
			err := c.patchNameSpace(ctx, ns, labelKey, newValue)
			if err != nil {
				log.Printf("Got error while patching namespace %v", err)
			}
		}
	}
}

func getDefaultValue(labelValue string) string {
	return strings.Split(labelValue, ",")[0]
}

func (c *Controller) handleLabels(ctx context.Context, ns *v1.Namespace, key string, labelValues string) string {
	if key == PodSecurityModeEnforcing {
		log.Printf("enforce label exists in namespace %s ensure warn label does not exist", ns.Name)
		if c != nil {
			err := c.patchNameSpace(ctx, ns, PodSecurityModeWarn, "")
			if err != nil {
				log.Printf("patching failed %v", err)
			}
		}
	}

	newValue := getDefaultValue(labelValues)
	val, ok := ns.Labels[key]
	if !ok {
		_, ok2 := ns.Labels[PodSecurityModeEnforcing]
		if key == PodSecurityModeWarn && ok2 {
			log.Printf("enforce label exists in namespace %s ignoring warn label", ns.Name)
			return ""
		}
		log.Printf("adding label %s with value %s to namespace %s", key, newValue, ns.Name)
		return newValue
	}
	allowedValues := strings.Split(labelValues, ",")
	if !Contains(allowedValues, val) {
		log.Printf("updating label %s to value %s in namespace %s. Not allowed value %s", key, newValue, ns.Name, val)
		return newValue
	}
	return ""
}

func (c *Controller) patchNameSpace(ctx context.Context, ns *v1.Namespace, label string, value string) error {
	oldJSON, err := json.Marshal(ns)
	if err != nil {
		return err
	}

	if len(ns.Labels) == 0 {
		ns.Labels = make(map[string]string)
	}

	if value != "" {
		ns.Labels[label] = value
	}

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

	_, err = c.kclient.CoreV1().Namespaces().Patch(ctx, ns.Name, types.JSONPatchType, pb, metav1.PatchOptions{})
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
