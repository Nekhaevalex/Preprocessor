# TD4 Preprocessor/Assembler syntax description

## Assembler opcodes
| #  | Binary | Instruction | Argument 1 | Argument 2 | Description                                                                             | FastAdd support | Commentary                           |
|----|--------|-------------|------------|------------|-----------------------------------------------------------------------------------------|-----------------|--------------------------------------|
| 0  | 0000   | add         | A          | Im         | A=A+Im                                                                                  | FALSE           |                                      |
| 1  | 0001   | mov         | A          | B          | A=B + FastAdd                                                                           | TRUE            |                                      |
| 2  | 0010   | in          | A          | -          | Takes value into A, sharing B as address                                                | TRUE            | Uses B as address only on TD4E/TD4E8 |
| 3  | 0011   | mov         | A          | Im         | A=Im                                                                                    | FALSE           |                                      |
| 4  | 0100   | mov         | B          | A          | B=A + FastAdd                                                                           | TRUE            |                                      |
| 5  | 0101   | add         | B          | Im         | B=B + Im                                                                                | FALSE           |                                      |
| 6  | 0110   | in          | B          | -          | Takes value into B, sharing A as address                                                | TRUE            | Uses A as address only on TD4E/TD4E8 |
| 7  | 0111   | mov         | B          | Im         | B=Im                                                                                    | FALSE           |                                      |
| 8  | 1000   | cmp         | A          | B          | Compare A with B and set C = 1 if (im = 0 => A = B, im = 1   => A > B, im = 2 => A < B) | TRUE            | Available only on TD4E/TD4E8         |
| 9  | 1001   | out         | B          | -          | Gives value from B, sharing    A as address                                             | TRUE            | Uses A as address only on TD4E/TD4E8 |
| 10 | 1010   | mov         | B          | PC         | B=PC + FastAdd                                                                          | TRUE            | Available only on TD4E/TD4E8         |
| 11 | 1011   | out         | Im         | -          | Gives Im value, sharing    A as address                                                 | FALSE           | Uses A as address only on TD4E/TD4E8 |
| 12 | 1100   | jnc         | B          | -          | PC=B if C != 1                                                                          | TRUE            | Available only on TD4E/TD4E8         |
| 13 | 1101   | jmp         | B          | -          | PC=B                                                                                    | FALSE           | Available only on TD4E/TD4E8         |
| 14 | 1110   | jnc         | Im         | -          | PC=Im if C!=1                                                                           | TRUE            | Supposed using label instead of Im   |
| 15 | 1111   | jmp         | Im         | -          | PC=Im                                                                                   | FALSE           | Supposed using label instead of Im   |

## Preprocessor commands
| Directive | Example               |
|-----------|-----------------------|
| #define   | #define a 5           |
| #else     |                       |
| #endif    |                       |
| #endmacro |                       |
| #error    | #error "Error at"     |
| #ifdef    | #ifdef a              |
| #ifndef   | #ifndef a             |
| #import   | #import "file.h"      |
| #line     | #line 10 asm.s        |
| #macro    | #macro test a, b      |
| #pext     | #pext io 12           |
| #resdef   | #resdef a b           |
| #return   | #return a             |
| #sumdef   | #sumdef a b           |
| #undef    | #undef a              |
| #warn     | #warn "Hello, World!" |

## Tree structure
```bash
program
+--section
   +--section_name                String
   +--section_content             Block
      +--define_directive         Directive
      |  +--name:                 Ident
      |  +--definition:           Ident, nullable
      +--import_directive         Directive
      |  +--name:                 Ident
      +--line_directive           Directive
      |  +--name:                 Ident
      |  +--line_number:          Ident
      +--warn_directive           Directive
      |  +--message:              Ident
      +--sumdef_directive         Directive
      |  +--def_1                 Ident
      |  +--def_2                 Ident
      +--resdef_directive         Directive
      |  +--def_1                 Ident
      |  +--def_2                 Ident
      +--pext_directive           Directive
      |  +--pext_name             Ident
      |  +--pext_address          Ident
      +--error_directive          Directive
      |  +--message               Ident
      +--undef_directive          Directive
      |  +--definition            Ident
      +--ifdef_directive          Directive
      |  +--definition            Ident
      |  +--body_true             Block
      |  |  +--...
      |  +--body_false            Block, nullable
      |     +--...
      +--ifndef_directive         Directive
      |  +--definition            Ident
      |  +--body_true             Block
      |  |  +--...
      |  +--body_false            Block, nullable
      |     +--...
      +--macro_directive          Directive
      |  +--macro_name            String
      |  +--args                  [Ident]
      |  +--body                  Block
      |     +--...
      |     +--return_directive   Directive
      |        +--definition      Ident
      +--add_opcode               Opcode
      |  +--reg                   Reg
      |  +--value                 Ident
      +--mov_opcode               Opcode
      |  +--reg1                  Reg
      |  +--reg2                  Reg, nullable
      |  +--fa                    Ident
      +--in_opcode                Opcode
      |  +--reg                   Reg
      +--out_opcode               Opcode
      |  +--reg                   Reg
      |  +--fa                    Ident
      +--cmp_opcode               Opcode
      |  +--reg_a                 Reg
      |  +--reg_b                 Reg
      |  +--operation             Ident
      +--jmp_opcode               Opcode
      |  +--reg_b                 Reg
      |  +--addr                  Ident
      +--jnc_opcode               Opcode
      |  +--reg_b                 Reg
      |  +--addr                  Ident
      +--macro_call               MacroCall
      |  +--macro_name            String
      |  +--args                  [Ident|MacroCall]
      +--comment                  Comment
```