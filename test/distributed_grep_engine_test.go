package test

import (
	"cs425_mp1/internal/distributed_engine"
	"cs425_mp1/internal/grep"
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

	engine := distributed_engine.CreateEngine("test", "8080", nil, 20, false, "test%d.json")
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
Tests running one grep query with 3 VMs and seeing if the outputs are correct.
Only tests one query since each query is independent of each other and don't have any effect
on the correctness of the output, only the speed (due to caching)

NOTE: When running this test case, it assumes you already have VM 2 and VM 3 booted up with the program running.
This program must run on VM 1.
*/
func TestExecute3VM(t *testing.T) {
	cmdArgs := []string{"-n", "3", "-f", "../vm1.log", "-t", "test_execute_data/actual/"}
	cmd := exec.Command("../main", cmdArgs...)
	cmd.Stdin = strings.NewReader("grep -c GET\nexit\n")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// read open the json file it outputted
	actual_query, actual_gOut := distributed_engine.DeserializeJson("test_execute_data/actual/test1.json")
	expec_query, expec_gOut := distributed_engine.DeserializeJson("test_execute_data/expected/test1_expected.json")

	if actual_query != expec_query {
		t.Error("Query does not match")
	}
	if len(actual_gOut) != len(expec_gOut) {
		t.Error("actual grep output is not same length as expected grep output")
	}
	for i := 0; i < len(actual_gOut); i++ {
		if !grep.GrepOutputsAreEqual(&(actual_gOut[i]), &(expec_gOut[i])) {
			t.Errorf("Grep Outputs [%d] are NOT equal", i)
		}
	}
}

/*
Tests running the same query twice, and checking if the speed of the second query is faster than the speed of the first
query's execution time, which essentially checks if it cached its results

NOTE: When running this test case, it assumes you already have VM 2 and VM 3 booted up with the program running.
This program must run on VM 1.
*/
func TestExecuteCaching(t *testing.T) {
	cmdArgs := []string{"-n", "3", "-f", "../vm1.log", "-t", "test_execute_data/actual/2/"}
	cmd := exec.Command("../main", cmdArgs...)
	cmd.Stdin = strings.NewReader("grep -c GET\ngrep -c GET\nexit\n")
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	// read open the json file it outputted
	_, actual_gOut1 := distributed_engine.DeserializeJson("test_execute_data/actual/2/test1.json")
	_, actual_gOut2 := distributed_engine.DeserializeJson("test_execute_data/actual/2/test2.json")

	// evaluate the execution time
	for i := 0; i < len(actual_gOut1); i++ {
		// we want gOut2 time to be less than gOut1 time. Otherwise, it's an error
		if actual_gOut1[i].ExecutionTime < actual_gOut2[i].ExecutionTime {
			t.Errorf("gOut1 Time (%d) >= gOut2 Time (%d)", actual_gOut1[i].ExecutionTime.Nanoseconds(), actual_gOut2[i].ExecutionTime.Nanoseconds())
		}
	}
}
