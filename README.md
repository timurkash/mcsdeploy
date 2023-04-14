# mcsdeploy

This util helps to configure local deploy considering

- configs
  - service1-v1
    - config.yaml 
  - service2-v1
      - config.yaml
  - service3-v1
      - config.yaml
  - ...
- envoy
  - Dockerfile
  - envoy.yaml
- docker-compose.yml
- Makefile

You have to set file `services.yaml` like
```yaml
services:
- name: service1
  version: v1
- name: service2
  version: v1
- name: service3
  version: v1
- ...
```

## -vup, -uvp, -upv

Use is to up git tag and image tag

## -env, -doc, mak

Use it to get slice of envoy, docker-compose.yml, Makefile



## -prt

To get slice of proto: rfc and messages you can run

```bash
mcsdeploy -prt shop_product
```

next parameter can be plural of 

then you can receive

```

  // ShopProduct
  rpc ActShopProduct (ShopProductRequest) returns (ShopProductReply);
  rpc ListShopProducts (ListShopProductsRequest) returns (ListShopProductsReply);

// ShopProduct
message ShopProductRequest {
  common.ActionId action_id = 1;
  ShopProductInfo shop_product = 2;
}

message ShopProductReply {
  ShopProductInfo shop_product = 1;
  common.IdTimestamps id_timestamps = 2;
}

message ShopProductInfo {
  string name = 1;
}

message ListShopProductsRequest {
  message Filter {
    common.String name = 1;
  }
  Filter filter = 1;
  common.OrderOffsetLimit ool = 2;
}

message ListShopProductsReply {
  repeated ShopProductReply shop_products = 1;
  common.Paging paging = 2;
}

```
