table "account_entries" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "account_id" {
    null = false
    type = character_varying(16)
  }
  column "payment_id" {
    null = false
    type = character_varying(16)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "amount" {
    null = false
    type = numeric(16,6)
  }
  column "balance" {
    null = false
    type = numeric(16,6)
  }
  column "title" {
    null = false
    type = character_varying(64)
  }
  column "comment" {
    null = false
    type = character_varying(128)
  }
  column "signature" {
    null = false
    type = character_varying(128)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_account_entries_accounts_account_id" {
    columns     = [column.account_id]
    ref_columns = [table.accounts.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_account_entries_payments_payment_id" {
    columns     = [column.payment_id]
    ref_columns = [table.payments.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_account_entries_account_id" {
    columns = [column.account_id]
  }
  index "ix_account_entries_payment_id" {
    columns = [column.payment_id]
  }
}
table "accounts" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "holder_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "valid_time" {
    null = true
    type = timestamp
  }
  column "flags" {
    null = false
    type = bigint
  }
  column "currency" {
    null = false
    type = character_varying(32)
  }
  column "balance" {
    null = true
    type = numeric(16,6)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_accounts_participants_holder_id" {
    columns     = [column.holder_id]
    ref_columns = [table.participants.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_accounts_holder_id" {
    columns = [column.holder_id]
  }
}
table "addresses" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "owner_id" {
    null = false
    type = character_varying(16)
  }
  column "region_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "archived" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "preferred" {
    null    = false
    type    = boolean
    default = false
  }
  column "summary" {
    null = false
    type = character_varying(255)
  }
  column "details" {
    null = false
    type = character_varying(255)
  }
  column "post_code" {
    null = true
    type = character_varying(32)
  }
  column "contact_name" {
    null = false
    type = character_varying(128)
  }
  column "email_address" {
    null = true
    type = character_varying(128)
  }
  column "phone_number" {
    null = true
    type = character_varying(128)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_addresses_regions_region_id" {
    columns     = [column.region_id]
    ref_columns = [table.regions.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_addresses_users_owner_id" {
    columns     = [column.owner_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_addresses_owner_id" {
    columns = [column.owner_id]
  }
  index "ix_addresses_region_id" {
    columns = [column.region_id]
  }
}
table "catalogs" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "parent_id" {
    null = true
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "name" {
    null = false
    type = character_varying(64)
  }
  column "title" {
    null = false
    type = character_varying(64)
  }
  column "description" {
    null = true
    type = character_varying(255)
  }
  column "template" {
    null = false
    type = jsonb
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_catalogs_catalogs_parent_id" {
    columns     = [column.parent_id]
    ref_columns = [table.catalogs.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_catalogs_parent_id_name" {
    unique  = true
    columns = [column.parent_id, column.name]
  }
}
table "categories" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "parent_id" {
    null = true
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "name" {
    null = false
    type = character_varying(64)
  }
  column "title" {
    null = false
    type = character_varying(64)
  }
  column "description" {
    null = true
    type = character_varying(255)
  }
  column "template" {
    null = false
    type = jsonb
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_categories_categories_parent_id" {
    columns     = [column.parent_id]
    ref_columns = [table.categories.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_categories_parent_id_title" {
    unique  = true
    columns = [column.parent_id, column.title]
  }
}
table "category_products" {
  schema = schema.public
  column "product_id" {
    null = false
    type = character_varying(16)
  }
  column "category_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  primary_key {
    columns = [column.category_id, column.product_id]
  }
  foreign_key "fk_category_products_categories_category_id" {
    columns     = [column.category_id]
    ref_columns = [table.categories.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_category_products_products_product_id" {
    columns     = [column.product_id]
    ref_columns = [table.products.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_category_products_category_id" {
    columns = [column.category_id]
  }
  index "ix_category_products_product_id" {
    columns = [column.product_id]
  }
}
table "chat_members" {
  schema = schema.public
  column "user_id" {
    null = false
    type = character_varying(16)
  }
  column "session_id" {
    null = false
    type = character_varying(16)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "read_cursor" {
    null = true
    type = character_varying(16)
  }
  column "permission" {
    null = false
    type = bigint
  }
  column "display_name" {
    null = true
    type = character_varying(64)
  }
  primary_key {
    columns = [column.user_id, column.session_id]
  }
  foreign_key "fk_chat_members_chat_sessions_session_id" {
    columns     = [column.session_id]
    ref_columns = [table.chat_sessions.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_chat_members_users_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_chat_members_session_id" {
    columns = [column.session_id]
  }
  index "ix_chat_members_user_id" {
    columns = [column.user_id]
  }
}
table "chat_records" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "session_id" {
    null = false
    type = character_varying(16)
  }
  column "creator_id" {
    null = false
    type = character_varying(16)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "version" {
    null = false
    type = character_varying(32)
  }
  column "headers" {
    null = false
    type = jsonb
  }
  column "content" {
    null = false
    type = jsonb
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_chat_records_chat_sessions_session_id" {
    columns     = [column.session_id]
    ref_columns = [table.chat_sessions.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_chat_records_users_creator_id" {
    columns     = [column.creator_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_chat_records_creator_id" {
    columns = [column.creator_id]
  }
  index "ix_chat_records_session_id" {
    columns = [column.session_id]
  }
}
table "chat_sessions" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "readonly" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "icon_url" {
    null = true
    type = character_varying(2048)
  }
  column "title" {
    null = false
    type = character_varying(32)
  }
  column "introduction" {
    null = true
    type = character_varying(128)
  }
  primary_key {
    columns = [column.id]
  }
}
table "client_users" {
  schema = schema.public
  column "client_id" {
    null = false
    type = character_varying(16)
  }
  column "user_id" {
    null = false
    type = character_varying(16)
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  primary_key {
    columns = [column.client_id, column.user_id]
  }
  foreign_key "fk_client_users_clients_client_id" {
    columns     = [column.client_id]
    ref_columns = [table.clients.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_client_users_users_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
}
table "clients" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "secret_key" {
    null = false
    type = character_varying(32)
  }
  column "secret_code" {
    null = true
    type = character_varying(64)
  }
  column "description" {
    null = true
    type = character_varying(255)
  }
  primary_key {
    columns = [column.id]
  }
  index "ix_clients_secret_key" {
    unique  = true
    columns = [column.secret_key]
  }
}
table "deliveries" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "region_id" {
    null = false
    type = character_varying(16)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "status" {
    null = false
    type = character_varying(32)
  }
  column "summary" {
    null = false
    type = character_varying(255)
  }
  column "details" {
    null = false
    type = character_varying(255)
  }
  column "contact_name" {
    null = false
    type = character_varying(128)
  }
  column "email_address" {
    null = true
    type = character_varying(128)
  }
  column "phone_number" {
    null = true
    type = character_varying(128)
  }
  column "provider_code" {
    null = true
    type = character_varying(64)
  }
  column "tracking_code" {
    null = true
    type = character_varying(128)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_deliveries_regions_region_id" {
    columns     = [column.region_id]
    ref_columns = [table.regions.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_deliveries_provider_code_tracking_code" {
    columns = [column.provider_code, column.tracking_code]
  }
  index "ix_deliveries_region_id" {
    columns = [column.region_id]
  }
}
table "devices" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "user_id" {
    null = true
    type = character_varying(16)
  }
  column "client_id" {
    null = false
    type = character_varying(16)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "trace_code" {
    null = false
    type = character_varying(64)
  }
  column "push_token" {
    null = true
    type = character_varying(128)
  }
  column "metadata" {
    null = true
    type = jsonb
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_devices_users_client_id" {
    columns     = [column.client_id]
    ref_columns = [table.clients.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_devices_users_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_devices_trace_code" {
    unique  = true
    columns = [column.trace_code]
  }
}
table "inventories" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "region_id" {
    null = true
    type = character_varying(16)
  }
  column "product_sku_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "valid_time" {
    null = true
    type = timestamp
  }
  column "weight" {
    null = false
    type = integer
  }
  column "quantity" {
    null = false
    type = integer
  }
  column "comment" {
    null = true
    type = character_varying(128)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_inventories_product_skus_product_sku_id" {
    columns     = [column.product_sku_id]
    ref_columns = [table.product_skus.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_inventories_regions_region_id" {
    columns     = [column.region_id]
    ref_columns = [table.regions.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_inventories_product_sku_id" {
    columns = [column.product_sku_id]
  }
  index "ix_inventories_region_id" {
    columns = [column.region_id]
  }
}
table "logins" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "user_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "provider" {
    null = false
    type = character_varying(16)
  }
  column "identifier" {
    null = false
    type = character_varying(64)
  }
  column "credential" {
    null = true
    type = character_varying(64)
  }
  column "metadata" {
    null = true
    type = jsonb
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_logins_users_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_logins_provider_identifier" {
    unique  = true
    columns = [column.provider, column.identifier]
  }
}
table "order_items" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "order_id" {
    null = false
    type = character_varying(16)
  }
  column "product_sku_id" {
    null = false
    type = character_varying(16)
  }
  column "delivery_id" {
    null = true
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "quantity" {
    null = false
    type = integer
  }
  column "currency" {
    null = false
    type = character_varying(32)
  }
  column "unit_amount" {
    null = false
    type = numeric(16,6)
  }
  column "comment" {
    null = true
    type = character_varying(128)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_order_items_deliveries_delivery_id" {
    columns     = [column.delivery_id]
    ref_columns = [table.deliveries.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_order_items_orders_order_id" {
    columns     = [column.order_id]
    ref_columns = [table.orders.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_order_items_product_skus_product_sku_id" {
    columns     = [column.product_sku_id]
    ref_columns = [table.product_skus.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_order_items_delivery_id" {
    columns = [column.delivery_id]
  }
  index "ix_order_items_order_id" {
    columns = [column.order_id]
  }
  index "ix_order_items_product_sku_id" {
    columns = [column.product_sku_id]
  }
}
table "order_payments" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "order_id" {
    null = false
    type = character_varying(16)
  }
  column "payment_id" {
    null = false
    type = character_varying(16)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "amount" {
    null = false
    type = numeric(16,6)
  }
  column "status" {
    null = false
    type = character_varying(32)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_order_payments_orders_order_id" {
    columns     = [column.order_id]
    ref_columns = [table.orders.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_order_payments_payments_payment_id" {
    columns     = [column.payment_id]
    ref_columns = [table.payments.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_order_payments_order_id" {
    columns = [column.order_id]
  }
  index "ix_order_payments_payment_id" {
    columns = [column.payment_id]
  }
}
table "orders" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "shop_id" {
    null = false
    type = character_varying(16)
  }
  column "owner_id" {
    null = false
    type = character_varying(16)
  }
  column "address_id" {
    null = true
    type = character_varying(16)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "status" {
    null = false
    type = character_varying(32)
  }
  column "currency" {
    null = false
    type = character_varying(32)
  }
  column "total_amount" {
    null = false
    type = numeric(16,6)
  }
  column "serial_code" {
    null = false
    type = character_varying(64)
  }
  column "comment" {
    null = true
    type = character_varying(128)
  }
  column "snapshot" {
    null = false
    type = jsonb
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_orders_addresses_address_id" {
    columns     = [column.address_id]
    ref_columns = [table.addresses.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_orders_shops_shop_id" {
    columns     = [column.shop_id]
    ref_columns = [table.shops.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_orders_users_owner_id" {
    columns     = [column.owner_id]
    ref_columns = [table.participants.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_orders_address_id" {
    columns = [column.address_id]
  }
  index "ix_orders_owner_id" {
    columns = [column.owner_id]
  }
  index "ix_orders_serial_code" {
    unique  = true
    columns = [column.serial_code]
  }
  index "ix_orders_shop_id" {
    columns = [column.shop_id]
  }
}
table "participants" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "owner_id" {
    null = true
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "common_name" {
    null = false
    type = character_varying(128)
  }
  column "short_name" {
    null = false
    type = character_varying(64)
  }
  column "license_type" {
    null = false
    type = character_varying(32)
  }
  column "license_code" {
    null = true
    type = character_varying(64)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_participants_users_owner_id" {
    columns     = [column.owner_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_participants_owner_id" {
    columns = [column.owner_id]
  }
}
table "payments" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "owner_id" {
    null = false
    type = character_varying(16)
  }
  column "payee_id" {
    null = false
    type = character_varying(16)
  }
  column "payer_id" {
    null = true
    type = character_varying(16)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "valid_time" {
    null = true
    type = timestamp
  }
  column "serial_code" {
    null = false
    type = character_varying(128)
  }
  column "currency" {
    null = false
    type = character_varying(32)
  }
  column "amount" {
    null = false
    type = numeric(16,6)
  }
  column "flags" {
    null = false
    type = bigint
  }
  column "status" {
    null = false
    type = character_varying(32)
  }
  column "title" {
    null = false
    type = character_varying(64)
  }
  column "comment" {
    null = true
    type = character_varying(128)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_payments_participants_owner_id" {
    columns     = [column.owner_id]
    ref_columns = [table.participants.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_payments_participants_payee_id" {
    columns     = [column.payee_id]
    ref_columns = [table.participants.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_payments_participants_payer_id" {
    columns     = [column.payer_id]
    ref_columns = [table.participants.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_payments_owner_id" {
    columns = [column.owner_id]
  }
  index "ix_payments_payee_id" {
    columns = [column.payee_id]
  }
  index "ix_payments_payer_id" {
    columns = [column.payer_id]
  }
  index "ix_payments_serial_code" {
    unique  = true
    columns = [column.serial_code]
  }
}
table "prices" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "product_sku_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "obsolete" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "valid_time" {
    null = true
    type = timestamp
  }
  column "weight" {
    null = false
    type = integer
  }
  column "currency" {
    null = false
    type = character_varying(32)
  }
  column "amount" {
    null = false
    type = numeric(16,6)
  }
  column "comment" {
    null = true
    type = character_varying(128)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_prices_product_skus_product_sku_id" {
    columns     = [column.product_sku_id]
    ref_columns = [table.product_skus.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_prices_product_sku_id" {
    columns = [column.product_sku_id]
  }
}
table "product_skus" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "product_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "valid_time" {
    null = true
    type = timestamp
  }
  column "name" {
    null = false
    type = character_varying(64)
  }
  column "title" {
    null = false
    type = character_varying(64)
  }
  column "description" {
    null = true
    type = character_varying(255)
  }
  column "template" {
    null = false
    type = jsonb
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_product_skus_products_product_id" {
    columns     = [column.product_id]
    ref_columns = [table.products.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_product_skus_product_id_name" {
    unique  = true
    columns = [column.product_id, column.name]
  }
}
table "products" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "catalog_id" {
    null = false
    type = character_varying(16)
  }
  column "shop_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "valid_time" {
    null = true
    type = timestamp
  }
  column "name" {
    null = false
    type = character_varying(64)
  }
  column "title" {
    null = false
    type = character_varying(64)
  }
  column "description" {
    null = true
    type = character_varying(255)
  }
  column "template" {
    null = false
    type = jsonb
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_products_catalogs_catalog_id" {
    columns     = [column.catalog_id]
    ref_columns = [table.catalogs.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_products_shops_shop_id" {
    columns     = [column.shop_id]
    ref_columns = [table.shops.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_products_catalog_id" {
    columns = [column.catalog_id]
  }
  index "ix_products_shop_id_name" {
    unique  = true
    columns = [column.shop_id, column.name]
  }
}
table "profiles" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "display_name" {
    null = false
    type = character_varying(64)
  }
  column "avatar_url" {
    null = true
    type = character_varying(2048)
  }
  column "gender" {
    null = true
    type = character_varying(32)
  }
  column "birthdate" {
    null = true
    type = date
  }
  column "introduction" {
    null = true
    type = character_varying(128)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_profiles_users_id" {
    columns     = [column.id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
}
table "realms" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "flags" {
    null = false
    type = bigint
  }
  column "name" {
    null = false
    type = character_varying(64)
  }
  column "title" {
    null = false
    type = character_varying(64)
  }
  column "description" {
    null = true
    type = character_varying(255)
  }
  primary_key {
    columns = [column.id]
  }
  index "ix_realms_name" {
    unique  = true
    columns = [column.name]
  }
}
table "regions" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "parent_id" {
    null = true
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "common_code" {
    null = false
    type = character_varying(64)
  }
  column "common_name" {
    null = false
    type = character_varying(128)
  }
  column "summary" {
    null = true
    type = character_varying(255)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_regions_regions_parent_id" {
    columns     = [column.parent_id]
    ref_columns = [table.regions.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_regions_common_code" {
    unique  = true
    columns = [column.common_code]
  }
  index "ix_regions_parent_id_common_name" {
    unique  = true
    columns = [column.parent_id, column.common_name]
  }
}
table "reviews" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "product_sku_id" {
    null = false
    type = character_varying(16)
  }
  column "author_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "is_anonymous" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "title" {
    null = true
    type = character_varying(64)
  }
  column "comment" {
    null = true
    type = character_varying(255)
  }
  column "rating" {
    null = false
    type = smallint
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_reviews_product_skus_product_sku_id" {
    columns     = [column.product_sku_id]
    ref_columns = [table.product_skus.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_reviews_users_author_id" {
    columns     = [column.author_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_reviews_author_id" {
    columns = [column.author_id]
  }
  index "ix_reviews_product_sku_id" {
    columns = [column.product_sku_id]
  }
}
table "role_users" {
  schema = schema.public
  column "role_id" {
    null = false
    type = character_varying(16)
  }
  column "user_id" {
    null = false
    type = character_varying(16)
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  primary_key {
    columns = [column.role_id, column.user_id]
  }
  foreign_key "fk_role_users_roles_role_id" {
    columns     = [column.role_id]
    ref_columns = [table.roles.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_role_users_users_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
}
table "roles" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "realm_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "name" {
    null = false
    type = character_varying(64)
  }
  column "description" {
    null = true
    type = character_varying(255)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_roles_realms_realm_id" {
    columns     = [column.realm_id]
    ref_columns = [table.realms.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_roles_realm_id_name" {
    unique  = true
    columns = [column.realm_id, column.name]
  }
}
table "shops" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "holder_id" {
    null = false
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = false
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "valid_time" {
    null = true
    type = timestamp
  }
  column "common_name" {
    null = false
    type = character_varying(128)
  }
  column "short_name" {
    null = true
    type = character_varying(64)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_shops_participants_holder_id" {
    columns     = [column.holder_id]
    ref_columns = [table.participants.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_shops_common_name" {
    unique  = true
    columns = [column.common_name]
  }
  index "ix_shops_holder_id" {
    columns = [column.holder_id]
  }
  index "ix_shops_short_name" {
    unique  = true
    columns = [column.short_name]
  }
}
table "user_devices" {
  schema = schema.public
  column "user_id" {
    null = false
    type = character_varying(16)
  }
  column "device_id" {
    null = false
    type = character_varying(16)
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  primary_key {
    columns = [column.user_id, column.device_id]
  }
  foreign_key "fk_user_devices_devices_device_id" {
    columns     = [column.device_id]
    ref_columns = [table.devices.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_user_devices_users_user_id" {
    columns     = [column.user_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
}
table "users" {
  schema = schema.public
  column "id" {
    null = false
    type = character_varying(16)
  }
  column "realm_id" {
    null = false
    type = character_varying(16)
  }
  column "creator_id" {
    null = true
    type = character_varying(16)
  }
  column "disabled" {
    null    = false
    type    = boolean
    default = false
  }
  column "approved" {
    null    = false
    type    = boolean
    default = false
  }
  column "verified" {
    null    = false
    type    = boolean
    default = false
  }
  column "immutable" {
    null    = false
    type    = boolean
    default = false
  }
  column "created_at" {
    null = false
    type = timestamp
  }
  column "updated_at" {
    null = true
    type = timestamp
  }
  column "deleted_at" {
    null = true
    type = timestamp
  }
  column "expires_at" {
    null = true
    type = timestamp
  }
  column "first_login_time" {
    null = true
    type = timestamp
  }
  column "last_active_time" {
    null = true
    type = timestamp
  }
  column "flags" {
    null = false
    type = bigint
  }
  column "attributes" {
    null = true
    type = jsonb
  }
  column "phone_number" {
    null = true
    type = character_varying(64)
  }
  column "email_address" {
    null = true
    type = character_varying(128)
  }
  column "description" {
    null = true
    type = character_varying(255)
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "fk_users_realms_realm_id" {
    columns     = [column.realm_id]
    ref_columns = [table.realms.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  foreign_key "fk_users_users_creator_id" {
    columns     = [column.creator_id]
    ref_columns = [table.users.column.id]
    on_update   = RESTRICT
    on_delete   = RESTRICT
  }
  index "ix_users_realm_id_email_address" {
    unique  = true
    columns = [column.realm_id, column.email_address]
  }
  index "ix_users_realm_id_phone_number" {
    unique  = true
    columns = [column.realm_id, column.phone_number]
  }
}
schema "public" {
  comment = "standard public schema"
}
