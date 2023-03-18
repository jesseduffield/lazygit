_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit 키 바인딩

## 글로벌 키 바인딩

<pre>
  <kbd>ctrl+r</kbd>: 최근에 사용한 저장소로 전환
  <kbd>pgup</kbd>: 메인 패널을 위로 스크롤 (fn+up/shift+k)
  <kbd>pgdown</kbd>: 메인 패널을 아래로로 스크롤 (fn+down/shift+j)
  <kbd>m</kbd>: view merge/rebase options
  <kbd>ctrl+p</kbd>: 커스텀 Patch 옵션 보기
  <kbd>R</kbd>: 새로고침
  <kbd>?</kbd>: 매뉴 열기
  <kbd>+</kbd>: 다음 스크린 모드 (normal/half/fullscreen)
  <kbd>_</kbd>: 이전 스크린 모드
  <kbd>ctrl+s</kbd>: view filter-by-path options
  <kbd>W</kbd>: Diff 메뉴 열기
  <kbd>ctrl+e</kbd>: Diff 메뉴 열기
  <kbd>@</kbd>: 명령어 로그 메뉴 열기
  <kbd>ctrl+w</kbd>: 공백문자를 Diff 뷰에서 표시 여부 전환
  <kbd>}</kbd>: diff 보기의 변경 사항 주위에 표시되는 컨텍스트의 크기를 늘리기
  <kbd>{</kbd>: diff 보기의 변경 사항 주위에 표시되는 컨텍스트 크기 줄이기
  <kbd>:</kbd>: execute custom command
  <kbd>z</kbd>: 되돌리기 (reflog) (실험적)
  <kbd>ctrl+z</kbd>: 다시 실행 (reflog) (실험적)
  <kbd>P</kbd>: 푸시
  <kbd>p</kbd>: 업데이트
</pre>

## List Panel Navigation

<pre>
  <kbd>,</kbd>: 이전 페이지
  <kbd>.</kbd>: 다음 페이지
  <kbd><</kbd>: 맨 위로 스크롤 
  <kbd>/</kbd>: 검색 시작
  <kbd>></kbd>: 맨 아래로 스크롤 
  <kbd>H</kbd>: 우 스크롤
  <kbd>L</kbd>: 좌 스크롤
  <kbd>]</kbd>: 이전 탭
  <kbd>[</kbd>: 다음 탭
</pre>

## Reflog

<pre>
  <kbd>ctrl+o</kbd>: 커밋 SHA를 클립보드에 복사
  <kbd>space</kbd>: 커밋을 체크아웃
  <kbd>y</kbd>: 커밋 attribute 복사
  <kbd>o</kbd>: 브라우저에서 커밋 열기
  <kbd>n</kbd>: 커밋에서 새 브랜치를 만듭니다.
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: 커밋을 복사 (cherry-pick)
  <kbd>C</kbd>: 커밋을 범위로 복사 (cherry-pick)
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>enter</kbd>: 커밋 보기
</pre>

## Stash

<pre>
  <kbd>space</kbd>: 적용
  <kbd>g</kbd>: pop
  <kbd>d</kbd>: drop
  <kbd>n</kbd>: 새 브랜치 생성
  <kbd>r</kbd>: rename stash
  <kbd>enter</kbd>: view selected item's files
</pre>

## Sub-commits

<pre>
  <kbd>ctrl+o</kbd>: 커밋 SHA를 클립보드에 복사
  <kbd>space</kbd>: 커밋을 체크아웃
  <kbd>y</kbd>: 커밋 attribute 복사
  <kbd>o</kbd>: 브라우저에서 커밋 열기
  <kbd>n</kbd>: 커밋에서 새 브랜치를 만듭니다.
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: 커밋을 복사 (cherry-pick)
  <kbd>C</kbd>: 커밋을 범위로 복사 (cherry-pick)
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>enter</kbd>: view selected item's files
</pre>

## 메인 패널 (Merging)

<pre>
  <kbd>e</kbd>: 파일 편집
  <kbd>o</kbd>: 파일 닫기
  <kbd>◀</kbd>: 이전 충돌을 선택
  <kbd>▶</kbd>: 다음 충돌을 선택
  <kbd>▲</kbd>: 이전 hunk를 선택
  <kbd>▼</kbd>: 다음 hunk를 선택
  <kbd>z</kbd>: 되돌리기
  <kbd>M</kbd>: git mergetool를 열기
  <kbd>space</kbd>: pick hunk
  <kbd>b</kbd>: pick all hunks
  <kbd>esc</kbd>: 파일 목록으로 돌아가기
</pre>

## 메인 패널 (Normal)

<pre>
  <kbd>mouse wheel ▼</kbd>: 아래로 스크롤 (fn+up)
  <kbd>mouse wheel ▲</kbd>: 위로 스크롤 (fn+down)
</pre>

## 메인 패널 (Patch Building)

<pre>
  <kbd>◀</kbd>: 이전 hunk를 선택
  <kbd>▶</kbd>: 다음 hunk를 선택
  <kbd>v</kbd>: 드래그 선택 전환
  <kbd>V</kbd>: 드래그 선택 전환
  <kbd>a</kbd>: toggle select hunk
  <kbd>ctrl+o</kbd>: 선택한 텍스트를 클립보드에 복사
  <kbd>o</kbd>: 파일 닫기
  <kbd>e</kbd>: 파일 편집
  <kbd>space</kbd>: line(s)을 패치에 추가/삭제
  <kbd>esc</kbd>: exit custom patch builder
</pre>

## 메인 패널 (Staging)

<pre>
  <kbd>◀</kbd>: 이전 hunk를 선택
  <kbd>▶</kbd>: 다음 hunk를 선택
  <kbd>v</kbd>: 드래그 선택 전환
  <kbd>V</kbd>: 드래그 선택 전환
  <kbd>a</kbd>: toggle select hunk
  <kbd>ctrl+o</kbd>: 선택한 텍스트를 클립보드에 복사
  <kbd>o</kbd>: 파일 닫기
  <kbd>e</kbd>: 파일 편집
  <kbd>esc</kbd>: 파일 목록으로 돌아가기
  <kbd>tab</kbd>: 패널 전환
  <kbd>space</kbd>: 선택한 행을 staged / unstaged
  <kbd>d</kbd>: 변경을 삭제 (git reset)
  <kbd>E</kbd>: edit hunk
  <kbd>c</kbd>: 커밋 변경내용
  <kbd>w</kbd>: commit changes without pre-commit hook
  <kbd>C</kbd>: Git 편집기를 사용하여 변경 내용을 커밋합니다.
</pre>

## 브랜치

<pre>
  <kbd>ctrl+o</kbd>: 브랜치명을 클립보드에 복사
  <kbd>i</kbd>: git-flow 옵션 보기
  <kbd>space</kbd>: 체크아웃
  <kbd>n</kbd>: 새 브랜치 생성
  <kbd>o</kbd>: 풀 리퀘스트 생성
  <kbd>O</kbd>: 풀 리퀘스트 생성 옵션
  <kbd>ctrl+y</kbd>: 풀 리퀘스트 URL을 클립보드에 복사
  <kbd>c</kbd>: 이름으로 체크아웃
  <kbd>F</kbd>: 강제 체크아웃
  <kbd>d</kbd>: 브랜치 삭제
  <kbd>r</kbd>: 체크아웃된 브랜치를 이 브랜치에 리베이스
  <kbd>M</kbd>: 현재 브랜치에 병합
  <kbd>f</kbd>: fast-forward this branch from its upstream
  <kbd>T</kbd>: 태그를 생성
  <kbd>g</kbd>: view reset options
  <kbd>R</kbd>: 브랜치 이름 변경
  <kbd>u</kbd>: set/unset upstream
  <kbd>enter</kbd>: 커밋 보기
</pre>

## 상태

<pre>
  <kbd>e</kbd>: 설정 파일 수정
  <kbd>o</kbd>: 설정 파일 열기
  <kbd>u</kbd>: 업데이트 확인
  <kbd>enter</kbd>: 최근에 사용한 저장소로 전환
  <kbd>a</kbd>: 모든 브랜치 로그 표시
</pre>

## 서브모듈

<pre>
  <kbd>ctrl+o</kbd>: 서브모듈 이름을 클립보드에 복사
  <kbd>enter</kbd>: 서브모듈 열기
  <kbd>d</kbd>: 서브모듈 삭제
  <kbd>u</kbd>: 서브모듈 업데이트
  <kbd>n</kbd>: 새로운 서브모듈 추가
  <kbd>e</kbd>: 서브모듈의 URL을 수정
  <kbd>i</kbd>: 서브모듈 초기화
  <kbd>b</kbd>: view bulk submodule options
</pre>

## 원격

<pre>
  <kbd>f</kbd>: 원격을 업데이트
  <kbd>n</kbd>: 새로운 Remote 추가
  <kbd>d</kbd>: Remote를 삭제
  <kbd>e</kbd>: Remote를 수정
</pre>

## 원격 브랜치

<pre>
  <kbd>ctrl+o</kbd>: 브랜치명을 클립보드에 복사
  <kbd>space</kbd>: 체크아웃
  <kbd>n</kbd>: 새 브랜치 생성
  <kbd>M</kbd>: 현재 브랜치에 병합
  <kbd>r</kbd>: 체크아웃된 브랜치를 이 브랜치에 리베이스
  <kbd>d</kbd>: 브랜치 삭제
  <kbd>u</kbd>: set as upstream of checked-out branch
  <kbd>esc</kbd>: 원격목록으로 돌아가기
  <kbd>g</kbd>: view reset options
  <kbd>enter</kbd>: 커밋 보기
</pre>

## 커밋

<pre>
  <kbd>ctrl+o</kbd>: 커밋 SHA를 클립보드에 복사
  <kbd>ctrl+r</kbd>: reset cherry-picked (copied) commits selection
  <kbd>b</kbd>: bisect 옵션 보기
  <kbd>s</kbd>: squash down
  <kbd>f</kbd>: fixup commit
  <kbd>r</kbd>: 커밋메시지 변경
  <kbd>R</kbd>: 에디터에서 커밋메시지 수정
  <kbd>d</kbd>: 커밋 삭제
  <kbd>e</kbd>: 커밋을 편집
  <kbd>p</kbd>: pick commit (when mid-rebase)
  <kbd>F</kbd>: create fixup commit for this commit
  <kbd>S</kbd>: squash all 'fixup!' commits above selected commit (autosquash)
  <kbd>ctrl+j</kbd>: 커밋을 1개 아래로 이동
  <kbd>ctrl+k</kbd>: 커밋을 1개 위로 이동
  <kbd>v</kbd>: 커밋을 붙여넣기 (cherry-pick)
  <kbd>A</kbd>: amend commit with staged changes
  <kbd>a</kbd>: reset commit author
  <kbd>t</kbd>: 커밋 되돌리기
  <kbd>T</kbd>: tag commit
  <kbd>ctrl+l</kbd>: 로그 메뉴 열기
  <kbd>space</kbd>: 커밋을 체크아웃
  <kbd>y</kbd>: 커밋 attribute 복사
  <kbd>o</kbd>: 브라우저에서 커밋 열기
  <kbd>n</kbd>: 커밋에서 새 브랜치를 만듭니다.
  <kbd>g</kbd>: view reset options
  <kbd>c</kbd>: 커밋을 복사 (cherry-pick)
  <kbd>C</kbd>: 커밋을 범위로 복사 (cherry-pick)
  <kbd>enter</kbd>: view selected item's files
</pre>

## 커밋 파일

<pre>
  <kbd>ctrl+o</kbd>: 커밋한 파일명을 클립보드에 복사
  <kbd>c</kbd>: checkout file
  <kbd>d</kbd>: discard this commit's changes to this file
  <kbd>o</kbd>: 파일 닫기
  <kbd>e</kbd>: 파일 편집
  <kbd>space</kbd>: toggle file included in patch
  <kbd>a</kbd>: toggle all files included in patch
  <kbd>enter</kbd>: enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: 파일 트리뷰로 전환
</pre>

## 태그

<pre>
  <kbd>space</kbd>: 체크아웃
  <kbd>d</kbd>: 태그 삭제
  <kbd>P</kbd>: 태그를 push
  <kbd>n</kbd>: 태그를 생성
  <kbd>g</kbd>: view reset options
  <kbd>enter</kbd>: 커밋 보기
</pre>

## 파일

<pre>
  <kbd>ctrl+o</kbd>: 파일명을 클립보드에 복사
  <kbd>d</kbd>: view 'discard changes' options
  <kbd>space</kbd>: Staged 전환
  <kbd>ctrl+b</kbd>: 파일을 필터하기 (Staged/unstaged)
  <kbd>c</kbd>: 커밋 변경내용
  <kbd>w</kbd>: commit changes without pre-commit hook
  <kbd>A</kbd>: 마지맛 커밋 수정
  <kbd>C</kbd>: Git 편집기를 사용하여 변경 내용을 커밋합니다.
  <kbd>e</kbd>: 파일 편집
  <kbd>o</kbd>: 파일 닫기
  <kbd>i</kbd>: ignore file
  <kbd>r</kbd>: 파일 새로고침
  <kbd>s</kbd>: 변경사항을 Stash
  <kbd>S</kbd>: Stash 옵션 보기
  <kbd>a</kbd>: 모든 변경을 Staged/unstaged으로 전환
  <kbd>enter</kbd>: stage individual hunks/lines for file, or collapse/expand for directory
  <kbd>g</kbd>: view upstream reset options
  <kbd>D</kbd>: view reset options
  <kbd>`</kbd>: 파일 트리뷰로 전환
  <kbd>M</kbd>: git mergetool를 열기
  <kbd>f</kbd>: fetch
</pre>
