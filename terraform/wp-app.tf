resource "yandex_compute_instance" "wp-app" {

	for_each = var.wp_hosts
	name = each.key
	zone = "${each.value == "wp-subnet-a" ? "ru-central1-a" : "ru-central1-b"}"

	resources {
	  cores  = 2
          memory = 2
        }

	boot_disk {
         initialize_params {
         image_id = "fd80c3igl4ivrjvjfljm"
	#fd80viupr3qjr5g6g9du"
        }
      }


  network_interface {
    # Указан id подсети default-ru-central1-a
    subnet_id = "${each.value == "wp-subnet-a" ? yandex_vpc_subnet.wp-subnet-a.id : yandex_vpc_subnet.wp-subnet-b.id}"
    nat       = true
  }

  metadata = {
    ssh-keys = "ubuntu:${file("~/.ssh/id_ed25519.pub")}"
  }
}






#resource "yandex_compute_instance" "wp-app-1" {
#  name = "wp-app-1"
#  zone = "ru-central1-a"

#  resources {
#    cores  = 2
#    memory = 2
#  }

#  boot_disk {
#    initialize_params {
#      image_id = "fd80viupr3qjr5g6g9du"
#    }
#  }

#  network_interface {
    # Указан id подсети default-ru-central1-a
#    subnet_id = yandex_vpc_subnet.wp-subnet-a.id
#    nat       = true
#  }

#  metadata = {
#    ssh-keys = "ubuntu:${file("~/.ssh/id_ed25519.pub")}"
#  }
#}

#resource "yandex_compute_instance" "wp-app-2" {
#  name = "wp-app-2"
#  zone = "ru-central1-b"

#  resources {
#    cores  = 2
#    memory = 2
#  }

#  boot_disk {
#    initialize_params {
#      image_id = "fd80viupr3qjr5g6g9du"
#    }
#  }

#  network_interface {
    # Указан id подсети default-ru-central1-b
#    subnet_id = yandex_vpc_subnet.wp-subnet-b.id
#    nat       = true
#  }

#  metadata = {
#    ssh-keys = "ubuntu:${file("~/.ssh/id_ed25519.pub")}"
#  }
#}
