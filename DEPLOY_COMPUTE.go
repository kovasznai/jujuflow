
sudo apt-get install -y software-properties-common
sudo add-apt-repository ppa:juju/stable
sudo apt-get update
sudo apt-get install -y juju

juju​ deploy​ ​--to=COMPUTE0_ID​ ​--config​ ​~/openstack-compute.yaml​ ​nova-compute
juju​ add-unit​ --to=COMPUTE1_ID​ nova-compute
juju​ add-unit​ --to=COMPUTE2_ID​ nova-compute


/usr/bin/juju

#!/bin/sh
export PATH=/usr/lib/juju-2.2/bin:"$PATH"
exec juju "$@"

ubuntu@vm3:~$ ls /usr/lib/juju-2.2/bin
juju  jujud  juju-metadata


vi /usr/lib/juju-2.2/bin/juju
<BINARY>


########################################################################################################

===================================================================================================
juju​ add-unit​ --to=COMPUTE1_ID​ nova-compute


juju/juju/cmd/juju/commands/main.go

package commands

import (
        "github.com/juju/juju/cmd/juju/application"


func registerCommands(r commandRegistry, ctx *cmd.Context) {

        r.Register(application.NewAddUnitCommand())


--------------------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/addunit.go

func (c *addUnitCommand) Info() *cmd.Info {
        return &cmd.Info{
                Name:    "add-unit",
                Args:    "<application name>",
                Purpose: usageAddUnitSummary,
                Doc:     usageAddUnitDetails,
        }
}


juju/cmd/cmd.go:233:type Info struct {

// Info holds some of the usage documentation of a Command.
type Info struct {
        // Name is the Command's name.
        Name string

        // Args describes the command's expected positional arguments.
        Args string

        // Purpose is a short explanation of the Command's purpose.
        Purpose string

        // Doc is the long documentation for the Command.
        Doc string

        // Aliases are other names for the Command.
        Aliases []string
}


--------------------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/addunit.go:133:func NewAddUnitCommand() cmd.Command {


package application

import (
        "github.com/juju/juju/cmd/modelcmd"



// NewAddUnitCommand returns a command that adds a unit[s] to an application.
func NewAddUnitCommand() cmd.Command {
        return modelcmd.Wrap(&addUnitCommand{})
}

--------------------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/addunit.go


// addUnitCommand is responsible adding additional units to an application.
type addUnitCommand struct {
        modelcmd.ModelCommandBase
        UnitCommandBase
        ApplicationName string
        api             serviceAddUnitAPI
}





########################################################################################################

=============================================================================================
juju​ deploy​ ​--to=COMPUTE0_ID​ ​--config​ ​~/openstack-compute.yaml​ ​nova-compute


juju/juju/cmd/juju/commands/main.go:389:	r.Register(application.NewDeployCommand())

package commands

import (
        "github.com/juju/juju/cmd/juju/application"

func registerCommands(r commandRegistry, ctx *cmd.Context) {

        r.Register(application.NewDeployCommand())


--------------------------------------------------------------------------------------------------

juju/juju/cmd/juju/application/deploy.go

import (
        "github.com/juju/cmd"


func (c *DeployCommand) Info() *cmd.Info {
        return &cmd.Info{
                Name:    "deploy",
                Args:    "<charm or bundle> [<application name>]",
                Purpose: "Deploy a new application or bundle.",
                Doc:     deployDoc,
        }
}


./cmd/cmd.go:233:type Info struct {

// Info holds some of the usage documentation of a Command.
type Info struct {
        // Name is the Command's name.
        Name string

        // Args describes the command's expected positional arguments.
        Args string

        // Purpose is a short explanation of the Command's purpose.
        Purpose string

        // Doc is the long documentation for the Command.
        Doc string

        // Aliases are other names for the Command.
        Aliases []string
}


--------------------------------------------------------------------------------------------------


juju/juju/cmd/juju/application/deploy.go:230:	deployCmd := &DeployCommand{

// NewDeployCommand returns a command to deploy services.
func NewDeployCommand() modelcmd.ModelCommand {
        steps := []DeployStep{
                &RegisterMeteredCharm{
                        RegisterURL: planURL + "/plan/authorize",
                        QueryURL:    planURL + "/charm",
                },
        }
        deployCmd := &DeployCommand{          <========
                Steps: steps,
        }
        deployCmd.NewAPIRoot = func() (DeployAPI, error) {
                apiRoot, err := deployCmd.ModelCommandBase.NewAPIRoot()
                if err != nil {
                        return nil, errors.Trace(err)
                }
                bakeryClient, err := deployCmd.BakeryClient()
                if err != nil {
                        return nil, errors.Trace(err)
                }
                cstoreClient := newCharmStoreClient(bakeryClient).WithChannel(deployCmd.Channel)

                return &deployAPIAdapter{
                        Connection:        apiRoot,
                        apiClient:         &apiClient{Client: apiRoot.Client()},
                        charmsClient:      &charmsClient{Client: apicharms.NewClient(apiRoot)},
                        applicationClient: &applicationClient{Client: application.NewClient(apiRoot)},
                        modelConfigClient: &modelConfigClient{Client: modelconfig.NewClient(apiRoot)},
                        charmstoreClient:  &charmstoreClient{Client: cstoreClient},
                        annotationsClient: &annotationsClient{Client: annotations.NewClient(apiRoot)},
                        charmRepoClient:   &charmRepoClient{CharmStore: charmrepo.NewCharmStoreFromClient(cstoreClient)},
                }, nil
        }

        return modelcmd.Wrap(deployCmd)
}


--------------------------------------------------------------------------------------------------


juju/juju/cmd/juju/application/deploy.go:515:		Name:    "deploy",

type DeployCommand struct {

        modelcmd.ModelCommandBase
        UnitCommandBase

        CharmOrBundle string

        BundleOverlayFile []string

        Channel params.Channel

        Series string

        Force bool

        DryRun bool

        ApplicationName string
        Config          cmd.FileVar
        ConstraintsStr  string
        Constraints     constraints.Value
        BindToSpaces    string

        Storage map[string]storage.Constraints

        BundleStorage map[string]map[string]storage.Constraints

        Resources map[string]string

        Bindings map[string]string
        Steps    []DeployStep

        UseExisting bool

        BundleMachines map[string]string

        NewAPIRoot func() (DeployAPI, error)        <======== deployCmd.NewAPIRoot = func() (DeployAPI, error) {

        NewAPIRoot func() (DeployAPI, error)

        machineMap string
        flagSet    *gnuflag.FlagSet
}


###################################################################################################


juju/juju/cmd/juju/commands/main.go:389:	r.Register(application.NewDeployCommand())

package commands

import (
        "github.com/juju/juju/cmd/juju/application"

func registerCommands(r commandRegistry, ctx *cmd.Context) {

        r.Register(application.NewDeployCommand())     <== modelcmd.Wrap(deployCmd) <== deployCmd := &DeployCommand{


--------------------------------------------------------------------------------------------------


juju/juju/cmd/juju/commands/main.go:262:type commandRegistry interface {

type commandRegistry interface {
        Register(cmd.Command)           <========
        RegisterSuperAlias(name, super, forName string, check cmd.DeprecationCheck)
        RegisterDeprecated(subcmd cmd.Command, check cmd.DeprecationCheck)
}

--------------------------------------------------------------------------------------------------

juju/cmd/supercommand.go

// Register makes a subcommand available for use on the command line. The
// command will be available via its own name, and via any supplied aliases.
func (c *SuperCommand) Register(subcmd Command) {
        info := subcmd.Info()
        c.insert(commandReference{name: info.Name, command: subcmd})
        for _, name := range info.Aliases {
                c.insert(commandReference{name: name, command: subcmd, alias: info.Name})
        }
}


./cmd/cmd.go:233:type Info struct {

// Info holds some of the usage documentation of a Command.
type Info struct {
        // Name is the Command's name.
        Name string

        // Args describes the command's expected positional arguments.
        Args string

        // Purpose is a short explanation of the Command's purpose.
        Purpose string

        // Doc is the long documentation for the Command.
        Doc string

        // Aliases are other names for the Command.
        Aliases []string
}


--------------------------------------------------------------------------------------------------

juju/cmd/supercommand.go

func (c *SuperCommand) insert(value commandReference) {
        if _, found := c.subcmds[value.name]; found {
                panic(fmt.Sprintf("command already registered: %q", value.name))
        }
        c.subcmds[value.name] = value
}


--------------------------------------------------------------------------------------------------


juju/cmd/supercommand.go:119:type SuperCommand struct {

// SuperCommand is a Command that selects a subcommand and assumes its
// properties; any command line arguments that were not used in selecting
// the subcommand are passed down to it, and to Run a SuperCommand is to run
// its selected subcommand.
type SuperCommand struct {
        CommandBase
        Name                string
        Purpose             string
        Doc                 string
        Log                 *Log
        Aliases             []string
        version             string
        usagePrefix         string
        userAliasesFilename string
        userAliases         map[string][]string
        subcmds             map[string]commandReference
        help                *helpCommand
        commonflags         *gnuflag.FlagSet
        flags               *gnuflag.FlagSet
        action              commandReference
        showHelp            bool
        showDescription     bool
        showVersion         bool
        noAlias             bool
        missingCallback     MissingCallback
        notifyRun           func(string)
        notifyHelp          func([]string)
}



###################################################################################################

juju/juju/cmd/juju/application/addunit.go:76:type UnitCommandBase struct {


type UnitCommandBase struct {

        PlacementSpec string

        Placement []*instance.Placement
        NumUnits  int

        AttachStorage []string
}

-----------------------------------------------------------------------------------

juju/juju/cmd/modelcmd/modelcommand.go:316:func Wrap(c ModelCommand, options ...WrapOption) ModelCommand {


func Wrap(c ModelCommand, options ...WrapOption) ModelCommand {
        wrapper := &modelCommandWrapper{
                ModelCommand:    c,
                skipModelFlags:  false,
                useDefaultModel: true,
        }
        for _, option := range options {
                option(wrapper)
        }
        // Define a new type so that we can embed the ModelCommand
        // interface one level deeper than cmd.Command, so that
        // we'll get the Command methods from WrapBase
        // and all the ModelCommand methods not in cmd.Command
        // from modelCommandWrapper.
        type embed struct {
                *modelCommandWrapper
        }
        return struct {
                embed
                cmd.Command
        }{
                Command: WrapBase(wrapper),
                embed:   embed{wrapper},
        }
}



type modelCommandWrapper struct {
        ModelCommand

        skipModelFlags  bool
        useDefaultModel bool
        modelName       string
}





juju/juju/cmd/modelcmd/modelcommand.go:38:type ModelCommand interface {

type ModelCommand interface {

        Command

        SetClientStore(jujuclient.ClientStore)

        ClientStore() jujuclient.ClientStore

        SetModelName(modelName string, allowDefault bool) error

        ModelName() (string, error)

        ControllerName() (string, error)

        initModel() error
}





juju/juju/cmd/juju/charmcmd/charm.go:20:type Command struct {

package charmcmd

import (
        "github.com/juju/cmd"


type Command struct {
        cmd.SuperCommand
}



juju/cmd/supercommand.go:119:type SuperCommand struct {

type SuperCommand struct {
        CommandBase
        Name                string
        Purpose             string
        Doc                 string
        Log                 *Log
        Aliases             []string
        version             string
        usagePrefix         string
        userAliasesFilename string
        userAliases         map[string][]string
        subcmds             map[string]commandReference
        help                *helpCommand
        commonflags         *gnuflag.FlagSet
        flags               *gnuflag.FlagSet
        action              commandReference
        showHelp            bool
        showDescription     bool
        showVersion         bool
        noAlias             bool
        missingCallback     MissingCallback
        notifyRun           func(string)
        notifyHelp          func([]string)
}



-----------------------------------------------------------------------------------

juju/juju/cmd/modelcmd/modelcommand.go:79:type ModelCommandBase struct {


type ModelCommandBase struct {
        CommandBase

        store jujuclient.ClientStore

        _modelName      string
        _controllerName string

        allowDefaultModel bool

        doneInitModel bool

        initModelError error
}



juju/juju/cmd/modelcmd/base.go:59:type CommandBase struct {

package modelcmd
import (
        "github.com/juju/cmd"

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


juju/cmd/cmd.go:86:type CommandBase struct{}

type CommandBase struct{}














#########################################################################################################

juju​ deploy​ ​--to=COMPUTE0_ID​ ​--config​ ​~/openstack-compute.yaml​ ​nova-compute


juju/juju/cmd/juju/commands/main.go:389:	r.Register(application.NewDeployCommand())

package commands

import (

        "github.com/juju/cmd"

        "github.com/juju/juju/cmd/juju/application"



type commandRegistry interface {
        Register(cmd.Command)
        RegisterSuperAlias(name, super, forName string, check cmd.DeprecationCheck)
        RegisterDeprecated(subcmd cmd.Command, check cmd.DeprecationCheck)
}




github.com/juju/cmd/cmd.go:63:type Command interface {


type Command interface {

        IsSuperCommand() bool

        Info() *Info

        SetFlags(f *gnuflag.FlagSet)

        Init(args []string) error

        Run(ctx *Context) error

        AllowInterspersedFlags() bool
}



github.com/juju/cmd/supercommand.go:185:func (c *SuperCommand) Register(subcmd Command) {

func (c *SuperCommand) Register(subcmd Command) {
        info := subcmd.Info()
        c.insert(commandReference{name: info.Name, command: subcmd})
        for _, name := range info.Aliases {
                c.insert(commandReference{name: name, command: subcmd, alias: info.Name})
        }
}



github.com/juju/cmd/supercommand.go:119:type SuperCommand struct {

type SuperCommand struct {
        CommandBase
        Name                string
        Purpose             string
        Doc                 string
        Log                 *Log
        Aliases             []string
        version             string
        usagePrefix         string
        userAliasesFilename string
        userAliases         map[string][]string
        subcmds             map[string]commandReference
        help                *helpCommand
        commonflags         *gnuflag.FlagSet
        flags               *gnuflag.FlagSet
        action              commandReference
        showHelp            bool
        showDescription     bool
        showVersion         bool
        noAlias             bool
        missingCallback     MissingCallback
        notifyRun           func(string)
        notifyHelp          func([]string)
}


github.com/juju/cmd/supercommand.go

func (c *SuperCommand) insert(value commandReference) {
        if _, found := c.subcmds[value.name]; found {
                panic(fmt.Sprintf("command already registered: %q", value.name))
        }
        c.subcmds[value.name] = value
}






github.com/juju/cmd/supercommand.go:429:func (c *SuperCommand) Run(ctx *Context) error {

// Run executes the subcommand that was selected in Init.
func (c *SuperCommand) Run(ctx *Context) error {
        if c.showDescription {
                if c.Purpose != "" {
                        fmt.Fprintf(ctx.Stdout, "%s\n", c.Purpose)
                } else {
                        fmt.Fprintf(ctx.Stdout, "%s: no description available\n", c.Info().Name)
                }
                return nil
        }
        if c.action.command == nil {
                panic("Run: missing subcommand; Init failed or not called")
        }
        if c.Log != nil {
                if err := c.Log.Start(ctx); err != nil {
                        return err
                }
        }
        if c.notifyRun != nil {
                name := c.Name
                if c.usagePrefix != "" && c.usagePrefix != name {
                        name = c.usagePrefix + " " + name
                }
                c.notifyRun(name)
        }
        if deprecated, replacement := c.action.Deprecated(); deprecated {
                ctx.Infof("WARNING: %q is deprecated, please use %q", c.action.name, replacement)
        }
        err := c.action.command.Run(ctx)
        if err != nil && !IsErrSilent(err) {
                WriteError(ctx.Stderr, err)
                logger.Debugf("error stack: \n%v", errors.ErrorStack(err))
                // Now that this has been logged, don't log again in cmd.Main.
                if !IsRcPassthroughError(err) {
                        err = ErrSilent
                }
        } else {
                logger.Infof("command finished")
        }
        return err




juju/juju/cmd/juju/application/deploy.go:813:func (c *DeployCommand) Run(ctx *cmd.Context) error {


func (c *DeployCommand) Run(ctx *cmd.Context) error {
        var err error
        c.Constraints, err = common.ParseConstraints(ctx, c.ConstraintsStr)
        if err != nil {
                return err
        }
        apiRoot, err := c.NewAPIRoot()
        if err != nil {
                return errors.Trace(err)
        }
        defer apiRoot.Close()

        deploy, err := findDeployerFIFO(
                c.maybeReadLocalBundle,
                func() (deployFn, error) { return c.maybeReadLocalCharm(apiRoot) },
                c.maybePredeployedLocalCharm,
                c.maybeReadCharmstoreBundleFn(apiRoot),
                c.charmStoreCharm, // This always returns a deployer
        )
        if err != nil {
                return errors.Trace(err)
        }

        return block.ProcessBlockedError(deploy(ctx, apiRoot), block.BlockChange)
}



#########################################################################################################

---------------------------------------------------------------------------------------------------------

func registerCommands(r commandRegistry, ctx *cmd.Context) {

        r.Register(application.NewDeployCommand())


---------------------------------------------------------------------------------------------------------


juju​ deploy​ ​--to=COMPUTE0_ID​ ​--config​ ​~/openstack-compute.yaml​ ​nova-compute


./cmd/juju/application/deploy.go:230:	deployCmd := &DeployCommand{


import (

        "github.com/juju/juju/cmd/modelcmd"


var planURL = "https://api.jujucharms.com/omnibus/v2"

// NewDeployCommand returns a command to deploy services.
func NewDeployCommand() modelcmd.ModelCommand {
        steps := []DeployStep{
                &RegisterMeteredCharm{
                        RegisterURL: planURL + "/plan/authorize",
                        QueryURL:    planURL + "/charm",
                },
        }
        deployCmd := &DeployCommand{
                Steps: steps,
        }
        deployCmd.NewAPIRoot = func() (DeployAPI, error) {
                apiRoot, err := deployCmd.ModelCommandBase.NewAPIRoot()
                if err != nil {
                        return nil, errors.Trace(err)
                }
                bakeryClient, err := deployCmd.BakeryClient()
                if err != nil {
                        return nil, errors.Trace(err)
                }
                cstoreClient := newCharmStoreClient(bakeryClient).WithChannel(deployCmd.Channel)

                return &deployAPIAdapter{
                        Connection:        apiRoot,
                        apiClient:         &apiClient{Client: apiRoot.Client()},
                        charmsClient:      &charmsClient{Client: apicharms.NewClient(apiRoot)},
                        applicationClient: &applicationClient{Client: application.NewClient(apiRoot)},
                        modelConfigClient: &modelConfigClient{Client: modelconfig.NewClient(apiRoot)},
                        charmstoreClient:  &charmstoreClient{Client: cstoreClient},
                        annotationsClient: &annotationsClient{Client: annotations.NewClient(apiRoot)},
                        charmRepoClient:   &charmRepoClient{CharmStore: charmrepo.NewCharmStoreFromClient(cstoreClient)},
                }, nil
        }

        return modelcmd.Wrap(deployCmd)
}


-------------------------------------------------------------------------------------------------------


juju/juju/cmd/modelcmd/modelcommand.go:38:type ModelCommand interface {

type ModelCommand interface {

        Command

        SetClientStore(jujuclient.ClientStore)

        ClientStore() jujuclient.ClientStore

        SetModelName(modelName string, allowDefault bool) error

        ModelName() (string, error)

        ControllerName() (string, error)

        initModel() error

}

-------------------------------------------------------------------------------------------------------

./cmd/juju/application/deploy.go:138:type deployAPIAdapter struct {

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


-------------------------------------------------------------------------------------------------------


./cmd/juju/application/deploy.go:515:		Name:    "deploy",

type DeployCommand struct {

        modelcmd.ModelCommandBase

        UnitCommandBase

        CharmOrBundle string

        BundleOverlayFile []string

        Channel params.Channel

        Series string

        Force bool

        DryRun bool

        ApplicationName string
        Config          cmd.FileVar
        ConstraintsStr  string
        Constraints     constraints.Value
        BindToSpaces    string

        Storage map[string]storage.Constraints

        BundleStorage map[string]map[string]storage.Constraints

        Resources map[string]string

        Bindings map[string]string
        Steps    []DeployStep

        UseExisting bool

        BundleMachines map[string]string

        NewAPIRoot func() (DeployAPI, error)

        NewAPIRoot func() (DeployAPI, error)

        machineMap string
        flagSet    *gnuflag.FlagSet
}


-------------------------------------------------------------------------------------------------------


./cmd/juju/application/deploy.go:491:type DeployStep interface {

type DeployStep interface {

        SetFlags(*gnuflag.FlagSet)

        RunPre(MeteredDeployAPI, *httpbakery.Client, *cmd.Context, DeploymentInfo) error

        RunPost(MeteredDeployAPI, *httpbakery.Client, *cmd.Context, DeploymentInfo, error) error

}


-------------------------------------------------------------------------------------------------------


./cmd/juju/application/register.go:29:type RegisterMeteredCharm struct {

type RegisterMeteredCharm struct {
        Plan           string
        IncreaseBudget int
        RegisterURL    string
        QueryURL       string
        credentials    []byte
}


-------------------------------------------------------------------------------------------------------


juju/juju/cmd/modelcmd/modelcommand.go:79:type ModelCommandBase struct {


type ModelCommandBase struct {
        CommandBase

        store jujuclient.ClientStore

        _modelName      string
        _controllerName string

        allowDefaultModel bool

        doneInitModel bool

        initModelError error
}


-------------------------------------------------------------------------------------------------------


juju/juju/cmd/modelcmd/base.go:59:type CommandBase struct {

package modelcmd
import (
        "github.com/juju/cmd"

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


-------------------------------------------------------------------------------------------------------


juju/cmd/cmd.go:86:type CommandBase struct{}

type CommandBase struct{}


-------------------------------------------------------------------------------------------------------

juju/juju/cmd/modelcmd/modelcommand.go:316:func Wrap(c ModelCommand, options ...WrapOption) ModelCommand {


// Wrap wraps the specified ModelCommand, returning a ModelCommand
// that proxies to each of the ModelCommand methods.
// Any provided options are applied to the wrapped command
// before it is returned.
func Wrap(c ModelCommand, options ...WrapOption) ModelCommand {
        wrapper := &modelCommandWrapper{
                ModelCommand:    c,
                skipModelFlags:  false,
                useDefaultModel: true,
        }
        for _, option := range options {
                option(wrapper)
        }
        // Define a new type so that we can embed the ModelCommand
        // interface one level deeper than cmd.Command, so that
        // we'll get the Command methods from WrapBase
        // and all the ModelCommand methods not in cmd.Command
        // from modelCommandWrapper.
        type embed struct {
                *modelCommandWrapper
        }
        return struct {
                embed
                cmd.Command
        }{
                Command: WrapBase(wrapper),
                embed:   embed{wrapper},
        }
}


====================================================================================================

juju/juju/cmd/juju/application/deploy.go:138:type deployAPIAdapter struct {

import (

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



juju/juju/api/interface.go:215:type Connection interface {

type Connection interface {

        Close() error

        Addr() string

        IPAddr() string

        APIHostPorts() [][]network.HostPort

        Broken() <-chan struct{}

        IsBroken() bool

        PublicDNSName() string

        Login(name names.Tag, password, nonce string, ms []macaroon.Slice) error
        ServerVersion() (version.Number, bool)

        base.APICaller

        ControllerTag() names.ControllerTag

        Ping() error

        AllFacadeVersions() map[string][]int

        AuthTag() names.Tag

        ModelAccess() string

        ControllerAccess() string

        CookieURL() *url.URL

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




juju/juju/cmd/juju/application/deploy.go:138:type deployAPIAdapter struct 

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



./cmd/juju/application/deploy.go:106:type apiClient struct {

import (
        "github.com/juju/juju/api"
        apicharms "github.com/juju/juju/api/charms"
        "github.com/juju/juju/api/application"
        "github.com/juju/juju/api/modelconfig"
        "gopkg.in/juju/charmrepo.v2"
        "gopkg.in/juju/charmrepo.v2/csclient"
        "github.com/juju/juju/api/annotations"

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


-----------------------------------------------------------------------------------------


juju/juju/api/client.go:35:type Client struct {

type Client struct {
        base.ClientFacade
        facade base.FacadeCaller
        st     *state
}


-----------------------------------------------------------------------------------------


juju/juju/api/charms/client.go:18:type Client struct {

type Client struct {
        base.ClientFacade
        facade base.FacadeCaller
}


-----------------------------------------------------------------------------------------


juju/juju/api/application/client.go:28:type Client struct {

type Client struct {
        base.ClientFacade
        st     base.APICallCloser
        facade base.FacadeCaller
}


-----------------------------------------------------------------------------------------


juju/juju/api/modelconfig/modelconfig.go:16:type Client struct {

type Client struct {
        base.ClientFacade
        facade base.FacadeCaller
}


-----------------------------------------------------------------------------------------


juju/juju/api/annotations/client.go:14:type Client struct {

type Client struct {
        base.ClientFacade
        facade base.FacadeCaller
}


-----------------------------------------------------------------------------------------


gopkg.in/juju/charmrepo.v2-unstable/charmstore.go:31:type CharmStore struct {

import (
        "gopkg.in/juju/charmrepo.v2-unstable/csclient"

type CharmStore struct {
        client *csclient.Client
}



gopkg.in/juju/charmrepo.v2-unstable/csclient/csclient.go:51:type Client struct {

type Client struct {
        params        Params
        bclient       httpClient
        header        http.Header
        statsDisabled bool
        channel       params.Channel
}


-----------------------------------------------------------------------------------------


gopkg.in/juju/charmrepo.v2-unstable/csclient/csclient.go:51:type Client struct {

type Client struct {
        params        Params
        bclient       httpClient
        header        http.Header
        statsDisabled bool
        channel       params.Channel
}


-----------------------------------------------------------------------------------------





