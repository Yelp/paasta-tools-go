package main

import (
    "fmt"
    "flag"
    "errors"
    "os"
    "path"
    "path/filepath"
    "strings"
    "github.com/xeipuuv/gojsonschema"
    "github.com/ghodss/yaml"
    "io/ioutil"
    "github.com/fatih/color"
)

type ValidateOptions struct {
    Service string
    SOAConfigPath string
}

func parseFlags(opts *ValidateOptions) error {
    flag.StringVar(&opts.Service, "service", "", "service to validate")
    flag.StringVar(&opts.SOAConfigPath, "yelpsoa-config-root", "", "yelpsoa-configs path")
    flag.Parse()
    return nil
}

func sanitiseKubernetesName(service string) string {
    name := strings.ReplaceAll(service, "_", "--")
    if strings.HasPrefix(name, "--") {
        name = strings.Replace(name, "--", "underscore-", 1)
    }
    return strings.ToLower(name)
}

// getServicePathDetermines the path of the directory containing the conf files
func getServicePath(service string, soa_dir string) (string, error) {
    if service != "" {
        return filepath.Join(soa_dir, service), nil
    }

    current_path, _ := os.Getwd()
    if soa_dir == current_path {
        return soa_dir, nil
    }
    return soa_dir,  errors.New("Unknown service")
}

// guessServiceName deduces the service name from the pwd
func guessServiceName(service string) string {
    if service != "" {
        return service
    }
    current_path, _ := os.Getwd()
    return filepath.Base(current_path)

}

// validateAutoscalingConfigs hasn't been implemented yet.
func validateAutoscalingConfigs(service_path string) bool {
    return true
}

// validateUniqueInstanceNames hasn't been implemented yet.
func validateUniqueInstanceNames(service_path string) bool {
    return true
}

// validatePaastaObjects hasn't been implemented yet.
func validatePaastaObjects(service_path string) bool {
    return true
}

// validateTron hasn't been implemented yet.
func validateTron(service_path string) bool {
    return true
}

// validateSchema checks if the specified config file has a valid schema.
// file_path is the path to file to validate.
// file_type is what schema type should we validate against
func validateSchema(file_name string, file_type string) bool {
    fileContents, err := ioutil.ReadFile(file_name)
    if err != nil {
        fmt.Printf("Failed to load config file: %s\n", err)
        return false
    }
    fileContentsJSON, err := yaml.YAMLToJSON(fileContents)
    if err != nil {
        fmt.Printf("Failed to convert yaml to json: %s\n", err)
        return false
    }

    schemaLoader := gojsonschema.NewReferenceLoader("file://schemas/" + file_type + "_schema.json")
    fileLoader := gojsonschema.NewStringLoader(string(fileContentsJSON))
    result, err := gojsonschema.Validate(schemaLoader, fileLoader)
    if err != nil {
        fmt.Printf("Schema invalid: %s\n", err)
        return false
    }
    if result.Valid() {
        green := color.New(color.FgGreen).SprintFunc()
        fmt.Printf("%s Successfully validated Schema: %s\n", green("Yes"), file_name)
        return true
    } else {
        red := color.New(color.FgRed).SprintFunc()
        blue := color.New(color.FgBlue).SprintFunc()
        paasta_document_link := "http://paasta.readthedocs.io/en/latest/yelpsoa_configs.html"
        fmt.Printf("%s Failed to validate schema. More info: %s: %s\n", red("No"), blue(paasta_document_link), file_name)
        fmt.Printf("  Validation Message:\n")
        for _, desc := range result.Errors() {
            fmt.Printf("  - %s\n", desc)
        }
        return false
    }
}

// validateAllSchemas Finds all recognized config files in service directory,and validates their schema.
// service_path is the path to location of configuration files.
func validateAllSchemas(service_path string) bool {
    matches, _ := filepath.Glob(service_path + "/*.yaml")
    return_code := true
    for _, file_name := range matches {
        file_info, _ := os.Lstat(file_name)
        if file_info.Mode() & os.ModeSymlink != 0 {
            continue
        }
        basename := path.Base(file_name)
        // This should be file_types := [4]string{"marathon", "adhoc", "tron", "kubernetes"}
        // But there is still some issues need to solve. For example,
        // https://groups.google.com/forum/#!topic/golang-nuts/7qgSDWPIh_E
        file_types := [2]string{"marathon", "adhoc"}
        for _, file_type := range file_types {
            if strings.HasPrefix(basename, file_type) == true {
                if validateSchema(file_name, file_type) == false {
                    return_code = false
                }
            }
        }
    }
    return return_code
}

func validateServiceName(service string) bool {
    sanitise_name := sanitiseKubernetesName(service)
    if len(sanitise_name) > 63 {
        fmt.Printf("Length of service name %s should be no more than 63.\n", sanitise_name)
        return false
    }
    return true
}

// checkServicePath heck that the specified path exists and has yaml files.
// service_path is the path to directory that should contain yaml files
func checkServicePath(service_path string) bool {
    if _, err := os.Stat(service_path); err != nil {
        fmt.Printf("%s is not a directory", service_path)
        return false
    }
    if matches, err := filepath.Glob(service_path + "/*.yaml"); len(matches) == 0 || err != nil {
        fmt.Printf("%s does not contain any .yaml files\n", service_path)
        return false
    }
    return true
}

func paastaValidateSoaConfigs(service string, service_path string) bool {
    if checkServicePath(service_path) == false {
        return false
    }
    if validateServiceName(service) == false {
        return false
    }
    if validateAllSchemas(service_path) == false {
        return false
    }
    if validateTron(service_path) == false {
        return false
    }
    if validatePaastaObjects(service_path) == false {
        return false
    }
    if validateUniqueInstanceNames(service_path) == false {
        return false
    }
    if validateAutoscalingConfigs(service_path) == false {
        return false
    }

    return true
}

func main() {
    options := &ValidateOptions{}
    err := parseFlags(options)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    service_path, err := getServicePath(options.Service, options.SOAConfigPath)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    service := guessServiceName(options.Service)

    result := paastaValidateSoaConfigs(service, service_path)
    if result == false {
        os.Exit(1)
    }
}
