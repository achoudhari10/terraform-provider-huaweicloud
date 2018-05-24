# Configure the OpenStack Provider
provider "huaweicloud" {
  user_name   = "Longyin_kk"
  domain_name = "longyin_kk"
  password    = "kuaiLE02"
  auth_url    = "https://iam.cn-north-1.myhuaweicloud.com/v3"
  region      = "cn-north-1"
  tenant_id   = "c59fd21fd2a94963b822d8985b884673"
}

resource "huaweicloud_cce_cluster_v3" "cluster_1" {
  kind="Cluster"
  api_version="v3"
  name = "test-c2d"
  cluster_type="VirtualMachine"
  flavor="cce.s1.small"
  cluster_version = "v1.7.3-r10"
  vpc_id="3305eb40-2707-4940-921c-9f335f84a2ca"
  subnet_id="00e41db7-e56b-4946-bf91-27bb9effd664"
  container_network_type="overlay_l2"
  description="update1111"
}
output "resource_stack" {
  value = "${huaweicloud_cce_cluster_v3.cluster_1.id}"
}
output "status" {
  value = "${huaweicloud_cce_cluster_v3.cluster_1.status}"
}




