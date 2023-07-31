_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go run scripts/cheatsheet/main.go generate` from the project root._

# Lazygit 키 바인딩

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## 글로벌 키 바인딩

<pre>
  <kbd>&lt;c-r&gt;</kbd>: 최근에 사용한 저장소로 전환
  <kbd>&lt;pgup&gt;</kbd>: 메인 패널을 위로 스크롤 (fn+up/shift+k)
  <kbd>&lt;pgdown&gt;</kbd>: 메인 패널을 아래로로 스크롤 (fn+down/shift+j)
  <kbd>@</kbd>: 명령어 로그 메뉴 열기
  <kbd>}</kbd>: Diff 보기의 변경 사항 주위에 표시되는 컨텍스트의 크기를 늘리기
  <kbd>{</kbd>: Diff 보기의 변경 사항 주위에 표시되는 컨텍스트 크기 줄이기
  <kbd>:</kbd>: Execute custom command
  <kbd>&lt;c-p&gt;</kbd>: 커스텀 Patch 옵션 보기
  <kbd>m</kbd>: View merge/rebase options
  <kbd>R</kbd>: 새로고침
  <kbd>+</kbd>: 다음 스크린 모드 (normal/half/fullscreen)
  <kbd>_</kbd>: 이전 스크린 모드
  <kbd>?</kbd>: 매뉴 열기
  <kbd>&lt;c-s&gt;</kbd>: View filter-by-path options
  <kbd>W</kbd>: Diff 메뉴 열기
  <kbd>&lt;c-e&gt;</kbd>: Diff 메뉴 열기
  <kbd>&lt;c-w&gt;</kbd>: 공백문자를 Diff 뷰에서 표시 여부 전환
  <kbd>z</kbd>: 되돌리기 (reflog) (실험적)
  <kbd>&lt;c-z&gt;</kbd>: 다시 실행 (reflog) (실험적)
  <kbd>P</kbd>: 푸시
  <kbd>p</kbd>: 업데이트
</pre>

## List panel navigation

<pre>
  <kbd>,</kbd>: 이전 페이지
  <kbd>.</kbd>: 다음 페이지
  <kbd>&lt;</kbd>: 맨 위로 스크롤 
  <kbd>&gt;</kbd>: 맨 아래로 스크롤 
  <kbd>/</kbd>: 검색 시작
  <kbd>H</kbd>: 우 스크롤
  <kbd>L</kbd>: 좌 스크롤
  <kbd>]</kbd>: 이전 탭
  <kbd>[</kbd>: 다음 탭
</pre>

## Reflog

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 커밋 SHA를 클립보드에 복사
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: 커밋을 체크아웃
  <kbd>y</kbd>: 커밋 attribute 복사
  <kbd>o</kbd>: 브라우저에서 커밋 열기
  <kbd>n</kbd>: 커밋에서 새 브랜치를 만듭니다.
  <kbd>g</kbd>: View reset options
  <kbd>c</kbd>: 커밋을 복사 (cherry-pick)
  <kbd>C</kbd>: 커밋을 범위로 복사 (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>&lt;enter&gt;</kbd>: 커밋 보기
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Stash

<pre>
  <kbd>&lt;space&gt;</kbd>: 적용
  <kbd>g</kbd>: Pop
  <kbd>d</kbd>: Drop
  <kbd>n</kbd>: 새 브랜치 생성
  <kbd>r</kbd>: Rename stash
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: View selected item's files
  <kbd>/</kbd>: Filter the current view by text
</pre>

## Sub-commits

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 커밋 SHA를 클립보드에 복사
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: 커밋을 체크아웃
  <kbd>y</kbd>: 커밋 attribute 복사
  <kbd>o</kbd>: 브라우저에서 커밋 열기
  <kbd>n</kbd>: 커밋에서 새 브랜치를 만듭니다.
  <kbd>g</kbd>: View reset options
  <kbd>c</kbd>: 커밋을 복사 (cherry-pick)
  <kbd>C</kbd>: 커밋을 범위로 복사 (cherry-pick)
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>&lt;enter&gt;</kbd>: View selected item's files
  <kbd>/</kbd>: 검색 시작
</pre>

## Worktrees

<pre>
  <kbd>n</kbd>: Create worktree
  <kbd>&lt;space&gt;</kbd>: Switch to worktree
  <kbd>&lt;enter&gt;</kbd>: Switch to worktree
  <kbd>o</kbd>: Open in editor
  <kbd>d</kbd>: Remove worktree
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 메뉴

<pre>
  <kbd>&lt;enter&gt;</kbd>: 실행
  <kbd>&lt;esc&gt;</kbd>: 닫기
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 메인 패널 (Merging)

<pre>
  <kbd>e</kbd>: 파일 편집
  <kbd>o</kbd>: 파일 닫기
  <kbd>&lt;left&gt;</kbd>: 이전 충돌을 선택
  <kbd>&lt;right&gt;</kbd>: 다음 충돌을 선택
  <kbd>&lt;up&gt;</kbd>: 이전 hunk를 선택
  <kbd>&lt;down&gt;</kbd>: 다음 hunk를 선택
  <kbd>z</kbd>: 되돌리기
  <kbd>M</kbd>: Git mergetool를 열기
  <kbd>&lt;space&gt;</kbd>: Pick hunk
  <kbd>b</kbd>: Pick all hunks
  <kbd>&lt;esc&gt;</kbd>: 파일 목록으로 돌아가기
</pre>

## 메인 패널 (Normal)

<pre>
  <kbd>mouse wheel down</kbd>: 아래로 스크롤 (fn+up)
  <kbd>mouse wheel up</kbd>: 위로 스크롤 (fn+down)
</pre>

## 메인 패널 (Patch Building)

<pre>
  <kbd>&lt;left&gt;</kbd>: 이전 hunk를 선택
  <kbd>&lt;right&gt;</kbd>: 다음 hunk를 선택
  <kbd>v</kbd>: 드래그 선택 전환
  <kbd>V</kbd>: 드래그 선택 전환
  <kbd>a</kbd>: Toggle select hunk
  <kbd>&lt;c-o&gt;</kbd>: 선택한 텍스트를 클립보드에 복사
  <kbd>o</kbd>: 파일 닫기
  <kbd>e</kbd>: 파일 편집
  <kbd>&lt;space&gt;</kbd>: Line(s)을 패치에 추가/삭제
  <kbd>&lt;esc&gt;</kbd>: Exit custom patch builder
  <kbd>/</kbd>: 검색 시작
</pre>

## 메인 패널 (Staging)

<pre>
  <kbd>&lt;left&gt;</kbd>: 이전 hunk를 선택
  <kbd>&lt;right&gt;</kbd>: 다음 hunk를 선택
  <kbd>v</kbd>: 드래그 선택 전환
  <kbd>V</kbd>: 드래그 선택 전환
  <kbd>a</kbd>: Toggle select hunk
  <kbd>&lt;c-o&gt;</kbd>: 선택한 텍스트를 클립보드에 복사
  <kbd>o</kbd>: 파일 닫기
  <kbd>e</kbd>: 파일 편집
  <kbd>&lt;esc&gt;</kbd>: 파일 목록으로 돌아가기
  <kbd>&lt;tab&gt;</kbd>: 패널 전환
  <kbd>&lt;space&gt;</kbd>: 선택한 행을 staged / unstaged
  <kbd>d</kbd>: 변경을 삭제 (git reset)
  <kbd>E</kbd>: Edit hunk
  <kbd>c</kbd>: 커밋 변경내용
  <kbd>w</kbd>: Commit changes without pre-commit hook
  <kbd>C</kbd>: Git 편집기를 사용하여 변경 내용을 커밋합니다.
  <kbd>/</kbd>: 검색 시작
</pre>

## 브랜치

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 브랜치명을 클립보드에 복사
  <kbd>i</kbd>: Git-flow 옵션 보기
  <kbd>&lt;space&gt;</kbd>: 체크아웃
  <kbd>n</kbd>: 새 브랜치 생성
  <kbd>o</kbd>: 풀 리퀘스트 생성
  <kbd>O</kbd>: 풀 리퀘스트 생성 옵션
  <kbd>&lt;c-y&gt;</kbd>: 풀 리퀘스트 URL을 클립보드에 복사
  <kbd>c</kbd>: 이름으로 체크아웃
  <kbd>F</kbd>: 강제 체크아웃
  <kbd>d</kbd>: 브랜치 삭제
  <kbd>r</kbd>: 체크아웃된 브랜치를 이 브랜치에 리베이스
  <kbd>M</kbd>: 현재 브랜치에 병합
  <kbd>f</kbd>: Fast-forward this branch from its upstream
  <kbd>T</kbd>: 태그를 생성
  <kbd>g</kbd>: View reset options
  <kbd>R</kbd>: 브랜치 이름 변경
  <kbd>u</kbd>: Set/Unset upstream
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: 커밋 보기
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 상태

<pre>
  <kbd>o</kbd>: 설정 파일 열기
  <kbd>e</kbd>: 설정 파일 수정
  <kbd>u</kbd>: 업데이트 확인
  <kbd>&lt;enter&gt;</kbd>: 최근에 사용한 저장소로 전환
  <kbd>a</kbd>: 모든 브랜치 로그 표시
</pre>

## 서브모듈

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 서브모듈 이름을 클립보드에 복사
  <kbd>&lt;enter&gt;</kbd>: 서브모듈 열기
  <kbd>&lt;space&gt;</kbd>: 서브모듈 열기
  <kbd>d</kbd>: 서브모듈 삭제
  <kbd>u</kbd>: 서브모듈 업데이트
  <kbd>n</kbd>: 새로운 서브모듈 추가
  <kbd>e</kbd>: 서브모듈의 URL을 수정
  <kbd>i</kbd>: 서브모듈 초기화
  <kbd>b</kbd>: View bulk submodule options
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 원격

<pre>
  <kbd>f</kbd>: 원격을 업데이트
  <kbd>n</kbd>: 새로운 Remote 추가
  <kbd>d</kbd>: Remote를 삭제
  <kbd>e</kbd>: Remote를 수정
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 원격 브랜치

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 브랜치명을 클립보드에 복사
  <kbd>&lt;space&gt;</kbd>: 체크아웃
  <kbd>n</kbd>: 새 브랜치 생성
  <kbd>M</kbd>: 현재 브랜치에 병합
  <kbd>r</kbd>: 체크아웃된 브랜치를 이 브랜치에 리베이스
  <kbd>d</kbd>: 브랜치 삭제
  <kbd>u</kbd>: Set as upstream of checked-out branch
  <kbd>g</kbd>: View reset options
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: 커밋 보기
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 커밋

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 커밋 SHA를 클립보드에 복사
  <kbd>&lt;c-r&gt;</kbd>: Reset cherry-picked (copied) commits selection
  <kbd>b</kbd>: Bisect 옵션 보기
  <kbd>s</kbd>: Squash down
  <kbd>f</kbd>: Fixup commit
  <kbd>r</kbd>: 커밋메시지 변경
  <kbd>R</kbd>: 에디터에서 커밋메시지 수정
  <kbd>d</kbd>: 커밋 삭제
  <kbd>e</kbd>: 커밋을 편집
  <kbd>p</kbd>: Pick commit (when mid-rebase)
  <kbd>F</kbd>: Create fixup commit for this commit
  <kbd>S</kbd>: Squash all 'fixup!' commits above selected commit (autosquash)
  <kbd>&lt;c-j&gt;</kbd>: 커밋을 1개 아래로 이동
  <kbd>&lt;c-k&gt;</kbd>: 커밋을 1개 위로 이동
  <kbd>v</kbd>: 커밋을 붙여넣기 (cherry-pick)
  <kbd>B</kbd>: Mark commit as base commit for rebase
  <kbd>A</kbd>: Amend commit with staged changes
  <kbd>a</kbd>: Set/Reset commit author
  <kbd>t</kbd>: 커밋 되돌리기
  <kbd>T</kbd>: Tag commit
  <kbd>&lt;c-l&gt;</kbd>: 로그 메뉴 열기
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;space&gt;</kbd>: 커밋을 체크아웃
  <kbd>y</kbd>: 커밋 attribute 복사
  <kbd>o</kbd>: 브라우저에서 커밋 열기
  <kbd>n</kbd>: 커밋에서 새 브랜치를 만듭니다.
  <kbd>g</kbd>: View reset options
  <kbd>c</kbd>: 커밋을 복사 (cherry-pick)
  <kbd>C</kbd>: 커밋을 범위로 복사 (cherry-pick)
  <kbd>&lt;enter&gt;</kbd>: View selected item's files
  <kbd>/</kbd>: 검색 시작
</pre>

## 커밋 파일

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 커밋한 파일명을 클립보드에 복사
  <kbd>c</kbd>: Checkout file
  <kbd>d</kbd>: Discard this commit's changes to this file
  <kbd>o</kbd>: 파일 닫기
  <kbd>e</kbd>: 파일 편집
  <kbd>&lt;space&gt;</kbd>: Toggle file included in patch
  <kbd>a</kbd>: Toggle all files included in patch
  <kbd>&lt;enter&gt;</kbd>: Enter file to add selected lines to the patch (or toggle directory collapsed)
  <kbd>`</kbd>: 파일 트리뷰로 전환
  <kbd>/</kbd>: 검색 시작
</pre>

## 커밋메시지

<pre>
  <kbd>&lt;enter&gt;</kbd>: 확인
  <kbd>&lt;esc&gt;</kbd>: 닫기
</pre>

## 태그

<pre>
  <kbd>&lt;space&gt;</kbd>: 체크아웃
  <kbd>d</kbd>: 태그 삭제
  <kbd>P</kbd>: 태그를 push
  <kbd>n</kbd>: 태그를 생성
  <kbd>g</kbd>: View reset options
  <kbd>w</kbd>: View worktree options
  <kbd>&lt;enter&gt;</kbd>: 커밋 보기
  <kbd>/</kbd>: Filter the current view by text
</pre>

## 파일

<pre>
  <kbd>&lt;c-o&gt;</kbd>: 파일명을 클립보드에 복사
  <kbd>d</kbd>: View 'discard changes' options
  <kbd>&lt;space&gt;</kbd>: Staged 전환
  <kbd>&lt;c-b&gt;</kbd>: 파일을 필터하기 (Staged/unstaged)
  <kbd>c</kbd>: 커밋 변경내용
  <kbd>w</kbd>: Commit changes without pre-commit hook
  <kbd>A</kbd>: 마지맛 커밋 수정
  <kbd>C</kbd>: Git 편집기를 사용하여 변경 내용을 커밋합니다.
  <kbd>e</kbd>: 파일 편집
  <kbd>o</kbd>: 파일 닫기
  <kbd>i</kbd>: Ignore file
  <kbd>r</kbd>: 파일 새로고침
  <kbd>s</kbd>: 변경사항을 Stash
  <kbd>S</kbd>: Stash 옵션 보기
  <kbd>a</kbd>: 모든 변경을 Staged/unstaged으로 전환
  <kbd>&lt;enter&gt;</kbd>: Stage individual hunks/lines for file, or collapse/expand for directory
  <kbd>g</kbd>: View upstream reset options
  <kbd>D</kbd>: View reset options
  <kbd>`</kbd>: 파일 트리뷰로 전환
  <kbd>M</kbd>: Git mergetool를 열기
  <kbd>f</kbd>: Fetch
  <kbd>/</kbd>: 검색 시작
</pre>

## 확인 패널

<pre>
  <kbd>&lt;enter&gt;</kbd>: 확인
  <kbd>&lt;esc&gt;</kbd>: 닫기/취소
</pre>
