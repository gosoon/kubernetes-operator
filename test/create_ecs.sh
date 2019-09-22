#!/bin/bash

region="default"
name="test-cluster"
request_body=$(cat<<EOF
{
  "name": "test-cluster",
  "clusterType": "kubernetes",
  "masterList": [
    {
      "ip": "192.168.72.81"
    }
  ],
  "nodeList": [
    {
      "ip": "192.168.72.81"
    },
    {
      "ip": "192.168.73.150"
    }
  ],
  "etcdList": [
    {
      "ip": "192.168.73.150"
    },
    {
      "ip": "192.168.72.81"
    }
  ],
  "externalLoadBalancer":"192.168.72.81",
  "privateSSHKey": "-----BEGIN OPENSSH PRIVATE KEY-----\nb3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABFwAAAAdzc2gtcn\nNhAAAAAwEAAQAAAQEAyYsDgmjs2Mvb7jbev++KWxRzS5DmonCuVyhbSDt6fL0nywOhHeDO\np5FjZcl5sWXOiDoxJap+K2vgZNzKAA3ejjdCt0i7DVWVDibJN0/R2AJmTQdaAek+Tuilck\nYLZAMCDBIh5QYqDfPUxYp7uEHTh6XznJmXeo+zwd5XhJdx6PxcjhxWnek4vv8QPAnJRlY5\nX8NkzZQtttKKczifXqIaYOCbbIMqRke/fn6xxt1OOkPpeXmcr1XoqxJHtuGVD/kcZmzS3L\nJYn+sksRZnxskCGTOBSFb8X+XmENs7ryzz1+ApRd9iuBn/OF91JsP6gaLguO6YkmVLFzQU\nW7Q0ZN30VwAAA9iUgbjPlIG4zwAAAAdzc2gtcnNhAAABAQDJiwOCaOzYy9vuNt6/74pbFH\nNLkOaicK5XKFtIO3p8vSfLA6Ed4M6nkWNlyXmxZc6IOjElqn4ra+Bk3MoADd6ON0K3SLsN\nVZUOJsk3T9HYAmZNB1oB6T5O6KVyRgtkAwIMEiHlBioN89TFinu4QdOHpfOcmZd6j7PB3l\neEl3Ho/FyOHFad6Ti+/xA8CclGVjlfw2TNlC220opzOJ9eohpg4JtsgypGR79+frHG3U46\nQ+l5eZyvVeirEke24ZUP+RxmbNLcslif6ySxFmfGyQIZM4FIVvxf5eYQ2zuvLPPX4ClF32\nK4Gf84X3Umw/qBouC47piSZUsXNBRbtDRk3fRXAAAAAwEAAQAAAQAYJAiVlE/aYADF9diU\nkPK3mil9Qav+hRS8596XNlijnFyp2pNv6r+WHroTNSDYeONWOfItGtDmDpPgQPoJK6Ae1M\nuu/I07OacS/N5ZO7xc7VynmVvUosWN2hwHJhCzOBEEtR9OOYDMDwrLZp0PIwNdWill9pfI\nXHIhpKpC/YjtaSZU1cTJRgpE18PH452rQo91RIvT6K1AXEQ+IgRzUU4KKU3QHB1MqPD/Ga\nbiLyH3VBvG2H7G866qYyhfY3s+NDR8o79+nTVxmd7H7lMFMAPQ3fXRoTdGxK2kbkd06E8Q\nGnGlWtzXSdT7G7cxsUfHSaH7fhEwZn95XO6o1JCHlGgBAAAAgQD1L4oltiT8W2c/aK5Z6a\n9QCr71JD3S6xRmKMpCLqst2ehqDi3etNLPAi6gJP5etJ39AgfUrqM37XWoV0J7DyeWBwgv\n3EFO1FFr8v8zBF4SX7Y+Yptf7CQztrg00qQB6hFI7O+n6ytyGZry0RbTjN2P5kCDnSdCQt\nuG4q4cINY7swAAAIEA+wnXjna7o1xopYb2NesYvmZ9FlNdSsFvRASwVfLOCyyiFf6kxqSL\njHwUAyWlQnvCWmRh5O3TtSGwEljbQi6yfIaHN7qpA3T5qzhYJmEZqObOZkKS0asUNsmMJc\nfxRpHhibIB+A5w0uapxamcbTbZurKQvoBfBVIPAIZa7e0HOlcAAACBAM2Gvm9Zw5x3Mmjk\ncq9fqvdcBaKI8DZYSA635lG56H7AsMrpJHU3eqlSjffhZdPVRDab5kqGiDbcWltf2w9zIe\nGoMFjaykKeNuypbb1Ffwz82q4dFRQ0djqO/dJLKYu8QodSO+sU7jxEpHTwIIZs2e9EIZBJ\nvC7yVME0JJiUPtYBAAAAIGZlaWxlaUBmZWlsZWlkZU1hY0Jvb2stUHJvLmxvY2FsAQI=\n-----END OPENSSH PRIVATE KEY-----\n",
  "serviceCIDR": "",
  "retry": false
}
EOF
)

curl -s -XPOST -d "${request_body}" \
    http://127.0.0.1:8080/api/v1/region/${region}/cluster/${name}
