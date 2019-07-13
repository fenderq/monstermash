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
		"W2FCO 3FGGQ QIVGJ R4O45",
		"NJOWB OWQ7E 2V4IV FBZYX",
		"GTS4Q 7YJ3M VQ7JR MGVEJ",
		"7N7K4 U3RVO JP623 443TL",
		"XOXOE TLBR6 MEAHZ DEVLH",
		"WHEWO I6FAI V7RH7 DV2VP",
		"GV7KT VA2ZK 65DDK 7Z2IF",
		"H5VVL P6NTZ 3WX6R O3IE7",
		"VT4NM VW6GF K4LDV MQVP6",
		"UM7KQ VOXRN G6WEG 235OI",
	},
	}
	if err := testItem(&td); err != nil {
		t.Error(err)
	}
}

func TestRandBar(t *testing.T) {
	td := testData{"randfile", "bar", []string{
		"OMNI4 XDEAW RZ2OQ G42QE",
		"BPAJA IJJLO PXG67 I6WKD",
		"CYLE4 SQOQN SXFSD 275EO",
		"7AANP JWSOL 2NT4M QEECK",
		"4XZMB MVGN3 PZX3D 6NH77",
		"IAO65 UFLFL 4UFLD 4UL7G",
		"RBUIZ 4PCCW MJUFK 43R3N",
		"MSJXG U552U 2MPSX YNA6F",
		"NDGSK D6NX2 7MBWD 43CLB",
		"6O7FO 4G2IQ 75YSU XJH5H",
	},
	}
	if err := testItem(&td); err != nil {
		t.Error(err)
	}
}

func TestZeroFoo(t *testing.T) {
	td := testData{"zerofile", "foo", []string{
		"3NIDU GY4WE XWQYJ QPMRZ",
		"JOS6E LGMR5 N5KIF IVWV5",
		"CYY77 G4KLN ZLVDI VBWUC",
		"PKE3K 4WYPY VOBY5 3GN22",
		"43R7C BR2ZW XNUTX DDSE5",
		"EZAP6 72DSV ARXDU CSLF5",
		"QPTMA IBZVD GJ7PG J2BJN",
		"4KDD4 GZR3O ELS7W B7LHT",
		"RSLKU EOCQE R4XGN WAQQN",
		"LQ2AG NJZBN 2CUOU DM6PK",
	},
	}
	if err := testItem(&td); err != nil {
		t.Error(err)
	}
}

func TestZeroBar(t *testing.T) {
	td := testData{"zerofile", "bar", []string{
		"JGU7G NO6SI NHNK7 TS6PM",
		"H5KMJ PRYFW I6HQ4 3WVYX",
		"LVUKB CSBP6 RYVG4 ZAOTZ",
		"XEO2H 4T6MY 6WXCH 7H5E4",
		"643VC 4I335 EJ4YZ COG6U",
		"64UTU TEFHS JJZOD UY2A2",
		"L3FL5 Y56U4 ZMHOH JDKUM",
		"DAN7A 6ILRO LMRC3 NE5IE",
		"NZSIJ Z3RGH UGC5Y AR43K",
		"TEMBZ U6OO5 3KIVY 5YCBB",
	},
	}
	if err := testItem(&td); err != nil {
		t.Error(err)
	}
}

func testItem(td *testData) error {
	path := filepath.Join("files", td.filename)
	salt, err := GetSaltFromFile(path)
	if err != nil {
		return err
	}
	s, err := MakePasswords(salt, []byte(td.passwd))
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
