for i in {1..100}; do ../cmds/swdocs/swdocs apply --file rabbitmq.json file$i & done
for i in {200..300}; do ../cmds/swdocs/swdocs create --name bus$i file$i & done
