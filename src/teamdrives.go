// Copyright 2018 Google Inc. All Rights Reserved.
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

// ListTeamDrives lists each team drive and if the user has
// requested recursively traverses the drive to a specified depth.
func (g *Commands) ListTeamDrives() error {
	req := g.rem.service.Teamdrives.List()
}

func teamDriveToFile(td *drive.TeamDrive) *File {
	return &File{
		ModTime: td.CreatedDate.Round(time.Second),
		IsDir:   true,
		Name:    td.Name,
		Id:      td.Id,
		TeamDriveCapabilities: td.TeamDriveCapabilities,
	}
}
