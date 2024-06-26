package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/go-test/deep"
)

func TestDecryptFiles(t *testing.T) {
	testCases := []struct {
		name            string
		initFile        string
		inFile          string
		expectedOutFile string
		hexKey          string
	}{
		{
			name:            "cenc",
			inFile:          "../../mp4/testdata/prog_8s_enc_dashinit.mp4",
			expectedOutFile: "../../mp4/testdata/prog_8s_dec_dashinit.mp4",
			hexKey:          "63cb5f7184dd4b689a5c5ff11ee6a328",
		},
		{
			name:            "cbcs",
			inFile:          "../../mp4/testdata/cbcs.mp4",
			expectedOutFile: "../../mp4/testdata/cbcsdec.mp4",
			hexKey:          "22bdb0063805260307ee5045c0f3835a",
		},
		{
			name:            "cbcs audio",
			inFile:          "../../mp4/testdata/cbcs_audio.mp4",
			expectedOutFile: "../../mp4/testdata/cbcs_audiodec.mp4",
			hexKey:          "5ffd93861fa776e96cccd934898fc1c8",
		},
		{
			name:            "PIFF audio",
			initFile:        "testdata/PIFF/audio/init.mp4",
			inFile:          "testdata/PIFF/audio/segment-1.0001.m4s",
			expectedOutFile: "testdata/PIFF/audio/segment-1.0001_dec.m4s",
			hexKey:          "602a9289bfb9b1995b75ac63f123fc86",
		},
		{
			name:            "PIFF video",
			inFile:          "testdata/PIFF/video/complseg-1.0001.mp4",
			expectedOutFile: "testdata/PIFF/video/complseg-1.0001_dec.mp4",
			hexKey:          "602a9289bfb9b1995b75ac63f123fc86",
		},
	}

	for _, tc := range testCases {
		ifh, err := os.Open(tc.inFile)
		if err != nil {
			t.Error(err)
		}
		buf := bytes.Buffer{}
		var initFH *os.File
		if tc.initFile != "" {
			initFH, err = os.Open(tc.initFile)
			if err != nil {
				t.Error(err)
			}
		}
		err = decryptFile(ifh, initFH, &buf, tc.hexKey)
		ifh.Close()
		if err != nil {
			t.Error(err)
		}
		expectedOut, err := os.ReadFile(tc.expectedOutFile)
		if err != nil {
			t.Error(err)
		}
		gotOut := buf.Bytes()
		diff := deep.Equal(expectedOut, gotOut)
		if diff != nil {
			t.Errorf("Mismatch for case %s: %s", tc.name, diff)
		}
	}
}

func BenchmarkDecodeCenc(b *testing.B) {
	inFile := "../../mp4/testdata/prog_8s_enc_dashinit.mp4"
	hexKey := "63cb5f7184dd4b689a5c5ff11ee6a328"
	raw, err := os.ReadFile(inFile)
	if err != nil {
		b.Error(err)
	}
	outData := make([]byte, 0, len(raw))
	outBuf := bytes.NewBuffer(outData)
	for i := 0; i < b.N; i++ {
		inBuf := bytes.NewBuffer(raw)
		outBuf.Reset()
		err = decryptFile(inBuf, nil, outBuf, hexKey)
		if err != nil {
			b.Error(err)
		}
	}
}
