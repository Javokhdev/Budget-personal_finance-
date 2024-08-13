# Olympy

**Project Description:**
Loyiha maqsadi:
Olimpiada o'yinlari uchun real vaqtda yangilanadigan medallar reytingi va musobaqa jarayonlarini live text shaklida ko'rsatadigan tizim yaratish. musobaqa jarayonlarini live text shaklida degani quyidagicha: Masalan Dzyudo musoboqasi bo’layotgan bo’lsa, olingan ochko, sariq kartochka yoki boshqalarni kuzatib turish imkoniyati. Yoki, futbol musobaqasida kiritilgan gol, kartochkalar, o’yinchi almashishi va hakozo kabi ma’lumotlar.

## Main Functionality:
1. Medal rating
2. Live text streaming
3. Microservices
    - communication between services through gRPC
    - Protocol Buffers for data serialization
    - Databases:
    - PostgreSQL (for basic data)
    - MongoDB (for live streaming data)
    - Redis (for caching)
    - Gin framework (for REST API)
    - JWT authentication
    - Casbin (for roles and authorization)
    - Kafka or RabbitMQ (for messaging)
    - Docker (for containerization)
    - API Gateway
    - Swagger (for API documentation)
    - Unit testing
    - Database migration
    - Submodule for protos
4. API endpoints
5. Data model
6. Security
7. Testing
8. Documentation

## Installation

1. Initialize a git repository and clone the project:
    ```sh
    git init
    git clone https://gitlab.com/olympy1
    ```
2. Create a database named `olympy` on port `5432`.
3. Update the `.env` file with the appropriate configuration.
   ```.env
   DB_HOST=localhost
   DB_USER=postgres
   DB_NAME=olympy
   DB_PASSWORD=pass
   DB_PORT=5432
   LOGPATH=logs/info.log
   ```

4. Use the following Makefile commands to manage the database migrations and set up the project:
    ```makefile
    # Set the database URL
    exp:
        export DBURL="postgres://postgres:1234@localhost:5432/olympy?sslmode=disable'"

    # Run migrations
    mig-up:
        migrate -path migrations -database ${DBURL} -verbose up

    # Rollback migrations
    mig-down:
        migrate -path migrations -database ${DBURL} -verbose down

    # Create a new migration
    mig-create:
        migrate create -ext sql -dir migrations -seq create_table

    # Create an insert migration
    mig-insert:
        migrate create -ext sql -dir migrations -seq insert_table

    # Generate Swagger documentation
    swag:
        swag init -g api/api.go -o api/docs

    # Clean up migrations (commented out by default)
    # mig-delete:
    #   rm -r db/migrations
    ```
5. Set the environment variable and run the project:
    ```sh
    make exp
    make mig-up
    go run main.go
    ```
6. Open the following URL to access the Swagger documentation:
    ```
    http://localhost:8088/swagger/index.html#/
    ```

## Features and Usages
Loyiha qo'shimcha maqsadlari:

1. Statistik tahlillar: Musobaqa davomida va yakunida sportchilar va jamoalar haqidagi keng qamrovli statistik ma'lumotlarni taqdim etish. Bu jamoalar va sportchilar o'rtasidagi solishtirishlarni, individual natijalarni va rekordlarni ko'rsatishni o'z ichiga oladi.

2. Foydalanuvchi tajribasini oshirish: Foydalanuvchilar uchun qulay va intuitiv interfeys yaratish, u orqali foydalanuvchilar osongina kerakli ma'lumotlarni topishi va kuzatishi mumkin. Shuningdek,yangilanishlarni real vaqtda olish imkoniyati.

3. Foydalanuvchilar o'rtasidagi muloqot: Foydalanuvchilar o'zaro musobaqalar haqida fikr almashishi va sharh qoldirishi mumkin bo'lgan forum yoki chat funksiyalarini borligi.

4. Xavfsizlik va maxfiylik: Foydalanuvchilarning shaxsiy ma'lumotlarini himoya qilish va tizimning xavfsizligini ta'minlash, shu jumladan kirish nazorati va ma'lumotlarni shifrlash.

5. Ma'lumotlar omborini boshqarish: Katta hajmdagi ma'lumotlarni samarali boshqarish va saqlash, shu jumladan arxivlash va qayta ishlash.

## Texnologiyalar

Loyiha quyidagi texnologiyalar yordamida amalga oshiriladi:
- **Backend**: Golang
- **Ma'lumotlar bazasi**: PostgreSQL yoki MongoDB
- **Real vaqt yangilanishi**: WebSocket 
- **Autentifikatsiya va avtorizatsiya**: JWT 


## Dependencies

- **Scheduling**: [github.com/go-co-op/gocron](https://github.com/go-co-op/gocron)
- **Swagger**: [github.com/swaggo/swag](https://github.com/swaggo/swag)
- **Database**:
    - [database/sql](https://golang.org/pkg/database/sql/)
    - [github.com/lib/pq](https://github.com/lib/pq)
- **Environment Variables**: [github.com/joho/godotenv](https://github.com/joho/godotenv)
- **API Framework**: [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- **Migrations**: [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate)
****
## Acknowledgments

- Mubina Bahodirova
- Feruza Mirjalilova
- Javohir Xasanov

## Special thanks to HUSAN MUSA