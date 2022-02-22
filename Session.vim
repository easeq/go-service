let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ~/dev/go-service
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
set shortmess=aoO
argglobal
%argdel
$argadd ~/dev/go-service
edit .gitignore
let s:save_splitbelow = &splitbelow
let s:save_splitright = &splitright
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
let &splitbelow = s:save_splitbelow
let &splitright = s:save_splitright
wincmd t
let s:save_winminheight = &winminheight
let s:save_winminwidth = &winminwidth
set winminheight=0
set winheight=1
set winminwidth=0
set winwidth=1
exe 'vert 1resize ' . ((&columns * 105 + 105) / 211)
exe 'vert 2resize ' . ((&columns * 105 + 105) / 211)
argglobal
balt server/grpc/error/error.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 1 - ((0 * winheight(0) + 51) / 103)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 1
normal! 0
wincmd w
argglobal
if bufexists("server/grpc/error/error_test.go") | buffer server/grpc/error/error_test.go | else | edit server/grpc/error/error_test.go | endif
if &buftype ==# 'terminal'
  silent file server/grpc/error/error_test.go
endif
balt ~/go/pkg/mod/github.com/stretchr/testify@v1.7.0/require/require.go
setlocal fdm=manual
setlocal fde=0
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=0
setlocal fml=1
setlocal fdn=20
setlocal fen
silent! normal! zE
let &fdl = &fdl
let s:l = 154 - ((86 * winheight(0) + 51) / 103)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 154
normal! 0
wincmd w
exe 'vert 1resize ' . ((&columns * 105 + 105) / 211)
exe 'vert 2resize ' . ((&columns * 105 + 105) / 211)
tabnext 1
badd +39 server/grpc/error/error.go
badd +0 ~/dev/go-service
badd +44 ~/go/pkg/mod/google.golang.org/grpc@v1.42.0/status/status.go
badd +150 server/grpc/error/error_test.go
badd +0 ~/go/pkg/mod/github.com/stretchr/testify@v1.7.0/require/require_forward.go
badd +38 server/grpc/config.go
badd +8 server/grpc/error/error.proto
badd +105 ~/go/pkg/mod/google.golang.org/grpc@v1.42.0/internal/status/status.go
badd +51 ~/go/pkg/mod/google.golang.org/genproto@v0.0.0-20200526211855-cb27e3aa2013/googleapis/rpc/status/status.pb.go
badd +134 /snap/go/9028/src/builtin/builtin.go
badd +389 ~/go/pkg/mod/github.com/stretchr/testify@v1.7.0/require/require.go
badd +0 .gitignore
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20 shortmess=filnxtToOF
let &winminheight = s:save_winminheight
let &winminwidth = s:save_winminwidth
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
set hlsearch
let g:this_session = v:this_session
let g:this_obsession = v:this_session
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
