module github.com/darkrockmountain/gomail

go 1.22.5

require (
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.7.0
	github.com/SparkPost/gosparkpost v0.2.0
	github.com/aws/aws-sdk-go v1.55.5
	github.com/mailgun/mailgun-go/v4 v4.12.0
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/microsoft/kiota-abstractions-go v1.6.1
	github.com/microsoft/kiota-authentication-azure-go v1.0.2
	github.com/microsoft/kiota-http-go v1.4.3
	github.com/microsoft/kiota-serialization-form-go v1.0.0
	github.com/microsoft/kiota-serialization-json-go v1.0.7
	github.com/microsoft/kiota-serialization-multipart-go v1.0.0
	github.com/microsoft/kiota-serialization-text-go v1.0.0
	github.com/microsoftgraph/msgraph-sdk-go v1.46.0
	github.com/pkg/errors v0.9.1
	github.com/sendgrid/rest v2.6.9+incompatible
	github.com/sendgrid/sendgrid-go v3.14.0+incompatible
	github.com/stretchr/testify v1.9.0
	golang.org/x/oauth2 v0.22.0
	google.golang.org/api v0.190.0
)

require (
	cloud.google.com/go/auth v0.7.3 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.3 // indirect
	cloud.google.com/go/compute/metadata v0.5.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.13.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/cjlapao/common-go v0.0.39 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-chi/chi/v5 v5.1.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.13.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/microsoftgraph/msgraph-sdk-go-core v1.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/std-uritemplate/std-uritemplate/go v1.0.5 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.53.0 // indirect
	go.opentelemetry.io/otel v1.28.0 // indirect
	go.opentelemetry.io/otel/metric v1.28.0 // indirect
	go.opentelemetry.io/otel/trace v1.28.0 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240725223205-93522f1f2a9f // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240730163845-b1a4ccb954bf // indirect
	google.golang.org/grpc v1.64.1 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract [v0.0.0, v0.5.9] // Retract versions up to v0.5.x due to migration from GitLab to GitHub. Previous versions will not be maintained.
