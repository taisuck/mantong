# CLAUDE.md

## 프로젝트 개요

victoria metrics에 수집된 시계열 데이터를 가공 처리하는 go 프로젝트
- 웹 백엔드 서비스
- 시계열 데이터의 실시간 가공 후 이벤트 처리

## 모듈 및 구조

- 모듈 경로: `github.com/taisuck/mantong` (저장소 루트에 단일 모듈)
- 개발 문서 경로 : `https://app.notion.com/p/taisuck/3f85839e81d7803b8d9be9c01bc94f32?v=3f85839e81d780509e95000c3c1427ed`
- DB 스키마 문서 경로 : `https://app.notion.com/p/taisuck/DB-3f95839e81d780a79a95feb489292c5d`
- 소스파일은 `src` 폴더 아래에 위치

```
src/main.go     웹 서비스를 실행하는 main
src/api         웹 api request handler
src/logging     로그 레벨 파싱, KST 날짜별 로그 파일
src/mongodb     MongoDB 연결
```

## 주요 명령어

```sh
go build ./...             # 전체 빌드
go run ./src -config config.yml -log-level debug
go vet ./...               # 정적 분석
gofmt -l .                 # 포맷되지 않은 파일 목록 (출력이 비어 있으면 정상)
gofmt -w .                 # 제자리에서 포맷 적용
go mod tidy                # import 변경 후 go.mod/go.sum 동기화
go test ./...              # 전체 테스트
go test -run '^TestName$' ./path/to/pkg   # 단일 테스트 실행
```

크로스 빌드(어느 OS에서든 상대 OS용 바이너리를 만들 수 있다. cgo를 쓰지 않는다):

```sh
scripts/build_linux.sh     scripts/build_windows.sh     # Linux 셸에서
scripts\build_linux.bat    scripts\build_windows.bat    # Windows cmd에서
```


## 언어 및 대화 지침 (Language)

- 모든 설명, 질문 답변, 커밋 메시지 작성 시 **한국어(Korean)**를 기본 언어로 사용한다.
- 전문 용어나 코드 상의 명칭(변수명, 함수명, 라이브러리 등)은 필요에 따라 영문 그대로 표기한다.

## 아키텍쳐

## 주의사항 & 금지 항목

- 기존 코드를 리팩토링할 때는 요청 받은 부분 외에는 건드리지 말 것.
