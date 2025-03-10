output "application_public_ip" {
  value = module.application.public_ip
}

output "database_public_ip" {
  value = module.database.public_ip
}

output "web_public_ip" {
  value = module.web.public_ip
}

