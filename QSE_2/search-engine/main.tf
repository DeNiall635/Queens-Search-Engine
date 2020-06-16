// Configure the Google Cloud provider
provider "google" {
credentials = "${file("servicekey.json")}"
project = "spider-indexer-40173800"
region = "us-west1"
}

// Terraform plugin for creating random ids
resource "random_id" "instance_id" {
byte_length = 8
}

// A single Google Cloud Engine instance
resource "google_compute_instance" "default" {
name = "flask-vm-${random_id.instance_id.hex}"
machine_type = "f1-micro"
zone = "us-west1-a"

boot_disk {
initialize_params {
image = "debian-cloud/debian-9"
}
}

// Make sure flask is installed on all new instances for later steps
metadata_startup_script = "sudo apt-get update; sudo apt-get install -yq build-essential python-pip rsync; pip install flask; pip install scrapy; pip install scrapyrt; pip install requests"

network_interface {
network = "default"

access_config {
// Include this section to give the VM an external ip address
}
}
metadata = {
ssh-keys = "linux:${file("id_rsa.pub")}"
}
}

output "ip" {
value = "${google_compute_instance.default.network_interface.0.access_config.0.nat_ip}"
}
