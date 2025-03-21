package crons

import (
	"encoding/json"
	"github.com/SwanHtetAungPhyo/common/model"
	"github.com/SwanHtetAungPhyo/common/protos"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"io"

	"os"
)

var log = logrus.New()
var LengthOfDAG = 0

func SaveDAGToFile() {
	// Check if the length of the DAG has changed
	if len(model.SwanDAG.Vertices) == LengthOfDAG {
		log.Println("DAG length has not changed. Skipping save.")
		return
	}

	// If the length has changed, proceed to save the DAG
	file, err := os.Create("./dag_data.json")
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(file)

	data, err := json.MarshalIndent(model.SwanDAG, "", "		")
	if err != nil {
		log.Fatalf("Error marshaling DAG: %v", err)
		return
	}

	_, err = file.Write(data)
	if err != nil {
		log.Fatalf("Error writing data to file: %v", err)
		return
	}

	// Update the LengthOfDAG after saving
	LengthOfDAG = len(model.SwanDAG.Vertices)

	log.Println("DAG data saved successfully")
}

func LoadDAGFromFile() *protos.DAG {
	file, err := os.Open("./dag_data.json")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
		return nil
	}
	defer file.Close()

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
		return nil
	}

	var protoDag protos.DAG
	err = proto.Unmarshal(data, &protoDag)
	if err != nil {
		log.Fatalf("Error unmarshalling file: %v", err) // Fix: Make sure it's binary format
		return nil
	}

	log.Println("DAG data loaded successfully")
	return &protoDag
}
func SetupCronJob() {
	c := cron.New()
	_, err := c.AddFunc("@every 5s", func() {
		SaveDAGToFile()
	})
	if err != nil {
		log.Fatalf("Error setting up cron job: %v", err)
		return
	}

	log.Println("[INFO] Cron job setup complete")
	c.Start()
}
