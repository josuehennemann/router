SKIP=${1}

if [[ "$SKIP" != "s" ]]; then
	echo "compilando"
	go build 
	#valida se compilou o projeto, e somente faz o scp se compilou
	if [[ $? != 0 ]]; then
		echo "Falha ao compilar"
		exit 1
	fi
fi

echo "starting....."

./router -fileConf=config.json