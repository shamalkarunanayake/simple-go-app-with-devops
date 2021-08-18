# Setup Guide

## Prerequisites
- [Helm](https://helm.sh/docs/intro/install/) >= 3.2 (chocolaty install default helm version is less than 3.2, Please makesure to install required version)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) you must use a version that is within one minor version difference of your cluster. For example, a v1.2 client should work with v1.1, v1.2, and v1.3 master.
- [Docker](https://docs.docker.com/desktop/) >= 18.0.0
---

1. Download the repository 
  
    ```git clone https://github.com/shamalkarunanayake/simple-go-app-with-devops.git```

2. First have a look on go-app-custom-metrics folder to get in touch with simple go-app. It consist of a Dockerfile. and the image also published in dockerhub.

   ```docker pull shamalskk/final-go-app-log1```

   In main.go file there are some implementations to expose prometheus  metrices in go-app. given below are example steps to build and push Docker image using dockerfile

        
        docker build -t shamalskk/example-app . 

        docker images

        docker run -d -p 8080:8080 shamalskk/example-app

        curl localhost:8080/metrics 

        docker login 

        docker push shamalskk/example-app
        

3. Create a kubernetes cluster in azure using terraform

    ### Prerequisites
* Terraform version >= 0.14
    * cd terraform 

    * Export credentials for Azure

    ```
    export ARM_CLIENT_ID="XXXXX"
    export ARM_SUBSCRIPTION_ID="XXXXX"
    export ARM_TENANT_ID="XXXXX"
    export ARM_CLIENT_SECRET="XXXXX"
    ```
    * apply the below commands to create and destroy aks cluster. please make sure to keep the tfstate file with you! 

    ```
    terraform init

    terraform apply

    terraform plan 

    terraform destroy
    ```

4. Login to azure portal and get the cluster config to local machine

    * Browse azure portal and navigate to specific resource group. Now you can see aks and network is created. (you can refer the connect option inside aks to connect to created cluster.example steps are given below)
    * check whether config is set to the current context 


    ```
    az account set --subscription xxxxxxxx

    az aks get-credentials --resource-group xxxxxxxx --name xxxxxxxx

    kubectl config current-context
    ```

5. Deploy our simple go-app in kubernetes 

    * cd kubernetes 
    * you can apply the below commands to deploy go-app & check whether your deployemnt, service, pod are  correctly created 

    ```
    kubectl apply -f deployment-goapp.yaml

    kubectl apply -f service-goapp.yaml

    kubectl get deployment

    kubectl get service

    kubectl get pods 

    ```
6. Deploy prometheus operator into our kubernetes cluster using Helm 

    * prometheus operator -  https://github.com/prometheus-community/helm-charts/tree/main/charts/kube-prometheus-stack 
    * First add the helm repo 
    * You can install chart with fixed version number as well. examples are shown below

    ```
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

    helm repo update

    helm install prometheus prometheus-community/kube-prometheus-stack

    helm install prometheus prometheus-community/kube-prometheus-stack --version "9.4.1"
    ```

    * After installing , you can refer all the resources created using below command (pod, service, daemonset, deployment, replicaset, statefulset)
    * We need to create a service monitor to our go-app . go inside kubernetes folder and apply the service-monitor-go-app.yaml.

    ```
    kubectl get all 

    kubectl apply -f service-monitor-go-app.yaml

    kubectl get servicemonitor
    ```

    * Now we can port-forward prometheus, grafana instances and see whether the metrics of our simple-go-app and kubernetes cluster 

    ```
    kubectl port-forward prometheus-prometheus-kube-prometheus-prometheus-0 9090

    kubectl port-forward prometheus-grafana-xxxxx 3000
    ```

7. Logging 

    * Latest helm charts can be gathered from https://github.com/elastic/helm-charts 

    ```
    helm repo add elastic https://helm.elastic.co 
    helm repo add fluent https://fluent.github.io/helm-charts
    helm repo update
    ```
    ```
    # Install elasticsearch

    helm install elasticsearch elastic/elasticsearch

    # Install kibana

    helm install kibana elastic/kibana

    # Install fluentd

    helm install fluentd fluent/fluentd
    helm show values fluent/fluentd

    or else 

    cd elastic-fluentd-kibana
    kubectl apply -f fluentd.yaml

    # Install metric-beat

    helm install metricbeat elastic/metricbeat

    # Install file-beat
    helm install filebeat elastic/filebeat

    # Install logstash
    helm install logstash elastic/logstash
    ```

    * Go inside elastic-fluentd-kibana folder apply the counter.yaml to generate simple counter logs
    ```
    kubectl apply -f counter.yaml
    ```
    * After following below steps, you can access kibana 

    ```
    kubectl port-forward --namespace default  elasticsearch-master-0  9200:9200

    kubectl port-forward kibana-kibana-b4dfc69c7-rtx9g 5601

    kubectl --namespace default port-forward fluentd-XXXX 24231:24231
    curl http://127.0.0.1:24231/metrics
    ```

8. Deploy java big-memory-app to test auto-scaling 

   * Go inside kubernetes folder and deploy java-bigmemoryapp.yaml

   ```
   kubectl apply -f java-bigmemoryapp.yaml

   kubectl apply -f vertical-pod-autoscaler.yaml
   ```