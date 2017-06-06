package sar

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/require"
)

func TestGather(t *testing.T) {
	s := Sar{
		InputsDir:   "testdata/inputs",
		OutputsDir:  "testdata/outputs",
		SarUtilPath: "/usr/local/bin/sar",
	}

	acc := &testutil.Accumulator{}
	acc.SetDebug(true)

	// tear down
	// put processed files back to input dir
	defer func() {
		files, err := ioutil.ReadDir("testdata/outputs")
		require.NoError(t, err)

		for _, f := range files {
			err = os.Rename("testdata/outputs/"+f.Name(), "testdata/inputs/"+f.Name())
			require.NoError(t, err)
		}
	}()

	err := s.Gather(acc)
	require.NoError(t, err)

	t.Log(acc.NMetrics())

	require.True(t, acc.HasTag("sar", "linux_kernel"))
	require.True(t, acc.HasTag("sar", "host"))
	require.True(t, acc.HasTag("sar", "date"))
	require.True(t, acc.HasTag("sar", "arch"))
	require.True(t, acc.HasTag("sar", "num_cpu"))
}

func TestParse(t *testing.T) {
	testdata := `Linux 3.10.0-229.14.1.el7.x86_64 (cocsvma5app101) 	05/01/2017 	_x86_64_	(32 CPU)

03:00:01 AM     CPU     %user     %nice   %system   %iowait    %steal     %idle
03:10:01 AM     all      0.43      0.00      0.16      0.00      0.00     99.42
03:20:01 AM     all      0.40      0.00      0.15      0.00      0.00     99.45
03:30:01 AM     all      0.42      0.00      0.15      0.00      0.00     99.43
03:40:01 AM     all      0.43      0.00      0.15      0.00      0.00     99.42
03:50:01 AM     all      0.44      0.00      0.15      0.00      0.00     99.41
04:00:01 AM     all      0.40      0.00      0.15      0.00      0.00     99.45
04:10:01 AM     all      0.43      0.00      0.15      0.00      0.00     99.41
04:20:01 AM     all      0.40      0.00      0.15      0.00      0.00     99.45
04:30:01 AM     all      0.42      0.00      0.15      0.00      0.00     99.43
04:40:01 AM     all      0.37      0.00      0.15      0.00      0.00     99.48`

	sar := Sar{}

	acc := &testutil.Accumulator{}

	r := bufio.NewReader(strings.NewReader(testdata))

	err := sar.parse(acc, r)
	require.NoError(t, err)
	require.Equal(t, acc.NMetrics(), uint64(10))

	require.True(t, acc.HasTag("sar", "linux_kernel"))
	require.True(t, acc.HasTag("sar", "host"))
	require.True(t, acc.HasTag("sar", "date"))
	require.True(t, acc.HasTag("sar", "arch"))
	require.True(t, acc.HasTag("sar", "num_cpu"))
}
