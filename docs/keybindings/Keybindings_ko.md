_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit 키 바인딩

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## 글로벌 키 바인딩

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | 최근에 사용한 저장소로 전환 |  |
| `` <pgup> (fn+up/shift+k) `` | 메인 패널을 위로 스크롤 |  |
| `` <pgdown> (fn+down/shift+j) `` | 메인 패널을 아래로로 스크롤 |  |
| `` @ `` | 명령어 로그 메뉴 열기 |  |
| `` } `` | Diff 보기의 변경 사항 주위에 표시되는 컨텍스트의 크기를 늘리기 |  |
| `` { `` | Diff 보기의 변경 사항 주위에 표시되는 컨텍스트 크기 줄이기 |  |
| `` : `` | Execute custom command |  |
| `` <c-p> `` | 커스텀 Patch 옵션 보기 |  |
| `` m `` | View merge/rebase options |  |
| `` R `` | 새로고침 |  |
| `` + `` | 다음 스크린 모드 (normal/half/fullscreen) |  |
| `` _ `` | 이전 스크린 모드 |  |
| `` ? `` | 매뉴 열기 |  |
| `` <c-s> `` | View filter-by-path options |  |
| `` W `` | Diff 메뉴 열기 |  |
| `` <c-e> `` | Diff 메뉴 열기 |  |
| `` <c-w> `` | 공백문자를 Diff 뷰에서 표시 여부 전환 |  |
| `` z `` | 되돌리기 (reflog) (실험적) | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | 다시 실행 (reflog) (실험적) | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` P `` | 푸시 |  |
| `` p `` | 업데이트 |  |

## List panel navigation

| Key | Action | Info |
|-----|--------|-------------|
| `` , `` | 이전 페이지 |  |
| `` . `` | 다음 페이지 |  |
| `` < `` | 맨 위로 스크롤  |  |
| `` > `` | 맨 아래로 스크롤  |  |
| `` v `` | 드래그 선택 전환 |  |
| `` <s-down> `` | Range select down |  |
| `` <s-up> `` | Range select up |  |
| `` / `` | 검색 시작 |  |
| `` H `` | 우 스크롤 |  |
| `` L `` | 좌 스크롤 |  |
| `` ] `` | 이전 탭 |  |
| `` [ `` | 다음 탭 |  |

## Reflog

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 커밋 SHA를 클립보드에 복사 |  |
| `` w `` | View worktree options |  |
| `` <space> `` | 커밋을 체크아웃 |  |
| `` y `` | 커밋 attribute 복사 |  |
| `` o `` | 브라우저에서 커밋 열기 |  |
| `` n `` | 커밋에서 새 브랜치를 만듭니다. |  |
| `` g `` | View reset options |  |
| `` C `` | 커밋을 복사 (cherry-pick) |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 커밋 보기 |  |
| `` / `` | Filter the current view by text |  |

## Stash

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 적용 |  |
| `` g `` | Pop |  |
| `` d `` | Drop |  |
| `` n `` | 새 브랜치 생성 |  |
| `` r `` | Rename stash |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | View selected item's files |  |
| `` / `` | Filter the current view by text |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 커밋 SHA를 클립보드에 복사 |  |
| `` w `` | View worktree options |  |
| `` <space> `` | 커밋을 체크아웃 |  |
| `` y `` | 커밋 attribute 복사 |  |
| `` o `` | 브라우저에서 커밋 열기 |  |
| `` n `` | 커밋에서 새 브랜치를 만듭니다. |  |
| `` g `` | View reset options |  |
| `` C `` | 커밋을 복사 (cherry-pick) |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View selected item's files |  |
| `` / `` | 검색 시작 |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | Create worktree |  |
| `` <space> `` | Switch to worktree |  |
| `` <enter> `` | Switch to worktree |  |
| `` o `` | Open in editor |  |
| `` d `` | Remove worktree |  |
| `` / `` | Filter the current view by text |  |

## 메뉴

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 실행 |  |
| `` <esc> `` | 닫기 |  |
| `` / `` | Filter the current view by text |  |

## 메인 패널 (Merging)

| Key | Action | Info |
|-----|--------|-------------|
| `` e `` | 파일 편집 |  |
| `` o `` | 파일 닫기 |  |
| `` <left> `` | 이전 충돌을 선택 |  |
| `` <right> `` | 다음 충돌을 선택 |  |
| `` <up> `` | 이전 hunk를 선택 |  |
| `` <down> `` | 다음 hunk를 선택 |  |
| `` z `` | 되돌리기 |  |
| `` M `` | Git mergetool를 열기 |  |
| `` <space> `` | Pick hunk |  |
| `` b `` | Pick all hunks |  |
| `` <esc> `` | 파일 목록으로 돌아가기 |  |

## 메인 패널 (Normal)

| Key | Action | Info |
|-----|--------|-------------|
| `` mouse wheel down (fn+up) `` | 아래로 스크롤 |  |
| `` mouse wheel up (fn+down) `` | 위로 스크롤 |  |

## 메인 패널 (Patch Building)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 이전 hunk를 선택 |  |
| `` <right> `` | 다음 hunk를 선택 |  |
| `` v `` | 드래그 선택 전환 |  |
| `` a `` | Toggle select hunk |  |
| `` <c-o> `` | 선택한 텍스트를 클립보드에 복사 |  |
| `` o `` | 파일 닫기 |  |
| `` e `` | 파일 편집 |  |
| `` <space> `` | Line(s)을 패치에 추가/삭제 |  |
| `` <esc> `` | Exit custom patch builder |  |
| `` / `` | 검색 시작 |  |

## 메인 패널 (Staging)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 이전 hunk를 선택 |  |
| `` <right> `` | 다음 hunk를 선택 |  |
| `` v `` | 드래그 선택 전환 |  |
| `` a `` | Toggle select hunk |  |
| `` <c-o> `` | 선택한 텍스트를 클립보드에 복사 |  |
| `` o `` | 파일 닫기 |  |
| `` e `` | 파일 편집 |  |
| `` <esc> `` | 파일 목록으로 돌아가기 |  |
| `` <tab> `` | 패널 전환 |  |
| `` <space> `` | 선택한 행을 staged / unstaged |  |
| `` d `` | 변경을 삭제 (git reset) |  |
| `` E `` | Edit hunk |  |
| `` c `` | 커밋 변경내용 |  |
| `` w `` | Commit changes without pre-commit hook |  |
| `` C `` | Git 편집기를 사용하여 변경 내용을 커밋합니다. |  |
| `` / `` | 검색 시작 |  |

## 브랜치

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 브랜치명을 클립보드에 복사 |  |
| `` i `` | Git-flow 옵션 보기 |  |
| `` <space> `` | 체크아웃 |  |
| `` n `` | 새 브랜치 생성 |  |
| `` o `` | 풀 리퀘스트 생성 |  |
| `` O `` | 풀 리퀘스트 생성 옵션 |  |
| `` <c-y> `` | 풀 리퀘스트 URL을 클립보드에 복사 |  |
| `` c `` | 이름으로 체크아웃 |  |
| `` F `` | 강제 체크아웃 |  |
| `` d `` | View delete options |  |
| `` r `` | 체크아웃된 브랜치를 이 브랜치에 리베이스 |  |
| `` M `` | 현재 브랜치에 병합 |  |
| `` f `` | Fast-forward this branch from its upstream |  |
| `` T `` | 태그를 생성 |  |
| `` s `` | Sort order |  |
| `` g `` | View reset options |  |
| `` R `` | 브랜치 이름 변경 |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream |
| `` w `` | View worktree options |  |
| `` <enter> `` | 커밋 보기 |  |
| `` / `` | Filter the current view by text |  |

## 상태

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | 설정 파일 열기 |  |
| `` e `` | 설정 파일 수정 |  |
| `` u `` | 업데이트 확인 |  |
| `` <enter> `` | 최근에 사용한 저장소로 전환 |  |
| `` a `` | 모든 브랜치 로그 표시 |  |

## 서브모듈

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 서브모듈 이름을 클립보드에 복사 |  |
| `` <enter> `` | 서브모듈 열기 |  |
| `` <space> `` | 서브모듈 열기 |  |
| `` d `` | 서브모듈 삭제 |  |
| `` u `` | 서브모듈 업데이트 |  |
| `` n `` | 새로운 서브모듈 추가 |  |
| `` e `` | 서브모듈의 URL을 수정 |  |
| `` i `` | 서브모듈 초기화 |  |
| `` b `` | View bulk submodule options |  |
| `` / `` | Filter the current view by text |  |

## 원격

| Key | Action | Info |
|-----|--------|-------------|
| `` f `` | 원격을 업데이트 |  |
| `` n `` | 새로운 Remote 추가 |  |
| `` d `` | Remote를 삭제 |  |
| `` e `` | Remote를 수정 |  |
| `` / `` | Filter the current view by text |  |

## 원격 브랜치

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 브랜치명을 클립보드에 복사 |  |
| `` <space> `` | 체크아웃 |  |
| `` n `` | 새 브랜치 생성 |  |
| `` M `` | 현재 브랜치에 병합 |  |
| `` r `` | 체크아웃된 브랜치를 이 브랜치에 리베이스 |  |
| `` d `` | Delete remote tag |  |
| `` u `` | Set as upstream of checked-out branch |  |
| `` s `` | Sort order |  |
| `` g `` | View reset options |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | 커밋 보기 |  |
| `` / `` | Filter the current view by text |  |

## 커밋

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 커밋 SHA를 클립보드에 복사 |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` b `` | Bisect 옵션 보기 |  |
| `` s `` | Squash down |  |
| `` f `` | Fixup commit |  |
| `` r `` | 커밋메시지 변경 |  |
| `` R `` | 에디터에서 커밋메시지 수정 |  |
| `` d `` | 커밋 삭제 |  |
| `` e `` | 커밋을 편집 |  |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick commit (when mid-rebase) |  |
| `` F `` | Create fixup commit for this commit |  |
| `` S `` | Squash all 'fixup!' commits above selected commit (autosquash) |  |
| `` <c-j> `` | 커밋을 1개 아래로 이동 |  |
| `` <c-k> `` | 커밋을 1개 위로 이동 |  |
| `` V `` | 커밋을 붙여넣기 (cherry-pick) |  |
| `` B `` | Mark commit as base commit for rebase | Select a base commit for the next rebase; this will effectively perform a 'git rebase --onto'. |
| `` A `` | Amend commit with staged changes |  |
| `` a `` | Set/Reset commit author |  |
| `` t `` | 커밋 되돌리기 |  |
| `` T `` | Tag commit |  |
| `` <c-l> `` | 로그 메뉴 열기 |  |
| `` w `` | View worktree options |  |
| `` <space> `` | 커밋을 체크아웃 |  |
| `` y `` | 커밋 attribute 복사 |  |
| `` o `` | 브라우저에서 커밋 열기 |  |
| `` n `` | 커밋에서 새 브랜치를 만듭니다. |  |
| `` g `` | View reset options |  |
| `` C `` | 커밋을 복사 (cherry-pick) |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View selected item's files |  |
| `` / `` | 검색 시작 |  |

## 커밋 파일

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 커밋한 파일명을 클립보드에 복사 |  |
| `` c `` | Checkout file |  |
| `` d `` | Discard this commit's changes to this file |  |
| `` o `` | 파일 닫기 |  |
| `` e `` | 파일 편집 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Toggle file included in patch |  |
| `` a `` | Toggle all files included in patch |  |
| `` <enter> `` | Enter file to add selected lines to the patch (or toggle directory collapsed) |  |
| `` ` `` | 파일 트리뷰로 전환 |  |
| `` / `` | 검색 시작 |  |

## 커밋메시지

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 확인 |  |
| `` <esc> `` | 닫기 |  |

## 태그

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 체크아웃 |  |
| `` d `` | View delete options |  |
| `` P `` | 태그를 push |  |
| `` n `` | 태그를 생성 |  |
| `` g `` | View reset options |  |
| `` w `` | View worktree options |  |
| `` <enter> `` | 커밋 보기 |  |
| `` / `` | Filter the current view by text |  |

## 파일

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 파일명을 클립보드에 복사 |  |
| `` <space> `` | Staged 전환 |  |
| `` <c-b> `` | 파일을 필터하기 (Staged/unstaged) |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | 커밋 변경내용 |  |
| `` w `` | Commit changes without pre-commit hook |  |
| `` A `` | 마지맛 커밋 수정 |  |
| `` C `` | Git 편집기를 사용하여 변경 내용을 커밋합니다. |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | 파일 편집 |  |
| `` o `` | 파일 닫기 |  |
| `` i `` | Ignore file |  |
| `` r `` | 파일 새로고침 |  |
| `` s `` | 변경사항을 Stash |  |
| `` S `` | Stash 옵션 보기 |  |
| `` a `` | 모든 변경을 Staged/unstaged으로 전환 |  |
| `` <enter> `` | Stage individual hunks/lines for file, or collapse/expand for directory |  |
| `` d `` | View 'discard changes' options |  |
| `` g `` | View upstream reset options |  |
| `` D `` | View reset options |  |
| `` ` `` | 파일 트리뷰로 전환 |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Git mergetool를 열기 |  |
| `` f `` | Fetch |  |
| `` / `` | 검색 시작 |  |

## 확인 패널

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 확인 |  |
| `` <esc> `` | 닫기/취소 |  |
