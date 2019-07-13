//
// Copyright (c) 2019 Steven Roberts <sroberts@fenderq.com>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
//

package main

import (
	"fmt"
	"path/filepath"
	"testing"
)

type testData struct {
	filename string
	passwd   string
	results  []string
}

func TestRandFoo(t *testing.T) {
	td := testData{"randfile", "foo", []string{
		"3TP3A JORV7 5IIAY 3CKP3",
		"6MLD2 HOTU7 CGQUY GAJ67",
		"4BVXO BLDY7 G3XC6 5OCGM",
		"DGUQF MSC23 H2VML 6BQG2",
		"H6AQN JKSSP 3LGEU 6JANJ",
		"EHR6D IGQ6A T6ADE 6PORU",
		"FYXOS ITX2K VYS3S JX7VR",
		"CNHVJ SDCIH WEF6Q R6G2V",
		"2N6P6 LKIC7 6FEEP SZRXB",
		"7MDBI 6TY6G EB2YV KZET6",
	},
	}
	if err := testItem(&td); err != nil {
		t.Error(err)
	}
}

func TestRandBar(t *testing.T) {
	td := testData{"randfile", "bar", []string{
		"INJMS SAR3E MKIWR HFX3A",
		"IMY5M UVDOQ GWEXS TUNGW",
		"GMNJX IJKRP AFBSP ITAOS",
		"PIPLD IRQQI EBOG2 HXA3Q",
		"QOCN3 5TO2K 5RMCR UJIPO",
		"FU7II 6RN4E SDR2Q RJMQD",
		"6ET3P XLFAO FZ4WV EQ33Y",
		"QSFO4 LUMPQ 4YXKA HC6HH",
		"LFD76 C2U56 KDX7A UW6U7",
		"STBHH GGTN4 KJYBO RVEBD",
	},
	}
	if err := testItem(&td); err != nil {
		t.Error(err)
	}
}

func TestZeroFoo(t *testing.T) {
	td := testData{"zerofile", "foo", []string{
		"QIPW4 WK7T2 AZCQP 2QUGZ",
		"SPPIJ XGQA3 RAUST WN2BO",
		"WTS2G Z5GIZ VNECD H2BMN",
		"UTBUD O54HH KUYW4 GFHFR",
		"DPTOU ABMQA MRNUB PEWTF",
		"Z6M3C Z677S XIR7M 4ESMP",
		"RBQ2G 3MXDD FIRC7 ELIMC",
		"AVGQJ 54ACY CJ4YX COWBS",
		"5GBNO L4CNN VSCVH 3P7UB",
		"AOMW4 R77AX ETA2K Z6ZXL",
	},
	}
	if err := testItem(&td); err != nil {
		t.Error(err)
	}
}

func TestZeroBar(t *testing.T) {
	td := testData{"zerofile", "bar", []string{
		"YGRXM KGIFK ZJE2Q XKTNU",
		"MLQJR HRFHR I2ZXX 7G3KB",
		"CUVDL RM2CF XY2QO MUXUM",
		"5OLHA EFVON Y4C7R 7TSLI",
		"54NLE PGOTX 4NT3Y 4FDDI",
		"PPG4L H7K5V TCU7Z ZL7UO",
		"KZ55J VYT27 RZSIK 35FDW",
		"QISFE DRLDG ZMQIJ AF2ZD",
		"VEA4Q LVFG2 UPCNL T2LCJ",
		"2NRZO 3UOL5 BJDDU Z33K6",
	},
	}
	if err := testItem(&td); err != nil {
		t.Error(err)
	}
}

func testItem(td *testData) error {
	path := filepath.Join("files", td.filename)
	data, err := GetDataFromFile(path)
	if err != nil {
		return err
	}
	s, err := MakePasswords(data, []byte(td.passwd))
	if err != nil {
		return err
	}
	for i, s1 := range s {
		s2 := td.results[i]
		if s1 != s2 {
			return fmt.Errorf("password #%d: %s should be %s",
				i+1, s1, s2)
		}
	}
	return nil
}
