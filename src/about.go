// Copyright 2015 Google Inc. All Rights Reserved.
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

package drive

import (
	"fmt"

	"github.com/odeke-em/log"
	drive "google.golang.org/api/drive/v3"
)

const Version = "0.3.6"

const (
	Barely = iota
	AlmostExceeded
	HalfwayExceeded
	Exceeded
	Unknown
)

const (
	AboutNone = 1 << iota
	AboutQuota
	AboutFileSizes
	AboutFeatures
)

func (g *Commands) About(mask int) (err error) {
	if mask == AboutNone {
		return nil
	}

	about, err := g.rem.About()
	if err != nil {
		return err
	}
	printSummary(g.log, about, mask)

	return nil
}

func quotaRequested(mask int) bool {
	return (mask & AboutQuota) != 0
}

func fileSizesRequested(mask int) bool {
	return (mask & AboutFileSizes) != 0
}

func featuresRequested(mask int) bool {
	return (mask & AboutFeatures) != 0
}

func printSummary(logy *log.Logger, about *drive.About, mask int) {
	if quotaRequested(mask) {
		quotaInformation(logy, about)
	}
	if fileSizesRequested(mask) {
		fileSizesInfo(logy, about)
	}

	if featuresRequested(mask) {
		featuresInformation(logy, about)
	}
}

func fileSizesInfo(logy *log.Logger, about *drive.About) {
	uploadSizes := about.MaxImportSizes
	if len(uploadSizes) >= 1 {
		logy.Logln("\n* Maximum upload sizes per file type *")
		logy.Logf("%-50s %-20s\n", "FileType", "Size")
		for typ, size := range uploadSizes {
			logy.Logf("%-50s %-20s\n", typ, size)
		}
		logy.Logln()
	}
	return
}

func featuresInformation(logy *log.Logger, about *drive.About) {
	// Deprecated, no longer in drive/v3, was in drive/v2
	/*
		if len(about.Features) >= 1 {
			logy.Logf("%-30s %-30s\n", "Feature", "Request limit (queries/second)")
			for _, feature := range about.Features {
				if feature.FeatureName == "" {
					continue
				}
				logy.Logf("%-30s %-30f\n", feature.FeatureName, feature.FeatureRate)
			}
			logy.Logln()
		}
	*/
}

func quotaInformation(logy *log.Logger, about *drive.About) {
	storageQuota := about.StorageQuota
	userInfo := about.User

	freeBytes := storageQuota.Usage - storageQuota.UsageInDrive
	logy.Logf(
		"Name: %s\nBytes Used:\t%-20d (%s)\n"+
			"Bytes Free:\t%-20d (%s)\nBytes InTrash:\t%-20d (%s)\n"+
			"Total Bytes:\t%-20d (%s)\n",
		userInfo.DisplayName,
		storageQuota.UsageInDrive, prettyBytes(storageQuota.UsageInDrive),
		freeBytes, prettyBytes(freeBytes),
		storageQuota.UsageInDriveTrash, prettyBytes(storageQuota.UsageInDriveTrash),
		storageQuota.Usage, prettyBytes(storageQuota.Usage))

	/*
		// Deprecated in v3
		if len(storageQuota.QuotaBytesByService) >= 1 {
			logy.Logln("\n* Space used by Google Services *")
			logy.Logf("%-36s %-36s\n", "Service", "Bytes")
			for _, quotaService := range about.QuotaBytesByService {
				logy.Logf("%-36s %-36s\n", quotaService.ServiceName, prettyBytes(quotaService.BytesUsed))
			}
			logy.Logf("%-36s %-36s\n", "Space used by all Google Apps",
				prettyBytes(storageQuota.Usage))
		}
	*/
	logy.Logln()
}

func (g *Commands) QuotaStatus(query int64) (status int, err error) {
	if query < 0 {
		return Unknown, err
	}

	about, err := g.rem.About()
	if err != nil {
		return Unknown, err
	}

	quota := about.StorageQuota
	// Sanity check
	if quota.Usage < 1 {
		return Unknown, fmt.Errorf("QuotaBytesTotal < 1")
	}

	toBeUsed := query + quota.UsageInDrive
	if toBeUsed >= quota.Usage {
		return Exceeded, nil
	}

	percentage := float64(toBeUsed) / float64(quota.Usage)
	if percentage < 0.5 {
		return Barely, nil
	}
	if percentage < 0.8 {
		return HalfwayExceeded, nil
	}
	return AlmostExceeded, nil
}
