module github.com/ElisaOyj/namespace-controller

go 1.15

replace k8s.io/api => k8s.io/api v0.0.0-20190819141258-3544db3b9e44

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190819141724-e14f31a72a77

require (
	github.com/imdario/mergo v0.3.8 // indirect
	github.com/mattbaird/jsonpatch v0.0.0-20171005235357-81af80346b1a
	github.com/stretchr/testify v1.3.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	golang.org/x/time v0.0.0-20190921001708-c4c64cad1fd0 // indirect
	gopkg.in/yaml.v2 v2.2.1
	k8s.io/api v0.0.0-20191016225839-816a9b7df678
	k8s.io/apimachinery v0.0.0-20191017185446-6e68a40eebf9
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20191010214722-8d271d903fe4 // indirect
)
