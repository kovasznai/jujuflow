
juju/juju/provider/gce/disks.go

func (v *volumeSource) ListVolumes() ([]string, error) {
        var volumes []string
        disks, err := v.gce.Disks()                  <==============
        if err != nil {
                return nil, errors.Trace(err)
        }
        for _, disk := range disks {
                if !isValidVolume(disk.Name) {
                        continue
                }
                if disk.Labels[tags.JujuModel] != v.modelUUID {
                        continue
                }
                volumes = append(volumes, disk.Name)
        }
        return volumes, nil
}


---------------------------------------------------------------------


./disks.go

type volumeSource struct {
        gce       gceConnection                             <==============
        envName   string // non-unique, informational only
        modelUUID string
}


---------------------------------------------------------------------


juju/juju/provider/gce/environ.go

type gceConnection interface {

        // Disks will return a list of all Disks found in the project.
        Disks() ([]*google.Disk, error)                       <==============


---------------------------------------------------------------------


juju/juju/provider/gce/google/conn_disks.go

// Disks implements storage section of gceConnection.
func (gce *Connection) Disks() ([]*Disk, error) {
        computeDisks, err := gce.raw.ListDisks(gce.projectID)     <==============
        if err != nil {
                return nil, errors.Annotate(err, "cannot list disks")
        }
        disks := make([]*Disk, len(computeDisks))
        for i, disk := range computeDisks {
                disks[i] = NewDisk(disk)
        }
        return disks, nil
}


---------------------------------------------------------------------


juju/juju/provider/gce/google/raw.go

package google


func (rc *rawConn) ListDisks(project string) ([]*compute.Disk, error) {
        ds := rc.Service.Disks              <======== [1]
        call := ds.AggregatedList(project)       <======== [2]  *DisksAggregatedListCall
        var results []*compute.Disk         <======== [3]
        for {
                diskList, err := call.Do()    <======== [4]  *DiskAggregatedList <- func (c *DisksAggregatedListCall) Do(
                if err != nil {
                        return nil, errors.Trace(err)
                }
                for _, list := range diskList.Items {        <========= [5]  diskList: *DiskAggregatedList
                        results = append(results, list.Disks...)   <============= list.Disks [6]
                }
                if diskList.NextPageToken == "" {
                        break
                }
                call = call.PageToken(diskList.NextPageToken)        <=========== [7]
        }
        return results, nil
}


---------------------------------------------------------------------


juju/juju/provider/gce/google/raw.go

import (
        "google.golang.org/api/compute/v1"


type rawConn struct {
        *compute.Service
}


---------------------------------------------------------------------
[1]
google.golang.org/api/compute/v1/compute-gen.go

package compute


type Service struct {

        Disks *DisksService

---------------------------------------------------------------------


google.golang.org/api/compute/v1/compute-gen.go


type DisksService struct {
        s *Service
}


---------------------------------------------------------------------
[2]

google.golang.org/api/compute/v1/compute-gen.go


// AggregatedList: Retrieves an aggregated list of persistent disks.
// For details, see https://cloud.google.com/compute/docs/reference/latest/disks/aggregatedList
func (r *DisksService) AggregatedList(project string) *DisksAggregatedListCall {
        c := &DisksAggregatedListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
        c.project = project
        return c
}


---------------------------------------------------------------------
[3]

google.golang.org/api/compute/v1/compute-gen.go

package compute

// Disk: A Disk resource.
type Disk struct {
        CreationTimestamp string `json:"creationTimestamp,omitempty"`

        Description string `json:"description,omitempty"`

        DiskEncryptionKey *CustomerEncryptionKey `json:"diskEncryptionKey,omitempty"`

        Id uint64 `json:"id,omitempty,string"`

        Kind string `json:"kind,omitempty"`

        LabelFingerprint string `json:"labelFingerprint,omitempty"`

        Labels map[string]string `json:"labels,omitempty"`

        LastAttachTimestamp string `json:"lastAttachTimestamp,omitempty"`

        Licenses []string `json:"licenses,omitempty"`

        Name string `json:"name,omitempty"`

        Options string `json:"options,omitempty"`

        SelfLink string `json:"selfLink,omitempty"`

        SizeGb int64 `json:"sizeGb,omitempty,string"`

        SourceImage string `json:"sourceImage,omitempty"`

        SourceImageEncryptionKey *CustomerEncryptionKey `json:"sourceImageEncryptionKey,omitempty"`

        SourceImageId string `json:"sourceImageId,omitempty"`

        SourceSnapshot string `json:"sourceSnapshot,omitempty"`

        SourceSnapshotEncryptionKey *CustomerEncryptionKey `json:"sourceSnapshotEncryptionKey,omitempty"`

        SourceSnapshotId string `json:"sourceSnapshotId,omitempty"`

        Status string `json:"status,omitempty"`

        Type string `json:"type,omitempty"`

        Users []string `json:"users,omitempty"`

        Zone string `json:"zone,omitempty"`

        googleapi.ServerResponse `json:"-"`

        ForceSendFields []string `json:"-"`

        NullFields []string `json:"-"`
}

---------------------------------------------------------------------

google.golang.org/api/compute/v1/compute-gen.go

type DisksAggregatedListCall struct {
        s            *Service
        project      string
        urlParams_   gensupport.URLParams
        ifNoneMatch_ string
        ctx_         context.Context
        header_      http.Header
}


---------------------------------------------------------------------
[4]

google.golang.org/api/compute/v1/compute-gen.go

// Do executes the "compute.disks.aggregatedList" call.
// Exactly one of *DiskAggregatedList or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *DiskAggregatedList.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *DisksAggregatedListCall) Do(opts ...googleapi.CallOption) (*DiskAggregatedList, error) {
        gensupport.SetOptions(c.urlParams_, opts...)
        res, err := c.doRequest("json")
        if res != nil && res.StatusCode == http.StatusNotModified {
                if res.Body != nil {
                        res.Body.Close()
                }
                return nil, &googleapi.Error{
                        Code:   res.StatusCode,
                        Header: res.Header,
                }
        }
        if err != nil {
                return nil, err
        }
        defer googleapi.CloseBody(res)
        if err := googleapi.CheckResponse(res); err != nil {
                return nil, err
        }
        ret := &DiskAggregatedList{
                ServerResponse: googleapi.ServerResponse{
                        Header:         res.Header,
                        HTTPStatusCode: res.StatusCode,
                },
        }
        target := &ret
        if err := json.NewDecoder(res.Body).Decode(target); err != nil {
                return nil, err
        }
        return ret, nil
}

---------------------------------------------------------------------
[5]

google.golang.org/api/compute/v1/compute-gen.go

package compute

import (
        googleapi "google.golang.org/api/googleapi"


type DiskAggregatedList struct {

        Id string `json:"id,omitempty"`

        Items map[string]DisksScopedList `json:"items,omitempty"`   <============ [5]

        Kind string `json:"kind,omitempty"`

        NextPageToken string `json:"nextPageToken,omitempty"`

        SelfLink string `json:"selfLink,omitempty"`

        Warning *DiskAggregatedListWarning `json:"warning,omitempty"`

        googleapi.ServerResponse `json:"-"`

        ForceSendFields []string `json:"-"`

        NullFields []string `json:"-"`
}

---------------------------------------------------------------------

google.golang.org/api/googleapi/googleapi.go

// ServerResponse is embedded in each Do response and
// provides the HTTP status code and header sent by the server.
type ServerResponse struct {
        // HTTPStatusCode is the server's response status code.
        // When using a resource method's Do call, this will always be in the 2xx range.
        HTTPStatusCode int
        // Header contains the response header fields from the server.
        Header http.Header
}


---------------------------------------------------------------------
[6]

google.golang.org/api/compute/v1/compute-gen.go

type DisksScopedList struct {

        // Disks: [Output Only] List of disks contained in this scope.
        Disks []*Disk `json:"disks,omitempty"`

---------------------------------------------------------------------
[7]

// PageToken sets the optional parameter "pageToken": Specifies a page
// token to use. Set pageToken to the nextPageToken returned by a
// previous list request to get the next page of results.
func (c *DisksAggregatedListCall) PageToken(pageToken string) *DisksAggregatedListCall {
        c.urlParams_.Set("pageToken", pageToken)
        return c
}

---------------------------------------------------------------------




