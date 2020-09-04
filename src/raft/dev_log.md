1. vote rpc的term，我认为不是有效的term，因为election不一定成功，接收到该vote的leader or candidate，
不会因为vote的term比自己的高而发生角色的变化 (个人的想法，没有raft的论文明确的支持，但是应该是对的)