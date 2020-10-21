# Tree structure
```
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