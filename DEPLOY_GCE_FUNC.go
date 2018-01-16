
juju add-credential google

juju bootstrap google/us-east1 

juju add-model k8s


juju deploy cs:~containers/kubernetes-master-17
juju deploy cs:~containers/etcd-29 --to 0
juju deploy cs:~containers/easyrsa-8 --to lxd:0
juju deploy cs:~containers/flannel-13
juju deploy cs:~containers/kubernetes-worker-22
juju expose kubernetes-master
juju expose kubernetes-worker

=============================================================================================================


juju/juju/cmd/juju/commands/main.go:389:	r.Register(application.NewDeployCommand())

        // Manage and control services

        r.Register(application.NewDeployCommand())


---------------------------------------------------------------------------------
juju/juju/cmd/juju/application/deploy.go

        "github.com/juju/errors"

        "gopkg.in/juju/charmrepo.v2/csclient"

        "github.com/juju/juju/api/application"
        apicharms "github.com/juju/juju/api/charms"
        "github.com/juju/juju/api/modelconfig"
        "github.com/juju/juju/api/annotations"
        "gopkg.in/juju/charmrepo.v2"




// NewDeployCommand returns a command to deploy applications.
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
                apiRoot, err := deployCmd.ModelCommandBase.NewAPIRoot()   =>  type Connection interface
                if err != nil {
                        return nil, errors.Trace(err)
                }
                bakeryClient, err := deployCmd.BakeryClient()
                if err != nil {
                        return nil, errors.Trace(err)
                }
                cstoreClient := newCharmStoreClient(bakeryClient).WithChannel(deployCmd.Channel)   =>  *csclient.Client  WithChannel(

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


---------------------------------------------------------------------------------


juju/juju/cmd/modelcmd/modelcommand.go

// NewAPIRoot returns a new connection to the API server for the environment
// directed to the model specified on the command line.
func (c *ModelCommandBase) NewAPIRoot() (api.Connection, error) {
        modelName, _, err := c.ModelDetails()
        if err != nil {
                return nil, errors.Trace(err)
        }
        return c.newAPIRoot(modelName)
}


---------------------------------------------------------------------------------


juju/juju/cmd/modelcmd/base.go

// BakeryClient returns a macaroon bakery client that
// uses the same HTTP client returned by HTTPClient.
func (c *CommandBase) BakeryClient(store jujuclient.CookieStore, controllerName string) (*httpbakery.Client, error) {
        c.assertRunStarted()
        ctx, err := c.getAPIContext(store, controllerName)
        if err != nil {
                return nil, errors.Trace(err)
        }
        return ctx.NewBakeryClient(), nil
}


---------------------------------------------------------------------------------

juju/errors/functions.go

// Trace adds the location of the Trace call to the stack.  The Cause of the
// resulting error is the same as the error parameter.  If the other error is
// nil, the result will be nil.
//
// For example:
//   if err := SomeFunc(); err != nil {
//       return errors.Trace(err)
//   }
//
func Trace(other error) error {
        if other == nil {
                return nil
        }
        err := &Err{previous: other, cause: Cause(other)}
        err.SetLocation(1)
        return err
}


---------------------------------------------------------------------------------

juju/juju/cmd/juju/application/store.go

        "gopkg.in/juju/charmrepo.v2/csclient"


// newCharmStoreClient is called to obtain a charm store client.
// It is defined as a variable so it can be changed for testing purposes.
var newCharmStoreClient = func(client *httpbakery.Client) *csclient.Client {
        return csclient.New(csclient.Params{
                BakeryClient: client,
        })
}


---------------------------------------------------------------------------------

juju/juju/cmd/modelcmd/modelcommand.go

        "github.com/juju/juju/api"


// NewAPIRoot returns a new connection to the API server for the environment
// directed to the model specified on the command line.
func (c *ModelCommandBase) NewAPIRoot() (api.Connection, error) {
        modelName, _, err := c.ModelDetails()
        if err != nil {
                return nil, errors.Trace(err)
        }
        return c.newAPIRoot(modelName)
}


---------------------------------------------------------------------------------


juju/juju/api/charms/client.go


// Client allows access to the charms API end point.
type Client struct {
        base.ClientFacade
        facade base.FacadeCaller
}


// NewClient creates a new client for accessing the charms API.
func NewClient(st base.APICallCloser) *Client {
        frontend, backend := base.NewClientFacade(st, "Charms")
        return &Client{ClientFacade: frontend, facade: backend}
}


---------------------------------------------------------------------------------

juju/juju/api/application/client.go


// Client allows access to the service API end point.
type Client struct {
        base.ClientFacade
        st     base.APICallCloser
        facade base.FacadeCaller
}


// NewClient creates a new client for accessing the application api.
func NewClient(st base.APICallCloser) *Client {
        frontend, backend := base.NewClientFacade(st, "Application")
        return &Client{ClientFacade: frontend, st: st, facade: backend}
}


---------------------------------------------------------------------------------

juju/juju/api/modelconfig/modelconfig.go


// Client provides methods that the Juju client command uses to interact
// with models stored in the Juju Server.
type Client struct {
        base.ClientFacade
        facade base.FacadeCaller
}


// NewClient creates a new `Client` based on an existing authenticated API
// connection.
func NewClient(st base.APICallCloser) *Client {
        frontend, backend := base.NewClientFacade(st, "ModelConfig")
        return &Client{ClientFacade: frontend, facade: backend}
}


---------------------------------------------------------------------------------

juju/juju/api/annotations/client.go


// Client allows access to the annotations API end point.
type Client struct {
        base.ClientFacade
        facade base.FacadeCaller
}


// NewClient creates a new client for accessing the annotations API.
func NewClient(st base.APICallCloser) *Client {
        frontend, backend := base.NewClientFacade(st, "Annotations")
        return &Client{ClientFacade: frontend, facade: backend}
}


---------------------------------------------------------------------------------

gopkg.in/juju/charmrepo.v2-unstable/charmstore.go

// CharmStore is a repository Interface that provides access to the public Juju
// charm store.
type CharmStore struct {
        client *csclient.Client
}


// NewCharmStoreFromClient creates and returns a charm store repository.
// The provided client is used for charm store requests.
func NewCharmStoreFromClient(client *csclient.Client) *CharmStore {
        return &CharmStore{
                client: client,
        }
}


---------------------------------------------------------------------------------


juju/juju/cmd/juju/application/deploy.go

func (c *DeployCommand) deployCharm(
        id charmstore.CharmID,
        csMac *macaroon.Macaroon,
        series string,
        ctx *cmd.Context,
        apiRoot DeployAPI,
) (rErr error) {
        charmInfo, err := apiRoot.CharmInfo(id.URL.String())
        if err != nil {
                return err
        }

        if len(c.AttachStorage) > 0 && apiRoot.BestFacadeVersion("Application") < 5 {
                // DeployArgs.AttachStorage is only supported from
                // Application API version 5 and onwards.
                return errors.New("this juju controller does not support --attach-storage")
        }

        numUnits := c.NumUnits
        if charmInfo.Meta.Subordinate {
                if !constraints.IsEmpty(&c.Constraints) {
                        return errors.New("cannot use --constraints with subordinate application")
                }
                if numUnits == 1 && c.PlacementSpec == "" {
                        numUnits = 0
                } else {
                        return errors.New("cannot use --num-units or --to with subordinate application")
                }
        }
        applicationName := c.ApplicationName
        if applicationName == "" {
                applicationName = charmInfo.Meta.Name
        }

        // Process the --config args.
        // We may have a single file arg specified, in which case
        // it points to a YAML file keyed on the charm name and
        // containing values for any charm settings.
        // We may also have key/value pairs representing
        // charm settings which overrides anything in the YAML file.
        // If more than one file is specified, that is an error.
        var configYAML []byte
        files, err := c.ConfigOptions.AbsoluteFileNames(ctx)
        if err != nil {
                return errors.Trace(err)
        }
        if len(files) > 1 {
                return errors.Errorf("only a single config YAML file can be specified, got %d", len(files))
        }
        if len(files) == 1 {
                configYAML, err = ioutil.ReadFile(files[0])
                if err != nil {
                        return errors.Trace(err)
                }
        }

        attr, err := c.ConfigOptions.ReadConfigPairs(ctx)
        if err != nil {
                return errors.Trace(err)
        }
        appConfig := make(map[string]string)
        for k, v := range attr {
                appConfig[k] = v.(string)
        }

        // Application facade V5 expects charm config to either all be in YAML
        // or config map. If config map is specified, that overrides YAML.
        // So we need to combine the two here to have only one.
        if apiRoot.BestFacadeVersion("Application") < 6 && len(appConfig) > 0 {
                var configFromFile map[string]map[string]string
                err := yaml.Unmarshal(configYAML, &configFromFile)
                if err != nil {
                        return errors.Annotate(err, "badly formatted YAML config file")
                }
                if configFromFile == nil {
                        configFromFile = make(map[string]map[string]string)
                }
                charmSettings, ok := configFromFile[applicationName]
                if !ok {
                        charmSettings = make(map[string]string)
                }
                for k, v := range appConfig {
                        charmSettings[k] = v
                }
                appConfig = nil
                configFromFile[applicationName] = charmSettings
                configYAML, err = yaml.Marshal(configFromFile)
                if err != nil {
                        return errors.Trace(err)
                }
        }

        bakeryClient, err := c.BakeryClient()
        if err != nil {
                return errors.Trace(err)
        }

        uuid, ok := apiRoot.ModelUUID()
        if !ok {
                return errors.New("API connection is controller-only (should never happen)")
        }

        deployInfo := DeploymentInfo{
                CharmID:         id,
                ApplicationName: applicationName,
                ModelUUID:       uuid,
                CharmInfo:       charmInfo,
        }

        for _, step := range c.Steps {
                err = step.RunPre(apiRoot, bakeryClient, ctx, deployInfo)
                if err != nil {
                        return errors.Trace(err)
                }
        }

        defer func() {
                for _, step := range c.Steps {
                        err = errors.Trace(step.RunPost(apiRoot, bakeryClient, ctx, deployInfo, rErr))
                        if err != nil {
                                rErr = err
                        }
                }
        }()

        if id.URL != nil && id.URL.Schema != "local" && len(charmInfo.Meta.Terms) > 0 {
                ctx.Infof("Deployment under prior agreement to terms: %s",
                        strings.Join(charmInfo.Meta.Terms, " "))
        }

        ids, err := resourceadapters.DeployResources(
                applicationName,
                id,
                csMac,
                c.Resources,
                charmInfo.Meta.Resources,
                apiRoot,
        )
        if err != nil {
                return errors.Trace(err)
        }

        if len(appConfig) == 0 {
                appConfig = nil
        }
        args := application.DeployArgs{
                CharmID:          id,
                Cons:             c.Constraints,
                ApplicationName:  applicationName,
                Series:           series,
                NumUnits:         numUnits,
                ConfigYAML:       string(configYAML),
                Config:           appConfig,
                Placement:        c.Placement,
                Storage:          c.Storage,
                AttachStorage:    c.AttachStorage,
                Resources:        ids,
                EndpointBindings: c.Bindings,
        }
        if len(appConfig) > 0 {
                args.Config = appConfig
        }
        return errors.Trace(apiRoot.Deploy(args))
}
