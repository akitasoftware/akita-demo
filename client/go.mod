module akitasoftware.com/demo-client

go 1.20

require (
	github.com/akitasoftware/akita-libs v0.0.0-20230417213824-238b07b431fa
	github.com/akitasoftware/go-utils v0.0.0-20221207014235-6f4c9079488d
	github.com/brianvoe/gofakeit/v6 v6.21.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/martian/v3 v3.0.1
	github.com/pkg/errors v0.9.1
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
)

require (
	github.com/dukex/mixpanel v1.0.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/segmentio/analytics-go/v3 v3.2.1 // indirect
	github.com/segmentio/backo-go v1.0.0 // indirect
)

replace github.com/google/martian/v3 v3.0.1 => github.com/akitasoftware/martian/v3 v3.0.1-0.20210608174341-829c1134e9de
