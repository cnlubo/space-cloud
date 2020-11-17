package syncman

import (
	"context"
	"errors"

	"github.com/stretchr/testify/mock"

	"github.com/spaceuptech/space-cloud/gateway/config"
	"github.com/spaceuptech/space-cloud/gateway/model"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/letsencrypt"
	"github.com/spaceuptech/space-cloud/gateway/modules/global/routing"
)

type mockIntegrationManager struct {
	mock.Mock
	skip bool
}

func (m *mockIntegrationManager) SetConfig(integrations config.Integrations, integrationHooks config.IntegrationHooks) error {
	return m.Called(integrations, integrationHooks).Error(0)
}

func (m *mockIntegrationManager) SetIntegrations(integrations config.Integrations) error {
	return m.Called(integrations).Error(0)
}

func (m *mockIntegrationManager) SetIntegrationHooks(integrationHooks config.IntegrationHooks) {
	m.Called(integrationHooks)
}

func (m *mockIntegrationManager) InvokeHook(_ context.Context, params model.RequestParams) config.IntegrationAuthResponse {
	if m.skip {
		return mockHookResponse{checkResponse: false}
	}
	return m.Called(params).Get(0).(config.IntegrationAuthResponse)
}

type mockHookResponse struct {
	checkResponse bool
	err           string
	result        interface{}
	status        int
}

// CheckResponse indicates whether the integration is hijacking the authentication of the request or not.
// Its a humble way of saying that I'm the boss for this request
func (r mockHookResponse) CheckResponse() bool {
	return r.checkResponse
}

// Error returns error generated by the module if CheckResponse() returns true.
func (r mockHookResponse) Error() error {
	if r.err == "" {
		return nil
	}
	return errors.New(r.err)
}

// Result returns the value received from the integration
func (r mockHookResponse) Result() interface{} {
	return r.result
}

// Status returns the status code received from the integration
func (r mockHookResponse) Status() int {
	return r.status
}

type mockAdminSyncmanInterface struct {
	mock.Mock
}

func (m *mockAdminSyncmanInterface) SetServices(eventType string, services model.ScServices) {
	m.Called(eventType, services)
}

func (m *mockAdminSyncmanInterface) SetSessionID(sessionID string) {
	m.Called(sessionID)
}

func (m *mockAdminSyncmanInterface) SetIntegrationConfig(integrations config.Integrations) {
	m.Called(integrations)
}

func (m *mockAdminSyncmanInterface) ValidateIntegrationSyncOperation(integrations config.Integrations) error {
	return m.Called(integrations).Error(0)
}

func (m *mockAdminSyncmanInterface) ParseLicense(license string) (map[string]interface{}, error) {
	v := m.Called(license)
	return v.Get(0).(map[string]interface{}), v.Error(1)
}

func (m *mockAdminSyncmanInterface) GetIntegrationToken(id string) (string, error) {
	v := m.Called(id)
	return v.String(0), v.Error(1)
}

func (m *mockAdminSyncmanInterface) IsRegistered() bool {
	return m.Called().Bool(0)
}

func (m *mockAdminSyncmanInterface) GetSessionID() (string, error) {
	return m.Called().String(0), nil
}

func (m *mockAdminSyncmanInterface) RenewLicense(b bool) error {
	return m.Called(b).Error(0)
}

func (m *mockAdminSyncmanInterface) ValidateProjectSyncOperation(projects *config.Config, projectID *config.ProjectConfig) bool {
	return m.Called(projects, projectID).Bool(0)
}

func (m *mockAdminSyncmanInterface) SetConfig(admin *config.License) error {
	return m.Called(admin).Error(0)
}

func (m *mockAdminSyncmanInterface) IsTokenValid(_ context.Context, token, resource, op string, attr map[string]string) (model.RequestParams, error) {
	c := m.Called(token, resource, op, attr)
	return c.Get(0).(model.RequestParams), c.Error(1)
}

func (m *mockAdminSyncmanInterface) GetInternalAccessToken() (string, error) {
	c := m.Called()
	return c.String(0), c.Error(1)
}

func (m *mockAdminSyncmanInterface) GetConfig() *config.License {
	return m.Called().Get(0).(*config.License)
}

type mockModulesInterface struct {
	mock.Mock
}

func (m *mockModulesInterface) SetInitialProjectConfig(ctx context.Context, config config.Projects) error {
	a := m.Called(ctx, config)
	return a.Error(0)
}

func (m *mockModulesInterface) SetDatabaseConfig(ctx context.Context, projectID string, crudConfigs config.DatabaseConfigs) error {
	return m.Called(ctx, projectID, crudConfigs).Error(0)
}

func (m *mockModulesInterface) SetDatabaseSchemaConfig(ctx context.Context, projectID string, schemaConfigs config.DatabaseSchemas) error {
	return m.Called(ctx, projectID, schemaConfigs).Error(0)
}

func (m *mockModulesInterface) SetDatabaseRulesConfig(ctx context.Context, projectID string, ruleConfigs config.DatabaseRules) error {
	return m.Called(ctx, projectID, ruleConfigs).Error(0)
}

func (m *mockModulesInterface) SetDatabasePreparedQueryConfig(ctx context.Context, projectID string, prepConfigs config.DatabasePreparedQueries) error {
	return m.Called(ctx, projectID, prepConfigs).Error(0)
}

func (m *mockModulesInterface) SetFileStoreConfig(ctx context.Context, projectID string, fileStore *config.FileStoreConfig) error {
	c := m.Called(ctx, projectID, fileStore)
	return c.Error(0)
}

func (m *mockModulesInterface) SetFileStoreSecurityRuleConfig(ctx context.Context, projectID string, fileRule config.FileStoreRules) error {
	return m.Called(ctx, projectID, fileRule).Error(0)
}

func (m *mockModulesInterface) SetRemoteServiceConfig(ctx context.Context, projectID string, services config.Services) error {
	return m.Called(ctx, projectID, services).Error(0)
}

func (m *mockModulesInterface) SetLetsencryptConfig(ctx context.Context, projectID string, c *config.LetsEncrypt) error {
	return m.Called(ctx, projectID, c).Error(0)
}

func (m *mockModulesInterface) SetIngressRouteConfig(ctx context.Context, projectID string, routes config.IngressRoutes) error {
	return m.Called(ctx, projectID, routes).Error(0)
}

func (m *mockModulesInterface) SetIngressGlobalRouteConfig(ctx context.Context, projectID string, c *config.GlobalRoutesConfig) error {
	return m.Called(ctx, projectID, c).Error(0)
}

func (m *mockModulesInterface) SetEventingConfig(ctx context.Context, projectID string, eventingConfigs *config.EventingConfig) error {
	c := m.Called(ctx, projectID, eventingConfigs)
	return c.Error(0)
}

func (m *mockModulesInterface) SetEventingSchemaConfig(ctx context.Context, projectID string, schemaObj config.EventingSchemas) error {
	return m.Called(ctx, projectID, schemaObj).Error(0)
}

func (m *mockModulesInterface) SetEventingTriggerConfig(ctx context.Context, projectID string, triggerObj config.EventingTriggers) error {
	return m.Called(ctx, projectID, triggerObj).Error(0)
}

func (m *mockModulesInterface) SetEventingRuleConfig(ctx context.Context, projectID string, secureObj config.EventingRules) error {
	return m.Called(ctx, projectID, secureObj).Error(0)
}

func (m *mockModulesInterface) SetUsermanConfig(ctx context.Context, projectID string, auth config.Auths) error {
	return m.Called(ctx, projectID, auth).Error(0)
}

func (m *mockModulesInterface) LetsEncrypt() *letsencrypt.LetsEncrypt {
	return m.Called().Get(0).(*letsencrypt.LetsEncrypt)
}

func (m *mockModulesInterface) Routing() *routing.Routing {
	return m.Called().Get(0).(*routing.Routing)
}

func (m *mockModulesInterface) Delete(projectID string) {
	m.Called(projectID)
}

func (m *mockModulesInterface) SetProjectConfig(ctx context.Context, config *config.ProjectConfig) error {
	return m.Called(ctx, config).Error(0)
}

func (m *mockModulesInterface) SetGlobalConfig(projectID, secretSource string, secrets []*config.Secret, aesKey string) error {
	c := m.Called(projectID, secretSource, secrets, aesKey)
	return c.Error(0)
}

func (m *mockModulesInterface) SetCrudConfig(projectID string, crudConfig config.Crud) error {
	c := m.Called(projectID, crudConfig)
	return c.Error(0)
}

func (m *mockModulesInterface) SetServicesConfig(projectID string, services *config.ServicesModule) error {
	c := m.Called(projectID, services)
	return c.Error(0)
}

func (m *mockModulesInterface) GetSchemaModuleForSyncMan(projectID string) (model.SchemaEventingInterface, error) {
	c := m.Called(projectID)
	return c.Get(0).(*mockSchemaEventingInterface), c.Error(1)
}

func (m *mockModulesInterface) GetAuthModuleForSyncMan(projectID string) (model.AuthSyncManInterface, error) {
	c := m.Called(projectID)
	return c.Get(0).(model.AuthSyncManInterface), c.Error(1)
}

type mockStoreInterface struct {
	mock.Mock
}

func (m *mockStoreInterface) WatchLicense(cb func(eventType string, resourceId string, resourceType config.Resource, resource *config.License)) {
	m.Called(cb)
}

func (m *mockStoreInterface) SetLicense(ctx context.Context, resourceID string, resource *config.License) error {
	c := m.Called(ctx, resourceID, resource)
	return c.Error(0)
}

func (m *mockStoreInterface) GetGlobalConfig() (*config.Config, error) {
	c := m.Called()
	return c.Get(0).(*config.Config), c.Error(1)
}

func (m *mockStoreInterface) WatchResources(cb func(eventType string, resourceId string, resourceType config.Resource, resource interface{})) error {
	panic("implement me")
}

func (m *mockStoreInterface) SetResource(ctx context.Context, resourceID string, resource interface{}) error {
	return m.Called(ctx, resourceID, resource).Error(0)
}

func (m *mockStoreInterface) DeleteResource(ctx context.Context, resourceID string) error {
	return m.Called(ctx, resourceID).Error(0)
}

func (m *mockStoreInterface) DeleteProject(ctx context.Context, projectID string) error {
	return m.Called(ctx, projectID).Error(0)
}

func (m *mockStoreInterface) GetProjectsConfig() (config.Projects, error) {
	c := m.Called()
	return c.Get(0).(config.Projects), c.Error(1)
}

func (m *mockStoreInterface) WatchProjects(cb func(projects []*config.Project)) error {
	c := m.Called(cb)
	return c.Error(0)
}

func (m *mockStoreInterface) WatchServices(cb func(evenType string, serviceId string, projects model.ScServices)) error {
	c := m.Called(cb)
	return c.Error(0)
}

func (m *mockStoreInterface) Register() {
	m.Called()
}

func (m *mockStoreInterface) SetAdminConfig(ctx context.Context, adminConfig *config.Admin) error {
	c := m.Called(ctx, adminConfig)
	return c.Error(0)
}

func (m *mockStoreInterface) GetAdminConfig(ctx context.Context) (*config.Admin, error) {
	c := m.Called(ctx)
	return c.Get(0).(*config.Admin), c.Error(1)
}

func (m *mockStoreInterface) WatchClusterConfig(cb func(clusters []*config.Admin)) error {
	c := m.Called(cb)
	return c.Error(0)
}

type mockSchemaEventingInterface struct {
	mock.Mock
}

func (m *mockSchemaEventingInterface) CheckIfEventingIsPossible(dbAlias, col string, obj map[string]interface{}, isFind bool) (findForUpdate map[string]interface{}, present bool) {
	c := m.Called(dbAlias, col, obj, isFind)
	return map[string]interface{}{}, c.Bool(1)
}

func (m *mockSchemaEventingInterface) Parser(dbSchemas config.DatabaseSchemas) (model.Type, error) {
	c := m.Called(dbSchemas)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaValidator(ctx context.Context, col string, collectionFields model.Fields, doc map[string]interface{}) (map[string]interface{}, error) {
	c := m.Called(ctx, col, collectionFields, doc)
	return nil, c.Error(1)
}

func (m *mockSchemaEventingInterface) SchemaModifyAll(ctx context.Context, dbAlias, logicalDBName string, dbSchemas config.DatabaseSchemas) error {
	c := m.Called(ctx, dbAlias, logicalDBName, dbSchemas)
	return c.Error(0)
}

func (m *mockSchemaEventingInterface) SchemaInspection(ctx context.Context, dbAlias, project, col string) (string, error) {
	c := m.Called(ctx, dbAlias, project, col)
	return c.String(0), c.Error(1)
}

func (m *mockSchemaEventingInterface) GetSchema(dbAlias, col string) (model.Fields, bool) {
	c := m.Called(dbAlias, col)
	return c.Get(0).(model.Fields), c.Bool(1)
}
func (m *mockSchemaEventingInterface) GetSchemaForDB(ctx context.Context, dbAlias, col, format string) ([]interface{}, error) {
	c := m.Called(ctx, dbAlias, col, format)
	return c.Get(0).([]interface{}), c.Error(1)
}
