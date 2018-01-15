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


---------------------------------------------------------------------------------


./cmd/juju/application/deploy.go:261:type DeployCommand struct {


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


./cmd/juju/application/deploy.go:706:func (c *DeployCommand) deployCharm(

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

