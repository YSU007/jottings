import re
import os
import sys

proto_path = sys.argv[1]
out_path = sys.argv[2]

content = ""
files = os.listdir(proto_path)
for f in files:
    if os.path.splitext(f)[1] == ".proto":
        fd = open(f, mode='r')
        content += fd.read()


class HANDLEDef:
    def __init__(self, n, r, q, p):
        self.num = n
        self.rpc = r
        self.req = q
        self.rsp = p

    def handle(self):
        return f'''
        package logic
        import (
            "Jottings/tiny_rpc/log"
            "Jottings/tiny_rpc/model"
            "Jottings/tiny_rpc/proto"
        )
        func {self.rpc}Handle(a *model.PlayerAccount, req *proto.{self.req}, rsp *proto.{self.rsp}) (code uint32) {{
        return
        }}
        '''


RPC_COMMENT_NUM = 4
PROTO_NUM = "PROTO_NUM:"
RPC_NAME = "RPC_NAME:"
REQ_NAME = "REQ_NAME:"
RSP_NAME = "RSP_NAME:"
HANDLE_SUFFIX = "_handle.go"
ENUM_FILE = "proto.go"

rex_block = "\/\*[\w\W]*?\*\/"
rex_field = f"{PROTO_NUM}.*|{RPC_NAME}.*|{REQ_NAME}.*|{RSP_NAME}.*"

def_arr = []
block = re.findall(rex_block, content)
for b in block:
    bf = re.findall(rex_field, b)
    if len(bf) != RPC_COMMENT_NUM:
        print("proto file comment err\n", b)
        exit()
    num = str(bf[0]).removeprefix(PROTO_NUM)
    rpc = str(bf[1]).removeprefix(RPC_NAME)
    req = str(bf[2]).removeprefix(REQ_NAME)
    rsp = str(bf[3]).removeprefix(RSP_NAME)
    def_arr.append(HANDLEDef(num, rpc, req, rsp))

enum_arr = dict()
for d in def_arr:
    fp = f"{out_path}/{str(d.rpc).lower()}{HANDLE_SUFFIX}"
    f = open(fp, mode='w')
    f.write(d.handle())
    f.flush()
    f.close()
    enum_arr[d.rpc] = d.num
enum_str = enum_arr.__str__().strip("{}").replace("'", "").replace(":", "=").replace(",", "\n")
fp = f"{out_path}/{ENUM_FILE}"
f = open(fp, mode='w')
f.write(f'''
package proto
const (
    {enum_str}
)''')
f.flush()
f.close()
