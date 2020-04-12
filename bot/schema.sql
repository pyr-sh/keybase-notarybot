CREATE TABLE "signatures" (
    "id" TEXT NOT NULL PRIMARY KEY,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL,
    "username" TEXT NOT NULL,
    "name" TEXT NOT NULL,
    "file_url" TEXT NOT NULL,
    "line_pos" DECIMAL(7, 4),

    CHECK(EXTRACT(TIMEZONE FROM "created_at") = '0'),
    CHECK(EXTRACT(TIMEZONE FROM "updated_at") = '0')
);
