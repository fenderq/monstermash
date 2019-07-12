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
	"path/filepath"
	"testing"
)

type testItem struct {
	filename string
	passwd   string
	results  []string
}

func TestMonsterMash(t *testing.T) {
	tests := []testItem{
		{"zerofile", "foo", []string{
			"6EKHK H5FEB BTCGQ UTJ4M",
			"2I7CX M6OGU 5UMEW M3YHE",
			"DAYVS 5C3OD LVOWO LBBGN",
			"XUABI 2XWFD VD4J7 CEX7B",
			"F5K66 JVHOS 3P2JX LB67P",
			"DAJZ3 MLCPE XTE25 5U2TG",
			"3M74W BQS2T UAOLA CJFNA",
			"6B5RH 7I2GU D3DMU F77ZN",
			"FMIRB SBUPH SOSTY 3LVA2",
			"555DL AW7JB XGEBF WVPYT",
		},
		},
		{"zerofile", "bar", []string{
			"55OQX N7JNG 52NBE GMXBF",
			"GWTJN 2JCPM SCLGW OGWJU",
			"W7FV2 H64IQ X2GTF FNXHK",
			"HIMCX MXKMH RYEVI T5Z36",
			"SZZFH JD52Q HBRMQ H6LF2",
			"4F2TP P2VAX XVXMJ VURQO",
			"MOX7R 4F7ZB 7R57E FIPJG",
			"II2MV JJXEL M6RT5 5SSIR",
			"2HKM6 DPTXQ LHN6Q PC2QR",
			"HCIXV 22ZBI LP2TP 77ABU",
		},
		},
		{"randfile", "foo", []string{
			"UBSAQ JNIDM 7ALY3 YRHCJ",
			"JVCOY F6F3K IXOZW GGE4Z",
			"CZOKP 2XPYJ SET6M 3DNJ5",
			"XMA3L W4UXF WYDID NN6JF",
			"G3KMZ K3UT6 JB2L4 BH6O6",
			"O5XH4 PT6KR 6YNRE MAPUT",
			"HHIGZ 5V4XQ XCULG ZWBZ3",
			"BMX2X MQENJ 6JAMS D2BVK",
			"G6RRJ ZPGMW FKPRA YNNSC",
			"RR2AJ OXU45 IPQZZ SCJZ7",
		},
		},
		{"randfile", "bar", []string{
			"AFWZB 6SM3J CPH5S 5H5DF",
			"LJPV4 ORSAA JKXTB CD76D",
			"QVVOG EV4MJ S24VZ TBXAR",
			"OJR25 7DNUF KRNXP 2XUH7",
			"UATE3 WZWMK 2AX4U IGEKZ",
			"5SE66 RNYEV B2NKC 2YEQL",
			"IJ4QM U5IWF N75QJ HQZ37",
			"22C2P ZCWQW E6XNJ LT6GP",
			"55ZU6 CE33Z M4V2P HU5U7",
			"SBA3A SAGPH WATHJ PLG5L",
		},
		},
	}
	for i, v := range tests {
		path := filepath.Join("files", v.filename)
		t.Log("GetDataFromFile:", path)
		data, err := GetDataFromFile(path)
		if err != nil {
			t.Error(err)
		}
		t.Log("MakePasswords:", v.passwd)
		s, err := MakePasswords(data, []byte(v.passwd))
		if err != nil {
			t.Error(err)
		}
		for j, s1 := range s {
			s2 := v.results[j]
			if s1 != s2 {
				t.Errorf("%s should be %s "+
					"in test %d, password #%d",
					s1, s2, i+1, j+1)
			}
		}
	}
}
