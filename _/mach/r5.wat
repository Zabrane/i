;; risc-v virtual machine in webassembly
;; - instantiate module
;; - grow exported memory as needed
;; - copy risc-v program to memory at 384 (first 32*4 bytes are i32 registers, then 32*8 f64 registers)
;; - call "jmp" with pc argument (linear memory addr)
(module
 (type (func (param i32 i32) (result i32)))
 (type (func (param f64 f64) (result f64)))
 (type (func (param i32 i32) (result f64)))
 (type (func (param i32 i32)))
 (type (func (param i32 f64)))
 (type (func (param f64 f64) (result i32)))
 (type (func (param f64) (result f64)))
 
 (import "ext" "getc" (func $getc (result i32)))
 (import "ext" "putc" (func $putc (param i32)))
 (import "ext" "draw" (func $draw (param i32 i32 i32)))
 (import "ext" "sin"  (func $sin  (param f64) (result f64)))
 (import "ext" "cos"  (func $cos  (param f64) (result f64)))
 (import "ext" "exp"  (func $exp  (param f64) (result f64)))
 (import "ext" "log"  (func $log  (param f64) (result f64)))
 (import "ext" "atan2" (func $atan2 (param f64 f64) (result f64)))
 (import "ext" "hypot" (func $hypot (param f64 f64) (result f64)))
 
 (func $jmp (param $pc i32) 
  (local $i  i32) (local $op i32) (local $rd i32) (local $r1 i32) (local $r2 i32) 
  (local $f3 i32) (local $f7 i32) (local $x  i32) (local $y  i32) (local $z  i32) (local $im i32) (local $is i32) (local $ib i32)
  
  ;; fetch
  (local.set $i (i32.load (local.get $pc)))
  
  ;; decode
  (local.set $op (i32.and (local.get $i) (i32.const 127)))
  (local.set $rd (i32.and (i32.shr_u (local.get $i) (i32.const  7)) (i32.const  31)))
  (local.set $r1 (i32.and (i32.shr_u (local.get $i) (i32.const 15)) (i32.const  31)))
  (local.set $r2 (i32.and (i32.shr_u (local.get $i) (i32.const 20)) (i32.const  31)))
  (local.set $f3 (i32.and (i32.shr_u (local.get $i) (i32.const 12)) (i32.const   7)))
  (local.set $f7 (i32.and (i32.shr_u (local.get $i) (i32.const 25)) (i32.const 127)))
  (local.set $im (i32.shr_s (local.get $i) (i32.const 20)))
  (local.set $is (i32.or (local.get $rd) (i32.shr_s (local.get $i) (i32.const 25))))
  (local.set $ib (i32.or (i32.and   (i32.const 30) (i32.shr_u (local.get $i)  (i32.const 7)))
                 (i32.or (i32.shl   (i32.and (i32.const 128)  (local.get $i)) (i32.const 4)) 
                 (i32.or (i32.shr_u (i32.and (i32.const 0x7e000000) (local.get $i)) (i32.const 20)) 
                         (i32.shr_s (i32.and (i32.const 0x80000000) (local.get $i)) (i32.const 19))))))
  
  ;; execute
  (if       (i32.eq (local.get $op) (i32.const   3)) (then (i32.store (local.get $rd) (call_indirect (type 0) (local.get $r1) (local.get $im) (local.get $f3)))) ;; ld 0..7
  (else (if (i32.eq (local.get $op) (i32.const   7)) (then (f64.store offset=128 (local.get $rd) (call $fld (local.get $r1) (local.get $im))))                   ;; fld
  (else (if (i32.eq (local.get $op) (i32.const  19)) (then (i32.store (local.get $rd) 
	    (if (result i32) (i32.and (i32.eq (local.get $f3) (i32.const 1)) (i32.eq (local.get $f7) (i32.const 48))) (then (i32.clz (local.get $r1))) ;; clz
	    (else (call_indirect (type 0) (i32.load (local.get $r1)) (local.get $im) (i32.add (local.get $f3) (i32.const 8)))))))                                ;; addi 8..15
  (else (if (i32.eq (local.get $op) (i32.const  35)) (then (call_indirect (type 3) (i32.add (local.get $r1) (local.get $is)) (local.get $r2) (i32.add (local.get $f3) (i32.const 15)))) ;; sb 15..18
  (else (if (i32.eq (local.get $op) (i32.const  39)) (then (call $fsd (i32.add (local.get $r1) (local.get $is)) (f64.load offset=128 (local.get $rd))))          ;; fsd
  (else (if (i32.eq (local.get $op) (i32.const  51)) (then (i32.store (local.get $rd) (call_indirect (type 0) (i32.load (local.get $r1)) (i32.load (local.get $r2))
                                                                            (i32.add (local.get $f3)
                                                                            (if (result i32) (i32.eqz (local.get $f7)) (then (i32.const 8))                      ;; add 8..15
	                                                                    (else (if (result i32) (i32.eq (local.get $f7) (i32.const 1)) (then (i32.const 19))  ;; mul 19..26
	                                                                    (else (i32.const 27)))))))))                                                         ;; sub 27..34
  (else (if (i32.eq (local.get $op) (i32.const  55)) (then (i32.store (local.get $rd) (i32.and (local.get $i) (i32.const 0xfffff000))))                        ;; lui
  (else (if (i32.eq (local.get $op) (i32.const  99)) (then (if (call_indirect (type 0) (local.get $r1) (local.get $r2) (i32.add (local.get $f3) (i32.const 35)))
              (then (local.set $pc (i32.sub (i32.add (local.get $pc) (local.get $ib)) (i32.const 4))))))                                                         ;; beq 35..42
	      
	      
  (else (if (i32.eq (local.get $op) (i32.const 103)) (then (i32.store (local.get $rd) (i32.add (local.get $pc) (i32.const 4)))
                                                           (local.set $pc (i32.add (i32.load (local.get $r1)) (local.get $im))))                                 ;; jalr
  (else (if (i32.eq (local.get $op) (i32.const 111)) (then (i32.store (local.get $rd) (i32.add (local.get $pc) (i32.const 4)))
                                                           (local.set $pc (i32.add (local.get $pc) (local.get $im))))                                            ;; jal
  
  (else (if (i32.eq (local.get $op) (i32.const  83)) (then
            (if (i32.eq (local.get $f7) (i32.const 81)) (then                                                     
                (i32.store (local.get $rd) (call_indirect (type 5) (f64.load offset=128 (local.get $r1)) (f64.load offset=128 (local.get $r2)) (i32.add (i32.const 47) (local.get $f3)))))
            (else (if (i32.eq (local.get $f7) (i32.const  97)) (then (i32.store (local.get $rd) (i32.trunc_f64_s (f64.load offset=128 (local.get $r1)))))       ;; fcvtwd
	    (else (if (i32.eq (local.get $f7) (i32.const 105)) (then (f64.store offset=128 (local.get $rd) (f64.convert_i32_s (i32.load (local.get $r1)))))     ;; fcvtdw
	    (else (if (i32.eq (local.get $f7) (i32.const  45)) (then (f64.store offset=128 (local.get $rd) (f64.sqrt (f64.load offset=128 (local.get $r1)))))
	    (else (f64.store offset=128 (local.get $rd) (call_indirect (type 1) (f64.load offset=128 (local.get $r1)) (f64.load offset=128 (local.get $r2))
	              (if (result i32) (i32.eq (local.get $f7) (i32.const 81)) (then (i32.add (i32.const 47) (local.get $f3)))                                  ;; fle  47..50
		      (else (i32.add (i32.const 43) (i32.and (i32.const 3) (i32.shr_u (local.get $i) (i32.const 27)))))))))))))))))                             ;; fadd 43..46

  (else (if (i32.eq (local.get $op) (i32.const 127))      (then
        (if       (i32.eqz (local.get $f3))               (then (i32.store (local.get $rd) (call $getc)))
        (else (if (i32.eq  (local.get $f3) (i32.const 1)) (then (call $putc (i32.load (local.get $r1))))
        (else (if (i32.eq  (local.get $f3) (i32.const 2)) (then (call $draw (i32.load (local.get $rd))(i32.load (local.get $r1))(i32.load (local.get $r2))))
        (else (if (i32.lt_u(local.get $f3) (i32.const 7)) (then (f64.store offset=128 (local.get $rd) (call_indirect (type 6) 
                                                       (f64.load offset=128 (local.get $r1)) (i32.add (i32.const 51) (local.get $f7)))))                           ;; sin 51..54
        (else (f64.store offset=128 (local.get $rd) (call_indirect (type 1)
	         (f64.load offset=128 (local.get $r1)) (f64.load offset=128 (local.get $r2)) (i32.add (i32.const 51 (local.get $f7))))))))))))))                   ;; atan2 55..56
       
  (else (unreachable)))))))))))))))))))))))))
  
  (local.set $pc (i32.add (local.get $pc) (i32.const 4)))
 )
 (func $xxx  (type 0) unreachable)
 (func $lb   (type 0) (i32.load8_s offset=384 (i32.add (local.get 0) (local.get 1))))
 (func $lw   (type 0) (i32.load    offset=384 (i32.add (local.get 0) (local.get 1))))
 (func $lbu  (type 0) (i32.load8_u offset=384 (i32.add (local.get 0) (local.get 1))))
 (func $sb   (type 3) (i32.store8  offset=384 (local.get 0) (local.get 1)))
 (func $sw   (type 3) (i32.store   offset=384 (local.get 0) (local.get 1)))
 (func $fsd  (type 4) (f64.store   offset=384 (local.get 0) (local.get 1)))
 (func $add  (type 0) (i32.add   (local.get 0) (local.get 1)))
 (func $sub  (type 0) (i32.sub   (local.get 0) (local.get 1)))
 (func $fadd (type 1) (f64.add   (local.get 0) (local.get 1)))
 (func $fsub (type 1) (f64.sub   (local.get 0) (local.get 1)))
 (func $fmul (type 1) (f64.mul   (local.get 0) (local.get 1)))
 (func $fdiv (type 1) (f64.div   (local.get 0) (local.get 1)))
 (func $shl  (type 0) (i32.shl   (local.get 0) (local.get 1)))
 (func $shr  (type 0) (i32.shr_u (local.get 0) (local.get 1)))
 (func $sra  (type 0) (i32.shr_s (local.get 0) (local.get 1)))
 (func $xor  (type 0) (i32.xor   (local.get 0) (local.get 1)))
 (func $or   (type 0) (i32.or    (local.get 0) (local.get 1)))
 (func $and  (type 0) (i32.and   (local.get 0) (local.get 1)))
 (func $mul  (type 0) (i32.mul   (local.get 0) (local.get 1)))
 (func $div  (type 0) (i32.div_s (local.get 0) (local.get 1)))
 (func $divu (type 0) (i32.div_u (local.get 0) (local.get 1)))
 (func $rem  (type 0) (i32.rem_s (local.get 0) (local.get 1)))
 (func $remu (type 0) (i32.rem_u (local.get 0) (local.get 1)))
 (func $beq  (type 0) (i32.eq    (local.get 0) (local.get 1)))
 (func $bne  (type 0) (i32.ne    (local.get 0) (local.get 1)))
 (func $blt  (type 0) (i32.lt_s  (local.get 0) (local.get 1)))
 (func $bltu (type 0) (i32.lt_u  (local.get 0) (local.get 1)))
 (func $bge  (type 0) (i32.ge_s  (local.get 0) (local.get 1)))
 (func $bgeu (type 0) (i32.ge_u  (local.get 0) (local.get 1)))
 (func $fld  (type 2) (f64.load offset=384 (i32.add (local.get 0) (local.get 1))))
 (func $fle  (type 5) (f64.le    (local.get 0) (local.get 1)))
 (func $flt  (type 5) (f64.lt    (local.get 0) (local.get 1)))
 (func $feq  (type 5) (f64.eq    (local.get 0) (local.get 1)))
 (memory 1)
 (export "mem" (memory 0)) (export "jmp" (func $jmp))
 (table 57 funcref)
 (elem (i32.const 0)  $lb   $xxx  $lw  $xxx $lbu  $xxx  $xxx $xxx                  ;;  0..7
                          $add  $shl  $xxx $xxx $xor  $shr  $or  $and $sb $xxx $sw     ;;  7..18
			  $mul  $xxx  $xxx $xxx $div  $divu $rem $remu                 ;; 19..26
			  $sub  $xxx  $xxx $xxx $xxx  $sra  $xxx $xxx                  ;; 27..34
			  $beq  $bne  $xxx $xxx $blt  $bge  $bltu $bgeu                ;; 35..42
			  $fadd $fsub $fmul $fdiv                                      ;; 43..46
 			  $fle  $flt  $feq $xxx                                        ;; 47..50
			  $sin  $cos  $exp $log $atan2 $hypot                          ;; 51..56
))
