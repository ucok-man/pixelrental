-- Drop foreign key constraints
ALTER TABLE "payments" DROP CONSTRAINT IF EXISTS "payments_order_id_fkey";
ALTER TABLE "order_details" DROP CONSTRAINT IF EXISTS "order_details_order_id_fkey";
ALTER TABLE "order_details" DROP CONSTRAINT IF EXISTS "order_details_game_id_fkey";
ALTER TABLE "tokens" DROP CONSTRAINT IF EXISTS "tokens_user_id_fkey";
ALTER TABLE "carts" DROP CONSTRAINT IF EXISTS "carts_user_id_fkey";
ALTER TABLE "carts" DROP CONSTRAINT IF EXISTS "carts_game_id_fkey";
ALTER TABLE "orders" DROP CONSTRAINT IF EXISTS "orders_user_id_fkey";

-- Drop tables
DROP TABLE IF EXISTS "payments";
DROP TABLE IF EXISTS "order_details";
DROP TABLE IF EXISTS "tokens";
DROP TABLE IF EXISTS "carts";
DROP TABLE IF EXISTS "orders";
DROP TABLE IF EXISTS "games";
DROP TABLE IF EXISTS "users";
