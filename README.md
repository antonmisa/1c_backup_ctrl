# How to configure?

1. Globally it uses environment variable CONFIG_PATH to find configuration file, if not found trying to find in local dir.
    In config file: 
    path_to_rac                         - Path to rac executable, which installs with 1C client. ("C:/Program Files/1cv8/8.3.14.1857/bin/rac.exe")
    path_to_1cs                         - Path to 1C client. ("C:/Program Files/1cv8/8.3.14.1857/bin/1cv8.exe")

2. Next executable uses flags:
	--clusterConnection localhost:1545  - cluster connection string
    --clusterName localhost:1541        - cluster host:port to make a backup in cli mode
    --clusterAdmin AdminName            - cluster admin name if needed
    --clusterPwd  AdminPwd              - cluster password if needed
    --infobase    basename              - Infobase name (lowercase) in cluster to make a backup
    --infobaseUser ibName               - infobase user name, which has permission for backup
    --infobasePwd  ibPwd                - infobase user password
    --output DirToPutBackup             - Directory for backup   

# How to start? 

1. Using Powershell (Windows):
$env:CONFIG_PATH = "./config/config.yml"; ctrl.exe --clusterConnection localhost:1545 --clusterName localhost:1541 --infobase test --infobaseUser robot --infobasePwd robot --output ./backup

2. Using Bash (Linux):
set CONFIG_PATH="./config/config.yml" && ctrl.exe --clusterConnection localhost:1545 --clusterName localhost:1541 --infobase test --infobaseUser robot --infobasePwd robot --output ./backup