# -*- mode: ruby -*-
# # vi: set ft=ruby :

require 'fileutils'
require 'open-uri'
require 'tempfile'
require 'yaml'

Vagrant.require_version ">= 1.6.0"

$update_channel = "alpha"

CLUSTER_IP="10.3.0.1"
NODE_IP = "172.17.4.99"
NODE_MEMORY_SIZE = 1024
USER_DATA_PATH = File.expand_path("contrib/vagrant/user-data")
SSL_TARBALL_PATH = File.expand_path("contrib/vagrant/ssl/controller.tar")
DEPLOYIT_CONFIG_PATH = File.expand_path("contrib/vagrant/deployit.yml")

system("mkdir -p ./contrib/vagrant/ssl && ./contrib/vagrant/lib/init-ssl-ca ./contrib/vagrant/ssl") or abort ("failed generating SSL CA artifacts")
system("./contrib/vagrant/lib/init-ssl ./contrib/vagrant/ssl apiserver controller IP.1=#{NODE_IP},IP.2=#{CLUSTER_IP}") or abort ("failed generating SSL certificate artifacts")
system("./contrib/vagrant/lib/init-ssl ./contrib/vagrant/ssl admin kube-admin") or abort("failed generating admin SSL artifacts")

Vagrant.configure("2") do |config|
  # always use Vagrant's insecure key
  config.ssh.insert_key = false

  config.vm.box = "coreos-%s" % $update_channel
  config.vm.box_version = ">= 1151.0.0"
  config.vm.box_url = "http://%s.release.core-os.net/amd64-usr/current/coreos_production_vagrant.json" % $update_channel

  ["vmware_fusion", "vmware_workstation"].each do |vmware|
    config.vm.provider vmware do |v, override|
      v.vmx['numvcpus'] = 1
      v.vmx['memsize'] = NODE_MEMORY_SIZE
      v.gui = false

      override.vm.box_url = "http://%s.release.core-os.net/amd64-usr/current/coreos_production_vagrant_vmware_fusion.json" % $update_channel
    end
  end

  config.vm.provider :virtualbox do |v|
    v.cpus = 1
    v.gui = false
    v.memory = NODE_MEMORY_SIZE

    # On VirtualBox, we don't have guest additions or a functional vboxsf
    # in CoreOS, so tell Vagrant that so it can be smarter.
    v.check_guest_additions = false
    v.functional_vboxsf     = false
  end

  # plugin conflict
  if Vagrant.has_plugin?("vagrant-vbguest") then
    config.vbguest.auto_update = false
  end

  config.vm.network :private_network, ip: NODE_IP
  config.vm.network :forwarded_port, guest: 2967, host: 2967

  config.vm.synced_folder ".", "/home/core/deployit",
    id: "core", :nfs => true, :mount_options => ['nolock,vers=3,udp']
  

  config.vm.provision :file, :source => SSL_TARBALL_PATH, :destination => "/tmp/ssl.tar"
  config.vm.provision :shell, :inline => "mkdir -p /etc/kubernetes/ssl && tar -C /etc/kubernetes/ssl -xf /tmp/ssl.tar", :privileged => true

  config.vm.provision :file, :source => USER_DATA_PATH, :destination => "/tmp/vagrantfile-user-data"
  config.vm.provision :shell, :inline => "mv /tmp/vagrantfile-user-data /var/lib/coreos-vagrant/", :privileged => true

  #kubernetes should be already ready so let build and run deployit
  config.vm.provision :file, :source => DEPLOYIT_CONFIG_PATH, :destination => "/home/core/.deployit/config.yml"
  config.vm.provision "docker" do |d|
    d.build_image "/home/core/deployit", args: "-t deployit -f /home/core/deployit/images/deployit/Dockerfile"
    d.run "deployit", name: "deployit-daemon", args:"-v /home/core/.deployit/config.yml:/etc/deployit/config.yml -p 2967:2967"
  end

end
