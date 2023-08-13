package helpers

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/marlosl/gpt-telegram-bot/consts"
)

type BinaryFile struct {
	FunctionName string
	Filename     string
	HandlerDir   string
}

var Binaries = []BinaryFile{
	{
		FunctionName: "chat-gpt-talk",
		Filename:     "chat-gpt-talk-handler",
		HandlerDir:   "chatgpttalk",
	},
	{
		FunctionName: "chat-gpt-send-image",
		Filename:     "chat-gpt-send-image-handler",
		HandlerDir:   "chatgptsendimage",
	},
}

func GetFunctionNames() []string {
	var names []string
	for _, binary := range Binaries {
		names = append(names, binary.FunctionName)
	}

	return names
}

func GetZipFilenameAndFullFunctionName(functionName string) (string, string) {
	outputDir := os.Getenv(consts.ProjectOutputDir)
	for _, binary := range Binaries {
		if binary.FunctionName == functionName {
			outputFilename := filepath.Join(outputDir, binary.HandlerDir, binary.Filename+".zip")
			fullFunctionName := binary.Filename
			return outputFilename, fullFunctionName
		}
	}
	return "", ""
}

func initBuilder() (string, string) {
	projectDir := os.Getenv(consts.ProjectDir)
	outputDir := os.Getenv(consts.ProjectOutputDir)

	if projectDir == "" {
		log.Fatal("Project directory not set.")
	}

	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		log.Fatal("Project directory does not set.")
	}

	if err := createDirIfNotExist(outputDir); err != nil {
		log.Fatalf("Can't create the output directory: %v\n", err)
	}

	os.Setenv("GOOS", "linux")
	os.Setenv("GOARCH", "amd64")
	os.Setenv("CGO_ENABLED", "0")
	return projectDir, outputDir
}

func BuildSingleBinaryFile(functionName string) error {
	projectDir, outputDir := initBuilder()
	for _, binary := range Binaries {
		if binary.FunctionName == functionName {
			return BuildBinaryFile(binary, projectDir, outputDir)
		}
	}
	return fmt.Errorf("function %s not found", functionName)
}

func BuildBinaryFiles() {
	projectDir, outputDir := initBuilder()
	for _, binary := range Binaries {
		BuildBinaryFile(binary, projectDir, outputDir)
	}
}

func BuildBinaryFile(binary BinaryFile, projectDir string, outputDir string) error {
	if err := buildBinaryFile(binary, projectDir, outputDir); err != nil {
		fmt.Printf("Can't build binary file: %v\n", err)
		return err
	}

	if err := zipBinaryFile(binary, outputDir); err != nil {
		fmt.Printf("Can't compress binary file: %v\n", err)
		return err
	}
	fmt.Printf("Binary file %s was built successfully.\n", binary.Filename)
	return nil
}

func buildBinaryFile(binary BinaryFile, projectDir string, outputDir string) error {
	fmt.Printf("Building %s binary file...\n", binary.Filename)

	outputDirName := filepath.Join(outputDir, binary.HandlerDir)
	if err := createDirIfNotExist(outputDirName); err != nil {
		return err
	}

	outputFilename := filepath.Join(outputDirName, "main")

	goPrg := "go"

	buildArg := "build"
	outputArg := "-o"

	sourceFile := filepath.Join(projectDir, "cmd", binary.HandlerDir, "main.go")

	cmd := exec.Command(goPrg, buildArg, outputArg, outputFilename, sourceFile)
	stdout, err := cmd.Output()

	if err != nil {
		return err
	}

	fmt.Print(string(stdout))

	return nil
}

func zipBinaryFile(binary BinaryFile, outputDir string) error {
	fmt.Printf("Compressing %s binary file...\n", binary.Filename)
	zipPrg := "zip"

	zipOpt := "-jrm"

	outputFilename := filepath.Join(outputDir, binary.HandlerDir, binary.Filename+".zip")
	inputFilename := filepath.Join(outputDir, binary.HandlerDir, "main")

	if err := deleteFileIfExist(outputFilename); err != nil {
		return err
	}

	cmd := exec.Command(zipPrg, zipOpt, outputFilename, inputFilename)
	stdout, err := cmd.Output()

	if err != nil {
		return err
	}

	fmt.Print(string(stdout))
	return nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteFileIfExist(file string) error {
	if _, err := os.Stat(file); !os.IsNotExist(err) {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}

	return nil
}
