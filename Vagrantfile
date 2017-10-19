# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.ssh.insert_key = false
  config.vm.provider "virtualbox" do |v|
    v.customize [ "guestproperty", "set", :id, "/VirtualBox/GuestAdd/VBoxService/--timesync-set-threshold", 1000 ]
  end
  config.vm.box = "puppetlabs/centos-7.2-64-nocm"

  [1, 2, 3].each do |idx|
    config.vm.define "node#{idx}" do |conf|
      conf.vm.provider "virtualbox" do |v|
        v.memory = 1024
        v.cpus = 2
      end
      conf.vm.hostname = "node#{idx}"
      conf.vm.network :private_network, ip: "192.168.50.10#{idx}"
      conf.vm.network :forwarded_port, guest: 9092, host: (9092 + idx*100)   # 19093, 29093, ...

      #
      # Doesn't run in parallel...
      #
      #config.vm.provision "ansible_local" do |ansible|
      #  ansible.playbook = "playbook.yml"
      #end
    end
  end

end
