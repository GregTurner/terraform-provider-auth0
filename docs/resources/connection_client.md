---
page_title: "Resource: auth0_connection_client"
description: |-
  With this resource, you can enable a single client on a connection.
---

# Resource: auth0_connection_client

With this resource, you can enable a single client on a connection.

!> This resource appends an enabled client to a connection. In contrast, the `auth0_connection_clients` resource
manages all the clients enabled for a connection. To avoid potential issues, it is recommended not to use this
resource in conjunction with the `auth0_connection_clients` resource when managing enabled clients for the
same connection id.

## Example Usage

```terraform
resource "auth0_connection" "my_conn" {
  name     = "My-Auth0-Connection"
  strategy = "auth0"
}

resource "auth0_client" "my_client" {
  name = "My-Auth0-Client"
}

# One connection to one client association.
# To prevent issues, avoid using this resource together with the `auth0_connection_clients` resource.
resource "auth0_connection_client" "my_conn_client_assoc" {
  connection_id = auth0_connection.my_conn.id
  client_id     = auth0_client.my_client.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `client_id` (String) ID of the client for which the connection is enabled.
- `connection_id` (String) ID of the connection on which to enable the client.

### Read-Only

- `id` (String) The ID of this resource.
- `name` (String) The name of the connection on which to enable the client.
- `strategy` (String) The strategy of the connection on which to enable the client.

## Import

Import is supported using the following syntax:

```shell
# This resource can be imported by specifying the
# connection ID and client ID separated by ":".
#
# Example:
terraform import auth0_connection_client.my_conn_client_assoc con_XXXXX:XXXXXXXX
```
