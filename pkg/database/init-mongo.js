db.createUser({
    user: "ayocodedb",
    pwd: "secret",
    roles: [
        { role: "readWrite", db: "acourse" }
    ],
    mechanisms: ["<SCRAM-SHA-1|SCRAM-SHA-256>"],
})