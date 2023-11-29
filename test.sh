#!/bin/bash 

cd backend/internal/services/implementation
go test -gcflags=all=-l -tags=unit -cover -coverprofile=serviceTest.out -json ./... > testServices.log
~/go/bin/test-report -f ./testServices.log -o .


cd -
cd backend/internal/repository/postgres_repo
go test -gcflags=all=-l -tags=unit -cover -coverprofile=repoTest.out -json ./... > testRepos.log
~/go/bin/test-report -f ./testRepos.log -o .

cd -
rm backend/internal/repository/postgres_repo/testRepos.log
rm backend/internal/services/implementation/testServices.log

mv backend/internal/services/implementation/report.html reportService.html
mv backend/internal/repository/postgres_repo/report.html reportRepos.html

# go tool cover -html=serviceTest.out
# go tool cover -html=repoTest.out

# go test recordPostgres_test.go record.go setup.go -gcflags=all=-l -tags=integration -cover -coverprofile=serviceTest.out -json  > testIntegration.log
go test petPostgres_test.go recordPostgres_test.go pet.go record.go setup.go  -gcflags=all=-l -tags=integration -cover -coverprofile=serviceIntegration.out -json  > testIntegration.log