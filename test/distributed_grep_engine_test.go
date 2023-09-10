package test

import (
	"cs425_mp1/internal/distributed_engine"
	"cs425_mp1/internal/grep"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestCreatingJson(t *testing.T) {

	packagedString := "grep;-w;sample"

	output1 := "Output file1"
	filename1 := "sample_text_file1.txt"
	numLines1 := 20
	exectionTime1 := time.Duration(50)
	grepOut1 := grep.GrepOutput{output1, filename1, numLines1, exectionTime1}

	output2 := "Output file2 from grep"
	filename2 := "example_text_file2.txt"
	numLines2 := 3
	exectionTime2 := time.Duration(5)
	grepOut2 := grep.GrepOutput{output2, filename2, numLines2, exectionTime2}

	output3 := "Output file3 from grep"
	filename3 := "test_text_file3.txt"
	numLines3 := 8
	exectionTime3 := time.Duration(12)
	grepOut3 := grep.GrepOutput{output3, filename3, numLines3, exectionTime3}

	outputs := []grep.GrepOutput{grepOut1, grepOut2, grepOut3}

	engine := distributed_engine.CreateEngine("test", "8080", nil, 0, false, "test%d.json")
	_, err := engine.CreateJson(packagedString, outputs)

	if err != nil {
		t.Errorf("Failed to Create Json file")
	}

	query, outputs := distributed_engine.DeserializeJson("test1.json")

	if query != packagedString {
		t.Errorf("Expected query %s but got %s", packagedString, query)
	}

	if len(outputs) != 3 {
		t.Errorf("Expected output of length 3 but got %d", len(outputs))
	}

	if !grep.GrepOutputsAreEqual(&outputs[0], &grepOut1) {
		t.Errorf("Expected grepOutput does not match output")
	}

	if !grep.GrepOutputsAreEqual(&outputs[1], &grepOut2) {
		t.Errorf("Expected grepOutput does not match output")
	}

	if !grep.GrepOutputsAreEqual(&outputs[2], &grepOut3) {
		t.Errorf("Expected grepOutput does not match output")
	}
}

/*
When this function is ran, we assume that the other VMs have already have their programs running, with
the [r] keyword already pressed to indicate
*/
func TestExecuteSmall(t *testing.T) {
	cmd := exec.Command("./main", "-n", "3", "-f", "test_logs/test_log_file1.log", "-t")
	cmd.Stdin = strings.NewReader("r\ngrep -c ERROR\nexit\n")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// read open the json file it outputted
	query, grepOutputs := distributed_engine.DeserializeJson("test1.json")

	if query != "grep -c ERROR" {
		t.Error("Expected query does not match")
	}

	for _, gOut := range grepOutputs {
		fmt.Println(gOut.ToString())
	}
}
