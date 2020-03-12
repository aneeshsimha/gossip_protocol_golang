# gossip_protocol_golang

---

### AWS Setup notes:

Basically just
1. https://us-west-1.console.aws.amazon.com/ec2/v2/home?region=us-west-1#LaunchInstanceWizard:
2. Select `Ubuntu Server 18.04 LTS (HVM), SSD Volume Type`, as a t2.micro
3. -> review and launch -> launch
4. Make a new keypair if you don't have one and download and save that
5. Launch it, then find its public DNS hostname in the management console: https://us-west-1.console.aws.amazon.com/ec2/v2/home?region=us-west-1#Instances:
6. Once it's booted up, `ssh -i "/path/to/keyfile" ubuntu@ec2-a-bunch-of-numbers.us-west-1.compute.amazonaws.com`
7. Add a new security group that has a bunch of ports exposed to actually run on (we use 8000-8100)
8. Run `gossip_setup.sh && source ~/.profile` 
9. **NOTE:** It's very important to have the source code in `$GOPATH/src/github.com/aneeshsimha/gossip_protocol_golang`, otherwise it will not compile.