_This file is auto-generated. To update, make the changes in the pkg/i18n directory and then run `go generate ./...` from the project root._

# Lazygit 키 바인딩

_Legend: `<c-b>` means ctrl+b, `<a-b>` means alt+b, `B` means shift+b_

## 글로벌 키 바인딩

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-r> `` | 최근에 사용한 저장소로 전환 |  |
| `` <pgup> (fn+up/shift+k) `` | 메인 패널을 위로 스크롤 |  |
| `` <pgdown> (fn+down/shift+j) `` | 메인 패널을 아래로로 스크롤 |  |
| `` @ `` | 명령어 로그 메뉴 열기 | View options for the command log e.g. show/hide the command log and focus the command log. |
| `` P `` | 푸시 | Push the current branch to its upstream branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` p `` | 업데이트 | Pull changes from the remote for the current branch. If no upstream is configured, you will be prompted to configure an upstream branch. |
| `` ) `` | Increase rename similarity threshold | Increase the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` ( `` | Decrease rename similarity threshold | Decrease the similarity threshold for a deletion and addition pair to be treated as a rename. |
| `` } `` | Diff 보기의 변경 사항 주위에 표시되는 컨텍스트의 크기를 늘리기 | Increase the amount of the context shown around changes in the diff view. |
| `` { `` | Diff 보기의 변경 사항 주위에 표시되는 컨텍스트 크기 줄이기 | Decrease the amount of the context shown around changes in the diff view. |
| `` : `` | Execute custom command | Bring up a prompt where you can enter a shell command to execute. Not to be confused with pre-configured custom commands. |
| `` <c-p> `` | 커스텀 Patch 옵션 보기 |  |
| `` m `` | View merge/rebase options | View options to abort/continue/skip the current merge/rebase. |
| `` R `` | 새로고침 | Refresh the git state (i.e. run `git status`, `git branch`, etc in background to update the contents of panels). This does not run `git fetch`. |
| `` + `` | 다음 스크린 모드 (normal/half/fullscreen) |  |
| `` _ `` | 이전 스크린 모드 |  |
| `` ? `` | 매뉴 열기 |  |
| `` <c-s> `` | View filter-by-path options | View options for filtering the commit log, so that only commits matching the filter are shown. |
| `` W `` | Diff 메뉴 열기 | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` <c-e> `` | Diff 메뉴 열기 | View options relating to diffing two refs e.g. diffing against selected ref, entering ref to diff against, and reversing the diff direction. |
| `` q `` | 종료 |  |
| `` <esc> `` | 취소 |  |
| `` <c-w> `` | 공백문자를 Diff 뷰에서 표시 여부 전환 | Toggle whether or not whitespace changes are shown in the diff view. |
| `` z `` | 되돌리기 (reflog) (실험적) | The reflog will be used to determine what git command to run to undo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |
| `` <c-z> `` | 다시 실행 (reflog) (실험적) | The reflog will be used to determine what git command to run to redo the last git command. This does not include changes to the working tree; only commits are taken into consideration. |

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
| `` <c-o> `` | 커밋 해시를 클립보드에 복사 |  |
| `` <space> `` | 체크아웃 | Checkout the selected commit as a detached HEAD. |
| `` y `` | 커밋 attribute 복사 | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 브라우저에서 커밋 열기 |  |
| `` n `` | 커밋에서 새 브랜치를 만듭니다. |  |
| `` g `` | View reset options | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | 커밋을 복사 (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 커밋 보기 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Stash

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 적용 | Apply the stash entry to your working directory. |
| `` g `` | Pop | Apply the stash entry to your working directory and remove the stash entry. |
| `` d `` | Drop | Remove the stash entry from the stash list. |
| `` n `` | 새 브랜치 생성 | Create a new branch from the selected stash entry. This works by git checking out the commit that the stash entry was created from, creating a new branch from that commit, then applying the stash entry to the new branch as an additional commit. |
| `` r `` | Rename stash |  |
| `` <enter> `` | View selected item's files |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## Sub-commits

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 커밋 해시를 클립보드에 복사 |  |
| `` <space> `` | 체크아웃 | Checkout the selected commit as a detached HEAD. |
| `` y `` | 커밋 attribute 복사 | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 브라우저에서 커밋 열기 |  |
| `` n `` | 커밋에서 새 브랜치를 만듭니다. |  |
| `` g `` | View reset options | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | 커밋을 복사 (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View selected item's files |  |
| `` w `` | View worktree options |  |
| `` / `` | 검색 시작 |  |

## Worktrees

| Key | Action | Info |
|-----|--------|-------------|
| `` n `` | New worktree |  |
| `` <space> `` | Switch | Switch to the selected worktree. |
| `` o `` | Open in editor |  |
| `` d `` | Remove | Remove the selected worktree. This will both delete the worktree's directory, as well as metadata about the worktree in the .git directory. |
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
| `` <space> `` | Pick hunk |  |
| `` b `` | Pick all hunks |  |
| `` <up> `` | 이전 hunk를 선택 |  |
| `` <down> `` | 다음 hunk를 선택 |  |
| `` <left> `` | 이전 충돌을 선택 |  |
| `` <right> `` | 다음 충돌을 선택 |  |
| `` z `` | 되돌리기 | Undo last merge conflict resolution. |
| `` e `` | 파일 편집 | Open file in external editor. |
| `` o `` | 파일 닫기 | Open file in default application. |
| `` M `` | Git mergetool를 열기 | Run `git mergetool`. |
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
| `` a `` | Toggle select hunk | Toggle hunk selection mode. |
| `` <c-o> `` | 선택한 텍스트를 클립보드에 복사 |  |
| `` o `` | 파일 닫기 | Open file in default application. |
| `` e `` | 파일 편집 | Open file in external editor. |
| `` <space> `` | Line(s)을 패치에 추가/삭제 |  |
| `` <esc> `` | Exit custom patch builder |  |
| `` / `` | 검색 시작 |  |

## 메인 패널 (Staging)

| Key | Action | Info |
|-----|--------|-------------|
| `` <left> `` | 이전 hunk를 선택 |  |
| `` <right> `` | 다음 hunk를 선택 |  |
| `` v `` | 드래그 선택 전환 |  |
| `` a `` | Toggle select hunk | Toggle hunk selection mode. |
| `` <c-o> `` | 선택한 텍스트를 클립보드에 복사 |  |
| `` <space> `` | Staged 전환 | 선택한 행을 staged / unstaged |
| `` d `` | 변경을 삭제 (git reset) | When unstaged change is selected, discard the change using `git reset`. When staged change is selected, unstage the change. |
| `` o `` | 파일 닫기 | Open file in default application. |
| `` e `` | 파일 편집 | Open file in external editor. |
| `` <esc> `` | 파일 목록으로 돌아가기 |  |
| `` <tab> `` | 패널 전환 | Switch to other view (staged/unstaged changes). |
| `` E `` | Edit hunk | Edit selected hunk in external editor. |
| `` c `` | 커밋 변경내용 | Commit staged changes. |
| `` w `` | Commit changes without pre-commit hook |  |
| `` C `` | Git 편집기를 사용하여 변경 내용을 커밋합니다. |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` / `` | 검색 시작 |  |

## 브랜치

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 브랜치명을 클립보드에 복사 |  |
| `` i `` | Git-flow 옵션 보기 |  |
| `` <space> `` | 체크아웃 | Checkout selected item. |
| `` n `` | 새 브랜치 생성 |  |
| `` o `` | 풀 리퀘스트 생성 |  |
| `` O `` | 풀 리퀘스트 생성 옵션 |  |
| `` <c-y> `` | 풀 리퀘스트 URL을 클립보드에 복사 |  |
| `` c `` | 이름으로 체크아웃 | Checkout by name. In the input box you can enter '-' to switch to the last branch. |
| `` F `` | 강제 체크아웃 | Force checkout selected branch. This will discard all local changes in your working directory before checking out the selected branch. |
| `` d `` | Delete | View delete options for local/remote branch. |
| `` r `` | 체크아웃된 브랜치를 이 브랜치에 리베이스 | Rebase the checked-out branch onto the selected branch. |
| `` M `` | 현재 브랜치에 병합 | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` f `` | Fast-forward this branch from its upstream | Fast-forward selected branch from its upstream. |
| `` T `` | 태그를 생성 |  |
| `` s `` | Sort order |  |
| `` g `` | View reset options |  |
| `` R `` | 브랜치 이름 변경 |  |
| `` u `` | View upstream options | View options relating to the branch's upstream e.g. setting/unsetting the upstream and resetting to the upstream. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 커밋 보기 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## 상태

| Key | Action | Info |
|-----|--------|-------------|
| `` o `` | 설정 파일 열기 | Open file in default application. |
| `` e `` | 설정 파일 수정 | Open file in external editor. |
| `` u `` | 업데이트 확인 |  |
| `` <enter> `` | 최근에 사용한 저장소로 전환 |  |
| `` a `` | 모든 브랜치 로그 표시 |  |

## 서브모듈

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 서브모듈 이름을 클립보드에 복사 |  |
| `` <enter> `` | Enter | 서브모듈 열기 |
| `` d `` | Remove | Remove the selected submodule and its corresponding directory. |
| `` u `` | Update | 서브모듈 업데이트 |
| `` n `` | 새로운 서브모듈 추가 |  |
| `` e `` | 서브모듈의 URL을 수정 |  |
| `` i `` | Initialize | 서브모듈 초기화 |
| `` b `` | View bulk submodule options |  |
| `` / `` | Filter the current view by text |  |

## 원격

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | View branches |  |
| `` n `` | 새로운 Remote 추가 |  |
| `` d `` | Remove | Remove the selected remote. Any local branches tracking a remote branch from the remote will be unaffected. |
| `` e `` | Edit | Remote를 수정 |
| `` f `` | Fetch | 원격을 업데이트 |
| `` / `` | Filter the current view by text |  |

## 원격 브랜치

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 브랜치명을 클립보드에 복사 |  |
| `` <space> `` | 체크아웃 | Checkout a new local branch based on the selected remote branch, or the remote branch as a detached head. |
| `` n `` | 새 브랜치 생성 |  |
| `` M `` | 현재 브랜치에 병합 | View options for merging the selected item into the current branch (regular merge, squash merge) |
| `` r `` | 체크아웃된 브랜치를 이 브랜치에 리베이스 | Rebase the checked-out branch onto the selected branch. |
| `` d `` | Delete | Delete the remote branch from the remote. |
| `` u `` | Set as upstream | Set the selected remote branch as the upstream of the checked-out branch. |
| `` s `` | Sort order |  |
| `` g `` | View reset options | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 커밋 보기 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## 커밋

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 커밋 해시를 클립보드에 복사 |  |
| `` <c-r> `` | Reset cherry-picked (copied) commits selection |  |
| `` b `` | Bisect 옵션 보기 |  |
| `` s `` | Squash | Squash the selected commit into the commit below it. The selected commit's message will be appended to the commit below it. |
| `` f `` | Fixup | Meld the selected commit into the commit below it. Similar to squash, but the selected commit's message will be discarded. |
| `` r `` | 커밋메시지 변경 | Reword the selected commit's message. |
| `` R `` | 에디터에서 커밋메시지 수정 |  |
| `` d `` | 커밋 삭제 | Drop the selected commit. This will remove the commit from the branch via a rebase. If the commit makes changes that later commits depend on, you may need to resolve merge conflicts. |
| `` e `` | Edit (start interactive rebase) | 커밋을 편집 |
| `` i `` | Start interactive rebase | Start an interactive rebase for the commits on your branch. This will include all commits from the HEAD commit down to the first merge commit or main branch commit.
If you would instead like to start an interactive rebase from the selected commit, press `e`. |
| `` p `` | Pick | Pick commit (when mid-rebase) |
| `` F `` | Create fixup commit | Create fixup commit for this commit |
| `` S `` | Apply fixup commits | Squash all 'fixup!' commits above selected commit (autosquash) |
| `` <c-j> `` | 커밋을 1개 아래로 이동 |  |
| `` <c-k> `` | 커밋을 1개 위로 이동 |  |
| `` V `` | 커밋을 붙여넣기 (cherry-pick) |  |
| `` B `` | Mark as base commit for rebase | Select a base commit for the next rebase. When you rebase onto a branch, only commits above the base commit will be brought across. This uses the `git rebase --onto` command. |
| `` A `` | Amend | Amend commit with staged changes |
| `` a `` | Amend commit attribute | Set/Reset commit author or set co-author. |
| `` t `` | Revert | Create a revert commit for the selected commit, which applies the selected commit's changes in reverse. |
| `` T `` | Tag commit | Create a new tag pointing at the selected commit. You'll be prompted to enter a tag name and optional description. |
| `` <c-l> `` | 로그 메뉴 열기 | View options for commit log e.g. changing sort order, hiding the git graph, showing the whole git graph. |
| `` <space> `` | 체크아웃 | Checkout the selected commit as a detached HEAD. |
| `` y `` | 커밋 attribute 복사 | Copy commit attribute to clipboard (e.g. hash, URL, diff, message, author). |
| `` o `` | 브라우저에서 커밋 열기 |  |
| `` n `` | 커밋에서 새 브랜치를 만듭니다. |  |
| `` g `` | View reset options | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` C `` | 커밋을 복사 (cherry-pick) | Mark commit as copied. Then, within the local commits view, you can press `V` to paste (cherry-pick) the copied commit(s) into your checked out branch. At any time you can press `<esc>` to cancel the selection. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | View selected item's files |  |
| `` w `` | View worktree options |  |
| `` / `` | 검색 시작 |  |

## 커밋 파일

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 파일명을 클립보드에 복사 |  |
| `` c `` | 체크아웃 | Checkout file |
| `` d `` | Remove | Discard this commit's changes to this file |
| `` o `` | 파일 닫기 | Open file in default application. |
| `` e `` | Edit | Open file in external editor. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <space> `` | Toggle file included in patch | Toggle whether the file is included in the custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` a `` | Toggle all files included in patch | Add/remove all commit's files to custom patch. See https://github.com/jesseduffield/lazygit#rebase-magic-custom-patches. |
| `` <enter> `` | Enter file to add selected lines to the patch (or toggle directory collapsed) | If a file is selected, enter the file so that you can add/remove individual lines to the custom patch. If a directory is selected, toggle the directory. |
| `` ` `` | 파일 트리뷰로 전환 | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` / `` | 검색 시작 |  |

## 커밋메시지

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 확인 |  |
| `` <esc> `` | 닫기 |  |

## 태그

| Key | Action | Info |
|-----|--------|-------------|
| `` <space> `` | 체크아웃 | Checkout the selected tag tag as a detached HEAD. |
| `` n `` | 태그를 생성 | Create new tag from current commit. You'll be prompted to enter a tag name and optional description. |
| `` d `` | Delete | View delete options for local/remote tag. |
| `` P `` | 태그를 push | Push the selected tag to a remote. You'll be prompted to select a remote. |
| `` g `` | Reset | View reset options (soft/mixed/hard) for resetting onto selected item. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` <enter> `` | 커밋 보기 |  |
| `` w `` | View worktree options |  |
| `` / `` | Filter the current view by text |  |

## 파일

| Key | Action | Info |
|-----|--------|-------------|
| `` <c-o> `` | 파일명을 클립보드에 복사 |  |
| `` <space> `` | Staged 전환 | Toggle staged for selected file. |
| `` <c-b> `` | 파일을 필터하기 (Staged/unstaged) |  |
| `` y `` | Copy to clipboard |  |
| `` c `` | 커밋 변경내용 | Commit staged changes. |
| `` w `` | Commit changes without pre-commit hook |  |
| `` A `` | 마지맛 커밋 수정 |  |
| `` C `` | Git 편집기를 사용하여 변경 내용을 커밋합니다. |  |
| `` <c-f> `` | Find base commit for fixup | Find the commit that your current changes are building upon, for the sake of amending/fixing up the commit. This spares you from having to look through your branch's commits one-by-one to see which commit should be amended/fixed up. See docs: <https://github.com/jesseduffield/lazygit/tree/master/docs/Fixup_Commits.md> |
| `` e `` | Edit | Open file in external editor. |
| `` o `` | 파일 닫기 | Open file in default application. |
| `` i `` | Ignore file |  |
| `` r `` | 파일 새로고침 |  |
| `` s `` | Stash | Stash all changes. For other variations of stashing, use the view stash options keybinding. |
| `` S `` | Stash 옵션 보기 | View stash options (e.g. stash all, stash staged, stash unstaged). |
| `` a `` | 모든 변경을 Staged/unstaged으로 전환 | Toggle staged/unstaged for all files in working tree. |
| `` <enter> `` | Stage individual hunks/lines for file, or collapse/expand for directory | If the selected item is a file, focus the staging view so you can stage individual hunks/lines. If the selected item is a directory, collapse/expand it. |
| `` d `` | View 'discard changes' options | View options for discarding changes to the selected file. |
| `` g `` | View upstream reset options |  |
| `` D `` | Reset | View reset options for working tree (e.g. nuking the working tree). |
| `` ` `` | 파일 트리뷰로 전환 | Toggle file view between flat and tree layout. Flat layout shows all file paths in a single list, tree layout groups files by directory. |
| `` <c-t> `` | Open external diff tool (git difftool) |  |
| `` M `` | Git mergetool를 열기 | Run `git mergetool`. |
| `` f `` | Fetch | Fetch changes from remote. |
| `` / `` | 검색 시작 |  |

## 확인 패널

| Key | Action | Info |
|-----|--------|-------------|
| `` <enter> `` | 확인 |  |
| `` <esc> `` | 닫기/취소 |  |
