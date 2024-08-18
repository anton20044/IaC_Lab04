#output "load_balancer_public_ip" {
#  description = "Public IP address of load balancer"
#  value = yandex_lb_network_load_balancer.wp_lb.listener.*.external_address_spec[0].*.address
#}
output "load_balancer_public_ip" {
  description = "Public IP address of load balancer"
  value = tolist(tolist(yandex_lb_network_load_balancer.wp_lb.listener).0.external_address_spec).0.address
}

output "vm_linux_public_ip_address" {
  description = "Virtual machine IP"
  value = toset([ for i in yandex_compute_instance.wp-app: i.network_interface[0].nat_ip_address])
}


output "db_user" {
  description = "DB User"
  value = local.dbuser
}

output "db_name" {
  description = "DB name"
  value = local.dbname
}

output "db_pass" {
  description = "DB User pass"
  value = local.dbpassword
  sensitive = true
}

output "database_host_fqdn" {
  description = "DB hostname"
  value = toset(local.dbhosts)
}

