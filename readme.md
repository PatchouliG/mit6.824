## current status

### lab1 finish

#### how to run
cd 
sh mr-test.sh

#### lab1 info in mit
https://pdos.csail.mit.edu/6.824/labs/lab-mr.html

#### issue
map job生成的中间文件，应该先生成一个临时文件，由master rename
因为如果有一个运行缓慢的map job 发生，超时，master会分配一个新的map job
这两个job均会写入同一个文件，尚不清楚会不会导致文件内容混乱
