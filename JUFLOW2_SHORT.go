
juju/juju/provider/gce/instance_information.go

func (env *environ) InstanceTypes(c constraints.Value) (instances.InstanceTypesWithCostMetadata, error) {

                machines, err := env.gce.ListMachineTypes(z.Name()
|
|
|
juju/juju/provider/gce/environ.go

type environ struct {

        gce   gceConnection
|
|
|
juju/juju/provider/gce/environ.go
import (
        "google.golang.org/api/compute/v1"
        "github.com/juju/juju/provider/gce/google"

type gceConnection interface {

        ListMachineTypes(zone string) ([]google.MachineType, error)
|
|
|
juju/juju/provider/gce/google/conn_machines.go

func (gce *Connection) ListMachineTypes(zone string) ([]MachineType, error)

        machines, err := gce.raw.ListMachineTypes(gce.projectID, zone)
|
|
|
juju/juju/provider/gce/google/conn.go

type Connection struct {

        raw       rawConnectionWrapper 
|
|
|
juju/juju/provider/gce/google/conn.go

type rawConnectionWrapper interface {

        ListMachineTypes(projectID, zone string) (*compute.MachineTypeList, error)
|
|
|
juju/juju/provider/gce/google/raw.go

type rawConn struct {
        *compute.Service

func (rc *rawConn) ListMachineTypes(projectID, zone string) (*compute.MachineTypeList, error)

        op := rc.MachineTypes.List(projectID, zone)
        machines, err := op.Do()    
|
|
|
google.golang.org/api/compute/v1/compute-gen.go

type Service struct {

        MachineTypes *MachineTypesService
|
|
|
google.golang.org/api/compute/v1/compute-gen.go

type MachineTypesService struct {
        s *Service              
|
|
|
google.golang.org/api/compute/v1/compute-gen.go

func (r *MachineTypesService) List(project string, zone string) *MachineTypesListCall

        c := &MachineTypesListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
|
|
|
google.golang.org/api/compute/v1/compute-gen.go

type MachineTypesListCall struct {
|
|
|
google.golang.org/api/compute/v1/compute-gen.go

func (c *MachineTypesListCall) Do(opts ...googleapi.CallOption) (*MachineTypeList, error) {
