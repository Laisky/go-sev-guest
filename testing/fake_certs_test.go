// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testing

import (
	"crypto/x509"
	"testing"
	"time"

	"github.com/google/go-sev-guest/abi"
	"github.com/google/go-sev-guest/kds"
	"github.com/pborman/uuid"
)

func TestCertificatesParse(t *testing.T) {
	signer, err := DefaultTestOnlyCertChain("Milan", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	certBytes, err := signer.CertTableBytes()
	if err != nil {
		t.Fatal(err)
	}
	entries, err := abi.ParseSnpCertTableHeader(certBytes)
	if err != nil {
		t.Fatal(err)
	}
	var hasVcek bool
	var hasVlek bool
	var hasAsk bool
	var hasAsvk bool
	var hasArk bool
	if len(entries) != 5 {
		t.Errorf("ParseSnpCertTableHeader(_) returned %d entries, want 5", len(entries))
	}
	for _, entry := range entries {
		if uuid.Equal(entry.GUID, uuid.Parse(abi.VlekGUID)) {
			hasVlek = true
		}
		if uuid.Equal(entry.GUID, uuid.Parse(abi.VcekGUID)) {
			hasVcek = true
		}
		if uuid.Equal(entry.GUID, uuid.Parse(abi.AskGUID)) {
			hasAsk = true
		}
		if uuid.Equal(entry.GUID, uuid.Parse(abi.AsvkGUID)) {
			hasAsvk = true
		}
		if uuid.Equal(entry.GUID, uuid.Parse(abi.ArkGUID)) {
			hasArk = true
		}
		der := certBytes[entry.Offset : entry.Offset+entry.Length]
		if _, err := x509.ParseCertificate(der); err != nil {
			t.Errorf("could not parse certificate of %v: %v", entry.GUID, err)
		}
	}
	if !hasVlek {
		t.Errorf("fake certs missing VLEK")
	}
	if !hasVcek {
		t.Errorf("fake certs missing VCEK")
	}
	if !hasAsk {
		t.Errorf("fake certs missing ASK")
	}
	if !hasArk {
		t.Errorf("fake certs missing ARK")
	}
	if !hasAsvk {
		t.Errorf("fake certs missing ASVK")
	}
	if _, err := kds.VcekCertificateExtensions(signer.Vcek); err != nil {
		t.Errorf("could not parse generated VCEK extensions: %v", err)
	}
}
