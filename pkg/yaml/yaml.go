package yaml

import (
	"fmt"
	"io"
	"os"

	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

// Reads the given file from disk and unmarshals it as YAML
func YamlFromFile(filename string, obj interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("Could not open file %s: %s", filename, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Could not read file %s: %s", filename, err)
	}

	if err := yaml.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("Could not parse YAML from file %s: %s", filename, err)
	}

	return nil
}

func YamlToFile(data interface{}, path string, mode os.FileMode) error {
	yaml, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, yaml, mode)
	if err != nil {
		return err
	}

	klog.Infof("Written to disk: %s", path)

	return nil
}
