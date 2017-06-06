package sar

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Sar struct {
	InputsDir   string `toml:"inputs_dir_path"`
	OutputsDir  string `toml:"outputs_dir_path"`
	SarUtilPath string `toml:"sar_path"`
}

func (*Sar) Description() string {
	return "Sar utility output collector"
}

var sampleConfig = `
  ## Path to the sar command. 
  #
  sar_path = "usr/bin/sar"
  #
  #
  ## Input files directory path
  #
  inputs_dir_path = "./sar_inputs"
  #
  #
  ## Output file path
  outputs_dir_path = "./sar_processed"
`

func (*Sar) SampleConfig() string {
	return sampleConfig
}

func (s *Sar) Gather(acc telegraf.Accumulator) error {
	files, err := ioutil.ReadDir(s.InputsDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() || !strings.HasPrefix(f.Name(), "sa") {
			continue
		}

		filePath := path.Join(s.InputsDir, f.Name())

		cmd := exec.Command(s.SarUtilPath, "-f", filePath)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			acc.AddError(err)
			continue
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			acc.AddError(err)
			continue
		}

		err = cmd.Start()
		if err != nil {
			acc.AddError(err)
			continue
		}

		r := bufio.NewReader(stdout)
		err = s.parse(acc, r)
		if err != nil {
			acc.AddError(err)
			continue
		}

		// read error message and return as an error
		slurp, _ := ioutil.ReadAll(stderr)

		err = cmd.Wait()
		if err != nil {
			return errors.New(string(slurp))
		}

		err = os.Rename(filePath, path.Join(s.OutputsDir, f.Name()))
		if err != nil {
			acc.AddError(err)
			continue
		}
	}

	return nil
}

func (s *Sar) parse(acc telegraf.Accumulator, r *bufio.Reader) error {
	// get system meta info
	m, _, err := r.ReadLine()
	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	meta := strings.Fields(string(m))

	if len(meta) == 6 {
		return errors.New("invalid sar output file")
	}

	tags := map[string]string{
		"linux_kernel": meta[0] + meta[1],
		"host":         meta[2],
		"date":         meta[3],
		"arch":         meta[4],
		"num_cpu":      strings.TrimPrefix(meta[5], "(") + strings.TrimSuffix(meta[6], "CPU)"),
	}

	// meta info is followed by emptry line
	// ignore it
	r.ReadLine()

	// get headers
	h, _, err := r.ReadLine()
	if err == io.EOF {
		return nil
	}

	if err != nil {
		return err
	}

	headers := strings.Fields(string(h))

	for {
		h, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}

		if err != nil {
			acc.AddError(err)
		}

		if len(h) < 1 {
			continue
		}

		record := strings.Fields(string(h))

		fields := map[string]interface{}{}

		// average metrics
		if len(record) > 1 && record[0] == "Average:" {
			// start from third elem, ignore the timestamp
			for i := 2; i < len(headers); i++ {
				f, err := strconv.ParseFloat(record[i-1], 64)
				if err != nil {
					fields[strings.TrimPrefix(string(headers[i]), "%")] = record[i-1]
					continue
				}

				fields[strings.TrimPrefix(string(headers[i]), "%")] = f
			}

			acc.AddFields("average", fields, tags)
			continue
		}

		timestamp, _ := time.Parse("3:04:05 PM", record[0]+" "+record[1])

		// start from third elem, ignore the timestamp
		for i := 2; i < len(headers); i++ {
			f, err := strconv.ParseFloat(record[i], 64)
			if err != nil {
				fields[strings.TrimPrefix(string(headers[i]), "%")] = record[i]
				continue
			}

			fields[strings.TrimPrefix(string(headers[i]), "%")] = f
		}

		acc.AddFields("sar", fields, tags, timestamp)
	}

	return nil
}

func init() {
	inputs.Add("sar", func() telegraf.Input {
		return &Sar{}
	})
}
