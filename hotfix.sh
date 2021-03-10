#!/bin/sh
curl -H "Content-Type: application/json" -X POST -d '{"Vendor": "juniper", "Type": "srx", "Version": "6.0", "Mode": "login", "fix_type": 2, "Prompts": ["asdfsf"]}' http://localhost:8189/api/operator/hotfix
curl -H "Content-Type: application/json" -X POST -d '{"Vendor": "juniper", "Type": "srx", "Version": "6.0", "Mode": "login", "fix_type": 2, "Prompts": ["asdfsf"]}' http://localhost:8189/api/operator/dump
curl -H "Content-Type: application/json" -X POST -d '{"Vendor": "juniper", "Type": "srx", "Version": "6.0", "Mode": "login", "fix_type": 1, "Prompts": ["asdfsf", "append_tet"]}' http://localhost:8189/api/operator/hotfix
curl -H "Content-Type: application/json" -X POST -d '{"Vendor": "juniper", "Type": "srx", "Version": "6.0", "Mode": "login", "fix_type": 1, "Prompts": ["asdfsf"]}' http://localhost:8189/api/operator/dump
curl -H "Content-Type: application/json" -X POST -d '{"Vendor": "juniper", "Type": "srx", "Version": "6.0", "Mode": "login", "fix_type": 1, "Prompts": ["[[:alnum:]]{1,}[[:alnum:]-_]{0,} (#|\\$) $", "append_tet"]}' http://localhost:8189/api/operator/hotfix
curl -H "Content-Type: application/json" -X POST -d '{"Vendor": "juniper", "Type": "srx", "Version": "6.0", "Mode": "login", "fix_type": 1, "Prompts": ["asdfsf"]}' http://localhost:8189/api/operator/dump
curl -H "Content-Type: application/json" -X POST -d '{"Vendor": "fortinet", "Type": "FortiGate-VM64-KVM", "Version": "v5.6.x", "Mode": "root", "fix_type": 1, "Prompts": ["[[:alnum:]]{1,}[[:alnum:]-_]{0,} \([[:alnum:]]+\) (#|\$) $"]}' http://localhost:8189/api/operator/hotfix
