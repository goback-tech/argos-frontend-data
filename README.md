# argos-frontend-grafana
그라파나 기반 ARGOS  데이터 시각화 도구 


BUILD INSTRUCTIONS.

1. INSTALL GO LANG

1.1 GO 다운로드 및 설치
1.2 .bashrc 에 등록

export GOROOT="/usr/local/go"
export PATH="$PATH:$GOROOT/bin

2. BUILDING THE BACKEND

2.1 cd (소스경로)/src/github.com/grafana/grafana
2.2 go run build.go setup
2.3 go run build.go build

3. BUIDING THE FRONTEND ASSETS

3.1 npm install -g yarn
3.2 yarn install --pure-lockfile
3.3 yarn watch

4. CREATE OPTIMIZED RELEASE PACKAGES

4.1 go run build.go build PACKAGES

5. CONFIGURATION

아래의 순서대로 오버라이드 적용

5.1 grafana.ini 
5.2 custom.ini
