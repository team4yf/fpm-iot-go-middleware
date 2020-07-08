INSERT INTO "fim_application" ("id", "created_at", "updated_at", "deleted_at", "app_id", "name", "status", "active_at", "company", "contact", "email", "mobile", "active_code", "home_page", "health", "description", "deploy_mode") VALUES
(1,	NULL,	NULL,	NULL,	'ceaa191a',	'Dev App',	1,	'2020-07-08 10:16:55.409488+08',	'YunPlus.IO',	'Wang3',	'support@yunplus.io',	'13770683580',	'abcd123',	'http://open.yunplus.io:19501',	'http://open.yunplus.io:19501/health',	'YunPlus.IO Core Project',	'Multi');


INSERT INTO "fim_project" ("id", "created_at", "updated_at", "deleted_at", "app_id", "name", "status", "project_id", "code", "entry_url", "setting") VALUES
(1,	NULL,	NULL,	NULL,	'ceaa191a',	'Demo Project',	1,	1,	'demo',	'http://open.yunplus.io:19501/demo',	'{"light":{"brand":"lt10","appid":"LT0314fbf27a4d2986"}}');


INSERT INTO "fim_client" ("id", "created_at", "updated_at", "deleted_at", "app_key", "secret_key", "expired", "name", "api_base_url", "environment", "enable_ssl", "cert_path", "type", "status") VALUES
(1,	NULL,	NULL,	NULL,	'LT0314fbf27a4d2986',	'1bc7b874c74623298a6',	3600,	'瓴泰科技智慧路灯 API 接口',	'http://101.132.142.5:8088/api',	'prod',	NULL,	NULL,	'light',	1);
