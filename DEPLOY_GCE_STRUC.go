
juju/juju/cmd/juju/application/register.go

// RegisterMeteredCharm implements the DeployStep interface.
type RegisterMeteredCharm struct {
        Plan           string
        IncreaseBudget int
        RegisterURL    string
        QueryURL       string
        credentials    []byte
}

---------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go


        "github.com/juju/juju/cmd/modelcmd"


type DeployCommand struct {
        modelcmd.ModelCommandBase
        UnitCommandBase

        // CharmOrBundle is either a charm URL, a path where a charm can be found,
        // or a bundle name.
        CharmOrBundle string

        // BundleOverlay refers to config files that specify additional bundle
        // configuration to be merged with the main bundle.
        BundleOverlayFile []string

        // Channel holds the charmstore channel to use when obtaining
        // the charm to be deployed.
        Channel params.Channel

        // Series is the series of the charm to deploy.
        Series string

        // Force is used to allow a charm to be deployed onto a machine
        // running an unsupported series.
        Force bool

        // DryRun is used to specify that the bundle shouldn't actually be
        // deployed but just output the changes.
        DryRun bool

        ApplicationName string
        ConfigOptions   common.ConfigFlag
        ConstraintsStr  string
        Constraints     constraints.Value
        BindToSpaces    string

        // TODO(axw) move this to UnitCommandBase once we support --storage
        // on add-unit too.
        //
        // Storage is a map of storage constraints, keyed on the storage name
        // defined in charm storage metadata.
        Storage map[string]storage.Constraints

        // BundleStorage maps application names to maps of storage constraints keyed on
        // the storage name defined in that application's charm storage metadata.
        BundleStorage map[string]map[string]storage.Constraints

        // Resources is a map of resource name to filename to be uploaded on deploy.
        Resources map[string]string

        Bindings map[string]string
        Steps    []DeployStep

        // UseExisting machines when deploying the bundle.
        UseExisting bool
        // BundleMachines is a mapping for machines in the bundle to machines
        // in the model.
        BundleMachines map[string]string

        // NewAPIRoot stores a function which returns a new API root.
        NewAPIRoot func() (DeployAPI, error)

        machineMap string
        flagSet    *gnuflag.FlagSet
}



---------------------------------------------------------------------------------

juju/juju/cmd/modelcmd/modelcommand.go

// ModelCommandBase is a convenience type for embedding in commands
// that wish to implement ModelCommand.
type ModelCommandBase struct {
        CommandBase

        // store is the client controller store that contains information
        // about controllers, models, etc.
        store jujuclient.ClientStore

        // _modelName and _controllerName hold the current
        // model and controller names. They are only valid
        // after initModel is called, and should in general
        // not be accessed directly, but through ModelName and
        // ControllerName respectively.
        _modelName      string
        _controllerName string

        allowDefaultModel bool

        // doneInitModel holds whether initModel has been called.
        doneInitModel bool

        // initModelError holds the result of the initModel call.
        initModelError error
}


---------------------------------------------------------------------------------

juju/juju/cmd/modelcmd/base.go

        "github.com/juju/cmd"

// CommandBase is a convenience type for embedding that need
// an API connection.
type CommandBase struct {
        cmd.CommandBase
        cmdContext    *cmd.Context
        apiContexts   map[string]*apiContext
        modelAPI_     ModelAPI
        apiOpenFunc   api.OpenFunc
        authOpts      AuthOpts
        runStarted    bool
        refreshModels func(jujuclient.ClientStore, string) error
}

---------------------------------------------------------------------------------

juju/cmd/cmd.go

// CommandBase provides the default implementation for SetFlags, Init, and Help.
type CommandBase struct{}


---------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go


// DeployStep is an action that needs to be taken during charm deployment.
type DeployStep interface {

        // Set flags necessary for the deploy step.
        SetFlags(*gnuflag.FlagSet)

        // RunPre runs before the call is made to add the charm to the environment.
        RunPre(MeteredDeployAPI, *httpbakery.Client, *cmd.Context, DeploymentInfo) error

        // RunPost runs after the call is made to add the charm to the environment.
        // The error parameter is used to notify the step of a previously occurred error.
        RunPost(MeteredDeployAPI, *httpbakery.Client, *cmd.Context, DeploymentInfo, error) error
}

---------------------------------------------------------------------------------
juju/juju/cmd/juju/application/deploy.go

        "github.com/juju/juju/api"

type deployAPIAdapter struct {
        api.Connection
        *apiClient
        *charmsClient
        *applicationClient
        *modelConfigClient
        *charmRepoClient
        *charmstoreClient
        *annotationsClient
}

---------------------------------------------------------------------------------


juju/juju/api/interface.go

// Connection represents a connection to a Juju API server.
type Connection interface {

        // This first block of methods is pretty close to a sane Connection interface.

        // Close closes the connection.
        Close() error

        // Addr returns the address used to connect to the API server.
        Addr() string

        // IPAddr returns the IP address used to connect to the API server.
        IPAddr() string

        // APIHostPorts returns addresses that may be used to connect
        // to the API server, including the address used to connect.
        //
        // The addresses are scoped (public, cloud-internal, etc.), so
        // the client may choose which addresses to attempt. For the
        // Juju CLI, all addresses must be attempted, as the CLI may
        // be invoked both within and outside the model (think
        // private clouds).
        APIHostPorts() [][]network.HostPort

        // Broken returns a channel which will be closed if the connection
        // is detected to be broken, either because the underlying
        // connection has closed or because API pings have failed.
        Broken() <-chan struct{}

        // IsBroken returns whether the connection is broken. It checks
        // the Broken channel and if that is open, attempts a connection
        // ping.
        IsBroken() bool

        // PublicDNSName returns the host name for which an officially
        // signed certificate will be used for TLS connection to the server.
        // If empty, the private Juju CA certificate must be used to verify
        // the connection.
        PublicDNSName() string

        // These are a bit off -- ServerVersion is apparently not known until after
        // Login()? Maybe evidence of need for a separate AuthenticatedConnection..?
        Login(name names.Tag, password, nonce string, ms []macaroon.Slice) error
        ServerVersion() (version.Number, bool)

        // APICaller provides the facility to make API calls directly.
        // This should not be used outside the api/* packages or tests.
        base.APICaller

        // ControllerTag returns the tag of the controller.
        // This could be defined on base.APICaller.
        ControllerTag() names.ControllerTag

        // All the rest are strange and questionable and deserve extra attention
        // and/or discussion.

        // Ping makes an API request which checks if the connection is
        // still functioning.
        // NOTE: This method is deprecated. Please use IsBroken or Broken instead.
        Ping() error

        // I think this is actually dead code. It's tested, at least, so I'm
        // keeping it for now, but it's not apparently used anywhere else.
        AllFacadeVersions() map[string][]int

        // AuthTag returns the tag of the authorized user of the state API
        // connection.
        AuthTag() names.Tag

        // ModelAccess returns the access level of authorized user to the model.
        ModelAccess() string

        // ControllerAccess returns the access level of authorized user to the controller.
        ControllerAccess() string

        // CookieURL returns the URL that HTTP cookies for the API will be
        // associated with.
        CookieURL() *url.URL

        // These methods expose a bunch of worker-specific facades, and basically
        // just should not exist; but removing them is too noisy for a single CL.
        // Client in particular is intimately coupled with State -- and the others
        // will be easy to remove, but until we're using them via manifolds it's
        // prohibitively ugly to do so.
        Client() *Client
        Uniter() (*uniter.State, error)
        Upgrader() *upgrader.State
        Reboot() (reboot.State, error)
        InstancePoller() *instancepoller.API
        CharmRevisionUpdater() *charmrevisionupdater.State
        Cleaner() *cleaner.API
        MetadataUpdater() *imagemetadata.Client
        UnitAssigner() unitassigner.API
}


---------------------------------------------------------------------------------


gopkg.in/juju/charmrepo.v2-unstable/csclient/csclient.go


// Client represents the client side of a charm store.
type Client struct {
        params        Params
        bclient       httpClient
        header        http.Header
        statsDisabled bool
        channel       params.Channel
}


// WithChannel returns a new client whose requests are done using the
// given channel.
func (c *Client) WithChannel(channel params.Channel) *Client {
        client := *c
        client.channel = channel
        return &client
}

---------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go

// The following structs exist purely because Go cannot create a
// struct with a field named the same as a method name. The DeployAPI
// needs to both embed a *<package>.Client and provide the
// api.Connection Client method.
//
// Once we pair down DeployAPI, this will not longer be a problem.

type apiClient struct {
        *api.Client
}

type charmsClient struct {
        *apicharms.Client
}

type applicationClient struct {
        *application.Client
}

type modelConfigClient struct {
        *modelconfig.Client
}

type charmRepoClient struct {
        *charmrepo.CharmStore
}

type charmstoreClient struct {
        *csclient.Client
}

type annotationsClient struct {
        *annotations.Client
}

---------------------------------------------------------------------------------
