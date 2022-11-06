/*
Lint package for current project

Includes checkers from https://staticcheck.io/:
  - all SA* https://staticcheck.io/docs/checks/#SA
  - ST1013 https://staticcheck.io/docs/checks/#ST1013
  - QF1006 https://staticcheck.io/docs/checks/#QF1006
  - S1011 https://staticcheck.io/docs/checks/#S1011
  - and custom CheckOSExit

For run:

	go run cmd/staticlint/main.go ./...
*/
package main
