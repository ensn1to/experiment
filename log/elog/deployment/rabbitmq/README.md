the steps for running benchmark:

1. install docker-compose
2. docker-compose up 
3. docker exec it container-id /bin/bash
4. 启动插件:rabbitmq-plugins enable rabbitmq_management 
<!-- 5. go test -bench . -->



> ref: https://xie.infoq.cn/article/a52209e098e24a41f737112ad