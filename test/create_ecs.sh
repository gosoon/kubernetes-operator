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
  "masterVIP":"192.168.72.81",
  "privateSSHKey": "-----BEGIN RSA PRIVATE KEY-----MIIEpAIBAAKCAQEAyWLL26/GmJQHtW00poxCgs1wyKUd4hvWiylP1VsGuQNaR3cFg7yggFJY8FCOs8ZCWUgJDe/bzXS/kFMc5NPUL6lhokw/2NJqtxOVnb+4Mvq/6TB3Uf/+h0FMRc6phk7bPmaZuDsHDrFbUSAvHjlH4OS3kuRT69r7OpviYw2y5pt/NhoNp9ydBcpStAvtv6NTS24Pb0PM+Pdn8zWA+6XXl91Q3lHvPd59ZhCutLx9MrMXJRiVZZ7xW5+qa/vP1QLcnlmsie/h3hZQc1EkyEwgc5eKXI/7VDGJWAfkB9pAaPKGu0iHM6B2MGSSm5wRWUh4A2CH/kt2TMpu+szllOo9jwIDAQABAoIBAF2tFDDDifitWwycmNIkCkg38g+TJtxnoJupAta2+eCT26nEho8p7eri4zYd8tNTFMfdB0ExYqgmd3lV/+m9U0U8YAsTtttPvY4dkQoJBVDJbP22qro8/xPBXw5VvGuaQMe9CCI3auf9vgF1nBVOBc5p9a5hgVwDx7sSifMTizVTKcqfsmkmUQM9eIpuIaSET1HPW4COSPRCNYKC6fY1TAnLtj7pyWHzLdD2i/PDUB1m7A8wR43FBRWJDhdpjdN/tuWnOeoXY3qswaRWdoilwQQQnQWZA6J6tg/HoWjusHJXXte1OZeZo8M3GYqvx73FNsvo1/1CnlFG0NHAF3eHaMECgYEA885gFMNd1af5UDBwNTdVRiUMvDRzDp1TTfaL9qWQSKzwGG/BloG6JQTm1N/7wD3KBE0pUWkMz7QaUI4qeo9rjS3d5R2j7+HuD7jgjL/a79KHaF2o5Asj7YgR7k88Inlumrt2whaBjSa6R5jOCOcJ6CE4l51eJ6NSQ21CoptIACECgYEA03VI0I1J6ilNZ6fBEpdJW8Z/2QD5hrx9Jal0n5B2QzzsMhHWbABM5ANgV+0pkw4gx3uVbA81aB0ejRRQAitEvfcUKjwO9zdVPGlATxUT/PXdZNjWES9Avs0IWVQUwTFLh9c95QSBaSkmG95oM0+uK18wxANsXbg/fn2vq8WJR68CgYEAqLB1WEhox2jmnq41JQz5CoSeECZ2KXl/ZyUcaHbbov1NTosctedihTSkBkxHoxbdjSZaXULDI38o6e3DxHMxZkiDDID6qgJ3thcj/x7L/D19hR+wuMBghnwsc+gM4omElrj4jYgG8UQHhXxbqls5Royd3IF84Q4m4BJcFag9JCECgYEAshEN5DFWkQ1+1U7600D60YHynzam6cNIT7LHNqdcL3raG7/RpNkL5ubA9soMPH2lNNbpGTolays6UutMBMeS97VdEcPJhnzeFiU7tly1nEseyJGgkpAMIaBe63pWj+mHBTlIMdb9cyTnpog/jxYGQRfD5QxM8Q76yPXmPOv3kpkCgYBTtBySpfsogsvDe3UvYFBCYJQi+qxaDX9iTrNEt7/9g7OpBXc5HoS8c8vsnwQeudTx5YDvqiyqdy1gS4wonyWW/1KLcmNnCmqQMENGGjoS5AF79qbqxCSuxahslE7wzAkwwaue5VvzgLepDEfysfGMm3edWZIz28BfUMh2nfpppg==-----END RSA PRIVATE KEY-----",
  "serviceCIDR": "",
  "retry": false
}
EOF
)

curl -s -XPOST -d "${request_body}" \
    http://127.0.0.1:8080/api/v1/region/${region}/cluster/${name}
