# ANSIBLE-TERRAFORM INFRA

---
## Project description
Set up compute for the web, application and database tiers. Manage configuration using ansible

---
## Project set up

### Ansible Setup
Change directory to the ansible directory and run a python virtual environment
```bash
python3 -m venv tf-ans
source tf-ans/bin/activate
pip install ansible boto boto3
pip install --upgrade pip
ansible --version  
```

Verify that ansible can reach all servers

```bash
ansible all -m ping
```

Generate Ansible roles using ansible-galaxy
```bash
ansible-galaxy init roles/database
ansible-galaxy init roles/application
ansible-galaxy init roles/web  
```

You can update the playbook and do a dry run
```bash
ansible-playbook -i inventory playbook.yml --check
```