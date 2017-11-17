
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
https://cloud.google.com/compute/docs/reference/latest/machineTypes

MachineTypes

A Machine Type resource.

{
  "kind": "compute#machineType",
  "id": unsigned long,
  "creationTimestamp": string,
  "name": string,
  "description": string,
  "guestCpus": integer,
  "memoryMb": integer,
  "imageSpaceGb": integer,
  "scratchDisks": [
    {
      "diskGb": integer
    }
  ],
  "maximumPersistentDisks": integer,
  "maximumPersistentDisksSizeGb": long,
  "deprecated": {
    "state": string,
    "replacement": string,
    "deprecated": string,
    "obsolete": string,
    "deleted": string
  },
  "zone": string,
  "selfLink": string,
  "isSharedCpu": boolean
}

aggregatedList
	Retrieves an aggregated list of machine types.
get
	Returns the specified machine type. 
        Get a list of available machine types by making a list() request.
list
	Retrieves a list of machine types available to the specified project.
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~



#############################################################################################################


juju/juju/provider/gce/instance_information.go

// InstanceTypes implements InstanceTypesFetcher
func (env *environ) InstanceTypes(c constraints.Value) (instances.InstanceTypesWithCostMetadata, error) {
        reg, err := env.Region()
        if err != nil {
                return instances.InstanceTypesWithCostMetadata{}, errors.Trace(err)
        }
        zones, err := env.gce.AvailabilityZones(reg.Region)
        if err != nil {
                return instances.InstanceTypesWithCostMetadata{}, errors.Trace(err)
        }
        resultUnique := map[string]instances.InstanceType{}

        for _, z := range zones {
                if !z.Available() {
                        continue
                }
                machines, err := env.gce.ListMachineTypes(z.Name())                          <========================
                if err != nil {
                        return instances.InstanceTypesWithCostMetadata{}, errors.Trace(err)
                }
                for _, m := range machines {
                        i := instances.InstanceType{
                                Id:       strconv.FormatUint(m.Id, 10),
                                Name:     m.Name,
                                CpuCores: uint64(m.GuestCpus),
                                Mem:      uint64(m.MemoryMb),
                                Arches:   []string{arch.AMD64},
                                VirtType: &virtType,
                        }
                        resultUnique[m.Name] = i
                }
        }

        result := make([]instances.InstanceType, len(resultUnique))
        i := 0
        for _, it := range resultUnique {
                result[i] = it
                i++
        }
        result, err = instances.MatchingInstanceTypes(result, "", c)
        if err != nil {
                return instances.InstanceTypesWithCostMetadata{}, errors.Trace(err)
        }
        return instances.InstanceTypesWithCostMetadata{InstanceTypes: result}, nil
}
    
---------------------------------------------------------------------------------------------------------------

juju/juju/provider/gce/environ.go

type environ struct {
        name  string
        uuid  string
        cloud environs.CloudSpec
        gce   gceConnection                                                   <========================

        lock sync.Mutex // lock protects access to ecfg
        ecfg *environConfig

        // namespace is used to create the machine and device hostnames.
        namespace instance.Namespace
}

---------------------------------------------------------------------------------------------------------------

juju/juju/provider/gce/environ.go


import (

        "google.golang.org/api/compute/v1"
        "github.com/juju/juju/provider/gce/google"



type gceConnection interface {
        VerifyCredentials() error

        // Instance gets the up-to-date info about the given instance
        // and returns it.
        Instance(id, zone string) (google.Instance, error)
        Instances(prefix string, statuses ...string) ([]google.Instance, error)
        AddInstance(spec google.InstanceSpec, zone string) (*google.Instance, error)
        RemoveInstances(prefix string, ids ...string) error
        UpdateMetadata(key, value string, ids ...string) error

        IngressRules(fwname string) ([]network.IngressRule, error)
        OpenPorts(fwname string, rules ...network.IngressRule) error
        ClosePorts(fwname string, rules ...network.IngressRule) error

        AvailabilityZones(region string) ([]google.AvailabilityZone, error)
        // Subnetworks returns the subnetworks that machines can be
        // assigned to in the given region.
        Subnetworks(region string) ([]*compute.Subnetwork, error)
        // Networks returns the available networks that exist across
        // regions.
        Networks() ([]*compute.Network, error)

        // Storage related methods.

        // CreateDisks will attempt to create the disks described by <disks> spec and
        // return a slice of Disk representing the created disks or error if one of them failed.
        CreateDisks(zone string, disks []google.DiskSpec) ([]*google.Disk, error)
        // Disks will return a list of all Disks found in the project.
        Disks() ([]*google.Disk, error)
        // Disk will return a Disk representing the disk identified by the
        // passed <name> or error.
        Disk(zone, id string) (*google.Disk, error)
        // RemoveDisk will destroy the disk identified by <name> in <zone>.
        RemoveDisk(zone, id string) error
        // SetDiskLabels sets the labels on a disk, ensuring that the disk's
        // label fingerprint matches the one supplied.
        SetDiskLabels(zone, id, labelFingerprint string, labels map[string]string) error
        // AttachDisk will attach the volume identified by <volumeName> into the instance
        // <instanceId> and return an AttachedDisk representing it or error.
        AttachDisk(zone, volumeName, instanceId string, mode google.DiskMode) (*google.AttachedDisk, error)
        // DetachDisk will detach <volumeName> disk from <instanceId> if possible
        // and return error.
        DetachDisk(zone, instanceId, volumeName string) error
        // InstanceDisks returns a list of the disks attached to the passed instance.
        InstanceDisks(zone, instanceId string) ([]*google.AttachedDisk, error)
        // ListMachineTypes returns a list of machines available in the project and zone provided.
        ListMachineTypes(zone string) ([]google.MachineType, error)                     <========================
}

---------------------------------------------------------------------------------------------------------------

juju/juju/provider/gce/google/conn_machines.go

// ListMachineTypes returns a list of MachineType available for the
// given zone.
func (gce *Connection) ListMachineTypes(zone string) ([]MachineType, error) {
        machines, err := gce.raw.ListMachineTypes(gce.projectID, zone)                   <========================
        if err != nil {
                return nil, errors.Trace(err)
        }
        res := make([]MachineType, len(machines.Items))
        for i, machine := range machines.Items {
                deprecated := false
                if machine.Deprecated != nil {
                        deprecated = machine.Deprecated.State != ""
                }
                res[i] = MachineType{
                        CreationTimestamp: machine.CreationTimestamp,
                        Deprecated:        deprecated,
                        Description:       machine.Description,
                        GuestCpus:         machine.GuestCpus,
                        Id:                machine.Id,
                        ImageSpaceGb:      machine.ImageSpaceGb,
                        Kind:              machine.Kind,
                        MaximumPersistentDisks:       machine.MaximumPersistentDisks,
                        MaximumPersistentDisksSizeGb: machine.MaximumPersistentDisksSizeGb,
                        MemoryMb:                     machine.MemoryMb,
                        Name:                         machine.Name,
                }
        }
        return res, nil
}

-------------------------------------------------------------------------------------------

juju/juju/provider/gce/google/conn.go

// TODO(ericsnow) Add specific error types for common failures
// (e.g. BadRequest, RequestFailed, RequestError, ConnectionFailed)?

// Connection provides methods for interacting with the GCE API. The
// methods are limited to those needed by the juju GCE provider.
//
// Before calling any of the methods, the Connect method should be
// called to authenticate and open the raw connection to the GCE API.
// Otherwise a panic will result.
type Connection struct {
        // TODO(ericsnow) name this something else?
        raw       rawConnectionWrapper                                  <=====================        
        region    string
        projectID string
}


-------------------------------------------------------------------------------------------

juju/juju/provider/gce/google/conn.go


// rawConnectionWrapper facilitates mocking out the GCE API during tests.
type rawConnectionWrapper interface {
        // GetProject sends a request to the GCE API for info about the
        // specified project. If the project does not exist then an error
        // will be returned.
        GetProject(projectID string) (*compute.Project, error)

        // GetInstance sends a request to the GCE API for info about the
        // specified instance. If the instance does not exist then an error
        // will be returned.
        GetInstance(projectID, id, zone string) (*compute.Instance, error)

        // ListInstances sends a request to the GCE API for a list of all
        // instances in project for which the name starts with the provided
        // prefix. The result is also limited to those instances with one of
        // the specified statuses (if any).
        ListInstances(projectID, prefix string, status ...string) ([]*compute.Instance, error)

        // AddInstance sends a request to GCE to add a new instance to the
        // given project, with the provided instance data. The call blocks
        // until the instance is created or the request fails.
        AddInstance(projectID, zone string, spec *compute.Instance) error

        // RemoveInstance sends a request to the GCE API to remove the instance
        // with the provided ID (in the specified zone). The call blocks until
        // the instance is removed (or the request fails).
        RemoveInstance(projectID, id, zone string) error

        // SetMetadata sends a request to the GCE API to update one
        // instance's metadata. The call blocks until the request is
        // completed or fails.
        SetMetadata(projectID, zone, instanceID string, metadata *compute.Metadata) error

        // GetFirewalls sends an API request to GCE for the information about
        // the firewalls with the namePrefix and returns them.
        // If no firewalls are not found, errors.NotFound is returned.
        GetFirewalls(projectID, namePrefix string) ([]*compute.Firewall, error)

        // AddFirewall requests GCE to add a firewall with the provided info.
        // If the firewall already exists then an error will be returned.
        // The call blocks until the firewall is added or the request fails.
        AddFirewall(projectID string, firewall *compute.Firewall) error

        // UpdateFirewall requests GCE to update the named firewall with the
        // provided info, overwriting the existing data. If the firewall does
        // not exist then an error will be returned. The call blocks until the
        // firewall is updated or the request fails.
        UpdateFirewall(projectID, name string, firewall *compute.Firewall) error

        // RemoveFirewall removed the named firewall from the project. If it
        // does not exist then this is a noop. The call blocks until the
        // firewall is added or the request fails.
        RemoveFirewall(projectID, name string) error

        // ListAvailabilityZones returns the list of availability zones for a given
        // GCE region. If none are found the the list is empty. Any failure in
        // the low-level request is returned as an error.
        ListAvailabilityZones(projectID, region string) ([]*compute.Zone, error)

        // CreateDisk will create a gce Persistent Block device that matches
        // the specified in spec.
        CreateDisk(project, zone string, spec *compute.Disk) error

        // ListDisks returns a list of disks available for a given project.
        ListDisks(project string) ([]*compute.Disk, error)

        // RemoveDisk will delete the disk identified by id.
        RemoveDisk(project, zone, id string) error

        // GetDisk will return the disk correspondent to the passed id.
        GetDisk(project, zone, id string) (*compute.Disk, error)

        // SetDiskLabels sets the labels on a disk, ensuring that the disk's
        // label fingerprint matches the one supplied.
        SetDiskLabels(project, zone, id, labelFingerprint string, labels map[string]string) error

        // AttachDisk will attach the disk described in attachedDisks (if it exists) into
        // the instance with id instanceId.
        AttachDisk(project, zone, instanceId string, attachedDisk *compute.AttachedDisk) error

        // Detach disk detaches device diskDeviceName (if it exists and its attached)
        // form the machine with id instanceId.
        DetachDisk(project, zone, instanceId, diskDeviceName string) error

        // InstanceDisks returns the disks attached to the instance identified
        // by instanceId
        InstanceDisks(project, zone, instanceId string) ([]*compute.AttachedDisk, error)

        // ListMachineTypes returns a list of machines available in the project and zone provided.
        ListMachineTypes(projectID, zone string) (*compute.MachineTypeList, error)                 <===============

        // ListSubnetworks returns a list of subnets available in the given project and region.
        ListSubnetworks(projectID, region string) ([]*compute.Subnetwork, error)

        // ListNetworks returns a list of Networks available in the given project.
        ListNetworks(projectID string) ([]*compute.Network, error)
}



-------------------------------------------------------------------------------------------


juju/juju/provider/gce/google/raw.go

import (
        "google.golang.org/api/compute/v1"


type rawConn struct {
        *compute.Service
}


// ListMachineTypes returns a list of machines available in the project and zone provided.
func (rc *rawConn) ListMachineTypes(projectID, zone string) (*compute.MachineTypeList, error) {     
        op := rc.MachineTypes.List(projectID, zone)                                               <==============
        machines, err := op.Do()                                                                  <==============
        if err != nil {
                return nil, errors.Annotatef(err, "listing machine types for project %q and zone %q", projectID, zone)
        }
        return machines, nil
}


-------------------------------------------------------------------------------------------


google.golang.org/api/compute/v1/compute-gen.go

package compute


type Service struct {
        client    *http.Client
        BasePath  string // API endpoint base URL
        UserAgent string // optional additional User-Agent fragment

        AcceleratorTypes *AcceleratorTypesService

        Addresses *AddressesService

        Autoscalers *AutoscalersService

        BackendBuckets *BackendBucketsService

        BackendServices *BackendServicesService

        DiskTypes *DiskTypesService

        Disks *DisksService

        Firewalls *FirewallsService

        ForwardingRules *ForwardingRulesService

        GlobalAddresses *GlobalAddressesService

        GlobalForwardingRules *GlobalForwardingRulesService

        GlobalOperations *GlobalOperationsService

        HealthChecks *HealthChecksService

        HttpHealthChecks *HttpHealthChecksService

        HttpsHealthChecks *HttpsHealthChecksService

        Images *ImagesService

        InstanceGroupManagers *InstanceGroupManagersService

        InstanceGroups *InstanceGroupsService

        InstanceTemplates *InstanceTemplatesService

        Instances *InstancesService

        InterconnectAttachments *InterconnectAttachmentsService

        InterconnectLocations *InterconnectLocationsService

        Interconnects *InterconnectsService

        Licenses *LicensesService

        MachineTypes *MachineTypesService                       <=======================

        Networks *NetworksService

        Projects *ProjectsService

        RegionAutoscalers *RegionAutoscalersService

        RegionBackendServices *RegionBackendServicesService

        RegionCommitments *RegionCommitmentsService

        RegionInstanceGroupManagers *RegionInstanceGroupManagersService

        RegionInstanceGroups *RegionInstanceGroupsService

        RegionOperations *RegionOperationsService

        Regions *RegionsService

        Routers *RoutersService

        Routes *RoutesService

        Snapshots *SnapshotsService

        SslCertificates *SslCertificatesService

        Subnetworks *SubnetworksService

        TargetHttpProxies *TargetHttpProxiesService

        TargetHttpsProxies *TargetHttpsProxiesService

        TargetInstances *TargetInstancesService

        TargetPools *TargetPoolsService

        TargetSslProxies *TargetSslProxiesService

        TargetTcpProxies *TargetTcpProxiesService

        TargetVpnGateways *TargetVpnGatewaysService

        UrlMaps *UrlMapsService

        VpnTunnels *VpnTunnelsService

        ZoneOperations *ZoneOperationsService

        Zones *ZonesService
}


-------------------------------------------------------------------------------------------


google.golang.org/api/compute/v1/compute-gen.go

type MachineTypesService struct {
        s *Service                                    <==================
}


-------------------------------------------------------------------------------------------


google.golang.org/api/compute/v1/compute-gen.go

// List: Retrieves a list of machine types available to the specified
// project.
// For details, see https://cloud.google.com/compute/docs/reference/latest/machineTypes/list
func (r *MachineTypesService) List(project string, zone string) *MachineTypesListCall {        <=============
        c := &MachineTypesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
        c.project = project
        c.zone = zone
        return c


-------------------------------------------------------------------------------------------


google.golang.org/api/compute/v1/compute-gen.go


// method id "compute.machineTypes.list":
type MachineTypesListCall struct {
        s            *Service
        project      string
        zone         string
        urlParams_   gensupport.URLParams
        ifNoneMatch_ string
        ctx_         context.Context
        header_      http.Header
}


-------------------------------------------------------------------------------------------


google.golang.org/api/compute/v1/compute-gen.go

// Do executes the "compute.machineTypes.list" call.
// Exactly one of *MachineTypeList or error will be non-nil. Any non-2xx
// status code is an error. Response headers are in either
// *MachineTypeList.ServerResponse.Header or (if a response was returned
// at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *MachineTypesListCall) Do(opts ...googleapi.CallOption) (*MachineTypeList, error) {
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
        ret := &MachineTypeList{
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


-------------------------------------------------------------------------------------------
