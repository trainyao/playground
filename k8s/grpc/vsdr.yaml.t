---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: dr-host3
spec:
  host: svc-host3
  subsets:
  - name: v1
    labels:
      app: host3
  - name: v2
    labels:
      app: host3
      type: stopping
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: vs-host3
spec:
  hosts:
  - svc-host3
  http:
  - fault:
      abort:
        httpStatus: 555
        percent: 100
    match:
    - uri:
        prefix: /
    # - uri:
        # exact: /login
    route:
    - destination:
        host: svc-host3
        port:
          number: 20003
        subset: v1
      weight: 10
    - destination:
        host: svc-host3
        port:
          number: 20003
        subset: v2
      weight: 90

