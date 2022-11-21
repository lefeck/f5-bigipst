# f5-bigipst
f5-bigipst is a stress testing tool that creates random configuration policies through multiple threads

## Basic Usage
```cgo
[root@localhost ~]# ./f5-bigipst  -h
  -P string
        the specfiy  the location of partition (default "Common")
  -a string
        the remote of host ip (default "127.0.0.1")
  -f string
        Specify the file location of the output results
  -m string
        Specify the ip addess of member, If you don't specify an IP address, an IP address of 10.0.0.0/8 will be generated randomly.
  -n int
        The total of task numbers (default 10)
  -p string
        the password of login host (default "admin")
  -t duration
        Set the timeout period for connecting to the host (default 1m0s)
  -u string
        the username of login host (default "admin")
  -vs string
        the specfiy  virtual server of ip addess, If you don't specify an IP address, an IP address of 192.0.0.0/8 will be generated randomly.
  -w int
        The Number of threads to start worker work (default 10)

[root@localhost ~]# ./f5-bigipst  -a 192.168.5.134  -n 300
pool name Pool_Name_9hxy0s85 create success.
pool name Pool_Name_qr0hku2c create success.
pool name Pool_Name_8h7khz7a create success.
pool name Pool_Name_r51bien7 create success.
pool name Pool_Name_2s2yqn0d create success.
pool name Pool_Name_yz42qmmk create success.
pool name Pool_Name_eyomeoix create success.
pool name Pool_Name_hqeqozwk create success.
pool name Pool_Name_k2edan8j create success.
pool name Pool_Name_s286prxt create success.
virtualserver name Virtual_Name_r51bien7 create success.
virtualserver name Virtual_Name_2s2yqn0d create success.
virtualserver name Virtual_Name_yz42qmmk create success.
virtualserver name Virtual_Name_9hxy0s85 create success.
virtualserver name Virtual_Name_qr0hku2c create success.
virtualserver name Virtual_Name_eyomeoix create success.
virtualserver name Virtual_Name_8h7khz7a create success.
virtualserver name Virtual_Name_hqeqozwk create success.
virtualserver name Virtual_Name_k2edan8j create success.
virtualserver name Virtual_Name_s286prxt create success.
pool name Pool_Name_qoipkj9i create success.
pool name Pool_Name_hltwoaah create success.
pool name Pool_Name_bqczxo84 create success.
pool name Pool_Name_0rfkkbo4 create success.
pool name Pool_Name_zlqkcl9c create success.
pool name Pool_Name_kacq7yd0 create success.
pool name Pool_Name_wxdl03vj create success.
pool name Pool_Name_fvq66zwd create success.
pool name Pool_Name_by7lbfk5 create success.
...
```