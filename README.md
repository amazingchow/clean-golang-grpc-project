Your best golang gRPC project template to create general gRPC services.

## Installation

Only for linux_amd64/darwin_amd64
```shell
# install boilr
curl -sL https://github.com/tmrts/boilr/releases/download/0.3.0/install | bash -s
mkdir -p ${HOME}/.local
mv ${HOME}/bin ${HOME}/.local
# initialize template registry
${HOME}/.local/bin/boilr init
# download golang-grpc-service template
cd /path/to
git clone https://github.com/amazingchow/clean-golang-grpc-project.git
# save the golang-grpc-service template
${HOME}/.local/bin/boilr template save -f clean-golang-grpc-project golang-grpc-service-template
# list all templates
${HOME}/.local/bin/boilr template list
```

## Quickstart

```shell
# ${RepoBase} is your git repository, ${RepoGroup} is your repository's parent group
mkdir -p ${HOME}/.${RepoBase}/${RepoGroup}
cd ${HOME}/.${RepoBase}/${RepoGroup}
# update project.json to meet your project requirements, and ${RepoName} is your repository's name
${HOME}/.local/bin/boilr template use golang-grpc-service-template ${RepoName}
cd ${HOME}/.${RepoBase}/${RepoGroup}/${RepoName}
# initialize service environment
git init && git add . && git commit -am "First commit"
touch .env.local .env.secret .env.shared
./scripts/compile_pb_using_protoc_docker.sh && make init
# do what you want do...
```

## Take a quick look at the project structure

```text
.
├── devops               // where we place the Dockerfile, etc...
├── etc                  // where we place the dev-config file and prod-config file
├── internal             // where we place the project-related source code
├── protos               // where we place the proto files
├── scripts              // where we place the scripts to build proto files, to test gRPC method, to benchtest gRPC method, etc...
├── .dockerignore        // docker ignore file
├── .gitignore           // git ignore file
├── docker-compose.yml   // docker compose file
├── LICENSE              // license file
├── Makefile             // where we place the make command to manage the projet
├── requirements.txt     // runtime requirements
└── server.py            // main script
```

## License

This project is open sourced under MIT license, see the [LICENSE](LICENSE) file for more details.
