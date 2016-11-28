
# (SSH) Metro
Metro is very simple SSH tunneling tool. Providing an ability to connect to 
machine via SSH and tunnel configured ports through the connection.





## Usage
./metro -host=192.168.0.100 -user=sshuser -password=sshuser -list=tunnels.csv -port=22 -timeout=20

*Arguments*
 - host,port - host and port of the SSH connection
 - user,password - credentials of the SSH connection
 - timeout - timeout for the SSH connection
 - list - file with tunnels (each line describes one tunnel in format ```<local_port>;<remote_host>:<remote_port>```)



 