package openapi_plugin_v1

import (
	"fmt"
	"io"
	"os"
	"path"
	"log"
	"flag"
	"io/ioutil"

	"github.com/golang/protobuf/proto"

	openapiv2 "github.com/googleapis/gnostic/OpenAPIv2"
	openapiv3 "github.com/googleapis/gnostic/OpenAPIv3"
)

// Environment contains the environment of a plugin call.
type Environment struct {
	Invocation      string      // string representation of call
	Document        interface{} // input document
	DocumentName	string 		// name of input document
	Response        *Response   // response message
	OutputPath      string      // output location
	RunningAsPlugin bool        // true if app is being run as a plugin
}

// NewEnvironment creates a plugin context from arguments and standard input.
func NewEnvironment() (env *Environment, err error) {
	env = &Environment{
		Invocation: os.Args[0],
		Response:   &Response{},
	}

	input := flag.String("input", "", "OpenAPI description (in binary protocol buffer form)")
	output := flag.String("output", "-", "Output file or directory")
	plugin := flag.Bool("plugin", false, "Run as a gnostic plugin (other flags are ignored).")
	flag.Parse()

	env.RunningAsPlugin = *plugin
	programName := path.Base(os.Args[0])

	if (*input == "") && !*plugin {
		flag.Usage = func() {
			fmt.Fprintf(os.Stderr, "\n")
			fmt.Fprintf(os.Stderr, programName+" is a gnostic plugin.\n")
			fmt.Fprintf(os.Stderr, `
When it is run from gnostic, the -plugin option is specified and gnostic
writes a binary request to stdin and waits for a binary response on stdout.

This program can also be run standalone using the other flags listed below.
When the -plugin option is specified, these flags are ignored.`)
			fmt.Fprintf(os.Stderr, "\n\nUsage:\n")
			flag.PrintDefaults()
		}
		flag.Usage()
		os.Exit(0)
	}

	if env.RunningAsPlugin {
		// handle invocation as a plugin.

		// Read the plugin input.
		pluginData, err := ioutil.ReadAll(os.Stdin)
		env.RespondAndExitIfError(err)
		if len(pluginData) == 0 {
			env.RespondAndExitIfError(fmt.Errorf("no input data"))
		}

		// Deserialize the request from the input.
		request := &Request{}
		err = proto.Unmarshal(pluginData, request)
		env.RespondAndExitIfError(err)

		env.OutputPath = request.OutputPath

		// Collect parameters passed to the plugin.
		parameters := request.Parameters
		for _, parameter := range parameters {
			env.Invocation += " " + parameter.Name + "=" + parameter.Value
		}

		// Log the invocation.
		log.Printf("Running plugin %s(input:%s)", env.Invocation, request.Wrapper.Version)

		// Read the document sent by the plugin.
		version := request.Wrapper.Version
		apiData := request.Wrapper.Value

		env.DocumentName = request.Wrapper.Name

		switch version {
		case "v2":
			documentv2 := &openapiv2.Document{}
			err = proto.Unmarshal(apiData, documentv2)
			env.RespondAndExitIfError(err)
			env.Document = documentv2
		case "v3":
			documentv3 := &openapiv3.Document{}
			err = proto.Unmarshal(apiData, documentv3)
			env.RespondAndExitIfError(err)
			env.Document = documentv3
		default:
			err = fmt.Errorf("Unsupported OpenAPI version %s", version)
			env.RespondAndExitIfError(err)
		}

	} else {
		// handle invocation from the command line.

		env.OutputPath = *output

		// Read the input document.
		apiData, err := ioutil.ReadFile(*input)
		if len(apiData) == 0 {
			env.RespondAndExitIfError(fmt.Errorf("no input data"))
		}

		env.DocumentName = path.Base(*input)

		// first try to unmarshal OpenAPI v2
		documentv2 := &openapiv2.Document{}
		err = proto.Unmarshal(apiData, documentv2)
		if err == nil {
			env.Document = documentv2
		} else {
			// ignore deserialization errors
		}

		// then try to unmarshal OpenAPI v3
		documentv3 := &openapiv3.Document{}
		err = proto.Unmarshal(apiData, documentv3)
		if err == nil {
			env.Document = documentv3
		} else {
			// ignore deserialization errors
		}

	}
	return env, err
}

// RespondAndExitIfError checks an error and if it is non-nil, records it and serializes and returns the response and then exits.
func (env *Environment) RespondAndExitIfError(err error) {
	if err != nil {
		env.Response.Errors = append(env.Response.Errors, err.Error())
		env.RespondAndExit()
	}
}

// RespondAndExit serializes and returns the plugin response and then exits.
func (env *Environment) RespondAndExit() {
	if env.RunningAsPlugin {
		responseBytes, _ := proto.Marshal(env.Response)
		os.Stdout.Write(responseBytes)
	} else {
		err := HandleResponse(env.Response, env.OutputPath)
		if err != nil {
			log.Printf("%s", err.Error())
		}
	}
	os.Exit(0)
}

func HandleResponse(response *Response, outputLocation string) error {
	if response.Errors != nil {
		return fmt.Errorf("Plugin error: %+v", response.Errors)
	}

	// Write files to the specified directory.
	var writer io.Writer
	switch {
	case outputLocation == "!":
		// Write nothing.
	case outputLocation == "-":
		writer = os.Stdout
		for _, file := range response.Files {
			writer.Write([]byte("\n\n" + file.Name + " -------------------- \n"))
			writer.Write(file.Data)
		}
	case isFile(outputLocation):
		return fmt.Errorf("unable to overwrite %s", outputLocation)
	default: // write files into a directory named by outputLocation
		if !isDirectory(outputLocation) {
			os.Mkdir(outputLocation, 0755)
		}
		for _, file := range response.Files {
			p := outputLocation + "/" + file.Name
			dir := path.Dir(p)
			os.MkdirAll(dir, 0755)
			f, _ := os.Create(p)
			defer f.Close()
			f.Write(file.Data)
		}
	}
	return nil
}

func isFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func isDirectory(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
