#!/bin/bash 

cd backend/internal/services/implementation
go test -gcflags=all=-l -tags=unit -cover -json ./... > testServices.log
~/go/bin/test-report -f ./testServices.log -o .


cd -
cd backend/internal/repository/postgres_repo
go test -gcflags=all=-l -tags=unit -cover -json ./... > testRepos.log
~/go/bin/test-report -f ./testRepos.log -o .

cd -
rm backend/internal/repository/postgres_repo/testRepos.log
rm backend/internal/services/implementation/testServices.log

mv backend/internal/services/implementation/report.html reportService.html
mv backend/internal/repository/postgres_repo/report.html reportRepos.html
