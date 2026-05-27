# ACL

This package is reserved for anti-corruption adapters to other SEMS microservices.
The Device Management Service stores `user_id` and `home_id` as external UUID references and does not create foreign keys to IAM or Home services.
