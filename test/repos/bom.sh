#!/bin/bash
set -ex; rm -rf repo; mkdir repo; cd repo

git init

cat <<EOT >> bomtest.txt ﻿A,B,C,D,E
F,G,H,I,J
K,L,M,N,O
P,Q,R,S,T
U,V,W,X,Y
Z,1,2,3,4
EOT
