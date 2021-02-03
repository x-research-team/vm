if exists("b:current_syntax")
  finish
endif


syn case match

syn keyword     vmDirective         include
syn keyword     vmDeclaration       let const


hi def link     vmDirective         Statement
hi def link     vmDeclaration       Type

" Linq Keywords
syn keyword     vmLinq              from where select group into orderby join in on equals by ascending descending

syn keyword     vmStatement         return let const spawn defer struct enum using async await service
syn keyword     vmException         try catch finally throw
syn keyword     vmConditional       if else elif unless and or case is
syn keyword     vmRepeat            do while for break continue grep map
syn keyword     vmBranch            break continue
syn keyword     vmClass             class new property get set default this parent static public private protected interface

hi def link     vmStatement         Statement
hi def link     vmClass             Statement
hi def link     vmConditional       Conditional
hi def link     vmBranch            Conditional
hi def link     vmLabel             Label
hi def link     vmRepeat            Repeat
hi def link     vmLinq              Keyword

syn match       vmDeclaration       /\<fn\>/
syn match       vmDeclaration       /^fn\>/


syn keyword vmCommentTodo contained TODO FIXME XXX NOTE
hi def link vmCommentTodo Todo

syn match comment "#.*$"    contains=@Spell,vmCommentTodo
syn match comment "\/\/.*$" contains=@Spell,vmCommentTodo

syn keyword     vmCast              int str float array


hi def link     vmCast              Type


syn keyword     vmBuiltins          len
syn keyword     vmBuiltins          println print stdin stdout stderr
syn keyword     vmBoolean           true false
syn keyword     vmNull              nil

hi def link     vmBuiltins          Keyword
hi def link     vmNull              Keyword
hi def link     vmBoolean           Boolean


" Comments; their contents
syn keyword     vmTodo              contained TODO FIXME XXX BUG
syn cluster     vmCommentGroup      contains=vmTodo
syn region      vmComment           start="#" end="$" contains=@vmCommentGroup,@Spell,@vmTodo


hi def link     vmComment           Comment
hi def link     vmTodo              Todo


" vm escapes
syn match       vmEscapeOctal       display contained "\\[0-7]\{3}"
syn match       vmEscapeC           display contained +\\[abfnrtv\\'"]+
syn match       vmEscapeX           display contained "\\x\x\{2}"
syn match       vmEscapeU           display contained "\\u\x\{4}"
syn match       vmEscapeBigU        display contained "\\U\x\{8}"
syn match       vmEscapeError       display contained +\\[^0-7xuUabfnrtv\\'"]+


hi def link     vmEscapeOctal       vmSpecialString
hi def link     vmEscapeC           vmSpecialString
hi def link     vmEscapeX           vmSpecialString
hi def link     vmEscapeU           vmSpecialString
hi def link     vmEscapeBigU        vmSpecialString
hi def link     vmSpecialString     Special
hi def link     vmEscapeError       Error
hi def link     vmException		Exception

" Strings and their contents
syn cluster     vmStringGroup       contains=vmEscapeOctal,vmEscapeC,vmEscapeX,vmEscapeU,vmEscapeBigU,vmEscapeError
syn region      vmString            start=+"+ skip=+\\\\\|\\"+ end=+"+ contains=@vmStringGroup
syn region      vmRegExString       start=+/[^/*]+me=e-1 skip=+\\\\\|\\/+ end=+/\s*$+ end=+/\s*[;.,)\]}]+me=e-1 oneline
syn region      vmRawString         start=+`+ end=+`+

hi def link     vmString            String
hi def link     vmRawString         String
hi def link     vmRegExString       String

" Characters; their contents
syn cluster     vmCharacterGroup    contains=vmEscapeOctal,vmEscapeC,vmEscapeX,vmEscapeU,vmEscapeBigU
syn region      vmCharacter         start=+'+ skip=+\\\\\|\\'+ end=+'+ contains=@vmCharacterGroup


hi def link     vmCharacter         Character


" Regions
syn region      vmBlock             start="{" end="}" transparent fold
syn region      vmParen             start='(' end=')' transparent


" Integers
syn match       vmDecimalInt        "\<\d\+\([Ee]\d\+\)\?\>"
syn match       vmHexadecimalInt    "\<0x\x\+\>"
syn match       vmOctalInt          "\<0\o\+\>"
syn match       vmOctalError        "\<0\o*[89]\d*\>"


hi def link     vmDecimalInt        Integer
hi def link     vmHexadecimalInt    Integer
hi def link     vmOctalInt          Integer
hi def link     Integer                 Number

" Floating point
syn match       vmFloat             "\<\d\+\.\d*\([Ee][-+]\d\+\)\?\>"
syn match       vmFloat             "\<\.\d\+\([Ee][-+]\d\+\)\?\>"
syn match       vmFloat             "\<\d\+[Ee][-+]\d\+\>"


hi def link     vmFloat             Float
"hi def link     vmImaginary         Number


if exists("vm_fold")
    syn match	vmFunction	"\<fn\>"
    syn region	vmFunctionFold	start="\<fn\>.*[^};]$" end="^\z1}.*$" transparent fold keepend

    syn sync match vmSync	grouphere vmFunctionFold "\<fn\>"
    syn sync match vmSync	grouphere NONE "^}"

    setlocal foldmethod=syntax
    setlocal foldtext=getline(v:foldstart)
else
    syn keyword vmFunction	fn
    syn match	vmBraces	"[{}\[\]]"
    syn match	vmParens	"[()]"
endif

syn sync fromstart
syn sync maxlines=100

hi def link vmFunction		Function
hi def link vmBraces		Function

let b:current_syntax = "vm"
