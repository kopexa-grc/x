module github.com/kopexa-grc/x

go 1.24.5

tool (
	github.com/hashicorp/copywrite
	gotest.tools/gotestsum
	mvdan.cc/gofumpt
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.19.1
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.13.0
	github.com/cockroachdb/errors v1.12.0
	github.com/google/uuid v1.6.0
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f
	github.com/labstack/echo/v4 v4.13.4
	github.com/labstack/gommon v0.4.2
	github.com/muesli/termenv v0.16.0
	github.com/planetscale/vtprotobuf v0.6.0
	github.com/redis/go-redis/v9 v9.17.0
	github.com/rs/zerolog v1.34.0
	github.com/stretchr/testify v1.11.1
	google.golang.org/protobuf v1.36.10
	sigs.k8s.io/yaml v1.6.0
)

require (
	github.com/AlecAivazis/survey/v2 v2.3.7 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.5.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20230217124315-7d5c6f04bbb8 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/bitfield/gotestdox v0.2.2 // indirect
	github.com/bmatcuk/doublestar/v4 v4.6.0 // indirect
	github.com/bradleyfalzon/ghinstallation/v2 v2.5.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/cli/go-gh/v2 v2.11.2 // indirect
	github.com/cli/safeexec v1.0.0 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/cockroachdb/logtags v0.0.0-20230118201751-21c54148d20b // indirect
	github.com/cockroachdb/redact v1.1.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dnephin/pflag v1.0.7 // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/getsentry/sentry-go v0.27.0 // indirect
	github.com/go-openapi/errors v0.20.2 // indirect
	github.com/go-openapi/strfmt v0.21.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/go-github/v45 v45.2.0 // indirect
	github.com/google/go-github/v53 v53.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/hashicorp/copywrite v0.22.0 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jedib0t/go-pretty v4.3.0+incompatible // indirect
	github.com/jedib0t/go-pretty/v6 v6.4.6 // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/knadh/koanf v1.5.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mergestat/timediff v0.0.3 // indirect
	github.com/mgutz/ansi v0.0.0-20200706080929-d51e80ef957d // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/samber/lo v1.37.0 // indirect
	github.com/spf13/cobra v1.6.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/thanhpk/randstr v1.0.4 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	go.mongodb.org/mongo-driver v1.10.0 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/crypto v0.41.0 // indirect
	golang.org/x/exp v0.0.0-20220303212507-bbda1eaf7a17 // indirect
	golang.org/x/mod v0.26.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/oauth2 v0.8.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/term v0.34.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	golang.org/x/tools v0.35.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gotest.tools/gotestsum v1.12.3 // indirect
	mvdan.cc/gofumpt v0.8.0 // indirect
)
